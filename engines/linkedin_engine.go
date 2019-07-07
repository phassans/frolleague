package engines

import (
	"fmt"

	"github.com/phassans/frolleague/clients/linkedin"
)

type (
	linkedInEngine struct {
		client   linkedin.Client
		tokenMap map[linkedin.UserID]linkedin.AccessToken
	}

	LinkedInEngine interface {
		GetAccessToken(authCode linkedin.AuthCode) (linkedin.AccessTokenResponse, error)
		LogIn(authCode linkedin.AuthCode, token linkedin.AccessToken) (linkedin.MeResponse, error)
		GetMe(userID linkedin.UserID) (linkedin.MeResponse, bool, error)
	}
)

func NewLinkedInEngine(client linkedin.Client) LinkedInEngine {
	m := make(map[linkedin.UserID]linkedin.AccessToken)
	return &linkedInEngine{client, m}
}

func (l *linkedInEngine) GetAccessToken(authCode linkedin.AuthCode) (linkedin.AccessTokenResponse, error) {
	return l.client.GetAccessToken(authCode)
}

func (l *linkedInEngine) LogIn(authCode linkedin.AuthCode, token linkedin.AccessToken) (linkedin.MeResponse, error) {
	if token == "" {
		resp, err := l.client.GetAccessToken(authCode)
		if err != nil {
			return linkedin.MeResponse{}, err
		}
		token = resp.AccessToken
	}
	fmt.Printf("accessToken: %s", token)

	meResp, err := l.client.GetMe(token)
	if err != nil {
		return linkedin.MeResponse{}, err
	}
	l.tokenMap[linkedin.UserID(meResp.ID)] = token

	return meResp, nil
}

func (l *linkedInEngine) GetMe(userID linkedin.UserID) (linkedin.MeResponse, bool, error) {
	if token, ok := l.tokenMap[userID]; ok {
		meResp, err := l.client.GetMe(token)
		if err != nil {
			return linkedin.MeResponse{}, false, err
		}
		if meResp.ID == string(userID) {
			return meResp, true, nil
		}
	}
	return linkedin.MeResponse{}, false, nil
}
