package reflector

import (
	"errors"
	"reflect"

	"github.com/eaciit/toolkit"
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
		destBuffer := reflect.New(dest.Type().Elem()).Interface()
		if e := toolkit.Serde(source.Interface(), destBuffer, ""); e != nil {
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
		destBuffer := reflect.New(destItemTypeVal).Interface()
		if e := toolkit.Serde(data.Interface(), destBuffer, ""); e != nil {
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

	//fmt.Println("Data:", toolkit.JsonString(newDest.Interface()))
	if copyIsNeeded {
		dest.Elem().Set(newDest)
	}

	return nil
}
