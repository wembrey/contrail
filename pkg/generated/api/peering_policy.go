package api

import (
	"database/sql"
	"net/http"

	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/generated/db"
	"github.com/Juniper/contrail/pkg/generated/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
)

//PeeringPolicyRESTAPI
type PeeringPolicyRESTAPI struct {
	DB *sql.DB
}

type PeeringPolicyCreateRequest struct {
	Data *models.PeeringPolicy `json:"peering-policy"`
}

//Path returns api path for collections.
func (api *PeeringPolicyRESTAPI) Path() string {
	return "/peering-policys"
}

//LongPath returns api path for elements.
func (api *PeeringPolicyRESTAPI) LongPath() string {
	return "/peering-policy/:id"
}

//SetDB sets db object
func (api *PeeringPolicyRESTAPI) SetDB(db *sql.DB) {
	api.DB = db
}

//Create handle a Create REST API.
func (api *PeeringPolicyRESTAPI) Create(c echo.Context) error {
	requestData := &PeeringPolicyCreateRequest{
		Data: models.MakePeeringPolicy(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "peering_policy",
		}).Debug("bind failed on create")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	model := requestData.Data
	if model == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	if model.UUID == "" {
		model.UUID = uuid.NewV4().String()
	}
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.CreatePeeringPolicy(tx, model)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "peering_policy",
		}).Debug("db create failed on create")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusCreated, requestData)
}

//Update handles a REST Update request.
func (api *PeeringPolicyRESTAPI) Update(c echo.Context) error {
	return nil
}

//Delete handles a REST Delete request.
func (api *PeeringPolicyRESTAPI) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			return db.DeletePeeringPolicy(tx, id)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//Show handles a REST Show request.
func (api *PeeringPolicyRESTAPI) Show(c echo.Context) error {
	id := c.Param("id")
	var result *models.PeeringPolicy
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ShowPeeringPolicy(tx, id)
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"peering_policy": result,
	})
}

//List handles a List REST API Request.
func (api *PeeringPolicyRESTAPI) List(c echo.Context) error {
	var result []*models.PeeringPolicy
	var err error
	if err := common.DoInTransaction(
		api.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListPeeringPolicy(tx, &common.ListSpec{
				Limit: 1000,
			})
			return err
		}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"peering-policys": result,
	})
}