package utils

import "strings"

// Role constants
const (
	RoleAdmin        = "admin"
	RoleHOD          = "hod"
	RoleManagement   = "management"
	RoleEngineer     = "engineer"
	RoleExternal     = "external"
)

// CanManageAssets checks if caller role can perform asset management & PM scheduling
func CanManageAssets(role string) bool {
	r := strings.ToLower(strings.TrimSpace(role))
	return r == RoleAdmin || r == RoleHOD || r == RoleManagement
}

// CanMutate checks if caller role can perform asset location mutations
func CanMutate(role string) bool {
	r := strings.ToLower(strings.TrimSpace(role))
	return r == RoleAdmin || r == RoleHOD
}

// CanAssignWorker checks if caller role can assign workers to work orders
func CanAssignWorker(role string) bool {
	r := strings.ToLower(strings.TrimSpace(role))
	return r == RoleAdmin || r == RoleHOD || r == RoleManagement
}

// CanCloseWO checks if caller role can close or verify completed work orders
func CanCloseWO(role string) bool {
	r := strings.ToLower(strings.TrimSpace(role))
	return r == RoleAdmin || r == RoleHOD || r == RoleManagement
}
