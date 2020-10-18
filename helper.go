package ilert

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}
