package response

import (
	"encoding/json"
	"fmt"
)

// BuildErrorFromResponse wrap the error that contains error, and translate to error interface with informatif description.
func BuildErrorFromResponse(r Response) error {
	if r.Error() == nil {
		return nil
	}

	errMessage := map[string]interface{}{
		"code":    r.HTTPStatusCode(),
		"data":    r.Data(),
		"message": r.Message(),
		"status":  r.Status(),
	}

	errMessageBuff, _ := json.Marshal(errMessage)
	return fmt.Errorf(string(errMessageBuff))
}
