package fsutil

func ToByte(s string) []byte {
	return []byte(s)
}

func ToStr(b []byte) string {
	return string(b[:])
}
