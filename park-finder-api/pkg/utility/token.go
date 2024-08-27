package utility

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func GenerateToken(email string, role string) string {
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(2 * 365 * 24 * time.Hour).Unix(),
		"user": map[string]string{
			"email": email,
			"role":  role,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, err := loadSecretKey()
	if err != nil {
		fmt.Println("loadSecretKey", err)
		return ""
	}

	signedToken, err := token.SignedString(key)
	if err != nil {
		fmt.Println("SignedString", err)
		return ""
	}

	return signedToken
}

func GetEmailFromToken(auth string) string {
	token, err := jwt.Parse(auth, KeyFunc, jwt.WithoutClaimsValidation())
	if err != nil {
		return ""
	}

	claims := token.Claims.(jwt.MapClaims)
	user := claims["user"].(map[string]interface{})
	email := user["email"].(string)
	return email
}

func GetRoleFromToken(auth string) string {
	token, err := jwt.Parse(auth, KeyFunc, jwt.WithoutClaimsValidation())
	if err != nil {
		return ""
	}

	claims := token.Claims.(jwt.MapClaims)
	user := claims["user"].(map[string]interface{})
	role := user["role"].(string)
	return role
}

func GetExpireDateFromToken(auth string) int64 {
	token, err := jwt.Parse(auth, KeyFunc, jwt.WithoutClaimsValidation())
	if err != nil {
		return 0
	}

	claims := token.Claims.(jwt.MapClaims)
	var expireDate int64
	switch exp := claims["exp"].(type) {
	case float64:
		expireDate = int64(exp)
	}
	fmt.Println(expireDate)
	return expireDate
}

func Skipper(c echo.Context) bool {
	skipPath := map[string]bool{
		"/user/login":    true,
		"/user/register": true,
	}
	path := c.Request().URL.Path
	return skipPath[path]
}

func KeyFunc(t *jwt.Token) (interface{}, error) {
	method := t.Method.Alg()
	if method != "HS256" {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}

	key, err := loadSecretKey()
	if err != nil {
		return nil, err
	}

	return key, nil
}

func loadSecretKey() ([]byte, error) {
	keyStr := os.Getenv("JWT_SECRET_KEY")
	if keyStr == "" {
		return nil, fmt.Errorf("cannot load jwt secret key")
	}

	key, err := secretToBytes(keyStr)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func secretToBytes(secret string) ([]byte, error) {
	if len(secret) != 64 {
		return nil, fmt.Errorf("key length not equal to 64 (32 bytes or 256 bits)")
	}

	match, err := regexp.MatchString("^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$", secret)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, fmt.Errorf("input string is not base64 encoded")
	}

	key, err := hex.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	return key, nil
}
