package dao

import (
	"mim/db"

	"gorm.io/gorm"
)

type Group struct {
	GroupID     int64
	GroupName   string
	Description string
	Members     []*User
	DeleteAt    gorm.DeletedAt
}

func (g *Group) TableName() string {
	return "groups"
}

func CreateGroup(group *Group) error {
	if err := db.DB.Create(group).Error; err != nil {
		return err
	}

	return nil
}

func FindGroupByID(id int64) (*Group, error) {

	return nil, nil
}
