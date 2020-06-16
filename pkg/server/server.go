package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "github.com/yrfg/timeman/pkg/server/api/v1"
	"log"
	"net/http"
)

type ServerMode int

const (
	Debug ServerMode = iota
	Release
)

func (m ServerMode) String() string {
	return [...]string{"debug", "release"}[m]
}

var (
	r *gin.Engine
)

type ServerSetupConfig struct {
	Mode ServerMode
	Port int
}

func Setup(setupCfg ServerSetupConfig) {
	var err error
	gin.SetMode(setupCfg.Mode.String())
	r = gin.Default()

	v1.Bind(r)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", setupCfg.Port),
		Handler: r,
	}
	go func() {
		err = s.ListenAndServe()
		if err != nil {
			log.Fatalf("listen error: %v",err)
		}
	}()
}
