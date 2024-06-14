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
	accessKeyId     string
	accessKeySecret string
	securityToken   string
	expiration      string
}

func GetOssCredentials(c *gin.Context) {
	accessKeyID := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")

	client, err := sts.NewClientWithAccessKey("cn-qingdao", accessKeyID, accessKeySecret)
	if err != nil {
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
		ResponseError(c, code.CodeServerBusy)
		return
	}

	token := tempToken{
		accessKeyId:     response.Credentials.AccessKeyId,
		accessKeySecret: response.Credentials.AccessKeySecret,
		securityToken:   response.Credentials.SecurityToken,
		expiration:      response.Credentials.Expiration,
	}
	zap.L().Info("temporary token:", zap.Any("info", token))

	ResponseSuccess(c, token)
}
