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
	ID int64 
	Username string
	Password string
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