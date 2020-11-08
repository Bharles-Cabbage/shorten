package main

import (
	"database/sql"
	"math/rand"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = "5432"
	user   = "postgres"
	dbname = "golang_practice"
)

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("static/html/index.html")

	psqlconn := "host=" + host + " port=" + port + " user=" + user + " password=" + os.Getenv("POSTGRESQL_PASSWORD") + " dbname=" + dbname + " sslmode=disable"

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "HTML Files load",
		})
	})

	router.POST("/shorten", func(c *gin.Context) {
		url := c.PostForm("url")

		shortURL := randString()

		c.String(200, url+shortURL)
	})

	// To be handled in future ... maybe ... not sure
	router.POST("/api", func(c *gin.Context) {
		c.String(200, "JSON post API")
	})

	router.Run(":8080")
}

func randString() string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randStr := make([]rune, 8)
	for i := range randStr {
		randStr[i] = chars[rand.Intn(len(chars))]
	}

	return string(randStr)
}
