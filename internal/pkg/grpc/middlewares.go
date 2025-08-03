package grpc

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func defaultMessageProducer(
	ctx context.Context,
	msg string,
	level zapcore.Level,
	code codes.Code,
	err error,
	duration zapcore.Field,
) {
	logger := ctxzap.Extract(ctx)
	params := []zap.Field{
		zap.String("grpc.code", code.String()),
		duration,
		zap.Any("request_id", ctx.Value(log.RequestIDKey)),
	}
	if err != nil {
		sts := status.Convert(err)
		msg = sts.Message()
		for _, v := range sts.Details() {
			errParams := errs.Params{}
			badRequest, ok := v.(*errdetails.BadRequest)
			if ok {
				for _, violation := range badRequest.GetFieldViolations() {
					errParams = append(
						errParams,
						errs.Param{Key: violation.GetField(), Value: violation.GetDescription()},
					)
				}
			}
			errorInfo, ok := v.(*errdetails.ErrorInfo)
			if ok {
				for key, value := range errorInfo.GetMetadata() {
					errParams = append(errParams, errs.Param{Key: key, Value: value})
				}
			}
			params = append(params, zap.Object("params", errParams))
		}
	}
	logger.Check(level, msg).Write(params...)
}
func unaryErrorServerInterceptor(ctx context.Context, req interface {
}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface {
}, error) {
	resp, err := handler(ctx, req)
	return resp, handleUnaryServerError(ctx, req, info, err)
}
func handleUnaryServerError(_ context.Context, _ any, _ *grpc.UnaryServerInfo, err error) error {
	if err == nil {
		return nil
	}
	var domainError *errs.Error
	if errors.As(err, &domainError) {
		stat := status.New(codes.Code(domainError.Code), domainError.Message)
		var withDetails *status.Status
		switch domainError.Code {
		case errs.ErrorCodeInvalidArgument:
			d := &errdetails.BadRequest{}
			for _, param := range domainError.Params {
				d.FieldViolations = append(
					d.FieldViolations,
					&errdetails.BadRequest_FieldViolation{
						Field:       param.Key,
						Description: param.Value,
					},
				)
			}
			withDetails, err = stat.WithDetails(d)
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		default:
			d := &errdetails.ErrorInfo{
				Reason:   domainError.Message,
				Domain:   "",
				Metadata: make(map[string]string),
			}
			for _, param := range domainError.Params {
				d.Metadata[param.Key] = param.Value
			}
			withDetails, err = stat.WithDetails(d)
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
		return withDetails.Err()
	}
	return status.Error(codes.Internal, err.Error())
}
