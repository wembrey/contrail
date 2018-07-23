package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Juniper/contrail/pkg/apisrv/keystone"
)

const (
	retryCount = 2
)

// HTTP represents API Server HTTP client.
type HTTP struct {
	httpClient *http.Client

	ID        string          `yaml:"id"`
	Password  string          `yaml:"password"`
	AuthURL   string          `yaml:"authurl"`
	Endpoint  string          `yaml:"endpoint"`
	AuthToken string          `yaml:"-"`
	Domain    string          `yaml:"domain"`
	InSecure  bool            `yaml:"insecure"`
	Debug     bool            `yaml:"debug"`
	Scope     *keystone.Scope `yaml:"scope"`
}

// Request represents API request to the server.
type Request struct {
	Method   string      `yaml:"method"`
	Path     string      `yaml:"path,omitempty"`
	Expected []int       `yaml:"expected,omitempty"`
	Data     interface{} `yaml:"data,omitempty"`
	Output   interface{} `yaml:"output,omitempty"`
}

// NewHTTP makes API Server HTTP client.
func NewHTTP(endpoint, authURL, id, password, domain string, insecure bool, scope *keystone.Scope) *HTTP {
	c := &HTTP{
		ID:       id,
		Password: password,
		AuthURL:  authURL,
		Endpoint: endpoint,
		Scope:    scope,
		Domain:   domain,
		InSecure: insecure,
	}
	c.Init()
	return c
}

//Init is used to initialize a client.
func (h *HTTP) Init() {
	tr := &http.Transport{
		Dial: (&net.Dialer{
			//Timeout: 5 * time.Second,
		}).Dial,
		//TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: h.InSecure},
	}
	client := &http.Client{
		Transport: tr,
		//Timeout:   time.Second * 10,
	}
	h.httpClient = client
}

// Login refreshes authentication token.
func (h *HTTP) Login() error {
	if h.AuthURL == "" {
		return nil
	}
	authURL := h.AuthURL + "/auth/tokens"
	authRequest := &keystone.AuthRequest{
		Auth: &keystone.Auth{
			Identity: &keystone.Identity{
				Methods: []string{"password"},
				Password: &keystone.Password{
					User: &keystone.User{
						Name:     h.ID,
						Password: h.Password,
						Domain: &keystone.Domain{
							ID: h.Domain,
						},
					},
				},
			},
			Scope: h.Scope,
		},
	}
	authResponse := &keystone.AuthResponse{}
	dataJSON, err := json.Marshal(authRequest)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", authURL, bytes.NewBuffer(dataJSON))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := h.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint: errcheck
	err = checkStatusCode([]int{201}, resp.StatusCode)
	if err != nil {
		output, _ := httputil.DumpResponse(resp, true) // nolint: gas
		log.WithError(err).WithField("output", string(output)).Error("Unexpected status code")
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return err
	}
	h.AuthToken = resp.Header.Get("X-Subject-Token")
	return nil
}

func checkStatusCode(expected []int, actual int) error {
	for _, expected := range expected {
		if expected == actual {
			return nil
		}
	}
	return errors.Errorf("unexpected return code: expected %v, actual %v", expected, actual)
}

// Create send a create API request.
func (h *HTTP) Create(path string, data interface{}, output interface{}) (*http.Response, error) {
	expected := []int{http.StatusCreated}
	return h.Do(echo.POST, path, data, output, expected)
}

// Read send a get API request.
func (h *HTTP) Read(path string, output interface{}) (*http.Response, error) {
	expected := []int{http.StatusOK}
	return h.Do(echo.GET, path, nil, output, expected)
}

// Update send an update API request.
func (h *HTTP) Update(path string, data interface{}, output interface{}) (*http.Response, error) {
	expected := []int{http.StatusOK}
	return h.Do(echo.PUT, path, data, output, expected)
}

// Delete send a delete API request.
func (h *HTTP) Delete(path string, output interface{}) (*http.Response, error) {
	expected := []int{http.StatusNoContent}
	return h.Do(echo.DELETE, path, nil, output, expected)
}

// EnsureDeleted send a delete API request.
func (h *HTTP) EnsureDeleted(path string, output interface{}) (*http.Response, error) {
	expected := []int{http.StatusNoContent, http.StatusNotFound}
	return h.Do(echo.DELETE, path, nil, output, expected)
}

// Do issues an API request.
func (h *HTTP) Do(method, path string, data interface{}, output interface{}, expected []int) (*http.Response, error) {
	request, err := h.prepareHTTPRequest(method, path, data)
	if err != nil {
		return nil, err
	}

	resp, err := h.doHTTPRequestRetryingOn401(request, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck

	err = checkStatusCode(expected, resp.StatusCode)
	if err != nil {
		output, _ := httputil.DumpResponse(resp, true) // nolint:  gas
		log.WithError(err).WithField("output", string(output)).Error("Unexpected status code")
		return resp, err
	}
	if method == echo.DELETE {
		return resp, nil
	}
	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return resp, errors.Wrap(err, "decoding response body failed")
	}
	if h.Debug {
		log.WithFields(log.Fields{
			"response": resp,
			"output":   output,
		}).Debug("API Server response")
	}
	return resp, err
}

func (h *HTTP) prepareHTTPRequest(method, path string, data interface{}) (*http.Request, error) {
	var request *http.Request
	if data == nil {
		var err error
		request, err = http.NewRequest(method, getURL(h.Endpoint, path), nil)
		if err != nil {
			return nil, errors.Wrap(err, "creating HTTP request failed")
		}
	} else {
		var dataJSON []byte
		dataJSON, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "encoding request data failed")
		}

		request, err = http.NewRequest(method, getURL(h.Endpoint, path), bytes.NewBuffer(dataJSON))
		if err != nil {
			return nil, errors.Wrap(err, "creating HTTP request failed")
		}
	}

	request.Header.Set("Content-Type", "application/json")
	if h.AuthToken != "" {
		request.Header.Set("X-Auth-Token", h.AuthToken)
	}
	return request, nil
}

func getURL(endpoint, path string) string {
	return endpoint + path
}

func (h *HTTP) doHTTPRequestRetryingOn401(request *http.Request, data interface{}) (*http.Response, error) {
	if h.Debug {
		log.WithFields(log.Fields{
			"method": request.Method,
			"url":    request.URL,
			"header": request.Header,
			"data":   data,
		}).Debug("Executing API Server request")
	}
	var resp *http.Response
	for i := 0; i < retryCount; i++ {
		var err error
		resp, err = h.httpClient.Do(request)
		if err != nil {
			return nil, errors.Wrap(err, "issuing HTTP request failed")
		}
		if resp.StatusCode != 401 {
			break
		}
		// token might be expired, refresh token and retry
		// skip refresh token after last retry
		if i < retryCount-1 {
			err = resp.Body.Close()
			if err != nil {
				return nil, errors.Wrap(err, "closing response body failed")
			}

			// refresh token and use the new token in request header
			err = h.Login()
			if err != nil {
				return nil, err
			}
			if h.AuthToken != "" {
				request.Header.Set("X-Auth-Token", h.AuthToken)
			}
		}
	}
	return resp, nil
}

// DoRequest requests based on request object.
func (h *HTTP) DoRequest(request *Request) (*http.Response, error) {
	return h.Do(request.Method, request.Path, request.Data, &request.Output, request.Expected)
}

// Batch execution.
func (h *HTTP) Batch(requests []*Request) error {
	for i, request := range requests {
		_, err := h.DoRequest(request)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("%dth request failed.", i))
		}
	}
	return nil
}