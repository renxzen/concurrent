package user

import (
	"gorm.io/gorm"

	db "backend/src/configs"
)

type userRepository struct {
	DB *gorm.DB
}

type UserRepository interface {
	GetByEmail(string) (User, error)
	Create(User) (User, error)
	Save(User) (User, error)
	DeleteByEmail(string) (User, error)
}

func NewUserRepository() UserRepository {
	db.Postgres.AutoMigrate(&User{})
	// db.Postgres.Migrator().DropTable(&User{})
	return &userRepository{
		DB: db.Postgres,
	}
}

func (u *userRepository) GetByEmail(email string) (user User, err error) {
	err = u.DB.Where("email=?", email).First(&user).Error
	return
}

func (u *userRepository) Create(user User) (User, error) {
	err := u.DB.Create(&user).Error
	return user, err
}

func (u *userRepository) Save(user User) (User, error) {
	err := u.DB.Save(&user).Error
	return user, err
}

func (u *userRepository) DeleteByEmail(email string) (user User, err error) {
	err = u.DB.Where("email=?", email).Delete(&user).Error
	return
}
