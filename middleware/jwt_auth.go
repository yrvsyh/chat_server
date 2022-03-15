package middleware

import (
	"chat_server/database"
	"chat_server/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type (
	JWTClaims struct {
		jwt.StandardClaims
		Name string `json:"name,omitempty"`
		// Role int    `json:"role,omitempty"`
	}
)

var (
	tokenSecretKey = []byte("secret_key")
	// token过期时间
	tokenExpireDuration = time.Hour * 24
	// token可刷新时间
	tokenRefreshDuration = time.Hour * 24 * 30
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := GetToken(c)
		claims, err := ParseToken(tokenString)
		if err != nil {
			utils.Error(c, -1, "Token解析失败")
			c.Abort()
		} else {
			if err := claims.Valid(); err == nil {
				log.WithFields(log.Fields{
					"claimsId": claims.Id,
					"name":     claims.Name,
				}).Info("TOKEN VALID")
				// // token黑名单判断
				// if _, err := database.RDB.Get(claims.StandardClaims.Id).Result(); err == redis.Nil {
				c.Next()
				return
				// }
			}
			utils.Error(c, -1, "无效的Token")
			c.Abort()
		}

	}
}

// 从http头获取token
func GetToken(c *gin.Context) string {
	tokenString := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		authHeaders := strings.Split(authHeader, " ")
		if len(authHeaders) == 2 && authHeaders[0] == "Bearer" {
			tokenString = authHeaders[1]
		}
	}
	return tokenString
}

// 生成token
// func GenToken(name string, role int) (string, error) {
func GenToken(name string) (string, error) {
	claims := &JWTClaims{
		Name: name,
		// Role: role,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewV4().String(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tokenSecretKey)
}

// 解析token
func ParseToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenSecretKey, nil
	})
	return claims, err
}

// 刷新token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err == nil {
		if err = claims.Valid(); err == nil {
			if time.Unix(claims.ExpiresAt, 0).Add(tokenRefreshDuration).After(time.Now()) {
				claims.StandardClaims.Id = uuid.NewV4().String()
				// 重置过期时间
				claims.ExpiresAt = time.Now().Add(tokenExpireDuration).Unix()
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				return token.SignedString(tokenSecretKey)
			}
		}
	}
	return "", err
}

// 使token失效, 添加至redis黑名单
func DelToken(tokenString string) error {
	claims, err := ParseToken(tokenString)
	if err == nil {
		if VerifyToken(tokenString) {
			tokenId := claims.StandardClaims.Id
			expireTime := time.Until(time.Unix(claims.ExpiresAt, 0))
			err = database.RDB.Set(tokenId, 1, expireTime).Err()
		}
	}
	return err
}

// 验证token是否过期
func VerifyToken(tokenString string) bool {
	claims, err := ParseToken(tokenString)
	if err == nil {
		if err = claims.Valid(); err == nil {
			return true
		}
	}
	return false
}
