package reportix_test

import (
	"context"
	"encoding/json"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/problem-company-toolkit/reportix"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// "google.golang.org/grpc/status"
)

var _ = Describe("Interceptor", func() {
	var (
		ctx     context.Context
		grpcErr error
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler

		reason    string = "MISSING_FIELD"
		msg       string = "missing required field"
		field     string = "id"
		object    string = "user"
		errDomain        = "foo-service"
		code             = codes.InvalidArgument

		debug []*errdetails.DebugInfo
	)

	BeforeEach(func() {
		ctx = context.Background()
		grpcErr = reportix.NewError(
			code,
			fmt.Sprintf("invalid uuid provided for field %s in object %s. %s", field, object, msg),
			&errdetails.ErrorInfo{
				Reason: reason,
				Domain: errDomain,
				Metadata: map[string]string{
					"field":  field,
					"object": object,
					"error":  msg,
				},
			},
			debug...,
		)

		info = &grpc.UnaryServerInfo{FullMethod: "/test.v1/TestMethod"}

	})

	Context("when the request is processed", func() {
		It("should mount grpc error to errorPayload json string", func() {
			handler = func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, reportix.ErrorToJSON(ctx, grpcErr)
			}

			var err error
			callback := func(ctx context.Context, grpcErr error) error {
				err = grpcErr
				return grpcErr
			}

			interceptor := reportix.NewErrInterceptor(reportix.ErrInterceptorOpts{Callback: callback})
			interceptor.UnaryServerInterceptor()(ctx, nil, info, handler)

			Expect(err).ShouldNot(BeNil())
			Expect(strings.Contains(err.Error(), "reason")).Should(BeTrue())
			Expect(strings.Contains(err.Error(), "domain")).Should(BeTrue())
			Expect(strings.Contains(err.Error(), "message")).Should(BeTrue())
			Expect(strings.Contains(err.Error(), "field")).Should(BeTrue())

			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus).ToNot(BeNil())

			errorPayload := reportix.ErrorPayload{}
			err = json.Unmarshal([]byte(grpcStatus.Message()), &errorPayload)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(errorPayload.Reason).To(BeEquivalentTo(reason))
			Expect(errorPayload.Domain).To(BeEquivalentTo(errDomain))
			Expect(errorPayload.Message).To(ContainSubstring(msg))
			Expect(errorPayload.Code).To(BeEquivalentTo(code))

			metadata := errorPayload.Metadata
			Expect(metadata).ToNot(BeNil())

			fieldName, ok := metadata["field"]
			Expect(ok).To(BeTrue())
			Expect(fieldName).To(BeEquivalentTo(field))
		})
	})
})
