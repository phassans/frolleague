package linkedin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *client) DoGet(accessToken string, path string) ([]byte, error) {
	logger := c.logger

	url := fmt.Sprintf("%s/%s", c.apiBaseURL, path)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Connection", "Keep-Alive")

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

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info().Msgf("doPost success!")
		return body, nil
	default:
		var e APIGetError
		err = json.Unmarshal(body, &e)
		if err != nil {
			// unmarshal error just return error
			return nil, fmt.Errorf("%s", string(body))
		}
		e.Status = resp.StatusCode

		logger = logger.With().Str("body", string(body)).Logger()
		logger.Error().Msgf("linkedIn get error with code %d for api %s", resp.StatusCode, path)

		return nil, e
	}

	return body, nil
}
