package controllers

func isAdmin(role string) bool {
	return role == "admin"
}

func canManageAssets(role string) bool {
	return role == "admin" || role == "hod" || role == "management"
}

func canMutate(role string) bool {
	return role == "admin" || role == "hod"
}

func canAssignEngineerRole(role string) bool {
	return role == "admin" || role == "hod" || role == "management"
}

func canUpdateProgressRole(role string) bool {
	return role == "admin" || role == "hod" || role == "management" || role == "engineer"
}

func canCloseOrDeleteWO(role string) bool {
	return role == "admin" || role == "hod" || role == "management"
}
