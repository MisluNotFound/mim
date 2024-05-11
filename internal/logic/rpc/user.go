package rpc

import (
	"context"
	"fmt"
	"mim/internal/logic/dao"
	"mim/pkg/code"
	"mim/pkg/jwt"
	"mim/pkg/proto"
	"mim/pkg/snowflake"

	"go.uber.org/zap"
)

func (r *LogicRpc) SignUp(ctx context.Context, req *proto.SignUpReq, resp *proto.SignUpResp) error {
	resp.Code = code.CodeSuccess

	_, ok, err := dao.FindUserByName(req.Username)
	if err != nil {
		zap.L().Error("logic SignUp() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	if ok {
		zap.L().Error("logic SignUp() failed: user exist")
		resp.Code = code.CodeUserExist
		return nil
	}

	u := dao.User{
		ID:       snowflake.GenID(),
		Username: req.Username,
		Password: req.Password,
	}

	if err := dao.CreateUser(&u); err != nil {
		zap.L().Error("logic SignUp() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	token, err := jwt.GenToken(u.ID, u.Username)
	if err != nil {
		zap.L().Error("logic SignUp() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}
	resp.Token = token
	zap.L().Info(resp.Token)
	fmt.Println("println: logic ", resp.Token)
	return nil
}

func (r *LogicRpc) SignIn(ctx context.Context, req *proto.SignInReq, resp *proto.SignInResp) error {
	resp.Code = code.CodeSuccess
	u, ok, err := dao.FindUserByName(req.Username)
	if err != nil {
		zap.L().Error("logic SignIn() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		zap.L().Error("logic SignIn() failed: user not exist")
		resp.Code = code.CodeUserNotExist
		return dao.ErrorUserNotExist
	}

	if u.Password != dao.Encrypt(req.Password) {
		zap.L().Error("logic SignIn() failed: invalid password")
		resp.Code = code.CodeInvalidPassword
		return dao.ErrorInvalidPassword
	}

	token, err := jwt.GenToken(u.ID, u.Username)
	if err != nil {
		zap.L().Error("logic SignIn() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}
	resp.Token = token
	zap.L().Info(resp.Token)
	return nil
}

func (r *LogicRpc) Auth(ctx context.Context, req *proto.AuthReq, resp *proto.AuthResp) error {
	resp.Code = code.CodeSuccess
	token := req.Token
	c, err := jwt.ParseToken(token)
	if err != nil {
		resp.Code = code.CodeInvalidToken
		return err
	}

	resp.UserID = c.UserID
	resp.Username = c.Username
	return nil
}
