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
	Avatar      string
	DeleteAt    gorm.DeletedAt
	Members     []*User `gorm:"-"`
	MyRemark    string  `gorm:"-"`
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
	if err := db.DB.Select("id", "username", "avatar").Where("id in ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetGroups(ids []int64) ([]Group, error) {
	var groups []Group

	if err := db.DB.Select("group_id, group_name, description, avatar").Where("group_id in ?", ids).Find(&groups).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []Group{}, nil
		}
		return groups, err
	}

	return groups, nil
}

func UpdateGroupName(groupID int64, name string) error {
	return db.DB.Model(&Group{}).Where("group_id = ?", groupID).Update("group_name", name).Error
}

func UpdateGroupPhoto(groupID int64, avatar string) error {
	return db.DB.Model(&Group{}).Where("group_id = ?", groupID).Update("avatar", avatar).Error
}

func DeleteGroup(groupID int64) error {
	return db.DB.Where("group_id = ?", groupID).Delete(&Group{}).Error
}
