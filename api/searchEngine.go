package api

import (
	"bookAPI/apiservice/data"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchBooks(c *gin.Context) {
	searchString := c.Query("searchString")
	if searchString == "" {
		c.JSON(200, gin.H{
			"systemCode":    400,
			"systemMessage": "請輸入要查詢的文字，謝謝。",
		})
		return
	}
	userId, err := strconv.Atoi(c.Query("userId"))
	fmt.Println(userId)
	fmt.Println(searchString)
	historySearchString, err := data.QueryHistorySearchString(searchString)
	if err != nil && err.Error() == "sql: no rows in result set" {
		data.InsertHistorySearchString(searchString, 1)
	} else if err != nil {
		c.JSON(200, err.Error())
		return
	} else {
		historySearchString.Count = historySearchString.Count + 1
		data.UpdateHistorySearchString(historySearchString)
	}
	searchResultList, err := data.QueryBookByInputString(userId, searchString)
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
		"data":          searchResultList,
	})
}

func SearchRecommendationString(c *gin.Context) {
	searchString := c.Query("searchString")
	resultStringList, err := data.QueryRecommendationString(searchString, 1)
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
		"data":          resultStringList,
	})
}
