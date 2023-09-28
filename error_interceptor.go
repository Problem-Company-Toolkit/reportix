package reportix

import (
	"context"

	"google.golang.org/grpc"
)

type CallbackFunc = func(context.Context, error) error

type ErrInterceptor struct {
	callback CallbackFunc
}

type ErrInterceptorOpts struct {
	Callback CallbackFunc
}

func NewErrInterceptor(opts ErrInterceptorOpts) *ErrInterceptor {
	if opts.Callback == nil {
		opts.Callback = func(ctx context.Context, err error) (formatedError error) {
			return formatedError
		}
	}

	return &ErrInterceptor{
		callback: opts.Callback,
	}
}

func (e *ErrInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		err = e.callback(ctx, err)

		return resp, err
	}
}
