package xerror

import (
	"errors"
	"fmt"
	"go-framework/util/xerror/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
)

func (e *Error) Error() string {
	return fmt.Sprintf("%s (code %d, status %d)", e.Message, e.Code, e.Status)
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(ToGRPCCode(int32(int(e.Code))), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   fmt.Sprintf("%d", e.Status),
			Metadata: e.Metadata,
		})
	return s
}

func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code
	}
	return false
}

func (e *Error) IsStatus(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Status == e.Status
	}
	return false
}

func New(code int32, message string) *Error {
	return &Error{
		Code:    code,
		Status:  code,
		Message: message,
	}
}

func Newf(code int32, format string, a ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, a...))
}

func (e *Error) WithMetadata(md map[string]string) *Error {
	err := proto.Clone(e).(*Error)
	err.Metadata = md
	return err
}

func (e *Error) WithStatus(status int32) *Error {
	e.Status = status
	return e
}

func Code(err error) int {
	if err == nil {
		return 200 // nolint:gomnd
	}
	if se := UnmarshalError(err); se != nil {
		return int(se.Code)
	}

	return http.StatusInternalServerError
}

func Status(err error) int32 {
	if err == nil {
		return 0 // nolint:gomnd
	}
	if se := UnmarshalError(err); se != nil {
		return se.Status
	}
	return http.StatusInternalServerError
}

func IsError(err error) bool {
	if err == nil {
		return false
	}
	if se := new(Error); errors.As(err, &se) {
		return true
	}
	return false
}

func UnmarshalError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		ret := New(FromGRPCCode(gs.Code()), gs.Message())
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				ret.Status = str2Int32(d.Reason)
				return ret.WithMetadata(d.Metadata)
			}
		}
		return ret
	}
	return New(http.StatusInternalServerError, err.Error())
}

func BadRequest(status int32, message string) *Error {
	return Newf(400, message).WithStatus(status)
}

func IsBadRequest(err error) bool {
	return Code(err) == 400
}

func Unauthorized(status int32, message string) *Error {
	return Newf(401, message).WithStatus(status)
}

func IsUnauthorized(err error) bool {
	return Code(err) == 401
}

func Forbidden(status int32, message string) *Error {
	return Newf(403, message).WithStatus(status)
}

func IsForbidden(err error) bool {
	return Code(err) == 403
}

func NotFound(status int32, message string) *Error {
	return Newf(404, message).WithStatus(status)
}

func IsNotFound(err error) bool {
	return Code(err) == 404
}

func Conflict(status int32, message string) *Error {
	return Newf(409, message).WithStatus(status)
}

func IsConflict(err error) bool {
	return Code(err) == 409
}

func InternalServer(status int32, message string) *Error {
	return Newf(500, message).WithStatus(status)
}

func IsInternalServer(err error) bool {
	return Code(err) == 500
}

func ServiceUnavailable(status int32, message string) *Error {
	return Newf(503, message).WithStatus(status)
}

func IsServiceUnavailable(err error) bool {
	return Code(err) == 503
}

func GatewayTimeout(status int32, message string) *Error {
	return Newf(504, message).WithStatus(status)
}

func IsGatewayTimeout(err error) bool {
	return Code(err) == 504
}

func ClientClosed(status int32, message string) *Error {
	return Newf(499, message).WithStatus(status)
}

func IsClientClosed(err error) bool {
	return Code(err) == 499
}

func IsStatus(err error, status int32) bool {
	return Status(err) == status
}

func str2Int32(str string) int32 {
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int32(res)
}
