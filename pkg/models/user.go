package models

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/hex-techs/blade/pkg/util/config"
	"github.com/hex-techs/blade/pkg/util/token"
	"golang.org/x/crypto/pbkdf2"
)

// User 用户表
type User struct {
	Base
	// 用户名，默认英文，唯一，不可为空
	Name string `gorm:"unique,index,size:64,not null" json:"name"`
	// 中文名称
	CnName string `gorm:"size:64" json:"cnName"`
	// 用户密码，加密后存储
	Password string `gorm:"size:1024" json:"password"`
	// 用户邮箱，唯一，不可为空
	Email string `gorm:"unique,index,not null" json:"email"`
	// 是否是管理员
	Admin bool `gorm:"default:false" json:"admin"`
	// 是否有效
	Enabled bool `gorm:"default:true" json:"enabled"`
	// 电话号码
	Phone string `gorm:"size:32" json:"phone"`
	// 社交账号，如微信，qq，钉钉，lark等
	IM string `gorm:"size:128" json:"im"`
	// 用户token，不存储在数据库中
	Token *Token `gorm:"-" json:"token"`
	// 用户角色，不存储在数据库中
	Roles []string `gorm:"-" json:"roles"`
}

// Token response user's token
type Token struct {
	// token信息
	Token string `json:"token"`
	// 过期时间
	Expired int64 `json:"expired"`
}

// EncodePasswd encodes password to safe format.
func (u *User) EncodePasswd() {
	newPasswd := pbkdf2.Key([]byte(u.Password), []byte(salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// ValidatePassword checks if given password matches the one belongs to the user.
func (u *User) ValidatePassword(Password string) bool {
	newUser := &User{Password: Password}
	newUser.EncodePasswd()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}

// GenUser generate User
func (u *User) GenUser() error {
	claim := &token.Claims{
		ID:             u.ID,
		Name:           u.Name,
		Admin:          u.Admin,
		StandardClaims: jwt.StandardClaims{},
	}
	t, e, err := token.GenerateJWTToken(claim, config.Read().Service.TokenExpired)
	if err != nil {
		return err
	}
	u.Token = &Token{
		Token:   t,
		Expired: e,
	}
	return nil
}

// TruncatePassword return null password to client
func (u *User) TruncatePassword() {
	u.Password = ""
}
