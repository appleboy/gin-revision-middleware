package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetRevision(t *testing.T) {
	content := []byte("temporary file's content")
	filename := "example"

	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		log.Fatalf("WriteFile %s: %v", filename, err)
	}

	// clean up
	defer os.Remove(filename)

	result, _ := GetRevision(filename)

	// TEST
	assert.Equal(t, result, string(content))
}

func TestRevisionMiddleware(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RevisionMiddleware())
	r.Handle("GET", "/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	// RUN
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// TEST
	assert.Equal(t, w.Body.String(), "Hello World")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	assert.Equal(t, w.HeaderMap.Get("X-Revision"), "1.0.0")
}
