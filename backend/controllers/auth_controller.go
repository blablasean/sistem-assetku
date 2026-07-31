package controllers

import (
	"errors"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"
	"strconv"
	"time"

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

	// Single Active Session Check:
	// Only reject secondary login if user has a valid active_token AND tab/browser is currently active (last_seen_at <= 30s)
	if user.IsLoggedIn && user.ActiveToken != "" {
		isAbandoned := user.LastSeenAt == nil || time.Since(*user.LastSeenAt) > 30*time.Second
		_, tokenErr := utils.ValidateToken(user.ActiveToken)

		if tokenErr == nil && !isAbandoned {
			return "", nil, errors.New("Akun sedang tersambung di perangkat lain")
		}
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	// Direct raw SQL update to guarantee is_logged_in, active_token, and last_seen_at are saved to MySQL
	now := time.Now()
	user.IsLoggedIn = true
	user.ActiveToken = token
	user.LastSeenAt = &now

	_ = c.db.Exec("UPDATE users SET is_logged_in = 1, active_token = ?, last_seen_at = ? WHERE id = ?", token, now, user.ID).Error

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
	return c.db.Exec("UPDATE users SET is_logged_in = 0, active_token = '', last_seen_at = NULL WHERE id = ?", userID).Error
}
