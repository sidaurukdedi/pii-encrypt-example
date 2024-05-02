package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pii-encrypt-example/pkg/middleware"
	"pii-encrypt-example/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UserHTTPHandler struct {
	logger      *logrus.Logger
	validator   *validator.Validate
	userUsecase UserUsecase
}

func NewUserHTTPHandler(logger *logrus.Logger, router *mux.Router, basicAuth middleware.RouteMiddleware, validator *validator.Validate, userUsecase UserUsecase) {
	handler := &UserHTTPHandler{
		logger:      logger,
		validator:   validator,
		userUsecase: userUsecase,
	}
	router.HandleFunc("/api/v1/user", basicAuth.Verify(handler.GetManyUsers)).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/user", basicAuth.Verify(handler.CreateUser)).Methods(http.MethodPost)
}

func (h UserHTTPHandler) GetManyUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var filter UserFilter

	queryString := r.URL.Query()
	nameQS := queryString.Get("name")
	if nameQS != "" {
		filter.Name = nameQS
	}
	resp := h.userUsecase.GetManyUsers(ctx, filter)
	response.JSON(w, resp)
}

func (h UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var payload UserRequest

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		resp = response.NewErrorResponse(err, http.StatusUnprocessableEntity, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	if err := h.validateRequestBody(payload); err != nil {
		resp = response.NewErrorResponse(err, http.StatusBadRequest, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	resp = h.userUsecase.CreateUser(ctx, payload)
	response.JSON(w, resp)
}

func (h UserHTTPHandler) validateRequestBody(body interface{}) (err error) {
	err = h.validator.Struct(body)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	errorField := errorFields[0]
	err = fmt.Errorf("invalid '%s' with value '%v'", errorField.Field(), errorField.Value())

	return
}
