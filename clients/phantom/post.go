package phantom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *client) DoPost(request interface{}) ([]byte, error) {
	logger := c.logger
	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", c.baseURL, apiPath)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Phantombuster-Key-1", "Hg2Dk5IHZbRPCHzzIUbEsXaIYKb1cxhY")

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
		logger.Error().Msgf("doPost non 200 response.json")
		return body, fmt.Errorf("post returned with errorCode: %d", resp.StatusCode)
	}
	logger.Info().Msgf("doPost success!")

	return body, nil
}
