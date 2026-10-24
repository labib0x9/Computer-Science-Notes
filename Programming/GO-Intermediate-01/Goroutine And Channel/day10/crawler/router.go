package main

import "github.com/gin-gonic/gin"

func homePageHandler(ctx *gin.Context) {

}

func abcContestHandler(ctx *gin.Context) {

}

func arcContestHandler(ctx *gin.Context) {

}

func runWeb() {
	router := gin.Default()

	router.GET("/", homePageHandler)
	router.GET("/abc", abcContestHandler)
	router.GET("/arc", arcContestHandler)

	router.Run(":8080")
}