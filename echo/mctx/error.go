package mctx

// varible global
var (
	OkStatus = BaseStatus{Code: "SUCCESS", Message: "Successful"}
)

// BaseStatus is class
type BaseStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ParseStatus is func init status
func ParseStatus(code string, message string) BaseStatus {
	return BaseStatus{
		Code:    code,
		Message: message,
	}
}

// HasError is func test err
func (e BaseStatus) HasError() bool {
	if e.Code == "SUCCESS" {
		return false
	}
	return e.Code != ""
}

// GetMsg is func test err
func (e BaseStatus) GetMsg() string {
	return "[" + e.Code + "] " + e.Message
}

// NewError is func return new error json
func NewError() BaseStatus {
	return BaseStatus{
		Code:    "ERROR",
		Message: "",
	}
}
