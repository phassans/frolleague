package rocket

import (
	"encoding/json"
	"fmt"
)

var adminCredentials = make(map[string]AdminCredentials)

const credentials = "credentials"

func (c *client) InitClient(username string, password string) error {
	resp, err := c.Login(UserLoginRequest{username, password})
	if err != nil {
		return err
	}

	adminCredentials[credentials] = AdminCredentials{resp.Data.AuthToken, resp.Data.UserID}
	return nil
}

func (c *client) GetAdminCredentials() AdminCredentials {
	if creds, ok := adminCredentials[credentials]; ok {
		return creds
	}

	panic("rocket client is not initialized")
	return AdminCredentials{}
}

func (c *client) Login(request UserLoginRequest) (UserLoginResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request, login, AdminCredentials{})
	if err != nil {
		var errResp ErrorResponse
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return UserLoginResponse{}, err
		}

		logger = logger.With().
			Str("code", errResp.Error).
			Str("error", errResp.Message).
			Str("status", errResp.Status).
			Logger()
		logger.Error().Msgf("login returned with error")
		return UserLoginResponse{}, fmt.Errorf("login returned with error: %s", errResp.Error)
	}

	// read response.json
	var resp UserLoginResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on UserLoginResponse")
		return UserLoginResponse{}, err
	}

	return resp, nil
}

func (c *client) CreateUser(request CreateUserRequest) (CreateUserResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request, createUser, c.GetAdminCredentials())
	if err != nil {
		var errResp Err
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return CreateUserResponse{}, err
		}

		logger = logger.With().Str("error", errResp.Error).Str("errorType", errResp.ErrorType).Bool("status", errResp.Success).Logger()
		logger.Error().Msgf("create user returned with error")

		// duplicate username error
		if errResp.Error == fmt.Sprintf("%s is already in use :( [error-field-unavailable]", request.Username) {
			return CreateUserResponse{}, ErrorDuplicateUserName{Username: request.Username}
		}

		return CreateUserResponse{}, ErrorUserCreation{ErrorMsg: errResp.Error}
	}

	// read response.json
	var resp CreateUserResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on CreateUserResponse")
		return CreateUserResponse{}, err
	}

	return resp, nil
}

func (c *client) DeleteUser(request DeleteUserRequest) (DeleteUserResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request, deleteUser, c.GetAdminCredentials())
	if err != nil {
		var errResp Err
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return DeleteUserResponse{}, err
		}

		logger = logger.With().Str("errorType", errResp.ErrorType).Str("error", errResp.Error).Bool("status", errResp.Success).Logger()
		logger.Error().Msgf("delete user returned with error")

		switch errResp.ErrorType {
		case ErrorUserParamNotProvidedType:
			return DeleteUserResponse{}, ErrorRequiredParam{ErrorMsg: fmt.Sprintf("missing required field!")}
		case ErrorInvalidUserType:
			return DeleteUserResponse{}, ErrorInvalidUser{ErrorMsg: fmt.Sprintf("user doesnt exist: %s", request.UserId)}
		}

		return DeleteUserResponse{}, fmt.Errorf("DeleteUser returned with error: %s", errResp.Error)
	}

	// read response.json
	var resp DeleteUserResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on DeleteUserResponse")
		return DeleteUserResponse{}, err
	}

	return resp, nil
}

func (c *client) InfoUser(request InfoUserRequest) (InfoUserResponse, error) {
	logger := c.logger

	urlParams := map[string]string{"username": request.Username}

	response, err := c.DoGet(urlParams, infoUser, c.GetAdminCredentials())
	if err != nil {
		var errResp Err
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return InfoUserResponse{}, err
		}

		logger = logger.With().Str("error", errResp.Error).Str("errorType", errResp.ErrorType).Bool("status", errResp.Success).Logger()
		logger.Error().Msgf("info user returned with error")

		switch errResp.ErrorType {
		case ErrorUserParamNotProvidedType:
			return InfoUserResponse{}, ErrorRequiredParam{ErrorMsg: fmt.Sprintf("missing required field!")}
		case ErrorInvalidUserType:
			return InfoUserResponse{}, ErrorInvalidUser{ErrorMsg: fmt.Sprintf("user doesnt exist: %s", request.Username)}
		}

		return InfoUserResponse{}, fmt.Errorf("InfoUser returned with error: %s", errResp.Error)

	}

	var resp InfoUserResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on InfoUserResponse")
		return InfoUserResponse{}, err
	}

	return resp, nil
}
