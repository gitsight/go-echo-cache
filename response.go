package cache

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

type ResponseRecorder struct {
	http.ResponseWriter
	r *httptest.ResponseRecorder

	headerCopied bool
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	r := httptest.NewRecorder()
	return &ResponseRecorder{
		ResponseWriter: w,
		r:              r,
	}
}

func (w *ResponseRecorder) Write(b []byte) (int, error) {
	w.copyHeaders()
	i, err := w.ResponseWriter.Write(b)
	if err != nil {
		return i, err
	}

	return w.r.Write(b[:i])
}

func (r ResponseRecorder) copyHeaders() {
	if r.headerCopied {
		return
	}

	r.headerCopied = true
	copyHeaders(r.ResponseWriter.Header(), r.r.Header())
}

func (w *ResponseRecorder) WriteHeader(statusCode int) {
	w.copyHeaders()

	w.r.WriteHeader(statusCode)
	w.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Result() ([]byte, error) {
	r.copyHeaders()
	r.ResponseWriter = nil

	var buf bytes.Buffer
	err := r.r.Result().Write(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func copyHeaders(src, dst http.Header) {
	for k, v := range src {
		for _, v := range v {
			dst.Set(k, v)
		}
	}
}
