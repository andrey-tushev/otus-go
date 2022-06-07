package hw09structvalidator

import (
	"errors"
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

func (v *ValidationErrors) add(field string, err error) {
	*v = append(*v, ValidationError{
		Field: field,
		Err:   err,
	})
}

func (v ValidationError) Error() string {
	return v.Field + " " + v.Err.Error()
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
	validationErrors := ValidationErrors{}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	valType := val.Type()
	for f := 0; f < valType.NumField(); f++ {
		structField := valType.Field(f)
		tag := structField.Tag.Get("validate")

		if tag == "" {
			continue
		}
		for _, rule := range strings.Split(tag, "|") {
			nameAndValue := strings.Split(rule, ":")
			if len(nameAndValue) != 2 {
				return errors.New("bad rule")
			}
			err := validateByRule(val.Field(f), nameAndValue[0], nameAndValue[1])
			if err != nil {
				validationErrors.add(structField.Name, err)
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateByRule(value reflect.Value, ruleName, ruleValue string) error {
	switch value.Kind() {
	case reflect.String:
		stringValue := value.String()

		switch ruleName {
		case "len":
			length, _ := strconv.Atoi(ruleValue)
			if len(stringValue) != length {
				return errors.New("has wrong length")
			}

		case "regexp":
			re, _ := regexp.Compile(ruleValue)
			if !re.Match([]byte(stringValue)) {
				return errors.New("has bad format")
			}

		case "in":
			ok := false
			for _, allowedValue := range strings.Split(ruleValue, ",") {
				if stringValue == allowedValue {
					ok = true
					break
				}
			}
			if !ok {
				return errors.New("contains illegal value")
			}
		}

	case reflect.Int:
		intValue := int(value.Int())

		switch ruleName {
		case "min":
			min, _ := strconv.Atoi(ruleValue)
			if intValue < min {
				return errors.New("is too small")
			}
		case "max":
			max, _ := strconv.Atoi(ruleValue)
			if intValue > max {
				return errors.New("is too big")
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
				return errors.New("contains illegal value")
			}
		}

	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)

			err := validateByRule(item, ruleName, ruleValue)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
