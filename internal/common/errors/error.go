package errors

import (
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrExists = errors.New("already exists.")
var ErrNotExists = errors.New("not found.")

func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters.")

	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}

func ErrNotFound() error {
	return status.Error(codes.NotFound, "not found.")
}

func ErrInternal() error {
	return status.Error(codes.Internal, "internal error.")
}

func LogErrAndCreateInternal(err error) error {
	LogErr(err)
	return status.Error(codes.Internal, "internal error.")
}

func LogErr(err error) {
	log.Err(err).Stack().Msg("")
}

func UnAuthErr() error {
	return status.Error(codes.Unauthenticated, "auth error.")
}
