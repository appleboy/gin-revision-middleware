package revision

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
	filename := "tempfile"

	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		log.Fatalf("WriteFile %s: %v", filename, err)
	}

	// clean up
	defer os.Remove(filename)

	result, _ := GetRevision(filename)

	// TEST
	assert.Equal(t, result, string(content))
}

func Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func TestRevisionMiddleware(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/v1")
	v1.Use(Middleware())
	{
		v1.GET("/hello", Hello)
	}

	// without middleware
	v2 := r.Group("/v2")
	{
		v2.GET("/hello", Hello)
	}

	// other revision file
	v3 := r.Group("/v3")
	v3.Use(Middleware("REVISION3"))
	{
		v3.GET("/hello", Hello)
	}

	// missin revision file
	v4 := r.Group("/v4")
	v4.Use(Middleware("REVISION4"))
	{
		v4.GET("/hello", Hello)
	}

	// RUN
	req, err := http.NewRequest("GET", "/v1/hello", nil)

	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// TEST
	assert.Equal(t, w.Body.String(), "Hello World")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	assert.Equal(t, w.HeaderMap.Get("X-Revision"), "1.0.0")

	// RUN
	req, err = http.NewRequest("GET", "/v2/hello", nil)

	if err != nil {
		log.Fatal(err)
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// TEST
	assert.Equal(t, w.Body.String(), "Hello World")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	assert.Empty(t, w.HeaderMap.Get("X-Revision"))

	// RUN
	req, err = http.NewRequest("GET", "/v3/hello", nil)

	if err != nil {
		log.Fatal(err)
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// TEST
	assert.Equal(t, w.Body.String(), "Hello World")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	assert.Equal(t, w.HeaderMap.Get("X-Revision"), "3.0.0")

	// RUN
	req, err = http.NewRequest("GET", "/v4/hello", nil)

	if err != nil {
		log.Fatal(err)
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// TEST
	assert.Equal(t, w.Body.String(), "Hello World")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")
	assert.Empty(t, w.HeaderMap.Get("X-Revision"))
}
