package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}
type TweetText struct {
	Tweettext string `json:"tweet_text,omitempty"  form:"Tweet"`
	User      string `json:"user,omitempty"  form:"user"`
	Number    int    `json:"number,omitempty"  form:"number"`
	Fav       int    `json:"fav,omitempty"  form:"fav"`
}

type TweetText2 struct {
	Tweettext string `json:"tweet_text,omitempty"  db:"text"`
	User      string `json:"user,omitempty"  db:"User"`
	Number    int    `json:"number,omitempty"  db:"number"`
	Fav       int    `json:"fav"  db:"fav"`
}

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	e := echo.New()

	e.GET("/cities/:cityName", getCityInfoHandler)
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.GET("/tweet", getTweetHandler)
	e.GET("/tweet/page/:number", getTweetPageHandler)
	e.POST("/tweet", postTweetHandler)
	e.DELETE("/tweet/:number", deleteTweetHandler)
	e.POST("/tweet/:number", postTweetfavHandler)

	e.Start(":4000")
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	fmt.Println(cityName)

	city := City{}
	db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)
	if city.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, city)
}

func getTweetHandler(c echo.Context) error {
	/*tweettext := c.Param("tweettext")
	fmt.Println(tweettext)*/

	text := []TweetText2{}
	db.Select(&text, "SELECT * FROM tweets")

	/*if TweetText2.Tweettext == "" {
		return c.NoContent(http.StatusNotFound)
	}*/
	return c.JSON(http.StatusOK, text)
}
func getTweetPageHandler(c echo.Context) error {
	/*tweettext := c.Param("tweettext")
	fmt.Println(tweettext)*/
	i := 1
	number := c.Param("number")
	i, _ = strconv.Atoi(number)
	text := []TweetText2{}
	fmt.Println(i)
	db.Select(&text, "SELECT * FROM tweets limit 10 offset ?", i)

	/*if TweetText2.Tweettext == "" {
		return c.NoContent(http.StatusNotFound)
	}*/
	return c.JSON(http.StatusOK, text)
}

func postTweetHandler(c echo.Context) error {
	req := TweetText{}
	c.Bind(&req)

	// もう少し真面目にバリデーションするべき
	if req.Tweettext == "" {
		// エラーは真面目に返すべき
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	db.Exec("INSERT INTO tweets (User,text,fav) VALUES (?,?,0)", req.User, req.Tweettext)

	return c.NoContent(http.StatusCreated)
}

func deleteTweetHandler(c echo.Context) error {
	req := TweetText2{}
	c.Bind(&req)

	db.Exec("DELETE FROM tweets WHERE number = (?)", req.Number)

	return c.NoContent(http.StatusCreated)
}

func postTweetfavHandler(c echo.Context) error {
	req := TweetText2{}
	c.Bind(&req)

	db.Exec("UPDATE tweets SET fav = fav + 1 WHERE number = (?)", req.Number)
	return c.NoContent(http.StatusCreated)
}
