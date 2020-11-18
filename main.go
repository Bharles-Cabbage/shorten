package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	user     = "postgres"
	password = "docker"
	dbname   = "golang_practice"
)

type urlShort struct {
	url  string
	slug string
}

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("static/html/index.html")

	psqlconn := "host=" + host + " port=5432 user=" + user + " password=" + os.Getenv("POSTGRESQL_PASSWORD") + " dbname=" + dbname + " sslmode=disable"

	db, err := sql.Open("postgres", psqlconn)
	checkError(err)
	defer db.Close()

	fmt.Printf("SUCCESFUL DB CONECT")
	err = db.Ping()
	checkError(err)

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "HTML Files load",
		})
	})

	router.POST("/shorten", func(c *gin.Context) {
		var shorturl urlShort

		url := c.PostForm("url")

		generatedSlug := randString()

		query := "SELECT * FROM urlshortner WHERE slug='AoI023iO';"
		err := db.QueryRow(query).Scan(&shorturl.url, &shorturl.slug)
		checkError(err)

		c.String(200, url+" "+generatedSlug+" "+shorturl.url+" "+shorturl.slug)
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
