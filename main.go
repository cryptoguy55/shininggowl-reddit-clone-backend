package main

import (
	"fmt"

	"first/articles"
	"first/common"
	"first/communities"
	"first/users"
	"first/websocket"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/jinzhu/gorm"
)

func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
		fmt.Printf("No .env file found")
	}
}

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
	db.AutoMigrate(&articles.ArticleModel{})
	db.AutoMigrate(&articles.TagModel{})
	db.AutoMigrate(&articles.FavoriteModel{})
	db.AutoMigrate(&articles.ArticleUserModel{})
	db.AutoMigrate(&articles.CommentModel{})
	db.AutoMigrate(&users.CommunityModel{})

}

func main() {
	db := common.Init()
	Migrate(db)
	defer db.Close()
	r := gin.Default()
	r.Use(cors.Default())
	v1 := r.Group("/api")
	// v1.Use(static.Serve("/public", static.LocalFile("./public", true)))
	v1.Static("/public", "./public")
	go websocket.Manager.Start()
	v1.GET("/ws/:id", func(c *gin.Context) {
		id := c.Param("id")
		websocket.WsPage(c.Writer, c.Request, id)
	})

	users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(false))
	articles.ArticlesAnonymousRegister(v1.Group("/articles"))
	articles.TagsAnonymousRegister(v1.Group("/tags"))
	communities.CommunitiesAnonymousRegister(v1.Group("/communities"))

	v1.Use(users.AuthMiddleware(true))
	users.UserRegister(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))
	communities.CommunitiesRegister(v1.Group("/communities"))

	articles.ArticlesRegister(v1.Group("/articles"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// test 1 to 1
	tx1 := db.Begin()
	userA := users.UserModel{
		Username: "AAAAAAAAAAAAAAAA",
		Email:    "aaaa@g.cn",
		Bio:      "hehddeda",
		//	Image:    nil,
	}
	tx1.Save(&userA)
	tx1.Commit()
	fmt.Println(userA)

	//db.Save(&ArticleUserModel{
	//    UserModelID:userA.ID,
	//})
	//var userAA ArticleUserModel
	//db.Where(&ArticleUserModel{
	//    UserModelID:userA.ID,
	//}).First(&userAA)
	//fmt.Println(userAA)
	fmt.Println("Hello go")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
