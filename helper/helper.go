package helper

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/types"
)

func isPrimaryType(dataType consts.DataType) bool {
	switch dataType {
	case consts.Bool,
		consts.Int, consts.Int8, consts.Int16, consts.Int32, consts.Int64,
		consts.Uint, consts.Uint8, consts.Uint16, consts.Uint32, consts.Uint64,
		consts.Float, consts.Double,
		consts.String:
		return true
	}

	return false
}

func ParseDataType(typeStr string) (consts.DataType, []string) {
	dataType := consts.DataType(typeStr)
	if isPrimaryType(dataType) {
		return dataType, nil
	}

	// array2:xxx
	if strings.HasPrefix(typeStr, string(consts.Array2)) {
		_, subType, _ := strings.Cut(typeStr, ":")
		return consts.Array2, []string{subType}
	}

	// array:xxx
	if strings.HasPrefix(typeStr, string(consts.Array)) {
		_, subType, _ := strings.Cut(typeStr, ":")
		return consts.Array, []string{subType}
	}

	// map:xxx,xxx
	if strings.HasPrefix(typeStr, string(consts.Map)) {
		_, subType, _ := strings.Cut(typeStr, ":")
		keyType, valueType, _ := strings.Cut(subType, ",")
		return consts.Map, []string{keyType, valueType}
	}

	if strings.HasPrefix(typeStr, string(consts.Enum)) {
		_, subType, _ := strings.Cut(typeStr, ":")
		return consts.Enum, []string{subType}
	}

	return consts.Unknown, nil
}

func ParseValue(rawValue string, dataType consts.DataType, params ...string) (any, error) {
	rawValue = originOrDefault(dataType, rawValue)

	switch dataType {
	case consts.Int, consts.Int64:
		return strconv.ParseInt(rawValue, 10, 64)
	case consts.Int32:
		return strconv.ParseInt(rawValue, 10, 32)
	case consts.Int16:
		return strconv.ParseInt(rawValue, 10, 16)
	case consts.Int8:
		return strconv.ParseInt(rawValue, 10, 8)
	case consts.Uint, consts.Uint64:
		return strconv.ParseUint(rawValue, 10, 64)
	case consts.Uint32:
		return strconv.ParseUint(rawValue, 10, 32)
	case consts.Uint16:
		return strconv.ParseUint(rawValue, 10, 16)
	case consts.Uint8:
		return strconv.ParseUint(rawValue, 10, 8)
	case consts.Float:
		return strconv.ParseFloat(rawValue, 32)
	case consts.Double:
		return strconv.ParseFloat(rawValue, 64)
	case consts.String:
		rawValue = strings.TrimPrefix(rawValue, "\"")
		rawValue = strings.TrimSuffix(rawValue, "\"")
		return rawValue, nil
	case consts.Bool:
		return strconv.ParseBool(rawValue)
	case consts.Array:
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
	case consts.Array2:
		arr := []any{}
		trimValue := strings.TrimPrefix(rawValue, "[")
		trimValue = strings.TrimSuffix(trimValue, "]")

		var (
			value  any
			remain = trimValue
			err    error
		)
		for len(remain) > 0 {
			value, remain, err = readValue(remain, consts.Array, ",", params...)
			if err != nil {
				return nil, err
			}

			arr = append(arr, value)
		}

		return arr, nil
	case consts.Map:
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

	case consts.Enum:
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

func readValue(rawValue string, dataType consts.DataType, valSep string, params ...string) (any, string, error) {
	if isPrimaryType(dataType) {
		beParse, remain, _ := strings.Cut(rawValue, valSep)
		value, err := ParseValue(beParse, dataType)
		return value, remain, err
	}

	if dataType == consts.Array || dataType == consts.Array2 {
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

	if dataType == consts.Map {
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

		value, err := ParseValue(beParse, consts.Map, params...)
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

func DefaultStringValue(dataType consts.DataType) string {
	switch dataType {
	case consts.Bool:
		return "false"
	case consts.Int, consts.Int8, consts.Int16, consts.Int32, consts.Int64,
		consts.Uint, consts.Uint8, consts.Uint16, consts.Uint32, consts.Uint64:
		return "0"
	case consts.Float, consts.Double:
		return "0.0"
	case consts.String:
		return "\"\""
	case consts.Array, consts.Array2:
		return "[]"
	case consts.Map:
		return "{}"
	case consts.Enum:
		return "0"
	}

	return ""
}

func originOrDefault(dataType consts.DataType, origin string) string {
	if len(origin) > 0 {
		return origin
	}

	return DefaultStringValue(dataType)
}
