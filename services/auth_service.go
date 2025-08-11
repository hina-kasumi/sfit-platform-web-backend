package services

import ()

type AuthService struct {
	userSer    *UserService
	jwtSer     *JwtService
	refreshSer *RefreshTokenService
}

func NewAuthService(userSer *UserService, jwtSer *JwtService, refreshSer *RefreshTokenService) *AuthService {
	return &AuthService{
		userSer:    userSer,
		jwtSer:     jwtSer,
		refreshSer: refreshSer,
	}
}

func (authSer *AuthService) Register(username, email, password string) (string, string, error) {
	user, err := authSer.userSer.CreateUser(username, email, password)
	if err != nil {
		return "", "", err
	}
	accessToken, err := authSer.jwtSer.GenerateToken(*user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := authSer.refreshSer.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
func (authSer *AuthService) Login(username, email, password string) (string, string, error) {
	user, err := authSer.userSer.GetUserByusernameOrEmail(username, email)
	if err != nil {
		return "", "", err
	}
	if err := user.IsValidPassword(password); err != nil {
		return "", "", err
	}

	accessToken, err := authSer.jwtSer.GenerateToken(*user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := authSer.refreshSer.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (authSer *AuthService) Logout(token string) error {
	err := authSer.jwtSer.BlacklistToken(token)
	if err != nil {
		return err
	}
	return nil
}
func (authSer *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	claim, err := authSer.refreshSer.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	sub, err := claim.GetSubject()
	if err != nil {
		return "", "", err
	}
	user, err := authSer.userSer.GetUserByID(sub)
	if err != nil {
		return "", "", err
	}
	accessToken, err := authSer.jwtSer.GenerateToken(*user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = authSer.refreshSer.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
