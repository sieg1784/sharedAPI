package service

import (
	"bookAPI/apiservice/data"
	"bookAPI/apiservice/utils"
	"time"

	"errors"

	"github.com/dgrijalva/jwt-go"
)

type UserToken struct {
	AccessToken string `json:"AccessToken"`
}

func LoginWithPassword(countryCode, phoneNum, password string) (UserToken, error) {

	user, err := data.QueryUserByUniqueKey(countryCode, phoneNum)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return UserToken{}, err
	}

	if password != user.Password {
		err = errors.New("error password")
		return UserToken{}, err
	}

	generateToken, err := GenerateToken(user)
	return generateToken, nil
}

func LoginAfterFirebaseAuth(countryCode, phoneNum string) (UserToken, error) {
	user, err := data.QueryUserByUniqueKey(countryCode, phoneNum)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return UserToken{}, err
	}

	generateToken, err := GenerateToken(user)
	return generateToken, nil
}

// 生成令牌  创建jwt风格的token
func GenerateToken(user data.User) (UserToken, error) {
	j := utils.NewJWT()
	claims := utils.CustomClaims{
		user.CountryCode,
		user.PhoneNum,
		user.Password,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),
			ExpiresAt: int64(time.Now().Unix() + 1209600),
			Issuer:    "cfun",
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		return UserToken{}, err
	}

	userToken := &UserToken{
		AccessToken: token,
	}

	return *userToken, err
}
