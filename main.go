package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

func RevisionMiddleware() gin.HandlerFunc {
	// Revision file contents will be only loaded once per process
	data, err := ioutil.ReadFile("REVISION")

	// If we cant read file, just skip to the next request handler
	// This is pretty much a NOOP middlware :)
	if err != nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Clean up the value since it could contain line breaks
	revision := strings.TrimSpace(string(data))

	// Set out header value for each response
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Revision", revision)
		c.Next()
	}
}
