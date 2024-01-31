package localapiclient

import (
	"net/http"
	"testing"
	"time"
)

type BadStatusHandler struct{}

func (b *BadStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func TestBadStatus(t *testing.T) {
	badStatusHandler := &BadStatusHandler{}

	_, err := CallLocalApi(badStatusHandler, "POST", "test")

	if err != ErrBadHttpStatus {
		t.Error("Expected fallback string, but got", err)
	}
}

type TimeoutHandler struct{}

var successfulReport = "successful bug report!"

func (b *TimeoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(6 * time.Second)
	w.Write([]byte(successfulReport))
}

func TestBugTimeout(t *testing.T) {
	timeoutHandler := &TimeoutHandler{}

	_, err := CallLocalApi(timeoutHandler, "GET", "test")

	if err != TimeoutStatus {
		t.Error("Expected fallback string, but got", err)
	}
}

type SuccessfulHandler struct{}

func (b *SuccessfulHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(successfulReport))
}

func TestSuccess(t *testing.T) {
	successfulHandler := &SuccessfulHandler{}

	w, err := CallLocalApi(successfulHandler, "GET", "test")

	if err != nil {
		t.Error("Expected no error, but got", err)
	}

	report := string(w.Body())
	if report != successfulReport {
		t.Error("Expected successful report, but got", report)
	}
}
