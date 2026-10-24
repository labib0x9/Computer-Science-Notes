package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type MyInt int

type Store struct {
	Key    string `map:"key"`
	Value  MyInt  `map:"value"`
	Whoami string `map:"wh0ami"`
}

// Decode decodes a map into a struct.
// Only supports string to string and string to MyInt conversions for simplicity.

func P(a ...any) {
	fmt.Println(a...)
}

func Decode(mp map[string]string, out interface{}) error {
	v := reflect.ValueOf(out)

	// check if out is a pointer
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("binding object must be a pointer, got %s", v.Kind())
	}

	v = v.Elem()

	// check if out is a struct
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("binding object must be a struct, got %s", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}

		tag := v.Type().Field(i).Tag.Get("map")
		if len(tag) == 0 {
			continue
		}

		value, ok := mp[tag]
		if !ok {
			continue
		}

		switch f.Type().Kind() {
		case reflect.String:
			f.SetString(value)
		case reflect.Int, reflect.Int16:
			n, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error converting feild %s (%s)", v.Type().Field(i).Name, f.Type().Kind())
			}
			f.Set(reflect.ValueOf(n).Convert(f.Type()))
		default:
			return fmt.Errorf("unsupported data type for field %s (%s)", v.Type().Field(i).Name, f.Type().Kind())
		}
	}

	return nil
}

func main() {

	mp := map[string]string{
		"key":   "Count",
		"value": "0",
	}

	var s Store
	if err := Decode(mp, &s); err != nil {
		panic(err)
	}

	fmt.Println(s)
}
