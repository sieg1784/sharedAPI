package api

import (
	"bookAPI/apiservice/data"
	"strconv"

	"github.com/gin-gonic/gin"
)

type fn func(int, int, int) ([]data.ResultBook, error)

var functionMap = map[string]fn{
	"cowork":              data.QueryCowork,
	"recommendationToday": data.QueryRecommendation,
	"bookList":            data.QueryBookList,
}

var displayTextMap = map[string]string{
	"cowork":              "合作說書",
	"recommendationToday": "今日推薦",
	"bookList":            "說書清單",
}

func Home(c *gin.Context) {

	userId, err := strconv.Atoi(c.Query("userId"))
	nowLearningList, err := data.QueryNowLearning(userId)
	if err != nil {
		c.JSON(200, gin.H{
			"systemCode":    500,
			"systemMessage": err.Error(),
		})
		return
	}

	recommendationList, err := data.QueryRecommendation(userId, 1, 6)
	if err != nil {
		c.JSON(200, gin.H{
			"systemCode":    500,
			"systemMessage": err.Error(),
		})
		return
	}

	coworkList, err := data.QueryCowork(userId, 1, 6)
	if err != nil {
		c.JSON(200, gin.H{
			"systemCode":    500,
			"systemMessage": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"systemCode":    200,
		"systemMessage": "ok",
		"data": gin.H{
			"learningRecently": gin.H{
				"displayText": "最近在學",
				"data":        nowLearningList,
			},
			"cowork": gin.H{
				"displayText": "合作說書",
				"data":        coworkList,
			},
			"recommendationToday": gin.H{
				"displayText": "今日推薦",
				"data":        recommendationList,
			},
		},
	})
	return
}

func BookList(c *gin.Context) {

	tempSize, err := data.QuerySystemConfig("pageSize")
	pageSize, _ := strconv.Atoi(tempSize)
	if err != nil {
		c.JSON(200, gin.H{
			"systemCode":    500,
			"systemMessage": err.Error(),
		})
		return
	}
	bookCategoryList, err := data.QueryBookCategory()
	moreType := c.Param("moreType")
	userId, _ := strconv.Atoi(c.Query("userId"))
	page, _ := strconv.Atoi(c.Query("page"))
	resultList, err := functionMap[moreType](userId, page, pageSize)

	if err != nil {
		c.JSON(200, gin.H{
			"systemCode":    500,
			"systemMessage": err.Error(),
		})
		return
	} else {
		c.JSON(200, gin.H{
			"systemCode":    200,
			"systemMessage": "ok",
			"data": gin.H{
				"bookCategory": bookCategoryList,
				moreType: gin.H{
					"displayText": displayTextMap[moreType],
					"data":        resultList,
				},
			},
		})
		return
	}
}
