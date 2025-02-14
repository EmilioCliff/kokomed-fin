package pkg

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

func GenerateAccessToken(consumerKey string, consumerSecret string) (string, error) {
	authString := consumerKey + ":" + consumerSecret
	encodedAuthString := base64.StdEncoding.EncodeToString([]byte(authString))

	url := "https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Basic "+encodedAuthString)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", Errorf(INTERNAL_ERROR, "unexpected response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "error reading body: %v", err)
	}

	var tokenResponse map[string]interface{}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", Errorf(INTERNAL_ERROR, "eror unmarshaling body: %v", err)
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", Errorf(INTERNAL_ERROR, "Access token not found in response")
	}

	return accessToken, nil
}