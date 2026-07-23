package controllers

import (
	"errors"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"

	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

func (c *AuthController) Login(username string, password string) (string, *models.User, error) {
	var user models.User
	if err := c.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("username atau password salah")
		}
		return "", nil, err
	}

	// Verify password with bcrypt
	if err := utils.VerifyPassword(user.Password, password); err != nil {
		return "", nil, errors.New("username atau password salah")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}
