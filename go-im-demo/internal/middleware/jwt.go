package middleware

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 读取JWT密钥的函数
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET_KEY")
	// 补充：判断生产环境/开发环境，处理secret为空的情况
	if secret == "" {
		//GO_ENV=production：生产环境（对外提供服务的正式环境）；
		//GO_ENV=development/test：开发 / 测试环境（本地调试用）。
		//这行代码的作用是「判断当前是否是生产环境」，为后续差异化处理做依据
		if os.Getenv("GO_ENV") == "production" {
			panic("错误：未设置JWT_SECRET_KEY环境变量！")
		}
		return "dev_default_secret_987654321"
	}
	return secret
}

var jwtSecret = []byte(getJWTSecret())

type Claims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// GenerateToken 生成token
func GenerateToken(userid uint, username string) (string, error) {
	//构建CustomClaims（设置过期时间、签发人、业务字段）
	//固定格式
	claims := Claims{
		UserID:   userid,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "Kunkun-go-forum",
			IssuedAt:  jwt.NewNumericDate(time.Now()),                            // 签发时间
			ID:        fmt.Sprintf("token_%d_%s", userid, time.Now().UnixNano()), // 唯一Token ID
		},
	}
	//生成Token（用HS256算法，传入Claims）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串并返回
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("生成Token失败：%w", err)
	}
	return tokenString, nil
}

// ParseToken jwt的验证与解析函数
func ParseToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if hmacAlg, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法类型：%v", token.Header["alg"])
		} else if hmacAlg.Name != "HS256" {
			return nil, fmt.Errorf("不支持的 HMAC 子算法：%s，预期：HS256", hmacAlg.Name)
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	//校验Token有效性，返回解析后的Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Token无效")
}
