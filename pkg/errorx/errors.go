package errorx

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(base *Error, err error) *Error {
	if err == nil {
		return base
	}
	return &Error{Code: base.Code, Message: base.Message, Err: err}
}

func WithDetail(base *Error, detail string) *Error {
	message := base.Message
	if detail != "" {
		message = message + ": " + detail
	}
	return &Error{Code: base.Code, Message: message}
}

var (
	ErrInvalidConfig  = New("invalid_config_error", "invalid configuration supplied")
	ErrMarshalRequest = New("marshal_request_error", "failed to marshal request payload")
	ErrBuildRequest   = New("build_request_error", "failed to build HTTP request")
	ErrSendRequest    = New("send_request_error", "failed to send HTTP request")
	ErrAPIResponse    = New("api_response_error", "received error response from WhatsApp Cloud API")
)
