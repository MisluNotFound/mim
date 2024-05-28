package redis

import "strconv"

const (
	prefixOnlineUser         = "online:user:"
	prefixSession            = "session:"              // 记录所有会话 prefix+userid: sender+target
	prefixMessage            = "message:list:"         // 记录会话中的所有消息 prefix+sender+target: seq
	prefixMessageOffline     = "session:offline:"      // 索引 记录哪些用户给目标用户发送了离线消息 prefix+uid: sender
	prefixMessageOfflineList = "message:offline:list:" // 索引 记录离线消息列表 prefix+senderid:
	prefixMessageList        = "message:list"          // 记录所有未接收的消息 unack
	prefixGroupUser          = "group:users:"          // 记录群在线成员
	prefixGroupOffline       = "group:offline:"        // 记录用户的哪些群聊有离线消息
	prefixGroupMessage       = "group:message:"        // 记录群消息缓存
	prefixEarlyMessage       = "early:message:"        // 记录用户入群时的时间 用于消息隔离
	prefixLastMessage        = "last:message:"         // 记录用户离线时最后一条群消息 prefix+user: group:seq
)

func getOnlineUser(userID interface{}) string {
	uidstr := ""
	switch v := userID.(type) {
	case int64:
		uidstr = strconv.FormatInt(v, 10)
	case string:
		uidstr = v // 如果已经是字符串类型，直接使用
	}

	return prefixOnlineUser + uidstr
}

func getGroupOffline(userID int64) string {
	uidstr := strconv.FormatInt(userID, 10)

	return prefixGroupOffline + uidstr
}

func getGroupUser(groupID int64) string {
	gidstr := strconv.FormatInt(groupID, 10)

	return prefixGroupUser + gidstr
}

func getGroupMessage(groupID int64) string {
	gidstr := strconv.FormatInt(groupID, 10)

	return prefixGroupMessage + gidstr
}

func getEarlyMessage(userID int64) string {
	uidstr := strconv.FormatInt(userID, 10)

	return prefixEarlyMessage + uidstr
}
func getLastMessage(userID int64) string {
	uidstr := strconv.FormatInt(userID, 10)

	return prefixLastMessage + uidstr
}
