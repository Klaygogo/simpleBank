package api

import (
	"os"
	"testing"

	db "github.com/Klaygogo/simplebank/db/sqlc"
	util "github.com/Klaygogo/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"time"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSyncKey:        util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}
