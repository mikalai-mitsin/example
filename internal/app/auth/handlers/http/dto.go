package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
)

type TokenPairDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenPairDTO(entity entities.TokenPair) (TokenPairDTO, error) {
	dto := TokenPairDTO{AccessToken: entity.Access.String(), RefreshToken: entity.Refresh.String()}
	return dto, nil
}

type ObtainTokenDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewObtainTokenDTO(r *http.Request) (ObtainTokenDTO, error) {
	update := ObtainTokenDTO{}
	if err := render.DecodeJSON(r.Body, &update); err != nil {
		return ObtainTokenDTO{}, err
	}
	return update, nil
}
func (dto ObtainTokenDTO) toEntity() (entities.Login, error) {
	login := entities.Login{Email: dto.Email, Password: dto.Password}
	return login, nil
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

func NewRefreshTokenDTO(r *http.Request) (RefreshTokenDTO, error) {
	dto := RefreshTokenDTO{}
	if err := render.DecodeJSON(r.Body, &dto); err != nil {
		return RefreshTokenDTO{}, err
	}
	return dto, nil
}
func (dto RefreshTokenDTO) toEntity() (entities.Token, error) {
	token := entities.Token(dto.RefreshToken)
	return token, nil
}
