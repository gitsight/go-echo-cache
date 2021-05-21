package cache

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestResponseRecorderResult(t *testing.T) {
	rw := httptest.NewRecorder()
	recorder := NewResponseRecorder(rw)
	c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), recorder)

	h := func(c echo.Context) error {
		c.Response().Header().Set("X-Foo", "42")
		return c.HTML(http.StatusOK, "foo")
	}

	assert.NoError(t, h(c))
	entry := recorder.Result()
	original := rw.Result()
	assert.Equal(t, original.StatusCode, entry.StatusCode)
	assert.Equal(t, original.Header, entry.Header)

	b, _ := io.ReadAll(original.Body)
	assert.Equal(t, b, entry.Body)
}
