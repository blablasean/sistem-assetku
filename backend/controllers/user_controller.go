package controllers

import (
	"errors"
	"strconv"
	"strings"
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
	r := strings.ToLower(strings.TrimSpace(callerRole))
	if r == "external" || r == "" {
		return errors.New("akses ditolak: anda harus login terlebih dahulu")
	}
	if u.Username == "" || u.Password == "" || u.Name == "" || u.Role == "" {
		return errors.New("username, password, nama, dan role wajib diisi")
	}

	var count int64
	c.db.Model(&models.User{}).Where("LOWER(username) = ?", strings.ToLower(u.Username)).Count(&count)
	if count > 0 {
		return errors.New("username '" + u.Username + "' sudah digunakan oleh akun lain")
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

	// Record in activity log
	utils.RecordActivity(c.db, "USER_MANAGEMENT", callerRole,
		"Menambahkan pengguna baru: "+u.Name+" ("+u.Username+") dengan role "+u.Role,
		strconv.Itoa(u.ID))

	return nil
}

// EditUser updates existing user credentials or profile fields
func (c *UserController) EditUser(u *models.User, callerRole string) error {
	r := strings.ToLower(strings.TrimSpace(callerRole))
	if r == "external" || r == "" {
		return errors.New("akses ditolak: anda harus login terlebih dahulu untuk mengedit data pengguna")
	}
	var existing models.User
	if err := c.db.First(&existing, u.ID).Error; err != nil {
		return errors.New("pengguna tidak ditemukan (ID #" + strconv.Itoa(u.ID) + ")")
	}

	oldRole := existing.Role
	if u.Name != "" {
		existing.Name = u.Name
	}
	if u.Role != "" {
		existing.Role = u.Role
	}
	if u.Avatar != "" {
		existing.Avatar = u.Avatar
	}

	// Validate & update username if changed
	if u.Username != "" && strings.ToLower(u.Username) != strings.ToLower(existing.Username) {
		var count int64
		c.db.Model(&models.User{}).Where("LOWER(username) = ? AND id != ?", strings.ToLower(u.Username), u.ID).Count(&count)
		if count > 0 {
			return errors.New("username '" + u.Username + "' sudah digunakan oleh akun lain")
		}
		existing.Username = u.Username
	}

	// Validate & update password if non-empty
	if u.Password != "" {
		if len(u.Password) < 4 {
			return errors.New("password baru minimal 4 karakter")
		}
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return errors.New("gagal memproses enkripsi password baru: " + err.Error())
		}
		existing.Password = hashedPassword
	}

	if err := c.db.Save(&existing).Error; err != nil {
		return errors.New("gagal menyimpan data pengguna: " + err.Error())
	}
	existing.Password = ""
	*u = existing

	// Record in activity log — note role changes specifically
	logMsg := "Mengedit pengguna: " + u.Name + " (" + u.Username + ")"
	if oldRole != u.Role {
		logMsg += " — role diubah dari " + oldRole + " menjadi " + u.Role
	}
	utils.RecordActivity(c.db, "USER_MANAGEMENT", callerRole, logMsg, strconv.Itoa(u.ID))

	return nil
}

// DeleteUser deletes a non-admin user account
func (c *UserController) DeleteUser(userID int, callerRole string) error {
	r := strings.ToLower(strings.TrimSpace(callerRole))
	if r == "external" || r == "" {
		return errors.New("akses ditolak: anda harus login terlebih dahulu")
	}
	var existing models.User
	if err := c.db.First(&existing, userID).Error; err != nil {
		return errors.New("pengguna tidak ditemukan")
	}
	if existing.Username == "admin" {
		return errors.New("akun Administrator utama (admin) tidak boleh dihapus")
	}
	if err := c.db.Delete(&existing).Error; err != nil {
		return err
	}
	utils.RecordActivity(c.db, "USER_MANAGEMENT", callerRole,
		"Menghapus akun pengguna: "+existing.Name+" ("+existing.Username+")",
		strconv.Itoa(userID))
	return nil
}
