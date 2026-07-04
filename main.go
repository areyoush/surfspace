package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Link struct {
	Id 	string
	Url string
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var linkMap = map[string]*Link{
	"example": &Link{Id: "example", Url: "http://example.com"},
}

func main() {
	r := gin.Default()

	r.GET("/", IndexHandler)
	r.GET("/:id", RedirectHandler)
	r.POST("/submit", SubmitHandler)

	r.Run(":8080")
}

func RedirectHandler(c *gin.Context) {
	id := c.Param("id")
	link, found := linkMap[id]

	if !found {
		c.String(http.StatusNotFound, "Link not found")
		return
	}

	c.Redirect(http.StatusMovedPermanently, link.Url)
	
}

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := ""
	for range length {
		index := seededRand.Intn(len(charset))
		result += string(charset[index])
	}

	return result
} 

func IndexHandler(c *gin.Context) {
	html := `<h1>Submit a new website</h1>
	<form action="/submit" method="POST">
	<input type="text" name="url">
	<input type="submit" value="Submit">
	</form>
	<h2>Existing Links</h2>
	<ul>`

	for _, link := range linkMap {
		html += `<li><a href="/` + link.Id + `">` + link.Id + `</a></li>`
	}
	html += `</ul>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func SubmitHandler(c *gin.Context) {
	url := c.PostForm("url")

	if url == "" {
		c.String(http.StatusBadRequest, "URL is required")
		return
	}

	if !(len(url) >= 4 && (url[:4] == "http" || url[:5] == "https")) {
		url = "https://" + url
	}

	id := generateRandomString(8)

	linkMap[id] = &Link{Id: id, Url: url}

	c.Redirect(http.StatusSeeOther, "/")
}