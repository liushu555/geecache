package geecache

//存储数据的抽象

type Byteview struct {
	bytes []byte
}

func (byteview Byteview) ByteSlice() []byte {
	return cloneByte(byteview.bytes)
}

func (byteview Byteview) Len() int {
	return len(byteview.bytes)
}

func cloneByte(bytes []byte) []byte {
	res := make([]byte, len(bytes))
	copy(res, bytes)
	return res
}
