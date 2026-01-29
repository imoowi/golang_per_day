package apisixtoken

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Key    string `json:"key"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}
type JWTGenerator struct {
	secretKey []byte
	issuer    string
	expiresIn time.Duration
}

func NewJWTGenerator(secret string, issuer string, expiresAt time.Duration) *JWTGenerator {
	return &JWTGenerator{
		secretKey: []byte(secret),
		issuer:    issuer,
		expiresIn: expiresAt,
	}
}

// 生成JWT令牌
func (g *JWTGenerator) GenerateToken(userID uint) (string, error) {
	return g.GenerateTokenWithRole(userID, "")
}

// 生成带角色的JWT令牌
func (g *JWTGenerator) GenerateTokenWithRole(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Key:    g.issuer, // 必须与APISIX consumer的key一致
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.issuer,
			Subject:   cast.ToString(userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        generateUUID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.secretKey)
}

// 验证JWT
func (g *JWTGenerator) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return g.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func generateUUID() string {
	// 使用 crypto/rand 生成 16 字节随机数
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// 如果读取失败，退而求其次用当前时间纳秒数
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	// 格式化为标准 UUID 形式：xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
