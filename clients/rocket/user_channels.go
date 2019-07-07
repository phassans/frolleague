package rocket

import (
	"encoding/json"
	"fmt"
)

func (c *client) AddUserToGroup(request AddUserToGroupRequest) (AddUserToGroupResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request, addUserToGroup, c.GetAdminCredentials())
	if err != nil {
		var errResp GroupErrorResponse
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return AddUserToGroupResponse{}, err
		}
		logger = logger.With().Bool("success", errResp.Success).Str("error", errResp.Error).Str("errorType", errResp.ErrorType).Logger()
		logger.Error().Msgf("AddUserToGroup returned with error")

		switch errResp.ErrorType {
		case ErrorGroupNotFoundType:
			return AddUserToGroupResponse{}, ErrorGroupNotFound{GroupName: request.RoomId}
		case ErrorInvalidUserType:
			return AddUserToGroupResponse{}, ErrorInvalidUser{ErrorMsg: fmt.Sprintf("user doesnt exist: %s", request.UserId)}
		case ErrorGroupParamNotProvidedType:
			return AddUserToGroupResponse{}, ErrorRequiredParam{fmt.Sprintf("required param missing!")}
		case ErrorUserParamNotProvidedType:
			return AddUserToGroupResponse{}, ErrorRequiredParam{ErrorMsg: fmt.Sprintf("missing required field!")}
		}

		return AddUserToGroupResponse{}, fmt.Errorf("AddUserToGroup returned with error: %s, type: %s", errResp.Error, errResp.ErrorType)
	}

	// read response.json
	var resp AddUserToGroupResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on AddUserToGroupResponse")
		return AddUserToGroupResponse{}, err
	}

	return resp, nil
}

func (c *client) RemoveUserFromGroup(request RemoveUserFromGroupRequest) (RemoveUserFromGroupResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request, removeUserFromGroup, c.GetAdminCredentials())
	if err != nil {
		var errResp GroupErrorResponse
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return RemoveUserFromGroupResponse{}, err
		}
		logger = logger.With().Bool("success", errResp.Success).Str("error", errResp.Error).Str("errorType", errResp.ErrorType).Logger()
		logger.Error().Msgf("RemoveUserFromGroupResponse returned with error")

		switch errResp.ErrorType {
		case ErrorGroupNotFoundType:
			return RemoveUserFromGroupResponse{}, ErrorGroupNotFound{GroupName: request.RoomId}
		case ErrorInvalidUserType:
			return RemoveUserFromGroupResponse{}, ErrorInvalidUser{ErrorMsg: fmt.Sprintf("user doesnt exist: %s", request.UserId)}
		case ErrorGroupParamNotProvidedType:
			return RemoveUserFromGroupResponse{}, ErrorRequiredParam{fmt.Sprintf("required param missing!")}
		case ErrorUserParamNotProvidedType:
			return RemoveUserFromGroupResponse{}, ErrorRequiredParam{ErrorMsg: fmt.Sprintf("missing required field!")}
		}

		return RemoveUserFromGroupResponse{}, fmt.Errorf("RemoveUserToGroup returned with error: %s, type: %s", errResp.Error, errResp.ErrorType)
	}

	// read response.json
	var resp RemoveUserFromGroupResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on RemoveUserFromGroupResponse")
		return RemoveUserFromGroupResponse{}, err
	}

	return resp, nil
}
