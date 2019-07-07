package engines

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/phassans/frolleague/clients/linkedin"
	"github.com/phassans/frolleague/common"
)

type (
	linkedInEngine struct {
		client linkedin.Client
		sql    *sql.DB
		logger zerolog.Logger
	}

	LinkedInEngine interface {
		// LinkedInMethods
		GetAccessToken(authCode linkedin.AuthCode) (linkedin.AccessTokenResponse, error)
		LogIn(authCode linkedin.AuthCode, token linkedin.AccessToken) (linkedin.MeResponse, error)
		GetMe(userID linkedin.UserID) (linkedin.MeResponse, bool, error)

		// DBMethods
		SaveUser(id LinkedInUserID, firstName FirstName, lastName LastName, linkedInImage LinkedInImage) error
		UpdateUserWithLinkedInURL(id LinkedInUserID, url LinkedInURL) error
	}
)

func NewLinkedInEngine(client linkedin.Client, psql *sql.DB, logger zerolog.Logger) LinkedInEngine {
	return &linkedInEngine{client, psql, logger}
}

func (l *linkedInEngine) GetAccessToken(authCode linkedin.AuthCode) (linkedin.AccessTokenResponse, error) {
	return l.client.GetAccessToken(authCode)
}

func (l *linkedInEngine) LogIn(authCode linkedin.AuthCode, token linkedin.AccessToken) (linkedin.MeResponse, error) {
	// get access token
	if token == "" {
		resp, err := l.client.GetAccessToken(authCode)
		if err != nil {
			return linkedin.MeResponse{}, err
		}
		token = resp.AccessToken
	}

	fmt.Printf("accessToken: %s", token)

	// get me
	meResp, err := l.client.GetMe(token)
	if err != nil {
		return linkedin.MeResponse{}, err
	}

	// save user token
	if err := l.SaveToken(LinkedInUserID(meResp.ID), AccessToken(token)); err != nil {
		return linkedin.MeResponse{}, err
	}

	// save user
	if err := l.SaveUser(LinkedInUserID(meResp.ID), FirstName(meResp.FirstName.Localized.EnUS), LastName(meResp.LastName.Localized.EnUS), LinkedInImage(meResp.ProfilePicture.DisplayImage)); err != nil {
		switch err.(type) {
		case common.DuplicateLinkedInUser:
			fmt.Printf("duplicate user!")
			return meResp, nil
		default:
			return meResp, err
		}
	}

	return meResp, nil
}

func (l *linkedInEngine) GetMe(userID linkedin.UserID) (linkedin.MeResponse, bool, error) {
	token, err := l.GetTokenByUserID(LinkedInUserID(userID))
	if err != nil {
		return linkedin.MeResponse{}, false, nil
	}

	if token != "" {
		meResp, err := l.client.GetMe(linkedin.AccessToken(token))
		if err != nil {
			return linkedin.MeResponse{}, false, err
		}
		if meResp.ID == string(userID) {
			return meResp, true, nil
		}
	}
	return linkedin.MeResponse{}, false, nil
}

func (l *linkedInEngine) SaveUser(id LinkedInUserID, firstName FirstName, lastName LastName, linkedInImage LinkedInImage) error {
	user, err := l.GetUserByID(id)
	if err != nil {
		if _, ok := err.(common.ErrorUserNotExist); !ok {
			return err
		}
	}
	if user.UserID != "" {
		return common.DuplicateLinkedInUser{LinkedInUserID: string(id), Message: fmt.Sprintf("user with userID: %v already exists", id)}
	}

	return l.doSaveUser(id, firstName, lastName, linkedInImage)
}

func (l *linkedInEngine) GetUserByID(id LinkedInUserID) (LinkedInUser, error) {
	var user LinkedInUser
	rows := l.sql.QueryRow("SELECT user_id FROM linkedin_user WHERE user_id = $1", id)

	switch err := rows.Scan(&user.UserID); err {
	case sql.ErrNoRows:
		return LinkedInUser{}, common.ErrorUserNotExist{Message: fmt.Sprintf("user doesnt exist")}
	case nil:
		return user, nil
	default:
		return LinkedInUser{}, common.DatabaseError{DBError: err.Error()}
	}
}

func (l *linkedInEngine) doSaveUser(id LinkedInUserID, firstName FirstName, lastName LastName, linkedInImage LinkedInImage) error {
	_, err := l.sql.Exec("INSERT INTO linkedin_user(user_id, first_name, last_name, picture, insert_time) "+
		"VALUES($1,$2,$3,$4,$5)", id, firstName, lastName, linkedInImage, time.Now())
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	l.logger.Info().Msgf("successfully saved a linkedIn user with ID: %d", id)
	return nil
}

func (l *linkedInEngine) UpdateUserWithLinkedInURL(id LinkedInUserID, url LinkedInURL) error {
	updateWithURL := `UPDATE linkedin_user SET url = $1 WHERE user_id=$2;`

	_, err := l.sql.Exec(updateWithURL, url, id)
	if err != nil {
		return err
	}
	return nil
}

func (l *linkedInEngine) SaveToken(userID LinkedInUserID, accessToken AccessToken) error {
	dbToken, err := l.GetTokenByUserID(userID)
	if err != nil {
		switch err.(type) {
		case common.ErrorNotExist:
			if err := l.doSaveToken(userID, accessToken); err != nil {
				return err
			}
		default:
			return err
		}
		return nil
	}

	if dbToken != "" {
		if err := l.UpdateUserWithToken(userID, accessToken); err != nil {
			return err
		}
	}

	return nil
}

func (l *linkedInEngine) doSaveToken(userID LinkedInUserID, token AccessToken) error {
	_, err := l.sql.Exec("INSERT INTO linkedin_user_token(user_id, token, insert_time) "+
		"VALUES($1,$2,$3)", userID, token, time.Now())
	if err != nil {
		return common.DatabaseError{DBError: err.Error()}
	}

	l.logger.Info().Msgf("successfully saved token for user with ID: %d", userID)
	return nil
}

func (l *linkedInEngine) GetTokenByUserID(userID LinkedInUserID) (AccessToken, error) {
	var token AccessToken
	rows := l.sql.QueryRow("SELECT token FROM linkedin_user_token WHERE user_id = $1", userID)

	switch err := rows.Scan(&token); err {
	case sql.ErrNoRows:
		return AccessToken(""), common.ErrorNotExist{Message: fmt.Sprintf("user token doesnt exist")}
	case nil:
		return token, nil
	default:
		return AccessToken(""), common.DatabaseError{DBError: err.Error()}
	}
}

func (l *linkedInEngine) UpdateUserWithToken(userID LinkedInUserID, token AccessToken) error {
	updateWithToken := `UPDATE linkedin_user_token SET token = $1 WHERE user_id=$2;`

	_, err := l.sql.Exec(updateWithToken, token, userID)
	if err != nil {
		return err
	}
	return nil
}
