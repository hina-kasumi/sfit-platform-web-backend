package authservice

import (
	userrepository "sfit-platform-web-backend/repositories/userRepository"
	jwtservice "sfit-platform-web-backend/services/jwtService"
)

func Register(username, email, password string) (string, string, error) {
	user, err := userrepository.CreateUser(username, email, password)
	if err != nil {
		return "", "", err
	}
	accessToken, err := jwtservice.GenerateToken(*user)
	if err != nil {
		return "", "", err
	}

	return accessToken, "", nil
}
func Login(username, email, password string) (string, string, error) {
	user, err := userrepository.GetUserByusernameOrEmail(username, "")
	if err != nil {
		return "", "", err
	}
	if err := user.IsValidPasswrod(password); err != nil {
		return "", "", err
	}

	accessToken, err := jwtservice.GenerateToken(*user)
	if err != nil {
		return "", "", err
	}

	return accessToken, "", nil
}

func Logout(token string) error {
	err := jwtservice.BlacklistToken(token)
	if err != nil {
		return err
	}
	return nil
}
func RefreshToken(token string) (string, string, error) {
	// Logic to refresh a user's token
	return "", "", nil
}
