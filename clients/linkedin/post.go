package linkedin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (c *client) DoPost(values url.Values, apiPath string) ([]byte, error) {
	logger := c.logger

	urlValues := fmt.Sprintf("%s/%s", c.baseURL, apiPath)
	req, err := http.NewRequest("POST", urlValues, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	logger = logger.With().Str("url", urlValues).Str("status", resp.Status).Logger()

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info().Msgf("doPost success!")
		return body, nil
	default:
		var e APIPostError
		err = json.Unmarshal(body, &e)
		if err != nil {
			// unmarshal error just return error
			return nil, fmt.Errorf("%s", string(body))
		}
		e.Code = resp.StatusCode

		logger = logger.With().Str("body", string(body)).Logger()
		logger.Error().Msgf("linkedIn post error with code %d for api %s", resp.StatusCode, apiPath)

		return nil, e
	}
}
