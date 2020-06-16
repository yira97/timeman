package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	v1 "github.com/yrfg/timeman/pkg/server/api/v1"
)

const (
	ServerSetupConfigModeDebug   = "debug"
	ServerSetupConfigModeRelease = "release"
)

var (
	r *gin.Engine
)

type ServerSetupConfig struct {
	Mode string
	Port int64
}

func Setup(setupCfg ServerSetupConfig) {
	var err error
	gin.SetMode(setupCfg.Mode)
	r = gin.Default()

	v1.Bind(r)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", setupCfg.Port),
		Handler: r,
	}
	go func() {
		err = s.ListenAndServe()
		if err != nil {
			os.Exit((3))
		}
	}()
}
