package api

import (
	"bookAPI/apiservice/data"
	"bookAPI/apiservice/service"
	"bookAPI/apiservice/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func QueryUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	user, _ := data.QueryUserByPK(userId)
	log.Println(user.ExpireDate)
	c.JSON(200, gin.H{
		"userId":       user.UserId,
		"countryCode":  user.CountryCode,
		"phoneNum":     user.PhoneNum,
		"roles":        user.Roles,
		"purchaseDate": utils.Time2MilliTimestamp(user.PurchaseDate),
		"expiraeDate":  utils.Time2MilliTimestamp(user.ExpireDate),
		"activate":     user.Activate,
		"nickname":     user.Nickname,
		"birthday":     utils.Time2MilliTimestamp(user.Birthday),
		"gender":       utils.GenderInt2String(user.Gender.Int64),
		"createDate":   utils.Time2MilliTimestamp(user.CreateDate),
	})
}

func DeleteUser(c *gin.Context) {
	countryCode := c.PostForm("countryCode")
	phoneNum := c.PostForm("phoneNum")
	logger := c.MustGet("logger").(*log.Logger)
	logger.Println(countryCode + ", " + phoneNum)
	result, err := data.DeleteUserByUniqueKey(countryCode, phoneNum)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, result)
	}
}

func LoginWithPassword(c *gin.Context) {

	countryCode := c.PostForm("countryCode")
	phoneNum := c.PostForm("phoneNum")
	password := c.PostForm("password")
	if strings.Trim(countryCode, " ") == "" || strings.Trim(phoneNum, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(200, gin.H{
			"systemCode":    410,
			"systemMessage": "請確實填寫手機號碼與密碼。",
		})
		return
	}
	token, err := service.LoginWithPassword(countryCode, phoneNum, password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(200, gin.H{
				"systemCode":    410,
				"systemMessage": "手機號碼或是密碼有錯喔。", // 查無此使用者, 但避免透落太多訊息給前端, 所以改成這樣
			})
			return
		} else if err.Error() == "error password" {
			c.JSON(200, gin.H{
				"systemCode":    410,
				"systemMessage": "手機號碼或是密碼有錯喔。",
			})
			return
		} else {
			c.JSON(500, err)
			return
		}
	}

	user, _ := data.QueryUserByUniqueKey(countryCode, phoneNum)
	c.Header("access_token", token.AccessToken)
	c.JSON(200, gin.H{
		"systemCode":    200,
		"systemMessage": "OK",
		"data": gin.H{
			"userId":   user.UserId,
			"roles":    user.Roles,
			"nickname": user.Nickname.ValueOrZero(),
			"birthday": utils.Time2MilliTimestamp(user.Birthday),
			"gender":   utils.GenderInt2String(user.Gender.Int64),
		},
	})
	return
}

func LoginWithIdToken(c *gin.Context) {
	countryCode := c.PostForm("countryCode")
	phoneNum := c.PostForm("phoneNum")
	if strings.Trim(countryCode, " ") == "" || strings.Trim(phoneNum, " ") == "" {
		c.JSON(200, gin.H{
			"systemCode":    410,
			"systemMessage": "請確實填寫手機號碼與密碼。",
		})
		return
	}

	idTokenFromFirebase := c.GetHeader("Authorization")
	if idTokenFromFirebase == "" {
		c.JSON(200, gin.H{
			"systemCode":    410,
			"systemMessage": "Id token is missing"})
		c.Abort()
		return
	}

	//verify token
	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)
	decodedIdToken, err := firebaseAuth.VerifyIDToken(context.Background(), idTokenFromFirebase)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
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
		return
	}

	token, err := service.LoginAfterFirebaseAuth(countryCodeFromFirebase, phoneNumFromFirebase)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(200, gin.H{
				"systemCode":    410,
				"systemMessage": "手機號碼有錯喔。", // 查無此使用者, 但避免透落太多訊息給前端, 所以改成這樣
			})
			return
		} else {
			c.JSON(500, err)
			return
		}
	}

	user, _ := data.QueryUserByUniqueKey(countryCode, phoneNum)
	c.Header("access_token", token.AccessToken)
	c.JSON(200, gin.H{
		"systemCode":    200,
		"systemMessage": "OK",
		"data": gin.H{
			"userId":   user.UserId,
			"roles":    user.Roles,
			"nickname": user.Nickname.ValueOrZero(),
			"birthday": utils.Time2MilliTimestamp(user.Birthday),
			"gender":   utils.GenderInt2String(user.Gender.Int64),
		},
	})
	return
}
