package reflector

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ariefdarmawan/serde"
)

func AssignValue(source reflect.Value, dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr {
		return errors.New("destination variable shoud be a pointer")
	}

	if source.Kind() == reflect.Ptr {
		source = source.Elem()
	}

	if dest.Elem().Type() == source.Type() {
		dest.Elem().Set(source)
	} else {
		var destBuffer interface{}
		destElemKind := dest.Elem().Kind()
		switch destElemKind {
		case reflect.Struct:
			destBuffer = reflect.New(dest.Type().Elem()).Interface()

		case reflect.Map:
			destBuffer = serde.CreatePtrFromType(dest.Type().Elem()).Interface()

		case reflect.Slice:
			rSlice := reflect.MakeSlice(dest.Type().Elem(), 0, 0)
			rBuffer := reflect.New(dest.Type().Elem())
			rBuffer.Elem().Set(rSlice)
			destBuffer = rBuffer.Interface()
		}
		if e := serde.Serde(source.Interface(), destBuffer); e != nil {
			return errors.New("fail to serializing source. " + e.Error())
		}
		//fmt.Println("dbf type is ", reflect.ValueOf(destBuffer).Type().String())
		dest.Elem().Set(reflect.ValueOf(destBuffer).Elem())
	}
	return nil
}

func AssignSliceItem(data reflect.Value, index int, dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr {
		return errors.New("destination variable shoud be a pointer")
	}
	if dest.Type().Elem().Kind() != reflect.Slice {
		return errors.New("source variable shoud be a slice or pointer of slice")
	}
	destItemType := dest.Type().Elem().Elem()
	destItemTypeVal := destItemType
	destItemIsPtr := destItemType.Kind() == reflect.Ptr
	if destItemIsPtr {
		destItemTypeVal = destItemTypeVal.Elem()
	}

	if data.Kind() == reflect.Ptr {
		data = data.Elem()
	}

	var item reflect.Value
	if destItemTypeVal == data.Type() {
		item = data
	} else {
		destBuffer := serde.CreatePtrFromType(destItemTypeVal).Interface()
		if e := serde.Serde(data.Interface(), destBuffer); e != nil {
			return errors.New("fail to serializing source. " + e.Error())
		}
		item = reflect.ValueOf(destBuffer).Elem()
	}

	var (
		newDest      = dest.Elem()
		copyIsNeeded = false
	)

	curLen := dest.Elem().Len()
	if (index + 1) > curLen {
		newDest = reflect.MakeSlice(dest.Elem().Type(), index+1, index+1)
		copyIsNeeded = true
		if curLen > 0 {
			reflect.Copy(newDest, dest.Elem())
		}
	}

	if destItemIsPtr {
		if item.Kind() != reflect.Ptr {
			ptrItem := reflect.New(item.Type())
			ptrItem.Elem().Set(item)
			newDest.Index(index).Set(ptrItem)
		} else {
			newDest.Index(index).Set(item)
		}
	} else {
		newDest.Index(index).Set(item)
	}

	//fmt.Println("Data:", codekit.JsonString(newDest.Interface()))
	if copyIsNeeded {
		dest.Elem().Set(newDest)
	}

	return nil
}

func CreateFromPtr[T any](ptr T, copyValue bool) (T, error) {
	vt := reflect.TypeOf(ptr)
	if vt.Kind() != reflect.Ptr {
		return ptr, fmt.Errorf("parameter should be ptr")
	}

	nv := reflect.New(vt.Elem())
	if copyValue {
		nv.Elem().Set(reflect.ValueOf(ptr).Elem())
	}
	return nv.Interface().(T), nil
}
