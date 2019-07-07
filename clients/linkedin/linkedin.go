package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

const (
	// linkedIn
	grantType    = "authorization_code"
	redirectURI  = "http://localhost:8000/projects/linkedInLogin/"
	clientID     = "86ex3hh85g80oi"
	clientSecret = "CffCzLefgz0f4X6S"

	// APIs
	apiPathAccessToken = "apiPathAccessToken"
	apiPathMe          = "apiPathMe"
)

type (
	client struct {
		baseURL    string
		apiBaseURL string
		logger     zerolog.Logger
	}

	Client interface {
		GetAccessToken(AuthCode) (AccessTokenResponse, error)
		GetMe(AccessToken) (MeResponse, error)
	}
)

var apiPath = map[string]string{apiPathAccessToken: "accessToken", apiPathMe: "me"}

// NewLinkedInClient returns a new linkedIn client
func NewLinkedInClient(baseURL string, apiBaseURL string, logger zerolog.Logger) Client {
	return &client{baseURL, apiBaseURL, logger}
}

func (c *client) GetAccessToken(code AuthCode) (AccessTokenResponse, error) {
	logger := c.logger

	request := AccessTokenRequest{
		Grant_type:    grantType,
		Code:          string(code),
		Redirect_uri:  redirectURI,
		Client_id:     clientID,
		Client_secret: clientSecret,
	}
	response, err := c.DoPost(structToMap(&request), apiPath[apiPathAccessToken])
	if err != nil {
		return AccessTokenResponse{}, err
	}

	// read response.json
	var resp AccessTokenResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on AccessTokenResponse")
		return AccessTokenResponse{}, err
	}

	return resp, err
}

func (c *client) GetMe(token AccessToken) (MeResponse, error) {
	logger := c.logger
	response, err := c.DoGet(string(token), apiPath[apiPathMe])
	if err != nil {
		return MeResponse{}, err
	}

	// read response.json
	var resp MeResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on AccessTokenResponse")
		return MeResponse{}, err
	}

	return resp, err
}

func structToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		if strings.ToLower(typ.Field(i).Name) == "code" {
			fmt.Println(v)
		}
		values.Set(strings.ToLower(typ.Field(i).Name), v)
	}
	return
}
