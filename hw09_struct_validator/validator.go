package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("%s: %s; ", err.Field, err.Err.Error()))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if !field.IsExported() {
			continue
		}

		value := val.Field(i)
		tag := field.Tag.Get("validate")

		if tag == "" {
			continue
		}

		tags := strings.Split(tag, "|")
		for _, t := range tags {
			tagParts := strings.SplitN(t, ":", 2)
			tagKey := tagParts[0]
			tagValue := ""
			if len(tagParts) > 1 {
				tagValue = tagParts[1]
			}

			var err error
			switch tagKey {
			case "len":
				err = validateLength(value, tagValue)
			case "regexp":
				err = validateRegexp(value, tagValue)
			case "min":
				err = validateMin(value, tagValue)
			case "max":
				err = validateMax(value, tagValue)
			case "in":
				err = validateIn(value, tagValue)
			}

			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateLength(value reflect.Value, expectedLength string) error {
	l, err := strconv.Atoi(expectedLength)
	if err != nil {
		return err
	}

	switch value.Kind() {
	case reflect.String:
		if len(value.String()) != l {
			return fmt.Errorf("length of value is %d but expected %d", len(value.String()), l)
		}
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.String {
			for i := 0; i < value.Len(); i++ {
				if len(value.Index(i).String()) != l {
					return fmt.Errorf("length of value is %d but expected %d", len(value.Index(i).String()), l)
				}
			}
		} else if value.Len() != l {
			return fmt.Errorf("length of value is %d but expected %d", value.Len(), l)
		}
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array,
		reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Struct,
		reflect.UnsafePointer:
		return fmt.Errorf("unsupported kind: %s", value.Kind())
	default:
		return fmt.Errorf("unexpected kind: %s", value.Kind())
	}

	return nil
}

func validateRegexp(value reflect.Value, pattern string) error {
	if value.Kind() != reflect.String || !value.IsValid() || value.String() == "" {
		return nil
	}

	re := regexp.MustCompile(pattern)
	if !re.MatchString(value.String()) {
		return fmt.Errorf("value %s does not match pattern %s", value.String(), pattern)
	}
	return nil
}

func validateMin(value reflect.Value, min string) error {
	if value.Kind() != reflect.Int {
		return fmt.Errorf("min validation expects an int but got %s", value.Kind())
	}
	minInt, err := strconv.Atoi(min)
	if err != nil {
		return err
	}
	if int(value.Int()) < minInt {
		return fmt.Errorf("value %d is less than minimum %d", value.Int(), minInt)
	}
	return nil
}

func validateMax(value reflect.Value, max string) error {
	if value.Kind() != reflect.Int {
		return fmt.Errorf("max validation expects an int but got %s", value.Kind())
	}
	maxInt, err := strconv.Atoi(max)
	if err != nil {
		return err
	}
	if int(value.Int()) > maxInt {
		return fmt.Errorf("value %d is greater than maximum %d", value.Int(), maxInt)
	}
	return nil
}

func validateIn(value reflect.Value, allowedValues string) error {
	allowed := strings.Split(allowedValues, ",")
	val := fmt.Sprintf("%v", value.Interface())
	for _, v := range allowed {
		if v == val {
			return nil
		}
	}
	return fmt.Errorf("value %s is not in the allowed list [%s]", val, strings.Join(allowed, ","))
}
