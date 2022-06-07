package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return v.Field + " has " + v.Err.Error()
}

func (v ValidationErrors) Error() string {
	str := ""
	for _, e := range v {
		if str != "" {
			str += ", "
		}
		str += e.Error()
	}
	return str
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	valType := val.Type()
	for f := 0; f < valType.NumField(); f++ {
		fieldType := valType.Field(f)
		tag := fieldType.Tag.Get("validate")

		if tag == "" {
			continue
		}
		for _, rule := range strings.Split(tag, "|") {
			nameAndValue := strings.Split(rule, ":")
			if len(nameAndValue) != 2 {
				return errors.New("bad rule")
			}
			err := validateValue(val.Field(f), nameAndValue[0], nameAndValue[1])
			if err != nil {
				fmt.Println(fieldType.Name, err)
			}

		}

		//
	}

	return nil
}

func validateValue(value reflect.Value, ruleName, ruleValue string) error {

	//fmt.Println(value.Type(), value.String(), ruleName, ruleValue)

	switch value.Kind() {
	case reflect.String:
		stringValue := value.String()

		switch ruleName {
		case "len":
			length, _ := strconv.Atoi(ruleValue)
			if len(stringValue) != length {
				return errors.New("wrong length")
			}

		case "regexp":
		case "in":
			ok := false
			for _, allowedValue := range strings.Split(ruleValue, ",") {
				if stringValue == allowedValue {
					ok = true
					break
				}
			}
			if !ok {
				return errors.New("illegal value")
			}
		}

	case reflect.Int:
		intValue := int(value.Int())

		switch ruleName {
		case "min":
			min, _ := strconv.Atoi(ruleValue)
			if intValue < min {
				return errors.New("too small")
			}
		case "max":
			max, _ := strconv.Atoi(ruleValue)
			if intValue > max {
				return errors.New("too big")
			}
		case "in":
			ok := false
			for _, allowedValue := range strings.Split(ruleValue, ",") {
				av, _ := strconv.Atoi(allowedValue)
				if intValue == av {
					ok = true
					break
				}
			}
			if !ok {
				return errors.New("illegal value")
			}
		}

	case reflect.Slice:

		for i := 0; i < value.Len(); i++ {
			//item := value.Index(i)
			//err := validateValue(value reflect.Value, ruleName, ruleValue string)

		}

	}

	return nil
}
