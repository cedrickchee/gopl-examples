// Package params provides a reflection-based parser for URL parameters.
package params

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Unpack populates the fields of the struct pointed to by ptr
// from the HTTP request parameters in req.
func Unpack(req *http.Request, ptr interface{}) error {
	// The `Unpack` function below does three things. First, it calls
	// `req.ParseForm()` to parse the request. Thereafter, `req.Form` contains
	// all the parameters, regardless of whether the HTTP client used the GET or
	// the POST request method.
	//
	// Next, `Unpack` builds a mapping from the _effective_ name of each field
	// to the variable for that field. The effective name may differ from the
	// actual name if the field has a tag. The `Field` method of `reflect.Type`
	// returns a `reflect.StructField` that provides information about the type
	// of each field such as its name, type, and optional tag. The `Tag` field
	// is a `reflect.StructTag`, which is a string type that provides a `Get`
	// method to parse and extract the substring for a particular key, such as
	// `http:"..."` in this case.
	//
	// Finally, `Unpack` iterates over the name/value pairs of the HTTP
	// parameters and updates the corresponding struct fields. Recall that the
	// same parameter name may appear more than once. If this happens, and the
	// field is a slice, then all the values of that parameter are accumulated
	// into the slice. Otherwise, the field is repeatedly overwritten so that
	// only the last value has any effect.

	if err := req.ParseForm(); err != nil {
		return err
	}

	// Build map of fields keyed by effective name.
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	// Update struct field for each parameter in the request.
	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue // ignore unrecognized HTTP parameters
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

// populate takes care of setting a single field `v` (or a single element of a
// slice field) from a parameter value. For now, it supports only strings,
// signed integers, and booleans. Supporting other types is left as an
// improvement.
func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
