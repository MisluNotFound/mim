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

func GetMessages(seqs []int64) ([]Message, error) {
	var messages []Message
	if err := db.DB.Where("seq IN ?", seqs).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
