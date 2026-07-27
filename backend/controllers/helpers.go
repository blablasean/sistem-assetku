package controllers

import "sistem-asetku-backend/utils"

func isAdmin(role string) bool {
	return role == utils.RoleAdmin
}

func canManageAssets(role string) bool {
	return utils.CanManageAssets(role)
}

func canMutate(role string) bool {
	return utils.CanMutate(role)
}

func canAssignEngineerRole(role string) bool {
	return utils.CanAssignWorker(role)
}

func canUpdateProgressRole(role string) bool {
	return utils.CanManageAssets(role) || role == utils.RoleEngineer
}

func canCloseOrDeleteWO(role string) bool {
	return utils.CanCloseWO(role)
}
