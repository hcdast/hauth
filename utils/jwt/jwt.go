/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:47
 * @LastEditors: hc
 * @LastEditTime: 2021-06-08 14:16:58
 * @Description: 鉴权
 */
package jwt

import (
	"errors"
	"time"

	"example-hauth/utils/logs"

	jwt "github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	*jwt.StandardClaims
	UserId      string
	DomainId    string
	OrgUnitId   string
	Authorities string `json:"authorities"`
}

var (
	key []byte = []byte("hc")
)

// 生成token
func GenToken(user_id, domain_id, org_id string, dt int64) string {
	claims := JwtClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + dt,
			Issuer:    "hc",
		},
		user_id,
		domain_id,
		org_id,
		"ROLE_ADMIN,AUTH_WRITE,ACTUATOR",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		logs.Error(err)
		return ""
	}
	return ss
}

// 销毁token
func DestoryToken() string {
	claims := JwtClaims{
		&jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() - 99999),
			Issuer:    "hc",
		},
		"exit",
		"exit",
		"exit",
		"ROLE_ADMIN,AUTH_WRITE,ACTUATOR",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		logs.Error(err)
		return ""
	}
	return ss
}

// 验证token
func CheckToken(token string) bool {
	_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return false
	}
	return true
}

// 解析token
func ParseJwt(token string) (*JwtClaims, error) {
	var jclaim = &JwtClaims{}
	_, err := jwt.ParseWithClaims(token, jclaim, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, errors.New("parase with claims failed.")
	}
	return jclaim, nil
}
