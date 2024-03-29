package nuonuo

import "fmt"

type Error struct {
	Code string
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

// 订单编号或流水号重复
func (e *Error) IsDuplicateOrderNo() bool {
	return e.Code == "E9106"
}

// 是否公共异常码
func (e *Error) IsCommon() bool {
	return len(e.Code) > 0 && e.Code[0] == '0'
}
