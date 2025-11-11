package schoolgql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eldarbr/schoolsubscriber/internal/myerrs"
)

const (
	graphqlEndpoint = `https://platform.21-school.ru/services/graphql`
	clientTimeout   = time.Second * 15
)

type IBaseResponse interface {
	GetErrorText() string
}

func (req *Request) MakeRequest(ctx context.Context, token string, resultPlaceholder IBaseResponse) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshalling JSON: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, graphqlEndpoint, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	/* meta */
	httpReq.Header.Set("Content-Type", "application/json")
	// X-Edu-Org-Unit-Id <== GET https://edu.21-school.ru/services/rest/edu-context/context-info
	httpReq.Header.Set("X-EDU-SCHOOL-ID", "6bfe3c56-0211-4fe1-9e59-51616caac4dd")
	httpReq.Header.Set("X-EDU-PRODUCT-ID", "96098f4b-5708-4c42-a62c-6893419169b3")
	httpReq.Header.Set("X-EDU-ROUTE-INFO", "v1")
	httpReq.Header.Set("X-Edu-Org-Unit-Id", "6bfe3c56-0211-4fe1-9e59-51616caac4dd")
	httpReq.Header.Set("schoolid", "6bfe3c56-0211-4fe1-9e59-51616caac4dd")
	httpReq.Header.Set("userrole", "STUDENT")
	// TODO: make this header part non-constant and actually gather the context-info

	/* auth */
	/* either of these */
	// httpReq.AddCookie(&http.Cookie{Name: "tokenId", Value: token})
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{ //nolint:exhaustruct // leave defaults.
		Timeout: clientTimeout,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &myerrs.StatusCodeError{StatusCode: resp.StatusCode}
	}

	err = json.NewDecoder(resp.Body).Decode(resultPlaceholder)
	if err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	if errText := resultPlaceholder.GetErrorText(); errText != "" {
		return &myerrs.PlatformError{Text: errText}
	}

	return nil
}
