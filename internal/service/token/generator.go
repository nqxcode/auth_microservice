package token

import (
	"time"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/auth_microservice/internal/utils"
)

type generator struct {
}

func NewGenerator() *generator {
	return &generator{}
}

func (g *generator) GenerateToken(info model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	token, err := utils.GenerateToken(info, secretKey, duration)
	if err != nil {
		return "", err
	}

	return token, nil
}
