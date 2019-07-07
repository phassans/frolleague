package rocket

import (
	"encoding/json"
	"fmt"
)

func (c *client) CreateGroup(request GroupCreateRequest) (GroupCreateResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request, createGroup, c.GetAdminCredentials())
	if err != nil {
		var errResp GroupErrorResponse
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return GroupCreateResponse{}, err
		}
		logger = logger.With().Bool("success", errResp.Success).Str("error", errResp.Error).Logger()
		logger.Error().Msgf("CreateGroup returned with error")

		switch errResp.ErrorType {
		case ErrorDuplicateGroupNameType:
			return GroupCreateResponse{}, ErrorDuplicateGroupName{GroupName: request.Name}
		case ErrorInvalidGroupNameType:
			return GroupCreateResponse{}, ErrorInvalidGroupName{GroupName: request.Name}
		}

		if errResp.Error == "Body param \"name\" is required" {
			return GroupCreateResponse{}, ErrorRequiredParam{ErrorMsg: fmt.Sprintf("missing required field name!")}
		}

		return GroupCreateResponse{}, fmt.Errorf("CreateGroup returned with error: %s, type: %s", errResp.Error, errResp.ErrorType)
	}

	// read response.json
	var resp GroupCreateResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on GroupCreateResponse")
		return GroupCreateResponse{}, err
	}

	return resp, nil
}

func (c *client) DeleteGroup(request DeleteGroupRequest) (DeleteGroupResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request, deleteGroup, c.GetAdminCredentials())
	if err != nil {
		var errResp GroupErrorResponse
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return DeleteGroupResponse{}, err
		}
		logger = logger.With().
			Bool("success", errResp.Success).
			Str("error", errResp.Error).
			Str("errorType", errResp.ErrorType).
			Logger()
		logger.Error().Msgf("DeleteGroup returned with error")

		switch errResp.ErrorType {
		case ErrorGroupNotFoundType:
			return DeleteGroupResponse{}, ErrorGroupNotFound{GroupName: request.RoomId}
		case ErrorGroupParamNotProvidedType:
			return DeleteGroupResponse{}, ErrorRequiredParam{fmt.Sprintf("required param missing!")}
		}

		return DeleteGroupResponse{}, fmt.Errorf("DeleteGroup returned with error: %s, type: %s", errResp.Error, errResp.ErrorType)
	}

	// read response.json
	var resp DeleteGroupResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on DeleteGroupResponse")
		return DeleteGroupResponse{}, err
	}

	return resp, nil
}

func (c *client) InfoGroup(request InfoGroupRequest) (InfoGroupResponse, error) {
	logger := c.logger

	urlParams := map[string]string{"roomName": request.RoomName}

	response, err := c.DoGet(urlParams, infoGroup, c.GetAdminCredentials())
	if err != nil {
		var errResp GroupErrorResponse
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return InfoGroupResponse{}, err
		}

		logger = logger.With().Str("error", errResp.Error).Str("errorType", errResp.ErrorType).Bool("status", errResp.Success).Logger()
		logger.Error().Msgf("InfoGroup returned with error")

		switch errResp.ErrorType {
		case ErrorGroupNotFoundType:
			return InfoGroupResponse{}, ErrorGroupNotFound{GroupName: request.RoomName}
		case ErrorGroupParamNotProvidedType:
			return InfoGroupResponse{}, ErrorRequiredParam{fmt.Sprintf("required param missing!")}
		}
		return InfoGroupResponse{}, fmt.Errorf("InfoGroup returned with error: %s", errResp.Error)

	}

	// read response.json
	var resp InfoGroupResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on InfoGroupResponse")
		return InfoGroupResponse{}, err
	}

	return resp, nil
}

// not needed so currently no tests
func (c *client) SetTypeGroup(request SetTypeGroupRequest) (SetTypeGroupResponse, error) {
	logger := c.logger

	response, err := c.DoPost(request, setTypeGroup, c.GetAdminCredentials())
	if err != nil {
		var errResp GroupErrorResponse
		err = json.Unmarshal(response, &errResp)
		if err != nil {
			logger = logger.With().Str("error", err.Error()).Logger()
			logger.Error().Msgf("unmarshal error on ErrorResponse")
			return SetTypeGroupResponse{}, err
		}
		logger = logger.With().
			Bool("success", errResp.Success).
			Str("error", errResp.Error).
			Str("errorType", errResp.ErrorType).
			Logger()
		logger.Error().Msgf("SetTypeGroup returned with error")
		return SetTypeGroupResponse{}, fmt.Errorf("SetTypeGroup returned with error: %s, type: %s", errResp.Error, errResp.ErrorType)
	}

	// read response.json
	var resp SetTypeGroupResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		logger = logger.With().Str("error", err.Error()).Logger()
		logger.Error().Msgf("unmarshal error on SetTypeGroupResponse")
		return SetTypeGroupResponse{}, err
	}

	return resp, nil
}
