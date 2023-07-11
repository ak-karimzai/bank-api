package grpcserver

import (
	errorhandler "github.com/ak-karimzai/bank-api/internel/error_handler"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errMap map[errorhandler.Error]codes.Code = map[errorhandler.Error]codes.Code{
	errorhandler.InternealServer: codes.Internal,
	errorhandler.Forbidden:       codes.Aborted,
	errorhandler.NotFound:        codes.NotFound,
}

func toGrpcError(err errorhandler.ResponseError) codes.Code {
	grpcErr, ok := errMap[err.Status]
	if !ok {
		return codes.Unknown
	}
	return grpcErr
}

func filedViolation(filed string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       filed,
		Description: err.Error(),
	}
}

func invalidArgumentError(
	violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{
		FieldViolations: violations,
	}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")
	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
