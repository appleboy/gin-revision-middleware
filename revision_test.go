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

func Result(t *testing.T, router *gin.Engine, path string, revision string) {
	// RUN
	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// TEST
	assert.Equal(t, w.Body.String(), "Hello World")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain; charset=utf-8")

	if revision == "" {
		assert.Empty(t, w.HeaderMap.Get("X-Revision"))
	} else {
		assert.Equal(t, w.HeaderMap.Get("X-Revision"), revision)
	}
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

	// missing revision file
	v4 := r.Group("/v4")
	v4.Use(Middleware("REVISION4"))
	{
		v4.GET("/hello", Hello)
	}

	Result(t, r, "/v1/hello", "1.0.0")
	Result(t, r, "/v2/hello", "")
	Result(t, r, "/v3/hello", "3.0.0")
	Result(t, r, "/v4/hello", "")
}
