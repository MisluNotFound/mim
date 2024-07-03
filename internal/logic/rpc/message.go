package rpc

import (
	"context"
	"fmt"
	"math"
	"mim/internal/logic/dao"
	"mim/internal/logic/redis"
	"mim/pkg/code"
	"mim/pkg/proto"

	"go.uber.org/zap"
)

// 获取历史记录
func (r *LogicRpc) PullMessage(ctx context.Context, req *proto.PullMessageReq, resp *proto.PullMessageResp) error {
	resp.Code = code.CodeSuccess

	fmt.Println("pull message req", req)
	// 判断是否为群
	if req.LastSeq == 0 {
		req.LastSeq = math.MaxInt64
	}
	if req.IsGroup {
		// 获取入群时间
		ug, ok, err := dao.IsJoined(req.UserID, req.SessionID)
		if err != nil {
			resp.Code = code.CodeServerBusy
			zap.L().Error("pull message failed: ", zap.Error(err))
			return err
		}

		if !ok {
			resp.Code = code.CodeNotJoinGroup
			return dao.ErrorNotJoinGroup
		}

		resp.Messages, err = dao.PullGroupMessage(req.SessionID, req.LastSeq, ug.JoinTime, req.Size)
		if err != nil {
			resp.Code = code.CodeServerBusy
			zap.L().Error("pull message failed: ", zap.Error(err))
			return err
		}
	} else {
		var err error
		resp.Messages, err = dao.PullSingleMessage(req.UserID, req.SessionID, req.LastSeq, req.Size)
		if err != nil {
			resp.Code = code.CodeServerBusy
			zap.L().Error("pull message failed: ", zap.Error(err))
			return err
		}
	}
	return nil
}

// 用户点开会话获取未读消息
func (r *LogicRpc) PullOfflineMessage(ctx context.Context, req *proto.PullOfflineMessageReq, resp *proto.PullMessageResp) error {
	resp.Code = code.CodeSuccess
	var messages []dao.Message
	var err error

	lastRead, err := redis.GetLastRead(req.UserID, req.SessionID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if req.IsGroup {
		var ok bool
		var ug *dao.UserGroup
		ug, ok, err = dao.IsJoined(req.UserID, req.SessionID)
		if err != nil {
			zap.L().Error("pull offline message failed: ", zap.Error(err))
			resp.Code = code.CodeServerBusy
			return err
		}

		if !ok {
			resp.Code = code.CodeNotJoinGroup
			return nil
		}

		messages, err = dao.PullOfflineMessage(req.UserID, req.SessionID, lastRead, true, ug.JoinTime)
	} else {
		messages, err = dao.PullOfflineMessage(req.UserID, req.SessionID, lastRead, false, 0)
	}

	if err != nil {
		zap.L().Error("pull offline message failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	err = redis.MarkAsRead(req.UserID, req.SessionID)
	if err != nil {
		zap.L().Error("pull offline message failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	resp.Messages = messages
	return nil
}

// 获取会未读消息数

// 获取session:lastRead
func (r *LogicRpc) GetUnReadCount(ctx context.Context, req *proto.GetUnReadCountReq, resp *proto.GetUnReadResp) error {
	resp.Code = code.CodeSuccess

	// 获取lastRead
	counts, err := redis.GetAllLastRead(req.UserID)
	fmt.Println("counts ", counts)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	// db获取最后一条消息
	var senders []int64
	for s := range counts {
		senders = append(senders, s)
	}
	lastMessages, err := dao.GetLastMessage(senders, req.UserID)

	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}
	// 构造返回值
	var infos []proto.UnReadInfo
	for s, c := range counts {
		info := proto.UnReadInfo{
			SessionID: s,
		}

		for _, m := range lastMessages {
			// 如果是群聊 查询群聊备注 头像
			if m.IsGroup {
				group, _, err := dao.FindGroupByID(m.TargetID)
				if err != nil {
					zap.L().Error("get group detail failed: ", zap.Error(err))
				}
				info.Avatar = group.Avatar
				info.Remark = group.GroupName
				info.LastMessage = m
			} else {
				// 如果是单聊 查询好友备注 头像
				friend, _ := dao.GetFriend(m.SenderID, m.TargetID)
				var friendID int64
				if friend.UserA == req.UserID {
					friendID = friend.UserB
					info.Remark = friend.AtoB
				} else {
					friendID = friend.UserA
					info.Remark = friend.BtoA
				}
				avatar, err := dao.GetUserPhoto(friendID)
				if err != nil {
					zap.L().Error("get user photo failed: ", zap.Error(err))
				}
				info.Avatar = avatar
				if m.SenderID == s || m.TargetID == s {
					info.LastMessage = m
				}
			}
			c, _ := dao.GetUnReadCount(c, m.Seq, s, req.UserID)
			info.Count = c
		}
		infos = append(infos, info)
	}
	resp.SessionInfo = infos
	return nil
}

func (r *LogicRpc) PullErrMessage(ctx context.Context, req *proto.PullErrMessageReq, resp *proto.PullErrMessageResp) error {
	resp.Code = code.CodeSuccess

	// 获取lastAck
	lastAck, err := redis.GetLastAck(req.UserID)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	if lastAck == -1 {
		return nil
	}

	// 查出lastAck之后的消息
	messages, err := dao.PullErrMessage(req.UserID, lastAck)
	if err != nil {
		resp.Code = code.CodeServerBusy
		return err
	}

	resp.Messages = messages
	return nil
}

func (r *LogicRpc) Online(ctx context.Context, req *proto.OnlineReq, resp *proto.OnlineResp) error {
	resp.Code = code.CodeSuccess
	if err := redis.AddOnlineUser(req.UserID, req.ServerID, req.BucketID); err != nil {
		zap.L().Info("logic Online() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (r *LogicRpc) Offline(ctx context.Context, req *proto.OfflineReq, resp *proto.OfflineResp) error {
	resp.Code = code.CodeSuccess

	if err := redis.RemoveOnlineUser(req.UserID); err != nil {
		zap.L().Error("logic Offline() failed: ", zap.Error(err))
		resp.Code = code.CodeServerBusy
		return err
	}

	return nil
}

func (r *LogicRpc) StoreOffline(ctx context.Context, req *proto.OfflineMessageReq, resp *proto.MessageResp) error {
	err := redis.AddUnReadCount(req.TargetID, req.SenderID, req.Seq)

	if err != nil {
		return err
	}

	return nil
}
