package usersvc

import (
  "encoding/json"

  "gitlab.com/abios/api-dash-backend/errorspec"
)

// parseResponse accepts an array of bytes and tries to parse it as either an
// user-svc Error Specification or as the specified target struct
func parseResponse(body []byte, target interface{}) error {
	var errspec errorspec.Error
  // Attempt to parse the payload as an user-svc Error Specification
	if err := json.Unmarshal(body, &errspec); err == nil {
		return errspec
	}
  // Otherwise attempt to parse the payload as the target struct
	return json.Unmarshal(body, target)
}
