package response

import (
	"encoding/json"
	"net/http"
)

// REST is a collection of behavior of REST.
type REST interface {
	JSON(w http.ResponseWriter)
}

type restObject struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Meta    interface{} `json:"meta,omitempty"` // will not be appeared if not set.
	// can add more
}

// RESTResponse is a model of REST response.
// type RESTResponse struct {
// 	Success bool        `json:"success"`
// 	Data    interface{} `json:"data"`
// 	Message string      `json:"message"`
// 	Status  string      `json:"status"`
// 	Code    int         `json:"code"`
// 	Meta    interface{} `json:"meta,omitempty"` // will not be appeared if not set.
// 	// can add more
// }

// NewRESTResponse will return object of rest response.
// func NewRESTResponse(response Response) REST {
// 	restResp := RESTResponse{}
// 	if response.Error() == nil {
// 		restResp.Success = true
// 	}
// 	restResp.Data = response.Data()
// 	restResp.Message = response.Message()
// 	restResp.Status = response.Status()
// 	restResp.Code = response.HTTPStatusCode()
// 	restResp.Meta = nil

// 	return restResp
// }

// // JSON will response json.
// func (r RESTResponse) JSON(w http.ResponseWriter) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(r.Code)
// 	json.NewEncoder(w).Encode(r)
// }

// JSON will response as json serialization.
func JSON(w http.ResponseWriter, resp Response) {
	var success bool
	if resp.Error() == nil {
		success = true
	}
	ro := restObject{
		Success: success,
		Data:    resp.Data(),
		Message: resp.Message(),
		Status:  resp.Status(),
		Code:    resp.HTTPStatusCode(),
		Meta:    resp.Meta(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ro.Code)
	json.NewEncoder(w).Encode(ro)
}
