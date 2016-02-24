package revision

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strings"
)

func GetRevision(fileName string) (string, error) {
	// Revision file contents will be only loaded once per process
	data, err := ioutil.ReadFile(fileName)

	// If we cant read file, just skip to the next request handler
	// This is pretty much a NOOP middlware :)
	if err != nil {
		log.Fatalf("Unable to read config file '%s'", fileName)

		return "", err
	}

	// Clean up the value since it could contain line breaks
	return strings.TrimSpace(string(data)), err
}

func RevisionMiddleware() gin.HandlerFunc {
	revision, err := GetRevision("REVISION")

	if err != nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Set out header value for each response
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Revision", revision)
		c.Next()
	}
}
