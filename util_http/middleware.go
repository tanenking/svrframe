package util_http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tanenking/svrframe/constants"
	"github.com/tanenking/svrframe/helper"
	"github.com/tanenking/svrframe/logx"

	"github.com/gin-gonic/gin"
)

type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func newResponseWriter(c *gin.Context) *customResponseWriter {
	return &customResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
}

func (w customResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w customResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

/*
	func FailedResponse(c *gin.Context, err error) {
		appErr, ok := err.(HTTP_ERROR)
		if !ok {
			appErr = ERR_CODE_FAILED
		}
		c.Set("error_code", appErr.GetCode())
		resp := gin.H{"error_code": appErr.GetCode(), "msg": appErr.Error(), "time_unix": helper.GetNowTimestamp()}
		c.JSON(appErr.GetHttpCode(), resp)
	}
*/
func FailedResponse(c *gin.Context, code int) {
	c.Set("error_code", code)
	//resp := gin.H{"error_code": code, "msg": cfgtable.GetErrorMsg(code), "time_unix": helper.GetNowTimestamp()}
	resp := gin.H{"error_code": code, "time_unix": helper.GetNowTimestamp()}
	c.JSON(http.StatusOK, resp)
}

func SuccessResponse(c *gin.Context, data interface{}) {
	resp := map[string]interface{}{
		"error_code": 0,
		"msg":        "success",
		"time_unix":  helper.GetNowTimestamp(),
		"data":       data,
	}
	c.JSON(http.StatusOK, resp)
}

func SupportOptionsMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, device_id, os_platform, token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method != "OPTIONS" {
			ctx := c.Request.Context()
			c.Request = c.Request.WithContext(ctx)

			c.Next()
		} else {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func SetRequestContext(c *gin.Context, i interface{}) {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, constants.RequestParam, i)
	c.Request = c.Request.WithContext(ctx)
}

func getRequestContext(c *gin.Context) []byte {
	ctx := c.Request.Context()
	i := ctx.Value(constants.RequestParam)
	if i == nil {
		return []byte{}
	}
	b := []byte{}
	if err := json.Unmarshal(b, i); err != nil {
		return []byte{}
	}
	return b
}

func Process() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := helper.GetNowTime()

		blw := newResponseWriter(c)
		c.Writer = blw

		c.Next()

		body_req := getRequestContext(c)

		latencyTime := time.Since(startTime)
		reqMethod := c.Request.Method

		reqUri := c.Request.RequestURI

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logx.InfoF("%s %s from %s status[%s], [%v], request = %s, response = %s", reqMethod, reqUri, clientIP, http.StatusText(statusCode), latencyTime, body_req, blw.body.String())
	}
}
