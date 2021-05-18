package cache

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
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

	original, err := httputil.DumpResponse(rw.Result(), true)
	assert.NoError(t, err)
	assert.NotEqual(t, "", string(original))

	copy, err := recorder.Result()
	assert.NoError(t, err)
	assert.Equal(t, string(original), string(copy))
}
