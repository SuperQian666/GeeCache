package geecache

// ByteView 使用[]byte可以缓存任意类型数据
type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	dst := make([]byte, len(b))
	copy(dst, b)
	return dst
}

func (v ByteView) String() string {
	return string(v.b)
}
