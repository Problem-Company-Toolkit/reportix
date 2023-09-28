package reportix

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorPayload struct {
	Reason   string            `json:"reason"`
	Domain   string            `json:"domain"`
	Code     uint32            `json:"code"`
	Message  string            `json:"message"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func NewError(
	code codes.Code,
	message string,
	errDetails *errdetails.ErrorInfo,
	debug ...*errdetails.DebugInfo,
) error {
	st, _ := status.New(code, message).WithDetails(errDetails)
	for _, debug := range debug {
		st, _ = st.WithDetails(debug)
	}
	err := st.Err()
	err = errors.WithStack(err)
	return err
}

func marshal(reason string, domain string, code uint32, message string, metadata map[string]string) error {
	appErr := ErrorPayload{
		Reason:   reason,
		Domain:   domain,
		Message:  message,
		Code:     code,
		Metadata: metadata,
	}

	encodedErr, err := json.Marshal(appErr)
	if err != nil {
		log.Printf("Failed to serialize error: %v", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	errMessage := status.Error(codes.Code(code), string(encodedErr))
	return errMessage
}

func ErrorToJSON(ctx context.Context, err error) error {
	if err == nil {
		return err
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	if st.Details() == nil || len(st.Details()) == 0 {
		return err
	}

	var errMetadata map[string]string
	detail := st.Details()[0].(*errdetails.ErrorInfo)
	errMetadata = func() map[string]string {
		if len(detail.GetMetadata()) > 0 {
			metadata := detail.GetMetadata()
			return metadata
		}
		return nil
	}()

	err = marshal(detail.GetReason(), detail.GetDomain(), uint32(st.Code()), st.Message(), errMetadata)
	return err
}
