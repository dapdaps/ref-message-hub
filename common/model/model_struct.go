package model

import (
	"fmt"
	"go/ast"
	"math/big"
	"ref-message-hub/common/ecode"
	"ref-message-hub/common/log"
	"ref-message-hub/common/referror"
	"reflect"
	"regexp"
	"strings"
)

const (
	Common            = "common"
	EthAddress        = "ethAddress"
	EthAddressNotZero = "ethAddressNotZero"
	PublicKey         = "publicKey"
	Number            = "number"
	Signature         = "signature"
	ListenKey         = "listenKey"
	Url               = "url"
)

func Validate(obj interface{}) (err error) {
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

	reflectType := reflect.ValueOf(obj).Type()
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
			switch fieldStruct.Type.Kind() {
			case reflect.String:
				str := fieldStruct.Tag.Get("degate")
				if len(str) > 0 {
					tags := strings.Split(str, ",")
					for _, value := range tags {
						if v.IsValid() {
							fieldFrom := v.Field(i)
							err = CheckTagValue(value, fieldFrom, fieldStruct.Tag.Get("json"))
							if err != nil {
								return
							}
						} else {
							err = &referror.Error{Code: ecode.RefUnknownError, Message: fmt.Sprintf("param error")}
							return
						}
					}
				}
			case reflect.Slice, reflect.Array:
				if v.IsValid() {
					err = CheckTagValue("", v.Field(i), fieldStruct.Tag.Get("json"))
					if err != nil {
						return
					}
				} else {
					err = &referror.Error{Code: ecode.RefUnknownError, Message: fmt.Sprintf("param error")}
					return
				}
			case reflect.Ptr:
				if v.IsValid() {
					err = CheckTagValue("", v.Field(i), fieldStruct.Tag.Get("json"))
					if err != nil {
						return
					}
				} else {
					err = &referror.Error{Code: ecode.RefUnknownError, Message: fmt.Sprintf("param error")}
					return
				}
			case reflect.Struct:
				if v.IsValid() {
					err = CheckTagValue("", v.Field(i), fieldStruct.Tag.Get("json"))
					if err != nil {
						return
					}
				} else {
					err = &referror.Error{Code: ecode.RefUnknownError, Message: fmt.Sprintf("param error")}
					return
				}
			}
		}
	}
	return
}

func traverseValue(reflectV reflect.Value) (err error) {
	reflectType := reflectV.Type()
	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
			switch fieldStruct.Type.Kind() {
			case reflect.String:
				str := fieldStruct.Tag.Get("degate")
				if len(str) > 0 {
					tags := strings.Split(str, ",")
					for _, value := range tags {
						fieldFrom := reflectV.Field(i)
						err = CheckTagValue(value, fieldFrom, fieldStruct.Tag.Get("json"))
						if err != nil {
							return
						}
					}
				}
			case reflect.Slice, reflect.Array:
				fmt.Println(fieldStruct.Type.Kind(), fieldStruct.Name)
				err = CheckTagValue("", reflectV.Field(i), fieldStruct.Tag.Get("json"))
				if err != nil {
					return
				}
			case reflect.Ptr:
				err = CheckTagValue("", reflectV.Field(i), fieldStruct.Tag.Get("json"))
				if err != nil {
					return
				}
			case reflect.Struct:
				err = CheckTagValue("", reflectV.Field(i), fieldStruct.Tag.Get("json"))
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func CheckTagValue(tag string, value reflect.Value, jsonName string) (err error) {
	switch value.Type().Kind() {
	case reflect.String:
		str := value.String()
		if tag == Common {
			err = CheckTagApiCommon(str, jsonName)
			return
		}
		if tag == Url {
			err = CheckTagApiUrl(str, jsonName)
			return
		}
		if tag == EthAddress {
			err = CheckTagEthAddress(str, jsonName)
			return
		}
		if tag == EthAddressNotZero {
			err = CheckEthAddressAndZeroAddress(str, jsonName)
			return
		}
		if tag == PublicKey {
			err = CheckTagPublicKey(str, jsonName)
			return
		}
		if tag == Number {
			err = CheckTagNumber(str, jsonName)
			return
		}
		if tag == Signature {
			err = CheckSignature(str, jsonName)
			return
		}
		if tag == ListenKey {
			err = CheckTagWsListenKey(str, jsonName)
			return
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			err = CheckTagValue(tag, value.Index(i), jsonName)
			if err != nil {
				return
			}
		}
	case reflect.Ptr:
		if value.IsNil() {
			return
		}
		v := value.Elem()
		err = traverseValue(v)
		if err != nil {
			return
		}
	case reflect.Struct:
		err = traverseValue(value)
		if err != nil {
			return
		}
	}
	return
}

var (
	reApiCommon    = regexp.MustCompile(`^[0-9a-zA-Z,\s]+$`) // api 常规字符串检查, 数字+字母+,+空格
	reWsSubscribe  = regexp.MustCompile("^[0-9a-zA-Z.@_]+$") // ws订阅检查, 数字+字母+ . + @ + _
	reWsListenKey  = regexp.MustCompile("^[0-9a-zA-Z.-_]+$") // listen key校验
	reETHAddr      = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	rePublicKey    = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)
	reSignatureKey = regexp.MustCompile(`^0x[0-9a-fA-F]{132,192}$`)
	reNumber       = regexp.MustCompile(`(^(0\.0*[1-9]+[0-9]*$|[1-9]+[0-9]*\.[0-9]*[0-9]$|[1-9]+[0-9]*$)|^0$)`)
	reApiUrl       = regexp.MustCompile(`^((https)?:\/\/)[^\s]+$`) // validate url
	reTest         = regexp.MustCompile(`^(\d+)+$`)
)

func CheckTagWsListenKey(value string, jsonName string) (err error) {
	if len(value) == 0 {
		return
	}
	if len(value) > 1000 {
		err = fmt.Errorf(fmt.Sprintf("%v Error: length exceeds limit", jsonName))
		return
	}
	if !reWsListenKey.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: there are illegal characters", jsonName))
		return
	}
	return
}

func CheckTagApiCommon(value string, jsonName string) (err error) {
	if len(value) == 0 {
		return
	}
	if len(value) > 1000 {
		err = fmt.Errorf(fmt.Sprintf("%v Error: length exceeds limit", jsonName))
		return
	}
	if !reApiCommon.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: there are illegal characters", jsonName))
		return
	}
	return
}

func CheckTagApiUrl(value string, jsonName string) (err error) {
	if len(value) == 0 {
		return
	}
	if len(value) > 1000 {
		err = fmt.Errorf(fmt.Sprintf("%v Error: length exceeds limit", jsonName))
		return
	}
	if !reApiUrl.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: there are illegal characters", jsonName))
		return
	}
	return
}

func CheckTagReWsSubscribe(value string, jsonName string) (err error) {
	if len(value) == 0 {
		return
	}
	if len(value) > 1000 {
		err = fmt.Errorf(fmt.Sprintf("%v Error: length exceeds limit", jsonName))
		return
	}
	if !reWsSubscribe.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: subscribe illegal", jsonName))
		return
	}
	return
}

func CheckTagEthAddress(value string, jsonName string) (err error) {
	if len(value) == 0 {
		return
	}
	if len(value) > 1000 {
		err = fmt.Errorf(fmt.Sprintf("%v Error: length exceeds limit", jsonName))
		return
	}
	if !reETHAddr.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: eth address illegal", jsonName))
		return
	}
	return
}

func CheckEthAddressAndZeroAddress(value string, jsonName string) (err error) {
	if !reETHAddr.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v error: %v not eth address", jsonName, value))
		return
	}
	var (
		b *big.Int
		o bool
	)
	b, o = new(big.Int).SetString(value[2:], 16)
	if !o {
		err = fmt.Errorf(fmt.Sprintf("%v error: %v to big int", jsonName, value))
		return
	}
	c := b.Cmp(big.NewInt(0))
	if c <= 0 {
		err = fmt.Errorf(fmt.Sprintf("%v error: %v is zero", jsonName, value))
		return
	}
	return
}

func CheckTagPublicKey(value string, jsonName string) (err error) {
	if len(value) == 0 {
		return
	}
	if len(value) > 1000 {
		err = fmt.Errorf(fmt.Sprintf("%v Error: length exceeds limit", jsonName))
		return
	}
	if !rePublicKey.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: public key illegal", jsonName))
		return
	}
	return
}

func CheckTagNumber(value string, jsonName string) (err error) {
	if len(value) == 0 {
		return
	}
	if len(value) > 1000 {
		err = fmt.Errorf(fmt.Sprintf("%v Error: length exceeds limit", jsonName))
		return
	}
	if !reNumber.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: number illegal", jsonName))
		return
	}
	return
}

func CheckSignature(value string, jsonName string) (err error) {
	if len(value) == 0 {
		return
	}
	if len(value) > 1000 {
		err = fmt.Errorf(fmt.Sprintf("%v Error: length exceeds limit", jsonName))
		return
	}
	if !reSignatureKey.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: signature illegal", jsonName))
		return
	}
	return
}

func CheckTest(value string, jsonName string) (err error) {
	if !reTest.MatchString(value) {
		err = fmt.Errorf(fmt.Sprintf("%v Error: signature illegal", jsonName))
		return
	}
	return
}
