package controllers

import (
	"errors"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"
	"strconv"

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

	// Save & update active session token (overwrites previous session cleanly)
	user.ActiveToken = token
	_ = c.db.Model(&user).Update("active_token", token)

	// Record login in activity log
	utils.RecordActivity(c.db, "AUTH", user.Username,
		"Login ke sistem sebagai "+user.Role+" ("+user.Name+")",
		strconv.Itoa(user.ID))

	return token, &user, nil
}

func (c *AuthController) Logout(userID int) error {
	var user models.User
	if err := c.db.First(&user, userID).Error; err == nil {
		// Record logout in activity log
		utils.RecordActivity(c.db, "AUTH", user.Username,
			"Logout dari sistem",
			strconv.Itoa(userID))
	}
	return c.db.Model(&models.User{}).Where("id = ?", userID).Update("active_token", "").Error
}
