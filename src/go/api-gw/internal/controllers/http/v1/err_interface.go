package v1

import "api-gateway/pkg/models"

type ResponseErrorInterface interface {
	Code() int
	Error() string
	GetPayload() *models.DtoErrorResponse
	IsClientError() bool
	IsCode(code int) bool
	IsRedirect() bool
	IsServerError() bool
	IsSuccess() bool
	String() string
}
