package main

import (
	"math/rand"
	"net/http"
	"strings"
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
	
	r.LoadHTMLGlob("templates/*")
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

	var sb strings.Builder
	sb.Grow(length)
	for range length {
		index := seededRand.Intn(len(charset))
		sb.WriteByte(charset[index])
	}
	return sb.String()
} 

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Links": linkMap,
	})
}

func SubmitHandler(c *gin.Context) {
	url := c.PostForm("url")

	if url == "" {
		c.String(http.StatusBadRequest, "URL is required")
		return
	}

	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	id := generateRandomString(8)

	linkMap[id] = &Link{Id: id, Url: url}

	c.Redirect(http.StatusSeeOther, "/")
}