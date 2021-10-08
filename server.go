package main

import (
	"bookAPI/apiservice/api"
	"bookAPI/apiservice/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	firebaseAuth := utils.SetupFirebase()

	router.Use(func(c *gin.Context) {
		c.Set("firebaseAuth", firebaseAuth)
	})

	logger := utils.SetupLog()

	router.Use(func(c *gin.Context) {
		c.Set("logger", logger)
	})
	authRouter := router.Group("/AudioBookApp/api/auth")
	authRouter.POST("/register", api.Register)
	authRouter.POST("/loginWithPassword", api.LoginWithPassword)
	authRouter.POST("/loginWithidToken", api.LoginWithIdToken)
	authRouter.POST("/loginWithIdTokenButNoExpire", api.LoginWithIdTokenButNoExpire)
	authRouter.POST("/resetPassword", api.ResetPassword)
	authRouter.POST("/checkRegisterAvailable", api.CheckRegisterAvailable)
	authRouter.GET("/user/:userId", api.QueryUser)
	authRouter.POST("/deleteUser", api.DeleteUser)

	userRouter := router.Group("/AudioBookApp/api/user")
	userRouter.Use(utils.JWTAuth())
	userRouter.POST("/updateUserInfo", api.UpdateUserInfo)

	// router.Use(utils.JWTAuth())
	generalRouter := router.Group("/AudioBookApp/api")
	generalRouter.GET("/home", api.Home)
	generalRouter.GET("/home/:moreType", api.BookList)
	generalRouter.GET("/player", api.QueryBookExtContent)
	generalRouter.POST("/player/updateProgressBar", api.UpdateUserReadingBook)
	generalRouter.GET("/search/books", api.SearchBooks)
	generalRouter.GET("/search/recommendationStrings", api.SearchRecommendationString)
	router.Run()

}
