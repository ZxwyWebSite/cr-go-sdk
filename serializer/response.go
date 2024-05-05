package serializer

// Response 基础序列化器
type Response[T any] struct {
	Code  int    `json:"code"`
	Data  T      `json:"data,omitempty"`
	Msg   string `json:"msg"`
	Error string `json:"error,omitempty"`
}

// 实现错误检测接口
func (c *Response[T]) Err() error {
	if c.Code != 0 {
		return &AppError{
			Code: c.Code,
			Msg:  c.Msg,
		} //errors.New(strconv.Itoa(c.Code) + `: ` + c.Msg)
	}
	return nil
}
