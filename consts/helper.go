package consts

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/lgynico/alpaca/types"
)

func isPrimaryType(dataType DataType) bool {
	switch dataType {
	case Bool,
		Int, Int8, Int16, Int32, Int64,
		Uint, Uint8, Uint16, Uint32, Uint64,
		Float, Double,
		String:
		return true
	}

	return false
}

func ParseDataType(typeStr string) (DataType, []string) {
	dataType := DataType(typeStr)
	if isPrimaryType(dataType) {
		return dataType, nil
	}

	// array2:xxx
	if strings.HasPrefix(typeStr, string(Array2)) {
		_, subType, _ := strings.Cut(typeStr, ":")
		return Array2, []string{subType}
	}

	// array:xxx
	if strings.HasPrefix(typeStr, string(Array)) {
		_, subType, _ := strings.Cut(typeStr, ":")
		return Array, []string{subType}
	}

	// map:xxx,xxx
	if strings.HasPrefix(typeStr, string(Map)) {
		_, subType, _ := strings.Cut(typeStr, ":")
		keyType, valueType, _ := strings.Cut(subType, ",")
		return Map, []string{keyType, valueType}
	}

	if strings.HasPrefix(typeStr, string(Enum)) {
		_, subType, _ := strings.Cut(typeStr, ":")
		return Enum, []string{subType}
	}

	return Unknown, nil
}

func ParseValue(rawValue string, dataType DataType, params ...string) (any, error) {
	switch dataType {
	case Int, Int64:
		return strconv.ParseInt(rawValue, 10, 64)
	case Int32:
		return strconv.ParseInt(rawValue, 10, 32)
	case Int16:
		return strconv.ParseInt(rawValue, 10, 16)
	case Int8:
		return strconv.ParseInt(rawValue, 10, 8)
	case Uint, Uint64:
		return strconv.ParseUint(rawValue, 10, 64)
	case Uint32:
		return strconv.ParseUint(rawValue, 10, 32)
	case Uint16:
		return strconv.ParseUint(rawValue, 10, 16)
	case Uint8:
		return strconv.ParseUint(rawValue, 10, 8)
	case Float:
		return strconv.ParseFloat(rawValue, 32)
	case Double:
		return strconv.ParseFloat(rawValue, 64)
	case String:
		rawValue = strings.TrimPrefix(rawValue, "\"")
		rawValue = strings.TrimSuffix(rawValue, "\"")
		return rawValue, nil
	case Bool:
		return strconv.ParseBool(rawValue)
	case Array:
		arr := []any{}
		dataType, params := ParseDataType(params[0])
		trimValue := strings.TrimPrefix(rawValue, "[")
		trimValue = strings.TrimSuffix(trimValue, "]")

		var (
			value  any
			remain = trimValue
			err    error
		)
		for len(remain) > 0 {
			value, remain, err = readValue(remain, dataType, ",", params...)
			if err != nil {
				return nil, err
			}

			arr = append(arr, value)
		}

		return arr, nil
	case Array2:
		arr := []any{}
		trimValue := strings.TrimPrefix(rawValue, "[")
		trimValue = strings.TrimSuffix(trimValue, "]")

		var (
			value  any
			remain = trimValue
			err    error
		)
		for len(remain) > 0 {
			value, remain, err = readValue(remain, Array, ",", params...)
			if err != nil {
				return nil, err
			}

			arr = append(arr, value)
		}

		return arr, nil
	case Map:
		m := map[any]any{}
		trimValue := strings.TrimPrefix(rawValue, "{")
		trimValue = strings.TrimSuffix(trimValue, "}")
		keyType, keyParams := ParseDataType(params[0])
		valueType, valueParams := ParseDataType(params[1])

		var (
			key    any
			value  any
			remain = trimValue
			err    error
		)

		for len(remain) > 0 {
			key, remain, err = readValue(remain, keyType, ":", keyParams...)
			if err != nil {
				return nil, err
			}
			value, remain, err = readValue(remain, valueType, ",", valueParams...)
			if err != nil {
				return nil, err
			}

			m[key] = value
		}

		return m, nil

	case Enum:
		enumName := params[0]
		enum, ok := types.GetEnum(enumName)
		if !ok {
			return nil, fmt.Errorf("enum %s not found", enumName)
		}

		for _, node := range enum.Nodes {
			if node.Key == rawValue {
				return node.Value, nil
			}
		}

		return nil, fmt.Errorf("unknow enum value: %s_%s", enumName, rawValue)
	}

	return nil, errors.New("unknown type")
}

func readValue(rawValue string, dataType DataType, valSep string, params ...string) (any, string, error) {
	if isPrimaryType(dataType) {
		beParse, remain, _ := strings.Cut(rawValue, valSep)
		value, err := ParseValue(beParse, dataType)
		return value, remain, err
	}

	if dataType == Array || dataType == Array2 {
		beParse := ""
		stack := arraystack.New()

		i := 0
		for ; i < len(rawValue); i++ {
			c := rawValue[i]
			beParse += string(c)

			if c == '[' {
				stack.Push(c)
			} else if c == ']' {
				stack.Pop()
			}

			if stack.Empty() {
				break
			}
		}

		if !stack.Empty() || i+1 < len(rawValue) && rawValue[i+1] != ',' {
			return nil, rawValue, fmt.Errorf("unparsable string for type %s: %s", dataType, rawValue)
		}

		value, err := ParseValue(beParse, dataType, params...)
		if err != nil {
			return nil, rawValue, err
		}
		remain := strings.TrimPrefix(rawValue, beParse)
		remain = strings.TrimPrefix(remain, ",")

		return value, remain, nil

	}

	if dataType == Map {
		beParse := ""
		stack := arraystack.New()
		i := 0
		for ; i < len(rawValue); i++ {
			c := rawValue[i]
			beParse += string(c)

			if c == '{' {
				stack.Push(c)
			} else if c == '}' {
				stack.Pop()
			}

			if stack.Empty() {
				break
			}
		}

		if !stack.Empty() || i+1 < len(rawValue) && rawValue[i+1] != ',' {
			return nil, rawValue, fmt.Errorf("unparsable string for type %s: %s", dataType, rawValue)
		}

		value, err := ParseValue(beParse, Map, params...)
		if err != nil {
			return nil, rawValue, err
		}
		remain := strings.TrimPrefix(rawValue, beParse)
		remain = strings.TrimPrefix(remain, ",")
		return value, remain, nil
	}

	return nil, rawValue, errors.New("unknown array value type")
}

func ToString(value reflect.Value) string {
	t := value.Type()

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'f', 4, 32)
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', 4, 64)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.String:
		return addQuote(value.String())
	case reflect.Array, reflect.Slice:
		str := "["
		for i := 0; i < value.Len(); i++ {
			if i > 0 {
				str += ","
			}
			str += ToString(value.Index(i))
		}
		str += "]"
		return str
	case reflect.Map:
		str := "{"
		for i, key := range value.MapKeys() {
			if i > 0 {
				str += ","
			}

			keyStr := addQuote(ToString(key))
			str += fmt.Sprintf("%s:%s", keyStr, ToString(value.MapIndex(key)))
		}
		str += "}"
		return str
	case reflect.Interface:
		return ToString(value.Elem())
	}

	return ""
}

func UnderlineToCamelCase(value string, capitalizeFirst bool) string {
	ss := strings.Split(value, "_")
	for i, v := range ss {
		if i == 0 && !capitalizeFirst {
			ss[i] = strings.ToLower(string(v[0])) + v[1:]
		} else {
			ss[i] = strings.ToUpper(string(v[0])) + v[1:]
		}
	}
	return strings.Join(ss, "")
}

func addQuote(value string) string {
	if len(value) == 0 {
		return `""`
	}

	if !strings.HasPrefix(value, `"`) {
		value = `"` + value
	}

	if !strings.HasSuffix(value, `"`) {
		value += `"`
	}

	return value
}
