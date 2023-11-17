package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"ref-message-hub/common/log"
	"ref-message-hub/common/model"
	"strconv"
	"time"
)

type Client struct {
	client     *http.Client
	sampleRate float64
}

func New(timeout time.Duration, sampleRate float64) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		sampleRate: sampleRate,
	}
}

func (client *Client) GetJSON(url string, header interface{}, params interface{}, result interface{}) (err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	var p *map[string]string
	if p, err = parseStruct(params); err != nil {
		return
	} else if p != nil {
		query := request.URL.Query()
		for k, v := range *p {
			query.Add(k, v)
		}
		request.URL.RawQuery = query.Encode()
	}

	if err = addHeader(&request.Header, header); err != nil {
		return
	}
	return client.req(request, params, result)
}

func (client *Client) PostJSON(url string, header interface{}, params interface{}, result interface{}) (err error) {
	body, err := encodeBody(params)
	if err != nil {
		return
	}
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return
	}
	if err = addHeader(&request.Header, header); err != nil {
		return
	}
	return client.req(request, params, result)
}

func (client *Client) DeleteJSON(url string, header interface{}, params interface{}, result interface{}) (err error) {
	body, err := encodeBody(params)
	if err != nil {
		return
	}
	request, err := http.NewRequest(http.MethodDelete, url, body)
	if err != nil {
		return
	}
	if err = addHeader(&request.Header, header); err != nil {
		return
	}
	return client.req(request, params, result)
}

func (client *Client) PutJSON(url string, header interface{}, params interface{}, result interface{}) (err error) {
	body, err := encodeBody(params)
	if err != nil {
		return
	}
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return
	}
	if err = addHeader(&request.Header, header); err != nil {
		return
	}
	return client.req(request, params, result)
}

func (client *Client) req(request *http.Request, params interface{}, result interface{}) (err error) {
	var httpResponse *http.Response
	start := time.Now()
	httpResponse, err = client.client.Do(request)
	log.Info("request %v %v, took %v", request.Method, request.URL, time.Since(start))
	if err != nil {
		log.Error("error getting %v: %v", request.URL.String(), err)
		return
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if result != nil {
		var body bytes.Buffer
		_, err = io.Copy(&body, httpResponse.Body)
		if err != nil {
			return
		}
		statusCode := httpResponse.StatusCode
		if statusCode != 200 {
			log.Info("raw response from server statusCode: %v", statusCode)
		}
		if v, ok := result.(*model.BaseResponse); ok {
			v.HttpStatusCode = statusCode
			v.HttpBodyText = strconv.Quote(body.String())
		} else {
			statusCodeJson := fmt.Sprintf("{\"http_status_code\":%v}", statusCode)
			_ = json.Unmarshal([]byte(statusCodeJson), result)
			err = json.Unmarshal(body.Bytes(), result)
			if err != nil {
				log.Error("Unmarshal %v: %v", request.URL.String(), err)
			}
		}
	}
	return
}

func encodeBody(v interface{}) (body io.Reader, err error) {
	if v == nil {
		return
	}

	switch b := v.(type) {
	case string:
		body = bytes.NewReader([]byte(b))
	case []byte:
		body = bytes.NewReader(b)
	default:
		var bs []byte
		bs, err = json.Marshal(v)
		if err != nil {
			return
		}
		body = bytes.NewReader(bs)
	}
	return
}

func parseStruct(v interface{}) (p *map[string]string, err error) {
	if v == nil {
		return
	}
	var bs []byte
	bs, err = json.Marshal(v)
	if err != nil {
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(bs, &m)
	if err != nil {
		return
	}
	p = &map[string]string{}
	for k, value := range m {
		switch value.(type) {
		case string:
			(*p)[k] = value.(string)
			break
		case int:
			(*p)[k] = strconv.Itoa(value.(int))
			break
		case float64:
			f := decimal.NewFromFloat(value.(float64))
			if decimal.NewFromBigInt(f.BigInt(), 0).Equal(f) {
				(*p)[k] = f.BigInt().String()
			} else {
				(*p)[k] = f.String()
			}
			break
		}
	}
	return
}

func addHeader(header *http.Header, v interface{}) (err error) {
	header.Add("content-type", "application/json")
	var m *map[string]string
	m, err = parseStruct(v)
	if err != nil {
		return
	}
	if m == nil {
		return
	}
	for k, v := range *m {
		header.Add(k, v)
	}
	return
}
