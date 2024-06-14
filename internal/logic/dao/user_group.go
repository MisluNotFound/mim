package dao

import (
	"errors"
	"mim/db"

	"gorm.io/gorm"
)

type UserGroup struct {
	UserID   int64
	GroupID  int64
	JoinTime int64 // 加入群聊时最后一条消息 用于消息隔离
	Role     string
	DeleteAt gorm.DeletedAt
}

func (ug *UserGroup) TableName() string {
	return "user_groups"
}

func CreateUserGroup(ug *UserGroup) error {
	lastMessage := Message{}
	if err := db.DB.Select("seq").Where("target_id = ?", ug.GroupID).Order("seq DESC").Limit(1).Error; err != nil {
		return err
	}

	ug.JoinTime = lastMessage.Seq
	if err := db.DB.Create(ug).Error; err != nil {
		return err
	}

	return nil
}

func IsJoined(uid, gid int64) (*UserGroup, bool, error) {
	var c int64
	var ug UserGroup
	if err := db.DB.Where("user_id = ? AND group_id = ?", uid, gid).First(&ug).Count(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	if c > 0 {
		return &ug, true, nil
	}

	return &ug, false, nil
}

func FindUserGroupsByG(gid int64) ([]*UserGroup, error) {
	var ugs []*UserGroup

	if err := db.DB.Where("group_id = ?", gid).Find(&ugs).Error; err != nil {
		return nil, err
	}

	return ugs, nil
}

func DelateUserGroup(uid, gid int64) error {
	var ug UserGroup
	if err := db.DB.Where("user_id = ? AND group_id = ?", uid, gid).Unscoped().Delete(&ug).Error; err != nil {
		return err
	}

	return nil
}

func FindUserGroupsByU(uid int64) ([]*UserGroup, error) {
	var ugs []*UserGroup

	if err := db.DB.Where("user_id = ?", uid).Find(&ugs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return ugs, nil
}

func UpdateMyName(uid int64, gid int64, name string) error {
	return db.DB.Model(&UserGroup{}).Where("user_id = ? AND group_id = ?", uid, gid).Update("nickname", name).Error
}
