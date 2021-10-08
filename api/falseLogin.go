package api

import (
	"bookAPI/apiservice/data"
	"bookAPI/apiservice/service"
	"bookAPI/apiservice/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoginWithIdTokenButNoExpire(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, "Id token is missing")
		return
	}

	decodedIdToken, err := utils.DecodeFirebaseIdToken(idTokenFromFirebase)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	str := decodedIdToken.Firebase.Identities["phone"] // str only used in next line, forgive this namming
	log.Println(str)
	phoneFromFirebase := utils.TrimBrackets(fmt.Sprintf("%v", str))
	countryCodeFromFirebase := utils.ExtractCountryCode(phoneFromFirebase)
	phoneNumFromFirebase := utils.ExtractPhoneNum(phoneFromFirebase)
	log.Println(countryCodeFromFirebase)
	log.Println(phoneNumFromFirebase)
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
