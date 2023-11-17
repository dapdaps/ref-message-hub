package referror

import "encoding/json"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
