package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Claims 表示项目当前使用的 JWT 载荷。
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwtv5.RegisteredClaims
}

// Manager 负责生成和解析 JWT。
type Manager struct {
	secret []byte
}

// NewManager 创建 JWT 管理器。
func NewManager(secret string) *Manager {
	return &Manager{
		secret: []byte(secret),
	}
}

// GenerateToken 生成登录后的 JWT。
func (m *Manager) GenerateToken(userID int64, username string) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwtv5.NewNumericDate(now),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken 解析并校验 JWT。
func (m *Manager) ParseToken(tokenText string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(tokenText, &Claims{}, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}

		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的 token")
	}

	return claims, nil
}
