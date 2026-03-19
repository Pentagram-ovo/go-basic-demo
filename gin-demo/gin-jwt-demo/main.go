package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义Claims结构体（核心：嵌入标准Claims + 扩展业务字段)
type CustomClaims struct {
	UserId   uint64 `json:"userid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

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

// 全局变量：JWT签名密钥
var jwtSecret = []byte(getJWTSecret())

// 模拟用户数据库
var userDB = map[string]string{
	"admin": "admin123",
	"user1": "123456",
}

// 登陆函数
func Login(username, password string) (string, error) {
	//校验用户信息，检查用户的用户名是否存在及密码是否正确
	if pwd, ok := userDB[username]; !ok || pwd != password {
		return "", errors.New("用户的信息错误")
	}
	//生成用户的Id标识
	userid := uint64(1000 + len(username))
	role := "user"
	if username == "admin" {
		role = "admin"
	}
	//构建CustomClaims（设置过期时间、签发人、业务字段）
	//固定格式
	claims := CustomClaims{
		UserId:   userid,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			Issuer:    "Kunkun-jwt-demo",
			IssuedAt:  jwt.NewNumericDate(time.Now()),                            // 签发时间
			ID:        fmt.Sprintf("token_%d_%s", userid, time.Now().UnixNano()), // 唯一Token ID
		},
	}
	//生成Token（用HS256算法，传入Claims）
	//这下面的几步都差不多是固定格式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串并返回
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("生成Token失败：%w", err)
	}
	return tokenString, nil
}

// jwt的验证与解析函数
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 校验签名算法
		//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//	return nil, fmt.Errorf("不支持的签名算法：%v", token.Header["alg"])
		//}此为大致校验，更详细的可以像下面一样，检查是否和上面用的算法一致

		// 先验证算法类型是 HMAC
		// 再验证具体算法是预期的 HS256
		if hmacAlg, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法类型：%v", token.Header["alg"])
		} else if hmacAlg.Name != "HS256" {
			return nil, fmt.Errorf("不支持的 HMAC 子算法：%s，预期：HS256", hmacAlg.Name)
		}
		return jwtSecret, nil
	})
	//会报错，在之后解决下
	////处理解析错误
	//if err != nil {
	//	// 区分具体错误类型（便于前端提示）
	//	var ve *jwt.ValidationError
	//	if errors.As(err, &ve) {
	//		if ve.Errors&jwt.ValidationErrorExpired != 0 {
	//			return nil, errors.New("Token已过期")
	//		}
	//		return nil, fmt.Errorf("Token验证失败：%w", ve)
	//	}
	//	return nil, fmt.Errorf("解析Token失败：%w", err)
	//}
	//校验Token有效性，返回解析后的Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Token无效")
}

// jwt刷新函数
func RefreshToken(oldToken string) (string, error) {
	// 第一步：解析旧Token（不校验过期，只提取Claims）
	token, err := jwt.ParseWithClaims(
		oldToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if hmacAlg, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("不支持的签名算法类型：%v", token.Header["alg"])
			} else if hmacAlg.Name != "HS256" {
				return nil, fmt.Errorf("不支持的 HMAC 子算法：%s，预期：HS256", hmacAlg.Name)
			}
			return jwtSecret, nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("解析旧Token失败：%w", err)
	}

	// 第二步：提取旧Claims（类型断言）
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", errors.New("解析Claims失败")
	}

	// 第三步：判断是否在刷新宽限期内（过期≤10分钟）
	expireTime := claims.ExpiresAt.Time
	now := time.Now()
	if now.After(expireTime.Add(10 * time.Minute)) {
		return "", errors.New("Token过期超过10分钟，无法刷新")
	}

	// 第四步：生成新Token（复用用户信息，延长过期时间）
	newClaims := CustomClaims{
		UserId:   claims.UserId,
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
			Issuer:    "wangjiaxing-jwt-demo",
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("token_%d_%s", claims.UserId, now.UnixNano()),
		},
	}
	// 生成新Token并返回
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signedToken, err := newToken.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("生成新Token失败：%w", err)
	}
	return signedToken, nil
}

func main() {
	// 测试1：登录生成Token
	fmt.Println("=== 测试1：用户登录 ===")
	token, err := Login("admin", "admin123")
	if err != nil {
		fmt.Printf("登录失败：%v\n", err)
		return
	}
	fmt.Printf("登录成功，Token：%s\n\n", token)

	// 测试2：验证Token
	fmt.Println("=== 测试2：验证Token ===")
	claims, err := ValidateToken(token)
	if err != nil {
		fmt.Printf("验证失败：%v\n", err)
		return
	}
	fmt.Printf("验证成功，用户信息：\n")
	fmt.Printf("用户ID：%d\n用户名：%s\n角色：%s\n过期时间：%s\n\n",
		claims.UserId, claims.Username, claims.Role, claims.ExpiresAt.Time.Format("2006-01-02 15:04:05"))

	// 测试3：刷新Token
	fmt.Println("=== 测试3：刷新Token ===")
	newToken, err := RefreshToken(token)
	if err != nil {
		fmt.Printf("刷新失败：%v\n", err)
		return
	}
	fmt.Printf("刷新成功，新Token：%s\n\n", newToken)

	// 测试4：验证错误场景（错误密码）
	fmt.Println("=== 测试4：错误密码登录 ===")
	_, err = Login("admin", "123456")
	fmt.Printf("预期失败：%v\n", err)
}
