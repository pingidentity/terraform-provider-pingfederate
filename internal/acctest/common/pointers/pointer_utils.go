package pointers

func String(val string) *string {
	return &val
}

func Bool(val bool) *bool {
	return &val
}

func Int64(val int64) *int64 {
	return &val
}

func Float64(val float64) *float64 {
	return &val
}
