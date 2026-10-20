package py

// MemoryView 返回与 b 共享底层数组的字节视图。
//
// Python: memoryview(b)
// DIFF: 只支持一维 []byte；不支持 format、shape、多维、只读视图。
func MemoryView(b []byte) []byte {
	if b == nil {
		return nil
	}
	return b[:len(b):len(b)]
}
