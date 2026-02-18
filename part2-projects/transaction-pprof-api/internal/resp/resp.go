package resp

// 统一错误码（当手册查：前后端约定、错误码分层）
const (
	CodeOK                 = 0
	CodeInternal           = 1
	CodeBadRequest         = 2
	CodeInsufficientBalance = 3
	CodeNotFound           = 4
)

// Resp 统一响应体（示例）：前后端约定可扩展 trace_id、error_detail 等
type Resp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func OK[T any](data T) Resp[T] {
	return Resp[T]{Code: CodeOK, Message: "ok", Data: data}
}

func Fail(code int, msg string) Resp[any] {
	if code == 0 {
		code = CodeInternal
	}
	return Resp[any]{Code: code, Message: msg}
}

