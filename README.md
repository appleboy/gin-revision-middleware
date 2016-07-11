# gin-revision-middleware

Revision middleware for Gin framework written in Go (Golang).

[![Build Status](https://travis-ci.org/appleboy/gin-revision-middleware.svg?branch=master)](https://travis-ci.org/appleboy/gin-revision-middleware) [![Build Status](https://drone.io/github.com/appleboy/gin-revision-middleware/status.png)](https://drone.io/github.com/appleboy/gin-revision-middleware/latest) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/gin-revision-middleware)](https://goreportcard.com/report/github.com/appleboy/gin-revision-middleware) [![Coverage Status](https://coveralls.io/repos/github/appleboy/gin-revision-middleware/badge.svg?branch=master)](https://coveralls.io/github/appleboy/gin-revision-middleware?branch=master) [![codecov](https://codecov.io/gh/appleboy/gin-revision-middleware/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/gin-revision-middleware)

## How to use

Please see the [demo](example/main.go) code and creat `REVISION` file and write your application version into file.

```go
package main

import (
  "github.com/appleboy/gin-revision-middleware"
  "github.com/gin-gonic/gin"
  "net/http"
  "os"
)

func rootHandler(context *gin.Context) {
  context.JSON(http.StatusOK, gin.H{
    "Hello": "world",
  })
}

func main() {
  port := os.Getenv("PORT")
  router := gin.Default()
  if port == "" {
    port = "8000"
  }
  router.Use(revision.Middleware())
  router.GET("/", rootHandler)
  router.Run(":" + port)
}
```

Screenshot for header

![header screenshot](screenshots/revision_header.png)


