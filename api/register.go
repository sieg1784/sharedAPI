package api

import (
	"bookAPI/apiservice/data"
	"bookAPI/apiservice/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"
)

func Register(c *gin.Context) {
	m := make(map[string]interface{})
	countryCode := c.PostForm("countryCode")
	phoneNum := c.PostForm("phoneNum")

	if strings.Trim(countryCode, " ") == "" || strings.Trim(phoneNum, " ") == "" {
		m["systemCode"] = 410
		m["systemMessage"] = "請確實填寫國碼與手機號碼。"
		c.JSON(200, m)
		return
	}

	password := c.PostForm("password")
	if strings.Trim(password, " ") == "" {
		m["systemCode"] = 412
		m["systemMessage"] = "請確實填寫密碼。"
		c.JSON(200, m)
		return
	}

	idTokenFromFirebase := c.GetHeader("Authorization")
	if strings.Trim(idTokenFromFirebase, " ") == "" {
		m["systemCode"] = 411
		m["systemMessage"] = "驗證碼有誤。"
		c.JSON(200, m)
		return
	}

	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)
	_, err := firebaseAuth.VerifyIDToken(context.Background(), idTokenFromFirebase)
	// _, err := utils.DecodeFirebaseIdToken(idTokenFromFirebase)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	roles := "freeUser"

	phoneNum = utils.NormalizePhoneNum(phoneNum)
	userId, err := data.InsertUser(countryCode, phoneNum, password, roles)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "uniqueUser"` {
		m["systemCode"] = 410
		m["systemMessage"] = "這個號碼已經有人註冊囉。"
		c.JSON(200, m)
	} else if err != nil {
		c.JSON(500, err.Error)
	} else {
		c.JSON(200, gin.H{
			"systemCode":    200,
			"systemMessage": "OK",
			"data": gin.H{
				"userId": userId,
				"roles":  roles,
			},
		})
	}
}

func UpdateUserInfo(c *gin.Context) {
	logger := c.MustGet("logger").(*log.Logger)
	userId := c.PostForm("userId")
	if strings.Trim(userId, " ") == "" {
		c.JSON(400, gin.H{
			"systemCode":    410,
			"systemMessage": "please do not update user data which not belongs to you",
		})
		return
	}

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"systemCode":    410,
			"systemMessage": "userId的格式不對喔。",
		})
	}

	nickname := c.PostForm("nickname")
	birthday := c.PostForm("birthday")
	gender := c.PostForm("gender")
	logger.Println(nickname + ", " + birthday + ", " + gender + ", " + userId)
	if nickname == "" && birthday == "" && gender == "" {
		c.JSON(200, gin.H{
			"systemCode":    412,
			"systemMessage": "至少要輸入一項資訊喔。",
		})
		return
	}

	user, err := data.QueryUserByPK(intUserId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		c.JSON(200, gin.H{
			"systemCode":    412,
			"systemMessage": "查無此使用者。",
		})
		return
	} else if err != nil {
		c.JSON(500, err)
		return
	}

	err = user.Nickname.UnmarshalText([]byte(nickname))
	var timeBirthday time.Time
	if birthday != "" {
		i, err := strconv.ParseInt(birthday, 10, 64)
		if err != nil {
			c.JSON(200, gin.H{
				"systemCode":    412,
				"systemMessage": "時間格式有問題。",
			})
			return
		}

		// timeBirthday = time.Unix(0, i*int64(time.Millisecond))
		timeBirthday = utils.MilliTimestamp2Time(i)
		user.Birthday = null.NewTime(timeBirthday, true)
	}

	if gender != "" {
		user.Gender = null.IntFrom(utils.GenderString2Int(gender))
	}

	intUserId, err = data.UpdateUserInfoByPK(user, intUserId)
	if err != nil {
		c.JSON(500, err)
		return
	} else {
		c.JSON(200, gin.H{
			"systemCode":    200,
			"systemMessage": "ok",
			"data": gin.H{
				"userId":   userId,
				"nickname": user.Nickname.ValueOrZero(),
				"birthday": utils.Time2MilliTimestamp(user.Birthday),
				"gender":   utils.GenderInt2String(user.Gender.Int64),
			},
		})
	}
	return
}

func ResetPassword(c *gin.Context) {
	idTokenFromFirebase := c.GetHeader("Authorization")
	if strings.Trim(idTokenFromFirebase, " ") == "" {
		c.JSON(200, gin.H{
			"systemCode":    411,
			"systemMessage": "驗證碼有誤",
		})
		return
	}

	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)
	decodedIdToken, err := firebaseAuth.VerifyIDToken(context.Background(), idTokenFromFirebase)
	// decodedIdToken, err := utils.DecodeFirebaseIdToken(idTokenFromFirebase)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	countryCode := c.PostForm("countryCode")
	phoneNum := c.PostForm("phoneNum")
	if strings.Trim(countryCode, " ") == "" || strings.Trim(phoneNum, " ") == "" {
		c.JSON(200, gin.H{
			"systemCode":    410,
			"systemMessage": "請確實填寫國碼與手機號碼。",
		})
		return
	}

	str := decodedIdToken.Firebase.Identities["phone"] // str only used in next line, forgive this namming
	phoneFromFirebase := utils.TrimBrackets(fmt.Sprintf("%v", str))
	countryCodeFromFirebase := utils.ExtractCountryCode(phoneFromFirebase)
	phoneNumFromFirebase := utils.ExtractPhoneNum(phoneFromFirebase)
	if countryCodeFromFirebase != countryCode || phoneNumFromFirebase != phoneNum {
		c.JSON(200, gin.H{
			"systemCode":    410,
			"systemMessage": "手機號碼有錯喔。",
		})
	}

	password := c.PostForm("password")
	if strings.Trim(password, " ") == "" {
		c.JSON(200, gin.H{
			"systemCode":    412,
			"systemMessage": "請確實填寫密碼",
		})
		return
	}

	user, err := data.QueryUserByUniqueKey(countryCode, phoneNum)
	if err != nil && err.Error() == "sql: no rows in result set" {
		c.JSON(200, gin.H{
			"systemCode":    410,
			"systemMessage": "查無此使用者。",
		})
		return
	} else if err != nil {
		c.JSON(500, err)
		return
	}

	user.Password = password
	_, err = data.UpdateUserByPK(user)
	if err != nil {
		c.JSON(500, err)
		return
	} else {
		c.JSON(200, gin.H{
			"systemCode":    200,
			"systemMessage": "ok",
			"data": gin.H{
				"userId": user.UserId,
			},
		})
	}
	return
}

func CheckRegisterAvailable(c *gin.Context) {
	countryCode := c.PostForm("countryCode")
	phoneNum := c.PostForm("phoneNum")
	if strings.Trim(countryCode, " ") == "" || strings.Trim(phoneNum, " ") == "" {
		c.JSON(200, gin.H{
			"systemCode":    410,
			"systemMessage": "請確實填寫國碼與手機號碼。",
		})
		return
	}

	_, err := data.QueryUserByUniqueKey(countryCode, phoneNum)
	if err != nil && err.Error() == "sql: no rows in result set" {
		c.JSON(200, gin.H{
			"systemCode":    200,
			"systemMessage": "此手機號碼可以使用。",
		})
		return
	} else if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, gin.H{
		"systemCode":    412,
		"systemMessage": "此手機號碼已被使用，要立即登入嗎？",
	})
}
