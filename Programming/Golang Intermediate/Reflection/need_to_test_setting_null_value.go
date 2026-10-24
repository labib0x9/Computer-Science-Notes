package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type Request struct {
	Name    *string `json:"name"`
	Status  *string `json:"status"`
	Persist *bool   `json:"persist"`
}

func encode(b []byte) io.Reader {
	return bytes.NewReader(b)
}

func decode(r io.Reader, req any) {
	var body map[string]json.RawMessage
	if err := json.NewDecoder(r).Decode(&body); err != nil {
		panic(err)
	}

	rv := reflect.ValueOf(req).Elem()
	fmt.Println("Field Number:", rv.NumField())

	for i := range rv.NumField() {
		f := rv.Type().Field(i).Tag.Get("json")
		fmt.Println("JSON FILED =", f)

		// if value, exits := body[f]; exits {
		// 	fmt.Println(f, "=", string(value))

		// }

		if value, exists := body[f]; exists {
			fmt.Println(f, "=", string(value))
			if field := rv.Field(i); field.CanSet() {
				newElem := reflect.New(field.Type().Elem()) // field.Type() is *string/*bool -> Elem() gives string/bool
				if err := json.Unmarshal(value, newElem.Interface()); err != nil {
					panic(err)
				}
				field.Set(newElem)
			}
		}
	}
}

func main() {

	var ra Request
	aj := `{}`
	a := encode([]byte(aj))
	// {}
	decode(a, &ra)
	fmt.Println()

	var rb Request
	bj := `{"name": "Labib"}`
	b := encode([]byte(bj))
	// {"name":"Labib"}
	decode(b, &rb)
	fmt.Println()

	var rc Request
	cj := `{"name": "Labib", "status": null}`
	c := encode([]byte(cj))
	// {"name":"Labib"}
	decode(c, &rc)
	fmt.Println()

	fmt.Println("Name", *rc.Name)
	fmt.Println("Status", *rc.Status)
}
