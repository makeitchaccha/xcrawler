package xcrawler

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement:false"`
	Username string `gorm:"primaryKey;"`
}

type History struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt      time.Time `gorm:"primaryKey;autoIncrement:false"`
	FavoritesCount int64     `gorm:"not null"`
	TweetCount     int64     `gorm:"not null"`
}

type Repository interface {
	SaveUser(user User) error
	SaveHistory(history History) error

	FindAllUsers() ([]User, error)
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) SaveUser(user User) error {
	return r.db.Create(&user).Error
}

func (r *repositoryImpl) SaveHistory(history History) error {
	return r.db.Create(&history).Error
}

func (r *repositoryImpl) FindAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
