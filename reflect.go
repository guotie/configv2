package config

import (
	"fmt"
	"reflect"
	"strings"
)

func scan(val interface{}, rv reflect.Value) (err error) {
	rv = indirect(rv)

	switch rv.Kind() {
	case reflect.Bool:
		v, ok := val.(bool)
		if !ok {
			return fmt.Errorf("Cannot convert val(%v) from type %s to bool.",
				val, reflect.TypeOf(val).String())
		}
		rv.SetBool(v)

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("Cannot convert val(%v) from type %s to float64.",
				val, reflect.TypeOf(val).String())
		}
		rv.SetInt(int64(v))

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("Cannot convert val(%v) from type %s to float64.",
				val, reflect.TypeOf(val).String())
		}
		rv.SetUint(uint64(v))

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("Cannot convert val(%v) from type %s to float64.",
				val, reflect.TypeOf(val).String())
		}
		rv.SetFloat(v)

	case reflect.String:
		v, ok := val.(string)
		if !ok {
			return fmt.Errorf("Cannot convert val(%v) from type %s to string.",
				val, reflect.TypeOf(val).String())
		}
		rv.SetString(v)
	// decode the primary types here

	case reflect.Struct:
		err = objectStruct(val, rv)
	case reflect.Map:
		err = objectMap(val, rv)
	// object

	case reflect.Array:
		fallthrough
	case reflect.Slice:
		// array
		err = array(val, rv)

	case reflect.Interface:
		// how to do with interface ?
		rv.Set(reflect.ValueOf(val))

	case 0:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.UnsafePointer:
		panic("Not support type " + rv.Kind().String())

	// should never arrive here
	case reflect.Ptr:
		panic("Should not be reflect.Ptr")

	// should never arrive here, too
	default:
		panic("Unknown reflect type, Should never arrive here")
	}

	return err
}

// code from json decode
//
func indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}

	return v
}

func objectMap(val interface{}, rv reflect.Value) (err error) {
	var subv reflect.Value

	mval, ok := val.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Cannot convert value(%v) to map[string]interface{}, check data!", val)
	}
	if rv.IsNil() {
		//rv = reflect.MakeMap(rv.Type())
		rv.Set(reflect.MakeMap(rv.Type()))
	}

	for key, valfd := range mval {
		elemType := rv.Type().Elem()
		subv = reflect.New(elemType).Elem()

		err = scan(valfd, subv)
		if err != nil {
			return
		}
		rv.SetMapIndex(reflect.ValueOf(key),
			reflect.ValueOf(subv.Interface()))
	}

	return
}

func objectStruct(val interface{}, rv reflect.Value) (err error) {
	var valElem interface{}

	mval, ok := val.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Cannot convert value(%v) to map[string]interface{}", val)
	}

	for i := 0; i < rv.NumField(); i++ {
		fdval := rv.Field(i)
		fdtyp := rv.Type().Field(i)

		tag := fdtyp.Tag.Get("json")
		if tag == "-" {
			continue
		}
		if tag != "" {
			valElem = mval[tag]
		} else {
			valElem = mval[fdtyp.Name]
			if valElem == nil {
				valElem = mval[strings.ToLower(fdtyp.Name)]
			}
		}
		if valElem == nil {
			continue
		}

		err = scan(valElem, fdval)
		if err != nil {
			return
		}
	}

	return
}

func array(val interface{}, rv reflect.Value) (err error) {
	sval, ok := val.([]interface{})
	if !ok {
		return fmt.Errorf("Cannot convert val(%v) to []interface{}", val)
	}

	if rv.Kind() == reflect.Slice && rv.IsNil() {
		rv.Set(reflect.New(rv.Type()).Elem())
	}
	for i, valfd := range sval {
		if rv.Kind() == reflect.Slice {
			if i >= rv.Cap() {
				newcap := rv.Cap() + rv.Cap()/2
				if newcap < 4 {
					newcap = 4
				}
				newv := reflect.MakeSlice(rv.Type(), rv.Len(), newcap)
				reflect.Copy(newv, rv)
				rv.Set(newv)
			}
			if i >= rv.Len() {
				rv.SetLen(i + 1)
			}
		} else {
			if i >= rv.Len() {
				fmt.Printf("array length oversized: len: %d, i: %d\n", rv.Len(), i)
				return
			}
		}

		subv := reflect.New(rv.Type().Elem())
		err = scan(valfd, subv)
		if err != nil {
			return
		}
		rv.Index(i).Set(subv.Elem())
	}

	return
}
