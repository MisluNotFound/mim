package dao

import (
	"errors"
	"mim/db"

	"gorm.io/gorm"
)

type Group struct {
	GroupID     int64
	GroupName   string
	Description string
	Members     []*User `gorm:"-"`
	DeleteAt    gorm.DeletedAt
}

const (
	Member string = "member"
	Admin  string = "admin"
	Owner  string = "owner"
)

func (g *Group) TableName() string {
	return "groups"
}

func CreateGroup(group *Group) error {
	if err := db.DB.Create(group).Error; err != nil {
		return err
	}

	return nil
}

// true表示存在 当查询出现除了RecordNotFound之外的错误时返回err
func FindGroupByID(id int64) (*Group, bool, error) {
	var g Group
	var c int64
	if err := db.DB.Where("group_id = ?", id).First(&g).Count(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	if c > 0 {
		return &g, true, nil
	}

	return &g, false, nil
}

func FindMembers(ids []int64) ([]*User, error) {
	var users []*User
	if err := db.DB.Select("id", "username").Where("id in ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

