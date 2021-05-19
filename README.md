# go-echo-cache [![GoDoc](https://godoc.org/github.com/gitsight/go-echo-cache?status.svg)](https://pkg.go.dev/github.com/gitsight/go-echo-cache) [![Test](https://github.com/gitsight/go-echo-cache/workflows/Test/badge.svg)](https://github.com/gitsight/go-echo-cache/actions?query=workflow%3ATest+branch%3Amaster) 

*go-echo-cache*, is a server-side HTTP cache middleware designed to work with [Echo framework](https://echo.labstack.com/).

The in-memory cache is managed by [freecache](https://github.com/coocood/freecache), a cache library with zero GC overhead and high concurrent performance.

Installation
------------

The recommended way to install *go-echo-cache* is:

```go
go get -u github.com/gitsight/go-echo-cache     
```


Examples
--------

### Basic example

A basic example that mimics the standard `git clone` command

```go
package main

import (
	"net/http"
	"time"

	"github.com/coocood/freecache"
	cache "github.com/gitsight/go-echo-cache"
	"github.com/labstack/echo/v4"
)

func main() {
	c := freecache.NewCache(1024 * 1024) // Pre-allocated cache of 1Mb)

	e := echo.New()
	e.Use(cache.New(&cache.Config{}, c))
	e.GET("/", func(c echo.Context) error {
		c.String(http.StatusOK, time.Now().String())
		return nil
	})

	e.Start(":8080")
}

```

License
-------
Apache License Version 2.0, see [LICENSE](LICENSE)
