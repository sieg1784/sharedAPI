package data

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

type User struct {
	UserId       int
	CountryCode  string
	PhoneNum     string
	Password     string
	Roles        null.String
	PurchaseDate null.Time
	ExpireDate   null.Time
	Activate     bool
	Nickname     null.String
	Birthday     null.Time
	Gender       null.Int
	CreateDate   null.Time
}

func InsertUser(countryCode string, phoneNum string, password string, roles string) (int, error) {
	userId := 0
	now := time.Now()
	sqlStatement :=
		`INSERT INTO "user" (country_code, phone_num, password, roles, create_date)
		VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
	err := Db.QueryRow(sqlStatement, countryCode, phoneNum, password, roles, now).Scan(&userId)
	if err != nil {
		return 0, err
	} else {
		return userId, nil
	}
}

func QueryUserByUniqueKey(countryCode, phoneNum string) (User, error) {
	user := User{}
	sqlStatement := `SELECT * FROM "user" WHERE country_code = $1 and phone_num = $2`
	row := Db.QueryRow(sqlStatement, countryCode, phoneNum)
	err := row.Scan(&user.UserId, &user.CountryCode, &user.PhoneNum, &user.Password, &user.Roles,
		&user.PurchaseDate, &user.ExpireDate, &user.Activate, &user.Nickname, &user.Birthday, &user.Gender, &user.CreateDate)
	return user, err
}

func QueryUserByPK(userId int) (User, error) {
	user := User{}
	sqlStatement := `SELECT user_id, country_code, phone_num, roles
		, purchase_date, expire_date, activate, nickname, birthday, gender, create_date FROM "user" WHERE user_id = $1`
	row := Db.QueryRow(sqlStatement, userId)
	err := row.Scan(&user.UserId, &user.CountryCode, &user.PhoneNum, &user.Roles,
		&user.PurchaseDate, &user.ExpireDate, &user.Activate, &user.Nickname, &user.Birthday, &user.Gender, &user.CreateDate)
	fmt.Println(user)
	return user, err
}

func UpdateUserByPK(user User) (int, error) {
	sqlStatement :=
		`UPDATE "user" SET password = $1, roles = $2, purchase_date = $3, expire_date = $4, activate = $5
		, nickname = $6, birthday = $7, gender = $8
		WHERE user_id = $9 RETURNING user_id`
	err := Db.QueryRow(sqlStatement, user.Password, user.Roles, user.PurchaseDate, user.ExpireDate, user.Activate,
		user.Nickname, user.Birthday, user.Gender, user.UserId).Scan(&user.UserId)
	if err != nil {
		return 0, err
	} else {
		return user.UserId, nil
	}
}

func UpdateUserInfoByPK(user User, userId int) (int, error) {
	sqlStatement :=
		`UPDATE "user" SET nickname = $1, birthday = $2, gender = $3
		WHERE user_id = $4 RETURNING user_id`
	err := Db.QueryRow(sqlStatement, user.Nickname, user.Birthday, user.Gender, userId).Scan(&userId)
	if err != nil {
		return 0, err
	} else {
		return userId, nil
	}
}

func UpdateUserByPKDeprecated(nickname string, birthday time.Time, gender int, userId int) (int, error) {
	sqlStatement :=
		`UPDATE "user" SET nickname = $1, birthday = $2, gender = $3
		WHERE user_id = $4 RETURNING user_id`
	err := Db.QueryRow(sqlStatement, nickname, birthday, gender, userId).Scan(&userId)
	if err != nil {
		return 0, err
	} else {
		return userId, nil
	}
}

func DeleteUserByUniqueKey(countryCode string, phoneNum string) (bool, error) {
	sqlStatement :=
		`DELETE FROM "user" WHERE country_code = $1 and phone_num = $2`

	res, err := Db.Exec(sqlStatement, countryCode, phoneNum)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			return count == 1, nil
		} else {
			return false, err
		}
	} else {
		return false, err
	}

}
