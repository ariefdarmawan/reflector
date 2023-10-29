package reflector

import (
	"errors"
	"reflect"

	"github.com/sebarcode/codekit"
)

func CopyAttributes[T any](source interface{}, dest T, excludeFieldNames ...string) (T, error) {
	rs := reflect.Indirect(reflect.ValueOf(source))
	if rs.Kind() != reflect.Struct {
		return dest, errors.New("source should be ptr of struct or struct")
	}
	rd := reflect.ValueOf(dest)
	if rd.Kind() != reflect.Ptr {
		return dest, errors.New("dest should be ptr of struct")
	}
	rde := rd.Elem()
	if rde.Kind() != reflect.Struct {
		return dest, errors.New("dest should be ptr of struct")
	}
	var err error
	func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
		sourceType := rs.Type()
		destType := rde.Type()
		destFieldNames := make([]string, rde.NumField())
		sourceFieldCount := rs.NumField()
		destFieldCount := rs.NumField()
		for i := 0; i < destFieldCount; i++ {
			destFieldNames[i] = destType.Field(i).Name
		}
		for i := 0; i < sourceFieldCount; i++ {
			fieldName := sourceType.Field(i).Name
			if codekit.HasMember(destFieldNames, fieldName) && !codekit.HasMember(excludeFieldNames, fieldName) {
				rde.FieldByName(fieldName).Set(rs.Field(i))
			}
		}
		rd.Set(rde)
	}()
	return dest, err
}

func CopyAttributeByNames[T any](source interface{}, dest T, copiedFieldNames ...string) (T, error) {
	rs := reflect.Indirect(reflect.ValueOf(source))
	if rs.Kind() != reflect.Struct {
		return dest, errors.New("source should be ptr of struct or struct")
	}
	rd := reflect.ValueOf(dest)
	if rd.Kind() != reflect.Ptr {
		return dest, errors.New("dest should be ptr of struct")
	}
	rde := rd.Elem()
	if rde.Kind() != reflect.Struct {
		return dest, errors.New("dest should be ptr of struct")
	}
	var err error
	func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
		sourceType := rs.Type()
		destType := rde.Type()
		destFieldNames := make([]string, rde.NumField())
		sourceFieldCount := rs.NumField()
		destFieldCount := rs.NumField()
		for i := 0; i < destFieldCount; i++ {
			destFieldNames[i] = destType.Field(i).Name
		}
		for i := 0; i < sourceFieldCount; i++ {
			fieldName := sourceType.Field(i).Name
			if codekit.HasMember(destFieldNames, fieldName) && codekit.HasMember(copiedFieldNames, fieldName) {
				rde.FieldByName(fieldName).Set(rs.Field(i))
			}
		}
		rd.Set(rde)
	}()
	return dest, err
}
