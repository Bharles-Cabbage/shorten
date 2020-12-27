package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
    "net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//const (
	//host     = "localhost"
	//user     = "postgres"
	//password = "docker"
	//dbname   = "golang_practice"
//)

type urlShort struct {
	url  string
	slug string
}

type URISlug struct{
    Slug string `uri:"url" binding:"required"`
}

func dbFunc(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        if _, err := db.Exec("CREATE TABLE IF NOT EXISTS urlshortner (URL char(2000) PRIMARY KEY, SLUG CHAR(9) NOT NULL"); err != nil {
            c.String(http.StatusInternalServerError,
                fmt.Sprintf("Error creating database table: %q", err))
            return
        }

        if _, err := db.Exec("INSERT INTO urlshortner VALUES ('test-url','test-slug')"); err != nil {
            c.String(http.StatusInternalServerError,
                fmt.Sprintf("Error inserting into table %q", err))
            return
        }

        rows, err := db.Query("SELECT * FROM urlshortner")
        if err != nil {
            c.String(http.StatusInternalServerError,
                fmt.Sprintf("Error reading urlshortner: %q", err))
            return
        }

        defer rows.Close()
        for rows.Next() {
            var url urlShort
            if err := rows.Scan(&url); err != nil {
                c.String(http.StatusInternalServerError,
                    fmt.Sprintf("Error scanning urlshortner: %q", err))
                return
            }
            c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", "urlshortner"))
        }
    }
}

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("static/html/index.tmpl")
    router.StaticFile("/static/main.css", "static/css/main.css")

	//psqlconn := "host=" + host + " port=5432 user=" + user + " password=" + os.Getenv("POSTGRESQL_PASSWORD") + " dbname=" + dbname + " sslmode=disable"

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	checkError(err)
	defer db.Close()

	fmt.Printf("SUCCESFUL DB CONECT")
	err = db.Ping()
	checkError(err)

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "HTML Files load",
		})
	})

    router.GET("/r/:url", func(c *gin.Context){
        var uri URISlug
        var shortURL urlShort
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

        query := "SELECT * FROM urlshortner WHERE slug='"+ uri.Slug +"'"
        err = db.QueryRow(query).Scan(&shortURL.url, &shortURL.slug)

        if err != nil {
            c.JSON(404, gin.H{"err": "Page not Found", "description": "No records associated with the URL found"})
        }
        c.Redirect(301, shortURL.url)

    })

	router.POST("/", func(c *gin.Context) {
		var shorturl urlShort
		var newshorturl urlShort
		var generatedSlug string
		var query string
		url := c.PostForm("url")

		generatedSlug = "Not generated this time"

		query = "SELECT * FROM urlshortner WHERE url='" + url + "'"
		err := db.QueryRow(query).Scan(&shorturl.url, &shorturl.slug)

		// IF Entered URL not in database generate a new url
		if err == sql.ErrNoRows {
			// Generate Random string until record for it is NOT found in the database
			for {
				generatedSlug = randString()

				query = "SELECT * FROM urlshortner WHERE slug='" + generatedSlug + "';"
				err = db.QueryRow(query).Scan(&newshorturl.url, &newshorturl.slug)

				if err != nil {
					break
				}
			}

			query = "INSERT INTO urlshortner VALUES($1, $2);"
			_, err := db.Exec(query, url, generatedSlug)
			checkError(err)

			query = "SELECT * FROM urlshortner WHERE slug='" + generatedSlug + "';"
			err = db.QueryRow(query).Scan(&newshorturl.url, &newshorturl.slug)
			checkError(err)
		} else if err != nil {
			checkError(err)
		}

        c.HTML(200, "index.tmpl", gin.H{
            "ShortURL": shorturl.slug,
        })
		//c.String(200, url+" | "+generatedSlug+" | "+shorturl.url+" | "+shorturl.slug)
        //c.Redirect(200, "/", gin.H{
            //"ShortURL": shorturl.slug,
        //})
	})

	// To be handled in future ... maybe ... not sure
	router.POST("/api", func(c *gin.Context) {
		c.String(200, "JSON post API")
	})

    router.GET("/db", dbFunc(db))

    router.Run(":" + os.Getenv("PORT"))
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
