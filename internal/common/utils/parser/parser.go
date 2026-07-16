package parserUtils

func MergeField[T comparable](dbField, input T) T {
	var zeroValue T
	if input != zeroValue {
		return input
	}
	return dbField
}

func StrOrEmpty(value string) *string {
	if value != "" {
		return &value
	}
	return nil
}

func StrOrDBNull(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}