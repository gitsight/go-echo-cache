package cache

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coocood/freecache"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	e := echo.New()

	var i int
	h := Cache(freecache.NewCache(42 * 1024 * 1024))(func(c echo.Context) error {
		i++
		return c.String(http.StatusOK, fmt.Sprintf("test_%d", i))
	})

	client := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := e.NewContext(r, w)
		assert.NoError(t, h(c))
	}))

	defer client.Close()

	resp, err := http.Get(client.URL)
	assert.NoError(t, err)
	assertRequest(t, resp, http.StatusOK, "test_1")

	resp, err = http.Get(client.URL)
	assert.NoError(t, err)
	assertRequest(t, resp, http.StatusOK, "test_1")

}

func assertRequest(t *testing.T, resp *http.Response, status int, content string) {
	if status != http.StatusOK {
		return
	}

	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, content, string(b))
	resp.Body.Close()
}
