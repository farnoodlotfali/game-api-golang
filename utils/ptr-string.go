package utils

func PtrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
