package authservice

import (
	jwtservice "sfit-platform-web-backend/services/jwtService"
	userservice "sfit-platform-web-backend/services/userService"
)

func Register(username, email, password string) (string, string, error) {
	user, err := userservice.CreateUser(username, email, password)
	if err != nil {
		return "", "", err
	}
	accessToken, err := jwtservice.GenerateToken(*user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwtservice.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
func Login(username, email, password string) (string, string, error) {
	user, err := userservice.GetUserByusernameOrEmail(username, email)
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
	refreshToken, err := jwtservice.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func Logout(token string) error {
	err := jwtservice.BlacklistToken(token)
	if err != nil {
		return err
	}
	return nil
}
func RefreshToken(refreshToken string) (string, string, error) {
	claim, err := jwtservice.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	sub, err := claim.GetSubject()
	if err != nil {
		return "", "", err
	}
	user, err := userservice.GetUserByID(sub)
	if err != nil {
		return "", "", err
	}
	accessToken, err := jwtservice.GenerateToken(*user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = jwtservice.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
