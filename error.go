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
