package dao

import (
	"mim/db"

	"gorm.io/gorm"
)

type Message struct {
	Seq       int64
	SenderID  int64
	TargetID  int64
	Content   []byte
	IsRead    bool // true代表未读
	DeletedAt gorm.DeletedAt
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

// kind = 1 获取单聊消息 从start往前size条
// kind = 2 获取群聊消息 从start往前size条 这两种情况下，如果start==0 则表示从最新的开始获取
// kind = 3 群聊的离线消息 从start往后到最新
func GetMessages(userID, targetID, start int64, size int, kind int) ([]Message, error) {
	var messages []Message

	if targetID == 0 {
		if err := db.DB.Select("seq, sender_id, target_id, content").
			Where("sender_id = ? OR target_id = ?", userID, userID).
			Order("seq DESC").
			Limit(size * 2).
			Find(&messages).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.DB.Select("seq, sender_id, target_id, content").
			Where("(sender_id = ? AND target_id = ?) OR (sender_id = ? AND target_id = ?)", userID, targetID, targetID, userID).
			Order("seq DESC").
			Limit(size * 2).
			Find(&messages).Error; err != nil {
			return nil, err
		}
	}

	return messages, nil
}
