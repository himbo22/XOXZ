package util

func GetString(m map[string]any, key string, def string) *string {
	if v, ok := m[key].(string); ok {
		return &v
	}
	return &def
}

func GetBool(m map[string]any, key string, def bool) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return def
}
