package rocket

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *client) DoGet(requestParams map[string]string, requestType string, params AdminCredentials) ([]byte, error) {
	logger := c.logger

	url := fmt.Sprintf("%s/%s/%s", c.baseURL, apiPath, requestType)
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("X-Auth-Token", params.AuthToken)
	req.Header.Set("X-User-Id", params.UserId)

	q := req.URL.Query()
	for key, value := range requestParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	logger = logger.With().Str("url", url).Str("status", resp.Status).Logger()

	if resp.StatusCode != 200 {
		logger = logger.With().Str("body", string(body)).Logger()
		logger.Error().Msgf("doGet non 200 response.json")
		return body, fmt.Errorf("get returned with errorCode: %d", resp.StatusCode)
	}
	logger.Info().Msgf("doGet success!")

	return body, nil
}
