package rpc

import (
	"context"
	"mim/internal/logic/dao"
	"mim/pkg/code"
	"mim/pkg/proto"
	"mim/pkg/snowflake"

	"go.uber.org/zap"
)

func getMembers(id int64) ([]*dao.User, error) {
	ugs, err := dao.FindUserGroupsByG(id)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, ug := range ugs {
		ids = append(ids, ug.UserID)
	}

	return dao.FindMembers(ids)
}

func (l *LogicRpc) NewGroup(ctx context.Context, req *proto.NewGroupReq, resp *proto.NewGroupResp) error {
	resp.Code = code.CodeSuccess

	g := dao.Group{
		GroupID:     snowflake.GenID(),
		GroupName:   req.GroupName,
		Description: req.Description,
	}

	if err := dao.CreateGroup(&g); err != nil {
		zap.L().Error("logic NewGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	ug := dao.UserGroup{
		UserID:  req.OwnerID,
		GroupID: g.GroupID,
		Role:    dao.Owner,
	}

	if err := dao.CreateUserGroup(&ug); err != nil {
		zap.L().Error("logic NewGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	users, err := getMembers(g.GroupID)

	if err != nil {
		resp.Code = code.CodeServerBusy
		zap.L().Error("logic NewGroup() failed: ", zap.Error(err))
		return err
	}

	g.Members = users

	resp.Group = &g
	return nil
}

func (l *LogicRpc) JoinGroup(ctx context.Context, req *proto.JoinGroupReq, resp *proto.JoinGroupResp) error {
	resp.Code = code.CodeSuccess

	g, ok, err := dao.FindGroupByID(req.GroupID)
	if err != nil {
		zap.L().Error("logic JoinGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		zap.L().Error("logic JoinGroup() failed: group not exists")
		resp.Code = code.CodeGroupNotExist
		return dao.ErrorGroupNotExist
	}

	_, ok, err = dao.IsJoined(req.UserID, req.GroupID)
	if err != nil {
		zap.L().Error("logic JoinGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	if ok {
		zap.L().Error("logic JoinGroup() failed: user already join group")
		resp.Code = code.CodeAlreadyJoined
		return dao.ErrorGroupAlreadyJoined
	}

	ug := dao.UserGroup{
		UserID:  req.UserID,
		GroupID: req.GroupID,
		Role:    dao.Member,
	}

	if err := dao.CreateUserGroup(&ug); err != nil {
		resp.Code = code.CodeServerBusy
		zap.L().Error("logic JoinGroup() failed: ", zap.Error(err))
		return err
	}

	users, err := getMembers(req.GroupID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		zap.L().Error("logic JoinGroup() failed: ", zap.Error(err))
		return err
	}

	g.Members = users

	resp.Group = g

	return nil
}

func (l *LogicRpc) FindGroup(ctx context.Context, req *proto.FindGroupReq, resp *proto.FindGroupResp) error {
	resp.Code = code.CodeSuccess

	g, ok, err := dao.FindGroupByID(req.GroupID)
	if err != nil {
		zap.L().Error("logic FindGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		zap.L().Error("logic JoinGroup() failed: group not exists")
		resp.Code = code.CodeGroupNotExist
		return dao.ErrorGroupNotExist
	}

	users, err := getMembers(req.GroupID)
	if err != nil {
		zap.L().Error("logic FindGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	g.Members = users

	resp.Group = g
	return nil
}

func (l *LogicRpc) LeaveGroup(ctx context.Context, req *proto.LeaveGroupReq, resp *proto.LeaveGroupResp) error {
	resp.Code = code.CodeSuccess

	_, ok, err := dao.IsJoined(req.UserID, req.GroupID)
	if err != nil {
		zap.L().Error("logic LeaveGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		zap.L().Error("logic LeaveGroup() failed: not join the group")
		resp.Code = code.CodeNotJoinGroup
		return dao.ErrorNotJoinGroup
	}

	if err := dao.DelateUserGroup(req.UserID, req.GroupID); err != nil {
		zap.L().Error("logic LeaveGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (l *LogicRpc) FindGroups(ctx context.Context, req *proto.FindGroupsReq, resp *proto.FindGroupsResp) error {
	resp.Code = code.CodeSuccess

	ugs, err := dao.FindUserGroupsByU(req.UserID)
	if err != nil {
		zap.L().Error("logic FindGroups() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	var gids []int64
	for _, ug := range ugs {
		gids = append(gids, ug.GroupID)
	}

	resp.Groups = &gids

	return nil
}
