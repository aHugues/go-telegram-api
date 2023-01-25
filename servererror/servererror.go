// package servererror handles errors from the telegram API
package servererror

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ahugues/go-telegram-api/baseclt"
)

type ServerError struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("Telegram API error (statuscode %d): %s", e.ErrorCode, e.Description)
}

// FromResponse builds a ServerError from a HTTP response
func FromResponse(r *http.Response) error {
	var parsedErr ServerError

	buf := bytes.NewBuffer([]byte{})
	maxRead := baseclt.HTTPMaxRead
	if r.ContentLength != -1 {
		if r.ContentLength > maxRead {
			return fmt.Errorf("response too big (%d)", r.ContentLength)
		}
		maxRead = r.ContentLength
	}

	if _, err := io.CopyN(buf, r.Body, maxRead); err != nil {
		return fmt.Errorf("error parsing body: %w", err)
	}
	if err := json.Unmarshal(buf.Bytes(), &parsedErr); err != nil {
		return fmt.Errorf("error parsing error %w", err)
	}
	return &parsedErr
}
