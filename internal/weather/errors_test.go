package weather_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ewohltman/weather-station/internal/weather"
)

func TestResponseError_Error(t *testing.T) {
	t.Parallel()

	status := http.StatusInternalServerError
	title := "testTitle"
	detail := "testDetail"

	message := (&weather.ResponseError{
		Status: status,
		Title:  title,
		Detail: detail,
	}).Error()

	expected := fmt.Sprintf("status: %d, %s: %s", status, title, detail)

	if message != expected {
		t.Errorf("%q != %q", message, expected)
	}
}
