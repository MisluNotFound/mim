package dao

import (
	"errors"
	"mim/db"

	"gorm.io/gorm"
)

type Friend struct {
	UserA int64
	UserB int64
	AtoB  string `gorm:"column:a_to_b"` // a给b的备注
	BtoA  string `gorm:"column:b_to_a"` // b给a的备注
}

type FriendInfo struct {
	Info   User
	Remark string
}

func (f *Friend) TableName() string {
	return "friends"
}

func GetFriend(userA, userB int64) (*Friend, error) {
	friend := Friend{}
	err := db.DB.Where("(user_a = ? AND user_B = ?) or (user_a = ? AND user_b = ?)", userA, userB, userB, userA).Find(&friend).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &friend, nil
}

func AddFriend(userA, userB int64) error {
	friend := Friend{
		UserA: userA,
		UserB: userB,
	}
	if err := db.DB.Create(&friend).Error; err != nil {
		return err
	}

	return nil
}

func DeleteFriend(userA, userB int64) error {
	if err := db.DB.Where("(user_a = ? AND user_B = ?) or (user_a = ? AND user_b = ?)", userA, userB, userB, userA).Delete(&Friend{}).Error; err != nil {
		return err
	}

	return nil
}

func GetFriends(userID int64) ([]Friend, error) {
	friends := []Friend{}

	if err := db.DB.Where("user_a = ? OR user_b = ?", userID, userID).Find(&friends).Error; err != nil {
		return nil, err
	}

	return friends, nil
}

func GetFriendsInfo(ids []int64) ([]User, error) {
	users := []User{}

	if err := db.DB.Select("id, username").Where("id in ?", ids).Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func UpdateFriendRemark(friend *Friend) error {
	return db.DB.Where("user_a = ? AND user_b = ?", friend.UserA, friend.UserB).Save(friend).Error
}
