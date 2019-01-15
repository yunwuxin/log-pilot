package pilot

import "fmt"

type FormatConverter func(info *LogInfoNode) (map[string]string, error)

var converters = make(map[string]FormatConverter)

var defaultProperties = []string{"types", "time_key", "null_value_pattern", "null_empty_string", "time_format", "time_type"}

func Register(format string, converter FormatConverter) {
	converters[format] = converter
}

func Convert(info *LogInfoNode) (map[string]string, error) {
	converter := converters[info.value]
	if converter == nil {
		return nil, fmt.Errorf("unsupported log format: %s", info.value)
	}
	return converter(info)
}

type SimpleConverter struct {
	properties map[string]bool
}

func init() {

	simpleConverter := func(properties []string) FormatConverter {

		properties = append(properties, defaultProperties...);
		return func(info *LogInfoNode) (map[string]string, error) {
			validProperties := make(map[string]bool)
			for _, property := range properties {
				validProperties[property] = true
			}
			ret := make(map[string]string)
			for k, v := range info.children {
				if _, ok := validProperties[k]; !ok {
					return nil, fmt.Errorf("%s is not a valid properties for format %s", k, info.value)
				}
				ret[k] = v.value
			}
			return ret, nil
		}
	}

	Register("csv", simpleConverter([]string{"keys", "delimiter"}))
	Register("json", simpleConverter([]string{"json_parser"}))
	Register("regexp", simpleConverter([]string{"expression", "ignorecase", "multiline"}))
	Register("apache", simpleConverter([]string{}))
	Register("apache_error", simpleConverter([]string{}))
	Register("nginx", simpleConverter([]string{}))
	Register("none", simpleConverter([]string{"message_key"}))

}
