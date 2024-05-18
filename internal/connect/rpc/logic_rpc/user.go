package logicrpc

import (
	"context"
	"errors"
	"mim/pkg/code"
	"mim/pkg/proto"
)

func Auth(req *proto.AuthReq) (int64, string, error) {
	resp := &proto.AuthResp{}
	err := logicSRpc.Call(context.Background(), "Auth", req, resp)
	if err != nil {
		return -1, "", err
	}

	if resp.Code != code.CodeSuccess {
		return -1, "", errors.New("invalid token")
	}

	return resp.UserID, resp.Username, nil
}

func GetGroup(req *proto.FindGroupsReq) (*[]int64, error) {
	resp := &proto.FindGroupsResp{}

	err := logicSRpc.Call(context.Background(), "FindGroups", req, resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != code.CodeSuccess {
		return nil, errors.New("server busy")
	}

	return resp.Groups, nil
}

func Online(req *proto.OnlineReq) error {
	resp := &proto.OnlineResp{}

	err := logicPRpc.Call(context.Background(), "Online", req, resp)
	if err != nil {
		return err
	}

	if resp.Code != code.CodeSuccess {
		return errors.New("server busy")
	}

	return nil
}
