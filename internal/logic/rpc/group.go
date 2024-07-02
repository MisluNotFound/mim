package rpc

import (
	"context"
	"fmt"
	"mim/internal/logic/dao"
	"mim/internal/logic/redis"
	"mim/pkg/code"
	"mim/pkg/proto"
	"mim/pkg/snowflake"

	"go.uber.org/zap"
)

func getMembers(id int64) ([]proto.UserInfo, error) {
	ugs, err := dao.FindUserGroupsByG(id)
	if err != nil {
		return nil, err
	}

	var ids []int64
	members := make([]proto.UserInfo, len(ugs))
	for i, ug := range ugs {
		members[i] = proto.UserInfo{
			Role: ug.Role,
		}
		ids = append(ids, ug.UserID)
	}

	users, err := dao.FindMembers(ids)
	if err != nil {
		return []proto.UserInfo{}, err
	}
	for i, u := range users {
		members[i].Member = *u
	}

	return members, nil
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

	if err := redis.JoinGroup(req.OwnerID, g.GroupID); err != nil {
		resp.Code = code.CodeServerBusy
		zap.L().Error("logic JoinGroup() failed: ", zap.Error(err))
		return err
	}

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

	if err := redis.JoinGroup(req.UserID, req.GroupID); err != nil {
		resp.Code = code.CodeServerBusy
		zap.L().Error("logic JoinGroup() failed: ", zap.Error(err))
		return err
	}

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

	resp.Group = g
	return nil
}

func (l *LogicRpc) LeaveGroup(ctx context.Context, req *proto.LeaveGroupReq, resp *proto.LeaveGroupResp) error {
	resp.Code = code.CodeSuccess

	ug, ok, err := dao.IsJoined(req.UserID, req.GroupID)
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

	if err := dao.DeleteUserGroup(req.UserID, req.GroupID); err != nil {
		zap.L().Error("logic LeaveGroup() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	redis.LeaveGroup(req.UserID, req.GroupID)
	if ug.Role == dao.Owner {
		err := dao.DeleteGroup(req.GroupID)
		redis.DelGroup(req.GroupID)
		if err != nil {
			resp.Code = code.CodeServerBusy
			return err
		}
	}

	return nil
}

func (l *LogicRpc) GetGroups(ctx context.Context, req *proto.GetGroupsReq, resp *proto.GetGroupsResp) error {
	resp.Code = code.CodeSuccess

	ugs, err := dao.FindUserGroupsByU(req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	var ids []int64
	for _, ug := range ugs {
		ids = append(ids, ug.GroupID)
	}

	groups, err := dao.GetGroups(ids)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	resp.Groups = groups
	fmt.Println(resp.Groups)
	return nil
}

func (r *LogicRpc) UpdateGroupName(ctx context.Context, req *proto.UpdateGroupNameReq, resp *proto.UpdateGroupNameResp) error {
	resp.Code = code.CodeSuccess

	ug, ok, err := dao.IsJoined(req.UserID, req.GroupID)

	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeNotJoinGroup
		return dao.ErrorNotJoinGroup
	}

	if ug.Role != "owner" {
		resp.Code = code.CodePermissionDenied
		return dao.ErrorPermissionDenied
	}

	err = dao.UpdateGroupName(req.GroupID, req.Name)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (r *LogicRpc) UpdateGroupPhoto(ctx context.Context, req *proto.UpdateGroupPhotoReq, resp *proto.UpdateGroupPhotoResp) error {
	resp.Code = code.CodeSuccess

	ug, ok, err := dao.IsJoined(req.UserID, req.GroupID)

	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeNotJoinGroup
		return dao.ErrorNotJoinGroup
	}

	if ug.Role != "owner" {
		resp.Code = code.CodePermissionDenied
		return dao.ErrorPermissionDenied
	}

	err = dao.UpdateGroupPhoto(req.GroupID, req.Avatar)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (r *LogicRpc) GetMembers(ctx context.Context, req *proto.GetMembersReq, resp *proto.GetMembersResp) error {
	resp.Code = code.CodeSuccess

	_, ok, err := dao.IsJoined(req.UserID, req.GroupID)

	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeNotJoinGroup
		return dao.ErrorNotJoinGroup
	}

	members, err := getMembers(req.GroupID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	resp.Members = members

	return nil
}

func (r *LogicRpc) GetRole(ctx context.Context, req *proto.GetRoleReq, resp *proto.GetRoleResp) error {
	resp.Code = code.CodeSuccess

	ug, ok, err := dao.IsJoined(req.UserID, req.GroupID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeNotJoinGroup
		return dao.ErrorNotJoinGroup
	}

	resp.Role = ug.Role
	return nil
}

func (r *LogicRpc) RemoveMember(ctx context.Context, req *proto.RemoveMemberReq, resp *proto.RemoveMemberResp) error {
	resp.Code = code.CodeSuccess

	ownerUg, ok, err := dao.IsJoined(req.UserID, req.GroupID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeNotJoinGroup
		return dao.ErrorNotJoinGroup
	}

	if ownerUg.Role != dao.Owner {
		resp.Code = code.CodePermissionDenied
		return dao.ErrorPermissionDenied
	}

	_, ok, err = dao.IsJoined(req.MemberID, req.GroupID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if !ok {
		resp.Code = code.CodeNotJoinGroup
		return dao.ErrorNotJoinGroup
	}

	err = dao.DeleteUserGroup(req.MemberID, req.GroupID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	redis.LeaveGroup(req.MemberID, req.UserID)

	return nil
}