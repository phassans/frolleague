package controller

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/phassans/frolleague/common"
)

func (rtr *router) newPostHandler(endpoint postEndpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer rtr.cleanup(&err, w)
		logger := common.GetLogger()

		request := reflect.New(reflect.TypeOf(endpoint.HTTPRequest()))
		err = json.NewDecoder(r.Body).Decode(request.Interface())
		if err != nil {
			err = json.NewEncoder(w).Encode(APIError{ErrorMessage: err.Error(), Code: GetErrorStatus(err)})
			return
		}
		r.Body.Close()

		var result interface{}
		var engineErr error
		err = hystrixCall(endpoint, func() error {
			result, engineErr = endpoint.Execute(r.Context(), rtr, request.Elem().Interface())
			if IsHardError(engineErr) {
				return engineErr
			}
			return nil
		})
		if err == nil {
			err = engineErr
		}
		logger = logger.With().Str("endpoint", endpoint.GetPath()).Int("status", GetErrorStatus(err)).Logger()
		if err != nil {
			logger.Error().Msgf(err.Error())
			return
		}
		logger.Info().Msgf("POST success")
		err = json.NewEncoder(w).Encode(result)
	}
}
