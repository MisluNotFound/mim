package dao

import (
	"errors"
	"mim/db"

	"gorm.io/gorm"
)

type UserGroup struct {
	UserID   int64
	GroupID  int64
	Role     string
	DeleteAt gorm.DeletedAt
}

func (ug *UserGroup) TableName() string {
	return "user_groups"
}

func CreateUserGroup(ug *UserGroup) error {
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
	if err := db.DB.Where("user_id = ? AND group_id = ?", uid, gid).Delete(&ug).Error; err != nil {
		return err
	}

	return nil
}

func FindUserGroupsByU(uid int64) ([]*UserGroup, error) {
	var ugs []*UserGroup

	if err := db.DB.Where("user_id = ?", uid).Find(&ugs).Error; err != nil {
		return nil, err
	}

	return ugs, nil
}
