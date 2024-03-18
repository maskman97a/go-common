package common

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-common/config/router"
	"go-common/pkg/utils"
	"net/http"
	"time"
)

type BaseController interface {
	InitRouter(routerGroup *gin.RouterGroup)
}

type baseController struct {
}

func NewBaseController() BaseController {
	return &baseController{}
}

func (baseController *baseController) InitRouter(routerGroup *gin.RouterGroup) {
	api := routerGroup.Group("/")
	router.Get(api, "/", GetService)
}

func GetService(c *gin.Context) {
	c.JSON(http.StatusOK, NewBaseResponse(Success,
		"success",
		"", "",
	))
}
func ValidateRequest(c *gin.Context) (*BaseRequest, error) {
	logrus.Info(fmt.Sprintf("--Start %s--", "baseController.ValidateRequest"))
	var baseRequest BaseRequest
	err := c.BindJSON(&baseRequest)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if baseRequest.RequestTime == "" {
		return nil, errors.New("requestTime is empty")
	}
	_, err = time.Parse(DateTimestampPattern, baseRequest.RequestTime)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if baseRequest.IpAddress == "" {
		return nil, errors.New("ipAddress is empty")
	}
	if baseRequest.Signature == "" {
		return nil, errors.New("signature is empty")
	}
	if baseRequest.Data == "" {
		return nil, errors.New("data is empty")
	}
	if !utils.IsJSON(baseRequest.Data) {
		return nil, errors.New("data is invalid")
	}
	logrus.Info(fmt.Sprintf("--Finish %s--", "baseController.ValidateRequest"))
	return &baseRequest, nil
}

//func SetSignature(data string) (string, error) {
//	signResp, err := service.NewHsmService().Sign(request.SignRequest{Data: data})
//	if err != nil {
//		logrus.Error(err)
//		return "", err
//	}
//	return signResp.Signature, nil
//}
