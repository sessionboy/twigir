package res

type Response struct {
	Ok    bool        `json:"ok"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func Ok(msg string, data interface{}, args ...int) Response {
	return Response{true, msg, data, nil}
}

func Err(msg string, args ...int) Response {
	return Response{false, msg, nil, nil}
}

func ErrWithError(msg string, _error interface{}, args ...int) Response {
	return Response{false, msg, nil, _error}
}

func ErrWithData(msg string, data interface{}, args ...int) Response {
	return Response{false, msg, data, nil}
}
