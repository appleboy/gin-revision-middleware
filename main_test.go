package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

func TestMiddlewareGeneralCase(t *testing.T) {

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
