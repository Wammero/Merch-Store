package validators

func IsValidPassword(password string) bool {
	return len(password) > 0
}

func IsValidUsername(username string) bool {
	return len(username) > 0
}
