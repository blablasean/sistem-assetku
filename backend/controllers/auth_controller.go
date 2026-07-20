package controllers

import (
	"errors"
	"fmt"
	"sistem-asetku-backend/models"
	"time"

	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

func (c *AuthController) Login(username string, password string) (string, error) {
	var user models.User
	if err := c.db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("username atau password salah")
		}
		return "", err
	}

	token := fmt.Sprintf("%s:%d:%d", username, user.ID, time.Now().Unix())
	return token, nil
}

func (c *AuthController) Logout(token string) error {
	if token == "" {
		return errors.New("token tidak valid")
	}
	return nil
}
