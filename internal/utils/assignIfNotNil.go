package utils

func AssignIfNotNil[T any](source *T, defaultValue T) T {
	if source != nil {
		return *source
	}
	return defaultValue
}
