package dao

import (
	"fmt"
	"mim/db"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	Seq       int64
	SenderID  int64
	TargetID  int64
	Content   []byte
	IsRead    bool // true代表未读
	DeletedAt gorm.DeletedAt
	Type      string
	URL       string
	Timer     time.Time `gorm:"default:null"`
	IsGroup   bool
	Extra     interface{} `gorm:"-"`
}

func (m *Message) TableName() string {
	return "messages"
}

func StoreMysqlMessage(msg *Message) error {
	if err := db.DB.Create(msg).Error; err != nil {
		return err
	}

	return nil
}

func GetLastMessage(senders []int64, targetID int64) ([]Message, error) {
	var lastMessages []Message
	var sqlStr = "SELECT m.seq, m.sender_id, m.target_id, m.content, m.type, m.url, m.timer, m.is_group " +
		"FROM messages m JOIN (SELECT target_id, MAX(seq) AS max_seq FROM messages WHERE " +
		"(sender_id IN ? AND target_id = ?) OR (is_group = 1 AND target_id IN ?) GROUP BY target_id) AS sub " +
		"ON m.target_id = sub.target_id AND m.seq = sub.max_seq;"

	if err := db.DB.Raw(sqlStr, senders, targetID, senders).Scan(&lastMessages).Error; err != nil {
		return nil, err
	}

	return lastMessages, nil
}

func PullOfflineMessage(userID, senderID int64, isGroup bool, count int, joinTime int64) ([]Message, error) {
	var messages []Message
	var err error
	fmt.Println(count)
	if isGroup {
		err = db.DB.Select("seq, sender_id, target_id, content, type, url, timer, is_group").
			Where("sender_id = ? AND target_id = ? AND seq > ?", senderID, userID, joinTime).
			Order("seq DESC").
			Limit(count).
			Find(&messages).Error
	} else {
		err = db.DB.Select("seq, sender_id, target_id, content, type, url, timer, is_group").
			Where("sender_id = ? AND target_id = ?", senderID, userID).
			Order("seq DESC").
			Limit(count).
			Find(&messages).Error
	}

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func PullErrMessage(userID int64, lastAck int64) ([]Message, error) {
	var singleMessages []Message
	// 获取单聊消息
	if err := db.DB.Select("seq, sender_id, target_id, content, type, url, timer, is_group").
		Where("target_id = ? AND seq > ?", userID, lastAck).
		Order("seq DESC").
		Find(&singleMessages).Error; err != nil {
		return nil, err
	}

	// 获取群聊消息
	var groupMessages []Message
	if err := db.DB.Table("messages").
		Select("seq, sender_id, target_id, content, type, url, timer, is_group").
		Joins("INNER JOIN user_groups ON messages.target_id = user_groups.group_id").
		Where("messages.is_group = true").
		Where("user_groups.user_id = ?", userID).
		Where("messages.seq > ?", lastAck).
		Order("seq DESC").
		Find(&groupMessages).Error; err != nil {
		return nil, err
	}

	singleMessages = append(singleMessages, groupMessages...)
	return singleMessages, nil
}

func PullSingleMessage(userID, sessionID int64, start int64, size int) ([]Message, error) {
	var messages []Message

	if err := db.DB.Select("seq, sender_id, target_id, content, type, url, timer, is_group").
		Where("(sender_id = ? AND target_id = ?) OR (sender_id = ? AND target_id = ?)", userID, sessionID, sessionID, userID).
		Where("seq < ?", start).
		Order("seq DESC").
		Limit(size).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func PullGroupMessage(groupID, start, joinTime int64, size int) ([]Message, error) {
	var messages []Message

	if err := db.DB.Select("seq, sender_id, target_id, content, type, url, timer, is_group").
		Where("is_group = true").
		Where("target_id = ?", groupID).
		Where("seq between ? and ?", joinTime, start).
		Order("seq DESC").
		Limit(size).
		Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
