package api

import (
	"bookAPI/apiservice/data"
	"bookAPI/apiservice/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryBookExtContent(c *gin.Context) {
	type Transcript struct {
		Id                string `json:"id"`
		Start             string `json:"start"`
		End               string `json:"end"`
		StartTimeInSecond int    `json:"startTimeInSecond"`
		Sentences         string `json:"sentences"`
	}
	bookId := c.Query("bookId")

	bookAudioUrl, err := data.QueryBookAudioUrl(bookId)
	if err != nil {
		c.JSON(200, gin.H{
			"systemCode":    500,
			"systemMessage": err.Error(),
		})
		return
	}

	bookExtContentList, err := data.QueryBookExtContent(bookId)
	if err != nil {
		c.JSON(200, gin.H{
			"systemCode":    500,
			"systemMessage": err.Error(),
		})
		return
	}

	bookExtContentMap := make(map[string]string)
	for _, bookExtContent := range bookExtContentList {
		bookExtContentMap[bookExtContent.ExtType.String] = bookExtContent.ExtContent.String
	}

	bytes := []byte(bookExtContentMap["transcript"])
	var transcripts []Transcript
	err = json.Unmarshal(bytes, &transcripts)

	c.JSON(200, gin.H{
		"systemCode":    200,
		"systemMessage": "ok",
		"data": gin.H{
			"bookAudioUrl": bookAudioUrl,
			"bookExtContent": gin.H{
				"mindMap": gin.H{
					"imageUrl": bookExtContentMap["mindMap"],
				},
				"quotes": gin.H{
					"text": bookExtContentMap["quotes"],
				},
				"transcript": transcripts,
			},
		},
	})
}

func UpdateUserReadingBook(c *gin.Context) {
	userId, err := strconv.Atoi(c.PostForm("userId"))

	bookId := c.PostForm("bookId")
	tmpReadingDate, err := strconv.Atoi(c.PostForm("readingDate"))

	readingDate := utils.MilliTimestamp2Time(int64(tmpReadingDate))
	readingTime, err := strconv.Atoi(c.PostForm("readingTime"))

	userReadingBook, err := data.QueryUserReadingBook(userId, bookId)
	userReadingBook = data.UserReadingBook{}
	userReadingBook.UserId = userId
	userReadingBook.BookId = bookId
	userReadingBook.ReadingDate = readingDate
	userReadingBook.ReadingTime = readingTime
	fmt.Println(userReadingBook)
	if err != nil && err.Error() == "sql: no rows in result set" {
		data.InsertUserReadingBook(userReadingBook)
	} else if err != nil {
		c.JSON(200, err.Error())
		return
	} else {
		data.UpdateUserReadingBook(userReadingBook)
	}
	c.JSON(200, gin.H{
		"systemCode":    200,
		"systemMessage": "ok",
	})
}
