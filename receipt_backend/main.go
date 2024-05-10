// Copyright 2024, Yair Zadok, All rights reserved.

package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/IntuitDeveloper/OAuth2-Go/config"
	"github.com/IntuitDeveloper/OAuth2-Go/handlers"
	"fmt"
   	"github.com/Yair-Zadok/godeeby"
   	"database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./valid.db")
	if err != nil { fmt.Println(err) }
	
    err = godeeby.Setup_db(db)
	if err != nil { fmt.Println(err) }


	router := gin.Default()
	router.Use(cors.Default())
	
	handlers.CallDiscoveryAPI()

	// Register handler routes using the wrap function
	router.GET("/getCompanyInfo", gin.WrapF(handlers.GetCompanyInfo))
	router.GET("/refreshToken", gin.WrapF(handlers.RefreshToken))
	router.GET("/revokeToken", gin.WrapF(handlers.RevokeToken))
	router.GET("/connectToQuickbooks", gin.WrapF(handlers.ConnectToQuickbooks))
	router.GET("/signInWithIntuit", gin.WrapF(handlers.SignInWithIntuit))
	router.GET("/getAppNow", gin.WrapF(handlers.GetAppNow))
	router.GET("/oauth2redirect", gin.WrapF(handlers.CallBackFromOAuth))
	router.GET("/connected/post", gin.WrapF(post_receipt))
	

	port := config.OAuthConfig.Port
	log.Printf("running server on %s", port)
	router.Run(port) 
}
