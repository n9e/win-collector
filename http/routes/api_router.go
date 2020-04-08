package routes

import (
	"fmt"
	"os"

	"github.com/n9e/win-collector/stra"
	"github.com/n9e/win-collector/sys/funcs"

	"github.com/n9e/win-collector/sys/identity"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/didi/nightingale/src/toolkits/http/render"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/errors"
	"github.com/toolkits/pkg/logger"
)

func ping(c *gin.Context) {
	c.String(200, "pong")
}

func addr(c *gin.Context) {
	c.String(200, c.Request.RemoteAddr)
}

func pid(c *gin.Context) {
	c.String(200, fmt.Sprintf("%d", os.Getpid()))
}

func pushData(c *gin.Context) {
	if c.Request.ContentLength == 0 {
		render.Message(c, "blank body")
		return
	}

	recvMetricValues := []*dataobj.MetricValue{}
	metricValues := []*dataobj.MetricValue{}

	errors.Dangerous(c.ShouldBind(&recvMetricValues))

	var msg string
	for _, v := range recvMetricValues {
		logger.Debug("->recv: ", v)
		if v.Endpoint == "" {
			v.Endpoint = identity.Identity
		}
		err := v.CheckValidity()
		if err != nil {
			msg += fmt.Sprintf("recv metric %v err:%v\n", v, err)
			logger.Warningf(msg)
			continue
		}
		metricValues = append(metricValues, v)
	}

	funcs.Push(metricValues)

	if msg != "" {
		render.Message(c, msg)
		return
	}

	render.Data(c, "ok", nil)
	return
}

func getStrategy(c *gin.Context) {
	var resp []interface{}

	port := stra.GetPortCollects()
	for _, stra := range port {
		resp = append(resp, stra)
	}

	proc := stra.GetProcCollects()
	for _, stra := range proc {
		resp = append(resp, stra)
	}

	render.Data(c, resp, nil)
}
