package controllers

import (
	"errors"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"

	"gorm.io/gorm"
)

// UserController handles business logic for user management and role operations
type UserController struct {
	db *gorm.DB
}

// NewUserController creates a new instance of UserController
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

// GetAllUsers retrieves all registered user accounts with sanitized passwords
func (c *UserController) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := c.db.Order("id ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	for i := range users {
		users[i].Password = ""
	}
	return users, nil
}

// GetEngineers retrieves all active staff engineers
func (c *UserController) GetEngineers() ([]models.User, error) {
	var engineers []models.User
	if err := c.db.Where("role = ?", "engineer").Order("name ASC").Find(&engineers).Error; err != nil {
		return nil, err
	}
	for i := range engineers {
		engineers[i].Password = ""
	}
	return engineers, nil
}

// CreateUser registers a new user with bcrypt password hashing
func (c *UserController) CreateUser(u *models.User, callerRole string) error {
	if !isAdmin(callerRole) {
		return errors.New("akses ditolak: hanya Admin yang dapat menambah pengguna baru")
	}
	if u.Username == "" || u.Password == "" || u.Name == "" || u.Role == "" {
		return errors.New("username, password, nama, dan role wajib diisi")
	}

	var existing models.User
	if err := c.db.Where("username = ?", u.Username).First(&existing).Error; err == nil {
		return errors.New("username sudah digunakan, silakan pilih username lain")
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return errors.New("gagal memproses enkripsi password")
	}
	u.Password = hashedPassword

	if err := c.db.Create(u).Error; err != nil {
		return err
	}
	u.Password = ""
	return nil
}

// EditUser updates existing user credentials or profile fields
func (c *UserController) EditUser(u *models.User, callerRole string) error {
	if !isAdmin(callerRole) {
		return errors.New("akses ditolak: hanya Admin yang dapat mengedit data pengguna")
	}
	var existing models.User
	if err := c.db.First(&existing, u.ID).Error; err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	existing.Name = u.Name
	existing.Role = u.Role
	existing.Username = u.Username
	if u.Avatar != "" {
		existing.Avatar = u.Avatar
	}

	if u.Password != "" {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return errors.New("gagal memproses enkripsi password baru")
		}
		existing.Password = hashedPassword
	}

	if err := c.db.Save(&existing).Error; err != nil {
		return err
	}
	existing.Password = ""
	*u = existing
	return nil
}

// DeleteUser deletes a non-admin user account
func (c *UserController) DeleteUser(userID int, callerRole string) error {
	if !isAdmin(callerRole) {
		return errors.New("akses ditolak: hanya Admin yang dapat menghapus pengguna")
	}
	var existing models.User
	if err := c.db.First(&existing, userID).Error; err != nil {
		return errors.New("pengguna tidak ditemukan")
	}
	if existing.Username == "admin" {
		return errors.New("akun Administrator utama tidak boleh dihapus")
	}
	return c.db.Delete(&existing).Error
}
