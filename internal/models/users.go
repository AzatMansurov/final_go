package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string `json:"name"`
	TelegramId int64  `json:"telegram_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ChatId     int64  `json:"chat_id"`
}

type UserModel struct {
	Db *gorm.DB
}

func (m *UserModel) Create(user User) error {
	return m.Db.Create(&user).Error
}

func (m *UserModel) FindOne(telegramId int64) (*User, error) {
	existUser := User{}

	result := m.Db.First(&existUser, User{TelegramId: telegramId})

	if result.Error != nil {
		return nil, result.Error
	}

	return &existUser, nil
}

func (m *UserModel) FindAll() ([]int, error) {
	var (
		allUsers []User
		usersId  []int
	)

	result := m.Db.Select("chat_id").Find(&allUsers)

	if result.Error != nil {
		return nil, result.Error
	}

	for i := 0; i < len(allUsers); i++ {
		usersId = append(usersId, int(allUsers[i].ChatId))
	}

	return usersId, nil
}
