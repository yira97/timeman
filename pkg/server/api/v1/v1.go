package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yrfg/timeman/pkg/store"
)

func Bind(r *gin.Engine) {
	sub := r.Group("/v1")
	{
		tm := sub.Group("/timeman")
		{
			tm.POST("/create-timeline-map", CreateTimelineMap)
		}
	}
}

func CreateTimelineMap(c *gin.Context) {
	var (
		body struct {
			Name string `json:"name"`
		}
		respBody struct {
			ID int64 `json:"id"`
		}
		respErrBody struct {
			Err string `json:"err"`
		}
		err   error
		newID int64
	)
	err = c.ShouldBind(&body)
	if err != nil {
		respErrBody.Err = err.Error()
		c.JSON(http.StatusBadRequest, respErrBody)
		return
	}

	newID, err = store.CreateTimeManMap(store.DefaultConn, body.Name)
	if err != nil {
		respErrBody.Err = err.Error()
		c.JSON(http.StatusBadRequest, respErrBody)
		return
	}
	respBody.ID = newID
	c.JSON(http.StatusOK, respBody)
}
