package polluter

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"
	"github.com/romanyx/jwalk"
	yaml "gopkg.in/yaml.v3"
)

type CustomMap = map[any]any

type yamlParser struct{}

func (p yamlParser) parse(r io.Reader) (jwalk.ObjectWalker, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read from input")
	}

	j, err := yamlToJSON(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed convert to json")
	}

	i, err := jwalk.Parse(j)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse")
	}

	obj, ok := i.(jwalk.ObjectWalker)
	if !ok {
		return nil, errors.New("unexpected format")
	}

	return obj, nil
}

func yamlToJSON(data []byte) ([]byte, error) {
	mapSlice := CustomMap{}

	err := yaml.Unmarshal(data, &mapSlice)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	buf := new(bytes.Buffer)
	handleMapSlice(mapSlice, buf)

	return buf.Bytes(), nil
}

func handleMapSlice(mapSlice CustomMap, buf *bytes.Buffer) {
	buf.WriteString("{")
	first := true
	indent := ""
	for k, v := range mapSlice {
		buf.WriteString(indent + "\"" + k.(string) + "\"" + ":")
		switch v := v.(type) {
		case CustomMap:
			handleMapSlice(v, buf)
		case []interface{}:
			buf.WriteString("[")
			first := true
			indent := ""
			for _, i := range v {
				switch v := i.(type) {
				case CustomMap:
					buf.WriteString(indent)
					handleMapSlice(v, buf)
				default:
					buf.WriteString(indent + formatValue(v))
				}
				if first {
					first = false
					indent = ","
				}
			}
			buf.WriteString("]")
		default:
			buf.WriteString(formatValue(v))
		}
		if first {
			first = false
			indent = ","
		}

	}
	buf.WriteString("}")
}

func formatValue(typedYAMLObj interface{}) string {
	switch typedVal := typedYAMLObj.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", typedVal)
	case int:
		return strconv.FormatInt(int64(typedVal), 10)
	case int64:
		return strconv.FormatInt(typedVal, 10)
	case float64:
		return strconv.FormatFloat(typedVal, 'g', -1, 32)
	case uint64:
		return strconv.FormatUint(typedVal, 10)
	case bool:
		if typedVal {
			return "true"
		}
		return "false"
	default:
		return "null"
	}

	return ""
}
