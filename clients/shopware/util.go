package shopware

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func checkHttpResp(resp *resty.Response, err error, expectedStatus ...int) error {
	if resp == nil {
		return fmt.Errorf("no response. err: %w", err)
	}

	if err != nil {
		return fmt.Errorf("request failed: %w, body [%s]", err, resp.Body())
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized")
	}
	if resp.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("not found")
	}

	// Check if status are expected
	if len(expectedStatus) == 0 { // if no expected status code is set, default to 2XX
		if resp.StatusCode() < 300 && resp.StatusCode() >= 200 {
			return nil
		}
	} else { // check custom expected status code
		for _, expectedStatus := range expectedStatus {
			if resp.StatusCode() == expectedStatus {
				return nil
			}
		}
	}

	return fmt.Errorf("unexpected response code %d : %s", resp.StatusCode(), resp.Body())
}

func tokenHeaders(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json",
		"Content-Type":  "application/json",
	}
}
