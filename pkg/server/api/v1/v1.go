package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yrfg/boast/errs"
	"github.com/yrfg/boast/utils"
	"github.com/yrfg/timeman/pkg/store"
)

func Bind(r *gin.Engine) {
	sub := r.Group("/v1")
	{
		tm := sub.Group("/timeman")
		{
			tm.POST("/create-task-map", CreateTimeManMap)
			tm.POST("/update-task-map-name", UpdateTimeManMapName)
			tm.GET("/get-task-map-list", ListTimeManMaps)
			tm.POST("/create-task",CreateTask)
			tm.POST("/get-task-list", ListTasks)
		}
	}
}

func ListTasks(c *gin.Context) {
	var (
		filterTag string
		taskMapID int64
		rTaskMapID string
		count int
		offset int
		respBody struct {
			Data []store.TimeManTaskDisplay `json:"data"`
		}
	)
	page, pageSize, ok := utils.PPInQuery(c.Request)
	if !ok {
		errs.WriteDis(c.Writer, http.StatusBadRequest, errs.RequestParseError, "page error")
		return
	}
	count, offset = utils.PPToLF(page, pageSize)
	filterTag = c.Query("filter_tag")
	rTaskMapID = c.Query("task_map_id")
	taskMapID,_ = strconv.ParseInt(rTaskMapID, 10,64)
	tasks := store.ListTimeManTaskByFilterTag(context.Background(), store.DefaultConn,taskMapID, filterTag,offset, count )
	respBody.Data = tasks
	c.JSON(http.StatusOK, respBody)
}

func CreateTask(c *gin.Context) {
	var (
		body struct {
			Name string `json:"name"`
			TaskMapID int64 `json:"task_map_id"`
		}
		respBody struct {
			ID int64 `json:"id"`
		}
		err         error
	)
	err = c.ShouldBind(&body)
	if err != nil {
		errs.WriteDis(c.Writer,http.StatusBadRequest,errs.RequestParseError, err.Error() )
		return
	}

	respBody.ID, err = store.CreateTask(context.Background(),store.DefaultConn,body.TaskMapID, body.Name)
	if err != nil {
		errs.WriteDis(c.Writer,http.StatusBadRequest,errs.RequestError, err.Error() )
		return
	}
	c.JSON(http.StatusOK, respBody)
}

func ListTimeManMaps(c *gin.Context) {
	var (
		respBody struct {
			Data []store.TimeManTaskMapDisplay `json:"data"`
		}
		offset int
		count  int
	)
	page, pageSize, ok := utils.PPInQuery(c.Request)
	if !ok {
		errs.WriteDis(c.Writer,http.StatusBadRequest,errs.RequestParseError, "page error" )
		return
	}
	count, offset = utils.PPToLF(page, pageSize)
	list := store.ListTimeManTaskMap(context.Background(),store.DefaultConn, offset, count)
	respBody.Data = list

	c.JSON(http.StatusOK, respBody)
}

func UpdateTimeManMapName(c *gin.Context) {
	var (
		body struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
		respBody struct {
		}
		err         error
	)
	err = c.ShouldBind(&body)
	if err != nil {
		errs.WriteDis(c.Writer,http.StatusBadRequest,errs.RequestParseError, err.Error() )
		return
	}

	err = store.UpdateTimeManTaskMapName(context.Background(),store.DefaultConn, body.ID, body.Name)
	if err != nil {
		errs.WriteDis(c.Writer,http.StatusBadRequest,errs.RequestError, err.Error() )
		return
	}
	c.JSON(http.StatusOK, respBody)
}

func CreateTimeManMap(c *gin.Context) {
	var (
		body struct {
			Name string `json:"name"`
		}
		respBody struct {
			ID int64 `json:"id"`
		}
		err         error
		newID       int64
	)
	err = c.ShouldBind(&body)
	if err != nil {
		errs.WriteDis(c.Writer,http.StatusBadRequest,errs.RequestParseError, err.Error() )
		return
	}

	newID, err = store.CreateTimeManTaskMap(context.Background(),store.DefaultConn, body.Name)
	if err != nil {
		errs.WriteDis(c.Writer,http.StatusBadRequest,errs.RequestParseError, err.Error() )
		return
	}
	respBody.ID = newID
	c.JSON(http.StatusOK, respBody)
}
