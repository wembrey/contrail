package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	integrationetcd "github.com/Juniper/contrail/pkg/testutil/integration/etcd"
)

// TODO(dfurman): move to ASF

const (
	dialTimeout      = 10 * time.Second
	shortDialTimeout = 10 * time.Millisecond
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		fails  bool
	}{
		{
			name: "succeeds when TLS disabled and correct credentials given",
			config: &Config{
				Config: *etcdConfig(dialTimeout),
			},
			fails: false,
		},
		{
			name: "fails when TLS enabled and no certificates given",
			config: &Config{
				Config: *etcdConfig(shortDialTimeout),
				TLSConfig: TLSConfig{
					Enabled: true,
				},
			},
			fails: true,
		},
		{
			name: "fails when TLS enabled invalid certificate paths given",
			config: &Config{
				Config: *etcdConfig(dialTimeout),
				TLSConfig: TLSConfig{
					Enabled:         true,
					CertificatePath: "invalid-path",
					KeyPath:         "invalid-path",
					TrustedCAPath:   "invalid-path",
				},
			},
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.config)

			if tt.fails {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				defer closeClient(t, c)
			}
		})
	}
}

func etcdConfig(dialTimeout time.Duration) *clientv3.Config {
	return &clientv3.Config{
		Endpoints:   []string{integrationetcd.Endpoint},
		DialTimeout: dialTimeout,
	}
}

func closeClient(t *testing.T, c *Client) {
	assert.NoError(t, c.Close())
}

func TestClient_DoInTransaction(t *testing.T) {
	testKey := "in_transaction_test_key"
	tests := []struct {
		name    string
		ctx     context.Context
		do      func(context.Context) error
		wantErr bool
	}{
		{
			name:    "transaction is already in context, function returns no error",
			ctx:     WithTxn(context.Background(), &stmTxn{}),
			do:      func(context.Context) error { return nil },
			wantErr: false,
		},
		{
			name:    "transaction is already in context, function returns error",
			ctx:     WithTxn(context.Background(), &stmTxn{}),
			do:      func(context.Context) error { return assert.AnError },
			wantErr: true,
		},
		{
			name: "get the key twice",
			ctx:  context.Background(),
			do: func(ctx context.Context) error {
				txn := GetTxn(ctx)

				txn.Put(testKey, []byte("some value"))
				v1 := txn.Get(testKey)
				if string(v1) != "some value" {
					return errors.New("value should be updated")
				}

				txn.Put(testKey, []byte("newer value"))
				v2 := txn.Get(testKey)
				if string(v2) != "newer value" {
					return errors.New("value should be updated again")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(&Config{
				Config:      *etcdConfig(shortDialTimeout),
				ServiceName: t.Name(),
			})
			require.NoError(t, err)
			defer func() {
				err = c.Delete(context.Background(), testKey)
				require.NoError(t, err)
			}()

			err = c.DoInTransaction(tt.ctx, tt.do)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
