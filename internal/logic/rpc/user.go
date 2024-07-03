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
		ID:       snowflake.GenerateUniqueID(),
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

func (r *LogicRpc) NearBy(ctx context.Context, req *proto.AuthReq, resp *proto.AuthResp) error {
	resp.Code = code.CodeSuccess

	// 在redis中存放并查询users

	return nil
}

func (r *LogicRpc) GetInfo(ctx context.Context, req *proto.GetInfoReq, resp *proto.GetInfoResp) error {
	resp.Code = code.CodeSuccess

	u, ok, err := dao.FindUserByID(req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeUserNotExist
		return dao.ErrorUserNotExist
	}

	resp.User = u
	return nil
}

func (r *LogicRpc) UpdatePhoto(ctx context.Context, req *proto.UpdatePhotoReq, resp *proto.UpdatePhotoResp) error {
	resp.Code = code.CodeSuccess

	err := dao.UpdatePhoto(req.UserID, req.Avatar)
	if err != nil {
		resp.Code = code.CodeServerBusy
	}

	return nil
}

func (r *LogicRpc) UpdatePassword(ctx context.Context, req *proto.UpdatePasswordReq, resp *proto.UpdatePasswordResp) error {
	resp.Code = code.CodeSuccess

	u, ok, err := dao.FindUserByID(req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeUserNotExist
		return dao.ErrorInvalidID
	}

	if u.Password != dao.Encrypt(req.OldPassword) {
		resp.Code = code.CodeInvalidPassword
		return dao.ErrorInvalidPassword
	}

	err = dao.UpdatePassword(req.UserID, req.NewPassword)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (r *LogicRpc) UpdateName(ctx context.Context, req *proto.UpdateNameReq, resp *proto.UpdateNameResp) error {
	resp.Code = code.CodeSuccess

	_, ok, _ := dao.FindUserByName(req.Name)

	fmt.Println("username ", req.Name)
	if ok {
		resp.Code = code.CodeUserExist
		return dao.ErrorUserExist
	}

	err := dao.UpdateName(req.UserID, req.Name)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}
