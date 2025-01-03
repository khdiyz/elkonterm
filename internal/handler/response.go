package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// const (
// 	createdMessage = "created"
// 	successMessage = "success"
// )

type BaseResponse struct {
	Message string `json:"message"`
}

type IdResponse struct {
	ID string `json:"id"`
}

// Helper function to return an error response
func errorResponse(c *gin.Context, status int, err error) {
	c.JSON(status, BaseResponse{
		Message: err.Error(),
	})
}

// Converts a gRPC error into an HTTP response
func fromError(c *gin.Context, serviceError error) {
	st, _ := status.FromError(serviceError)
	err := st.Message()

	switch st.Code() {
	case codes.NotFound:
		errorResponse(c, http.StatusNotFound, errors.New(err))
	case codes.InvalidArgument:
		errorResponse(c, http.StatusBadRequest, errors.New(err))
	case codes.Unavailable:
		errorResponse(c, http.StatusUnavailableForLegalReasons, errors.New(err))
	case codes.AlreadyExists:
		errorResponse(c, http.StatusBadRequest, errors.New(err))
	case codes.Unauthenticated:
		errorResponse(c, http.StatusUnauthorized, errors.New(err))
	default:
		errorResponse(c, http.StatusInternalServerError, errors.New(err))
	}
}
