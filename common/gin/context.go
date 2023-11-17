package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ref-message-hub/common/ecode"
	"ref-message-hub/common/log"
	"ref-message-hub/common/model"
	"ref-message-hub/common/referror"
	"reflect"
	"strings"
	"sync/atomic"
)

var findClientIPError = atomic.Uint64{}

func GetRealIP(c *gin.Context) string {
	IPValue := c.Request.Header.Get("X-Original-Forwarded-For")
	IPList := strings.Split(IPValue, ",")
	if l := len(IPList); l > 0 {
		return strings.TrimSpace(IPList[l-1])
	}

	findClientIPError.Add(1)
	log.Error("No sufficient IPs found in Header for key X-Original-Forwarded-For,  IPValue : %v", IPValue)
	return ""
}

func ShouldBind(c *gin.Context, obj interface{}) (err error) {
	err = c.ShouldBind(obj)
	if err != nil {
		return
	}

	var (
		v = reflect.ValueOf(obj).Elem()
	)
	for v.Kind() == reflect.Slice || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		log.Error("from and to must be pointer")
		err = fmt.Errorf("from and to must be pointer")
		return
	}

	if !v.IsValid() {
		err = &referror.Error{Code: ecode.RefUnknownError, Message: fmt.Sprintf("param error")}
		return
	}

	err = model.Validate(obj)
	if err != nil {
		return
	}
	return
}
