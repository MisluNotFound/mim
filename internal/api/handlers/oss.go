package handlers

import (
	"fmt"
	"mim/pkg/code"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type tempToken struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      string
}

func GetOssCredentials(c *gin.Context) {
	accessKeyID := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")

	if accessKeyID == "" || accessKeySecret == "" {
		zap.L().Error("Environment variables not set")
		ResponseError(c, code.CodeServerBusy)
		return
	}
	client, err := sts.NewClientWithAccessKey("cn-qingdao", accessKeyID, accessKeySecret)
	if err != nil {
		zap.L().Error("get temp token failed:", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}

	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = "acs:ram::1192702365393983:role/sts"

	sessionID := uuid.New().String()
	timestamp := time.Now().Format("20060102150405")
	request.RoleSessionName = fmt.Sprintf("session-%s-%s", timestamp, sessionID)

	response, err := client.AssumeRole(request)
	if err != nil {
		zap.L().Error("get temp token failed:", zap.Error(err))
		ResponseError(c, code.CodeServerBusy)
		return
	}

	token := tempToken{
		AccessKeyId:     response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		SecurityToken:   response.Credentials.SecurityToken,
		Expiration:      response.Credentials.Expiration,
	}
	zap.L().Info("temporary token:", zap.Any("info", token))

	ResponseSuccess(c, token)
}
