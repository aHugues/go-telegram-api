// package baseclt contains the base HTTP client to contact the Telegram API
package baseclt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseTelegramAPIURL = "https://api.telegram.org"

const HTTPMaxRead int64 = 100_000

const HTTPTimeout = 5 * time.Second

func ParseJSONBody(response *http.Response, dest interface{}) error {
	buf := bytes.NewBuffer([]byte{})
	maxRead := HTTPMaxRead
	if response.ContentLength != -1 {
		if response.ContentLength > maxRead {
			return fmt.Errorf("response too big (%d)", response.ContentLength)
		}
		maxRead = response.ContentLength
	}

	if _, err := io.CopyN(buf, response.Body, maxRead); err != nil {
		return fmt.Errorf("error parsing body: %w", err)
	}
	if err := json.Unmarshal(buf.Bytes(), dest); err != nil {
		return fmt.Errorf("error parsing error %w", err)
	}
	return nil
}
