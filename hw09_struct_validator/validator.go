package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const validateTag = "validate"

// Общие ошибки
var (
	ErrNotAStruct = errors.New("not a struct")
	ErrBadRule    = errors.New("bad rule")
)

// Ошибки валидации
var (
	ErrWrongLength  = errors.New("has wrong length")
	ErrBadFormat    = errors.New("has bad format")
	ErrIllegalValue = errors.New("contains illegal value")
	ErrTooSmall     = errors.New("is too small")
	ErrTooBig       = errors.New("is too big")
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
		return ErrNotAStruct
	}

	valType := val.Type()
	for f := 0; f < valType.NumField(); f++ {
		structField := valType.Field(f)
		tag := structField.Tag.Get(validateTag)

		// Если тек валидации не пустой, то разбираем его и запускаем валидацию для каждого из правил
		if tag == "" {
			continue
		}
		for _, rule := range strings.Split(tag, "|") {
			nameAndValue := strings.Split(rule, ":")
			if len(nameAndValue) != 2 {
				return ErrBadRule
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

// validateByRule валидирует одно значение по одному правилу
func validateByRule(value reflect.Value, ruleName, ruleValue string) error {
	switch value.Kind() {
	// Валидация для строк
	case reflect.String:
		stringValue := value.String()

		switch ruleName {
		case "len":
			length, _ := strconv.Atoi(ruleValue)
			if len(stringValue) != length {
				return ErrWrongLength
			}

		case "regexp":
			re, _ := regexp.Compile(ruleValue)
			if !re.Match([]byte(stringValue)) {
				return ErrBadFormat
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
				return ErrIllegalValue
			}
		}

	// Валидация для интов
	case reflect.Int:
		intValue := int(value.Int())

		switch ruleName {
		case "min":
			min, _ := strconv.Atoi(ruleValue)
			if intValue < min {
				return ErrTooSmall
			}
		case "max":
			max, _ := strconv.Atoi(ruleValue)
			if intValue > max {
				return ErrTooBig
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
				return ErrIllegalValue
			}
		}

	// Валидация для слайсов
	case reflect.Slice:

		// Валидируем каждый элемент слайса
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
