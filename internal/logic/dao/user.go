package dao

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"mim/db"

	"gorm.io/gorm"
)

var secret = "mislu&mim"

type User struct {
	ID       int64
	Username string
	Password string
	Avatar   string
	DeleteAt gorm.DeletedAt
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Migrate() {
	db.DB.AutoMigrate(u)
}

// true表示存在 当查询出现除了RecordNotFound之外的错误时返回err
func FindUserByName(username string) (*User, bool, error) {
	var u User
	var c int64
	if err := db.DB.Where("username = ?", username).First(&u).Count(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	if c > 0 {
		return &u, true, nil
	}

	return &u, false, nil
}

func CreateUser(u *User) error {
	u.Password = Encrypt(u.Password)
	if err := db.DB.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func Encrypt(oPassword string) string {
	h := md5.New()
	h.Write([]byte(oPassword))
	return hex.EncodeToString(h.Sum([]byte(secret)))
}

func FindUserByID(id int64) (*User, bool, error) {
	var u User
	var c int64
	if err := db.DB.Where("id = ?", id).First(&u).Count(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	if c > 0 {
		return &u, true, nil
	}

	return &u, false, nil
}

func UpdatePhoto(id int64, avatar string) error {
	return db.DB.Model(&User{}).Where("id = ?", id).Update("avatar", avatar).Error
}

func UpdatePassword(id int64, password string) error {
	password = Encrypt(password)
	return db.DB.Model(&User{}).Where("id = ?", id).Update("password", password).Error
}

func UpdateName(id int64, name string) error {
	return db.DB.Model(&User{}).Where("id = ?", id).Update("username", name).Error
}

func GetUserPhoto(id int64) (string, error) {
	var avatar string
	err := db.DB.Model(&User{}).Select("avatar").Where("id = ?", id).Find(&avatar).Error
	return avatar, err
}
