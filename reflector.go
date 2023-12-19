package reflector

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
)

type Reflector struct {
	ptr reflect.Value
	v   reflect.Value
	t   reflect.Type

	err error
}

func (r *Reflector) setError(msg string) *Reflector {
	r.err = errors.New(msg)
	return r
}

func (r *Reflector) Error() error {
	return r.err
}

func From(obj interface{}) *Reflector {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return new(Reflector).setError("source object should be pointer of struct")
	}

	if v.Elem().Kind() != reflect.Struct {
		return new(Reflector).setError("source object should be pointer of struct")
	}

	r := new(Reflector)
	r.ptr = v
	r.v = v.Elem()
	r.t = reflect.TypeOf(obj).Elem()
	return r
}

func (r *Reflector) Get(name string) (interface{}, error) {
	v, err := r.getValue(r.v, name)
	if err != nil {
		return nil, err
	}
	return v.Interface(), nil
}

func (r *Reflector) GetTo(name string, dest interface{}) error {
	v, err := r.getValue(r.v, name)
	if err != nil {
		return err
	}
	dv := reflect.ValueOf(dest)
	if dv.Kind() != reflect.Ptr {
		return errors.New("dest should be a pointer")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic on GetTo: %s, stack: %s", r, prettyStack(string(debug.Stack())))
			}
		}()
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return
			}
			dv.Elem().Set(v.Elem())
		} else {
			dv.Elem().Set(v)
		}
	}()
	return err
}

func (r *Reflector) getValue(rv reflect.Value, name string) (reflect.Value, error) {
	names := strings.Split(name, ".")
	fv := rv.FieldByName(names[0])
	if !fv.IsValid() {
		return rv, errors.New("invalidField: " + name)
	}

	if len(names) > 1 {
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				fv = reflect.New(fv.Type().Elem())
			}
			return r.getValue(fv.Elem(), strings.Join(names[1:], "."))
		} else {
			return r.getValue(fv, strings.Join(names[1:], "."))
		}
	}

	return fv, nil
}

func (r *Reflector) Set(name string, value interface{}) *Reflector {
	return r.setValue(r.v, name, value)
}

func (r *Reflector) setValue(rv reflect.Value, name string, value interface{}) *Reflector {
	if r.err != nil {
		return r
	}

	func() {
		defer func() {
			if rec := recover(); rec != nil {
				// surpress error
			}
		}()

		names := strings.Split(name, ".")
		fieldName := names[0]
		v := rv.FieldByName(fieldName)
		if len(names) > 1 {
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					newPtr := reflect.New(v.Type().Elem())
					v.Set(newPtr)
				}
				r.setValue(v.Elem(), strings.Join(names[1:], "."), value)
			} else {
				r.setValue(v, strings.Join(names[1:], "."), value)
			}
			return
		}

		v.Set(reflect.ValueOf(value))
	}()

	return r
}

func (r *Reflector) Flush() error {
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

func (r *Reflector) FieldNames(tag string) ([]string, error) {
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

func prettyStack(stack string) string {
	stacks := strings.Split(stack, "\n")
	filtereds := []string{}
	for _, s := range stacks {
		if strings.Contains(s, ".go:") {
			filtereds = append(filtereds, strings.Trim(s, " \t\n"))
		}
	}
	return strings.Join(filtereds, "\n")
}
