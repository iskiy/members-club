package main

import (
	"github.com/gin-gonic/gin"
	"members-club/pkg/handlers"
)

func main(){
	r := gin.Default()
	r.GET("/members/", handlers.GetMembers)
	r.StaticFile("/", "html/home.html")
	r.POST("/member/", handlers.AddMember)
	r.Run()
}
