package redis

import "strconv"

const (
	prefixOnlineUser          = "online:user:"
	prefixSession             = "session:"       // 记录所有会话 prefix+userid: sender score 1/0
	prefixGroupUser           = "group:users:"   // 记录群在线成员
	prefixEarlyMessage        = "early:message:" // 记录用户入群时的时间 用于消息隔离
	prefixAckMessage          = "ack:"           // 记录用户ack err消息的seq
	prefixSessionOfflineCount = "offline:count:" // 记录会话中的未读消息数， prefix:session: string
)

func GetSessionID(senderID, targetID int64) string {
	senderStr := strconv.FormatInt(senderID, 10)
	targetStr := strconv.FormatInt(targetID, 10)
	var sessionID string
	if senderID > targetID {
		sessionID += targetStr + senderStr
	} else {
		sessionID += senderStr + targetStr
	}

	return sessionID
}

// prefixSession
func getUserSession(uid int64) string {
	uidStr := strconv.FormatInt(uid, 10)

	return prefixSession + uidStr
}

func getOnlineUser(userID interface{}) string {
	uidstr := ""
	switch v := userID.(type) {
	case int64:
		uidstr = strconv.FormatInt(v, 10)
	case string:
		uidstr = v
	}

	return prefixOnlineUser + uidstr
}

func getGroupUser(groupID int64) string {
	gidstr := strconv.FormatInt(groupID, 10)

	return prefixGroupUser + gidstr
}

func getAckMessage(userID int64) string {
	uidstr := strconv.FormatInt(userID, 10)

	return prefixAckMessage + uidstr
}

// prefixSessionOfflineCount 我id+会话id
func getSessionOfflineCount(userID interface{}, senderID interface{}) string {
	uidstr := ""
	switch v := userID.(type) {
	case int64:
		uidstr = strconv.FormatInt(v, 10)
	case string:
		uidstr = v
	}

	sessionStr := ""
	switch v := senderID.(type) {
	case int64:
		sessionStr = strconv.FormatInt(v, 10)
	case string:
		sessionStr = v
	}
	return prefixSessionOfflineCount + uidstr + sessionStr
}
