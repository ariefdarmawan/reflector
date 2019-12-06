package reflector

import (
	"errors"
	"reflect"
)

type reflector struct {
	ptr reflect.Value
	v   reflect.Value
	t   reflect.Type

	err error
}

func (r *reflector) setError(msg string) *reflector {
	r.err = errors.New(msg)
	return r
}

func From(obj interface{}) *reflector {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return new(reflector).setError("source object should be pointer of struct")
	}

	if v.Elem().Kind() != reflect.Struct {
		return new(reflector).setError("source object should be pointer of struct")
	}

	r := new(reflector)
	r.ptr = v
	r.v = v.Elem()
	r.t = reflect.TypeOf(obj).Elem()
	return r
}

func (r *reflector) Set(name string, value interface{}) *reflector {
	if r.err != nil {
		return r
	}

	func() {
		defer func() {
			if rec := recover(); rec != nil {
			}
		}()

		v := r.v.FieldByName(name)
		v.Set(reflect.ValueOf(value))
	}()
	return r
}

func (r *reflector) Flush() error {
	if r.err != nil {
		return r.err
	}

	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(r.(string))
			}
		}()

		r.ptr.Elem().Set(r.v)
	}()
	return err
}

func (r *reflector) FieldNames(tag string) ([]string, error) {
	if r.err != nil {
		return []string{}, r.err
	}

	fieldNum := r.t.NumField()
	fields := make([]string, fieldNum)
	//fmt.Println("num of fields:", fieldNum)

	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(r.(string))
			}
		}()

		fieldIdx := 0
		for {
			f := r.t.Field(fieldIdx)
			fn := f.Name
			if tag != "" {
				if tn := f.Tag.Get(tag); tn != "" {
					fn = tn
				}
			}
			fields[fieldIdx] = fn
			//fmt.Println(fieldIdx, fn)

			fieldIdx++
			if fieldIdx >= fieldNum {
				break
			}
		}
	}()
	return fields, err
}
