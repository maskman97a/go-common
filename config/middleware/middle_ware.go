package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-common/common"
	"io"
	"net/http"
	"strings"
	"time"
)

func NewBaseResponse(code int, msg string, signature string, data string) *common.BaseResponse {
	return &common.BaseResponse{Code: code, Message: msg, Signature: signature, Data: data, ResponseTime: time.Now().Format(common.DateTimestampPattern)}
}

var excludePath = [1]string{"/actuator"}

func MiddleWare(c *gin.Context) {
	if containsExcludePath(c.Request.URL.Path) {
		c.Next()
		return
	}
	contentType := c.GetHeader("content-type")
	if contentType != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, common.BaseResponse{Code: common.ErrorContentType,
			Message: "Content-Type must be application/json",
		})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	logrus.Info("---[Inbound] Start ", c.Request.Method, c.Request.URL.Path, " ---")
	logrus.Info("[Inbound] Request body: \n", string(body))

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	c.Next()

	logrus.Info("---[Inbound] Finish ", c.Request.Method, c.Request.URL.Path, " ---")
}
func containsExcludePath(path string) bool {
	for _, p := range excludePath {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}
