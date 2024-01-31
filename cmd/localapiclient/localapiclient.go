package localapiclient

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"
)

// LocalAPIResponseWriter substitutes for http.ResponseWriter in order to write byte streams directly
// to a receiver function in the application.
type LocalApiResponseWriter struct {
	headers http.Header
	body    bytes.Buffer
	status  int
}

func newLocalApiResponseWriter() *LocalApiResponseWriter {
	return &LocalApiResponseWriter{headers: http.Header{}, status: http.StatusOK}
}

func (w *LocalApiResponseWriter) Header() http.Header {
	return w.headers
}

// Write writes the data to the response body, which will be sent to Java. If WriteHeader is not called
// explicitly, the first call to Write will trigger an implicit WriteHeader(http.StatusOK).
func (w *LocalApiResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(data)
}

func (w *LocalApiResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *LocalApiResponseWriter) Body() []byte {
	return w.body.Bytes()
}

func (w *LocalApiResponseWriter) StatusCode() int {
	return w.status
}

var ErrBadHttpStatus = errors.New("bad http status for localapi response")
var TimeoutStatus = errors.New("localapi request timed out")

func CallLocalApi(h http.Handler, method string, endpoint string) (*LocalApiResponseWriter, error) {
	done := make(chan *LocalApiResponseWriter, 1)
	var responseError error
	go func() {
		req, err := http.NewRequest(method, "/localapi/v0/"+endpoint, nil)
		if err != nil {
			log.Printf("error creating new request for %s: %v", endpoint, err)
			responseError = err
			close(done)
			return
		}
		w := newLocalApiResponseWriter()
		h.ServeHTTP(w, req)
		if w.StatusCode() > 300 {
			log.Printf("%s bad http status: %v", endpoint, w.StatusCode())
			responseError = ErrBadHttpStatus
		}
		done <- w
	}()

	select {
	case w := <-done:
		return w, responseError
	case <-time.After(2 * time.Second):
		responseError = TimeoutStatus
		return nil, responseError
	}
}
