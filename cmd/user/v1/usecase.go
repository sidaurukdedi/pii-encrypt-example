package user

import (
	"context"
	"net/http"
	"pii-encrypt-example/entity"
	"pii-encrypt-example/pkg/crypto"
	"pii-encrypt-example/pkg/exception"
	"pii-encrypt-example/pkg/response"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserUsecase interface {
	GetManyUsers(ctx context.Context, filter UserFilter) (resp response.Response)
	CreateUser(ctx context.Context, userRequest UserRequest) (resp response.Response)
}

type userUsecase struct {
	logger         *logrus.Logger
	location       *time.Location
	crypto         crypto.Crypto
	userRepository UserRepository
}

func NewUserUsecase(logger *logrus.Logger, location *time.Location, crypto crypto.Crypto, userRepository UserRepository) UserUsecase {
	return &userUsecase{
		logger:         logger,
		location:       location,
		crypto:         crypto,
		userRepository: userRepository,
	}
}

// GetManyUsers implements Usecase
func (u *userUsecase) GetManyUsers(ctx context.Context, filter UserFilter) (resp response.Response) {

	if filter.Name != "" {
		filter.NameHashed = u.crypto.Hash(filter.Name)
	}

	result, err := u.userRepository.FindManyUser(ctx, filter)
	if err != nil {
		if err == exception.ErrNotFound {
			return response.NewErrorResponse(err, http.StatusNotFound, nil, response.StatNotFound, "")
		}
		u.logger.WithContext(ctx).Error(err)
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}

	totalDataOnPage := len(result)
	usersResponse := make([]UserResponse, totalDataOnPage)
	for i, v := range result {
		var userResponse UserResponse
		userResponse.UUID = v.UUID

		name, err := u.crypto.Decrypt(v.NameCrypt)
		if err != nil {
			u.logger.WithContext(ctx).Error(err)
			// return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
		}
		userResponse.Name = string(name)

		email, err := u.crypto.Decrypt(v.EmailCrypt)
		if err != nil {
			u.logger.WithContext(ctx).Error(err)
			// return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
		}
		userResponse.Email = string(email)

		userResponse.CreatedAt = v.CreatedAt

		usersResponse[i] = userResponse
	}

	return response.NewSuccessResponse(usersResponse, response.StatOK, "")
}

// CreateUser implements Usecase
func (u *userUsecase) CreateUser(ctx context.Context, userRequest UserRequest) (resp response.Response) {
	var user entity.User

	uuid := uuid.New().String()
	user.UUID = uuid
	nameCrypt, err := u.crypto.Encrypt(userRequest.Name)
	if err != nil {
		u.logger.WithContext(ctx).Error(err)
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}
	user.NameCrypt = nameCrypt
	user.NameHash = u.crypto.Hash(userRequest.Name)

	emailCrypt, err := u.crypto.Encrypt(userRequest.Email)
	if err != nil {
		u.logger.WithContext(ctx).Error(err)
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}
	user.EmailCrypt = emailCrypt

	createdAt := time.Now().In(u.location)
	user.CreatedAt = createdAt

	_, err = u.userRepository.SaveUser(ctx, user, nil)
	if err != nil {
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}

	userResponse := UserResponse{
		UUID:      uuid,
		Name:      userRequest.Name,
		Email:     userRequest.Email,
		CreatedAt: createdAt,
	}

	return response.NewSuccessResponse(userResponse, response.StatOK, "")
}
