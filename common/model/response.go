package model

type BaseResponse struct {
	HttpStatusCode int    `json:"http_status_code"`
	HttpBodyText   string `json:"http_body_text"`
}
