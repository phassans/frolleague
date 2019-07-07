package linkedin

import (
	"encoding/json"
	"fmt"
)

type (
	AuthCode    string
	AccessToken string
	UserID      string

	AccessTokenRequest struct {
		Grant_type    string `json:"grant_type"`
		Code          string `json:"code"`
		Redirect_uri  string `json:"redirect_uri"`
		Client_id     string `json:"client_id"`
		Client_secret string `json:"client_secret"`
	}

	AccessTokenResponse struct {
		AccessToken AccessToken `json:"access_token"`
		Expires     int         `json:"expires_in"`
	}

	MeResponse struct {
		LocalizedLastName string `json:"localizedLastName"`
		LastName          struct {
			Localized struct {
				EnUS string `json:"en_US"`
			} `json:"localized"`
			PreferredLocale struct {
				Country  string `json:"country"`
				Language string `json:"language"`
			} `json:"preferredLocale"`
		} `json:"lastName"`
		FirstName struct {
			Localized struct {
				EnUS string `json:"en_US"`
			} `json:"localized"`
			PreferredLocale struct {
				Country  string `json:"country"`
				Language string `json:"language"`
			} `json:"preferredLocale"`
		} `json:"firstName"`
		ProfilePicture struct {
			DisplayImage string `json:"displayImage"`
		} `json:"profilePicture"`
		ID                 string `json:"id"`
		LocalizedFirstName string `json:"localizedFirstName"`
	}

	APIPostError struct {
		Code             int    `json:code,omitempty`
		ErrorType        string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	APIGetError struct {
		ServiceErrorCode int    `json:"serviceErrorCode"`
		Message          string `json:"message"`
		Status           int    `json:"status"`
	}
)

func (v APIPostError) Error() string {
	b, _ := json.Marshal(fmt.Sprintf("%s:%s", v.ErrorType, v.ErrorDescription))
	return fmt.Sprintf("%s", string(b))
}

func (v APIGetError) Error() string {
	b, _ := json.Marshal(fmt.Sprintf("%s", v.Message))
	return fmt.Sprintf("%s", string(b))
}
