package weather

import "fmt"

// ResponseError is the parsed error result from the weather.gov API.
type ResponseError struct {
	Type          string `json:"type,omitempty"`
	Title         string `json:"title,omitempty"`
	Status        int    `json:"status,omitempty"`
	Detail        string `json:"detail,omitempty"`
	Instance      string `json:"instance,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
}

// Error satisfies the error interface for *ResponseError.
func (err *ResponseError) Error() string {
	return fmt.Sprintf("status: %d, %s: %s", err.Status, err.Title, err.Detail)
}
