package cache

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/coocood/freecache"
	"github.com/labstack/echo/v4"
)

func Cache(cache *freecache.Cache) echo.MiddlewareFunc {
	m := &CacheMiddleware{cache: cache}
	return m.Handler
}

type CacheMiddleware struct {
	cache *freecache.Cache
}

func (m *CacheMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := m.readCache(c)
		if err == nil {
			return nil
		}

		if err != freecache.ErrNotFound {
			c.Logger().Errorf("error reading cache: %s", err)
		}

		recorder := NewResponseRecorder(c.Response().Writer)
		c.Response().Writer = recorder

		err = next(c)
		if err := m.cacheResult(recorder); err != nil {
			c.Logger().Error(err)
		}

		return err
	}
}

func (m *CacheMiddleware) readCache(c echo.Context) error {
	value, err := m.cache.Get([]byte("foo"))
	if err != nil {
		return err
	}

	b := bufio.NewReader(bytes.NewBuffer(value))
	r, err := http.ReadResponse(b, c.Request())
	if err != nil {
		return err
	}

	defer r.Body.Close()
	copyHeaders(r.Header, c.Response().Header())
	c.Response().WriteHeader(r.StatusCode)

	_, err = io.Copy(c.Response(), r.Body)
	return err
}

func (m *CacheMiddleware) cacheResult(r *ResponseRecorder) error {
	b, err := r.Result()
	if err != nil {
		return fmt.Errorf("unable to read recorded response: %s", err)
	}

	key := []byte("foo")
	return m.cache.Set(key, b, 1)
}
