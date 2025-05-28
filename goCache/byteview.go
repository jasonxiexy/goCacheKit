package gocache

// A byteview holds an immutable view of bytes
type Byteview struct {
	b []byte
}

// len returns the view's length
func (v Byteview) Len() int {
	return len(v.b)
}

// byteslice returns a copy of the data as a byte slice
func (v Byteview) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// String returns the data as a string, making a copy if necessary
func (v Byteview) String() string {
	return string(v.b)
}
