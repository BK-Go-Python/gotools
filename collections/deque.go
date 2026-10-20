package collections

// Deque 是 Python collections.deque 风格的泛型双端队列。
//
// Deque 使用环形缓冲区实现，两端插入、删除均为摊还 O(1)；
// 零值可直接使用，无需先调用 NewDeque。
//
// Deque 不是并发安全的。nil *Deque 上的只读方法返回零值或空结果，
// 修改方法为无操作，因此调用方无需额外判空。
type Deque[T any] struct {
	buf  []T
	head int // buf 中左端元素的下标
	n    int // 元素个数
}

// NewDeque 创建一个 Deque，items 按从左到右的顺序入队。
func NewDeque[T any](items ...T) *Deque[T] {
	d := &Deque[T]{
		buf: make([]T, len(items)),
		n:   len(items),
	}
	copy(d.buf, items)
	return d
}

// Append 依次把 items 追加到右端，等价于 Python 的 deque.extend。
// 每次插入摊还 O(1)，共计 O(len(items))。
func (d *Deque[T]) Append(items ...T) {
	if d == nil {
		return
	}
	for _, item := range items {
		if d.n == len(d.buf) {
			d.grow()
		}
		pos := d.head + d.n
		if pos >= len(d.buf) {
			pos -= len(d.buf)
		}
		d.buf[pos] = item
		d.n++
	}
}

// AppendLeft 按参数顺序依次把 items 加到左端。
//
// 与 Python 的 deque.extendleft 一致：因为每一项都插入到当前最左端，
// 所以最终顺序与入参相反。例如 AppendLeft(1, 2, 3) 后从左到右为 [3 2 1]。
// 每次插入摊还 O(1)，共计 O(len(items))。
func (d *Deque[T]) AppendLeft(items ...T) {
	if d == nil {
		return
	}
	for _, item := range items {
		if d.n == len(d.buf) {
			d.grow()
		}
		d.head--
		if d.head < 0 {
			d.head += len(d.buf)
		}
		d.buf[d.head] = item
		d.n++
	}
}

// Extend 等价于 Python 的 deque.extend，按顺序把 items 追加到右端。
func (d *Deque[T]) Extend(items ...T) {
	d.Append(items...)
}

// ExtendLeft 等价于 Python 的 deque.extendleft，按顺序把 items 加到左端。
// 最终顺序与入参相反，详见 AppendLeft。
func (d *Deque[T]) ExtendLeft(items ...T) {
	d.AppendLeft(items...)
}

// Pop 删除并返回右端元素。队列为空或 d 为 nil 时返回零值和 false。
// 复杂度 O(1)。
func (d *Deque[T]) Pop() (T, bool) {
	var zero T
	if d == nil || d.n == 0 {
		return zero, false
	}
	pos := d.head + d.n - 1
	if pos >= len(d.buf) {
		pos -= len(d.buf)
	}
	item := d.buf[pos]
	d.buf[pos] = zero // 清除引用，便于垃圾回收
	d.n--
	if d.n == 0 {
		d.head = 0
	}
	return item, true
}

// PopLeft 删除并返回左端元素。队列为空或 d 为 nil 时返回零值和 false。
// 复杂度 O(1)。
func (d *Deque[T]) PopLeft() (T, bool) {
	var zero T
	if d == nil || d.n == 0 {
		return zero, false
	}
	item := d.buf[d.head]
	d.buf[d.head] = zero // 清除引用，便于垃圾回收
	d.head++
	if d.head == len(d.buf) {
		d.head = 0
	}
	d.n--
	if d.n == 0 {
		d.head = 0
	}
	return item, true
}

// Peek 返回右端元素但不删除。队列为空或 d 为 nil 时返回零值和 false。
// 复杂度 O(1)。
func (d *Deque[T]) Peek() (T, bool) {
	var zero T
	if d == nil || d.n == 0 {
		return zero, false
	}
	pos := d.head + d.n - 1
	if pos >= len(d.buf) {
		pos -= len(d.buf)
	}
	return d.buf[pos], true
}

// PeekLeft 返回左端元素但不删除。队列为空或 d 为 nil 时返回零值和 false。
// 复杂度 O(1)。
func (d *Deque[T]) PeekLeft() (T, bool) {
	var zero T
	if d == nil || d.n == 0 {
		return zero, false
	}
	return d.buf[d.head], true
}

// Len 返回元素个数。d 为 nil 时返回 0。复杂度 O(1)。
func (d *Deque[T]) Len() int {
	if d == nil {
		return 0
	}
	return d.n
}

// Items 按从左到右的顺序返回元素的新切片。
//
// 返回的切片与 Deque 不共享底层数组，修改它不会影响 Deque；
// 空 Deque（包括 nil）返回非 nil 的空切片。复杂度 O(n)。
func (d *Deque[T]) Items() []T {
	if d == nil || d.n == 0 {
		return []T{}
	}
	items := make([]T, d.n)
	if d.head+d.n <= len(d.buf) {
		copy(items, d.buf[d.head:d.head+d.n])
		return items
	}
	k := copy(items, d.buf[d.head:])
	copy(items[k:], d.buf[:d.n-k])
	return items
}

// Rotate 旋转队列：n > 0 时右旋 n 位，n < 0 时左旋 |n| 位，n = 0 不变。
// 旋转位数对 Len 取模，因此超过长度的 n 等价于 n % Len。
//
// 复杂度 O(min(k, Len-k))，其中 k = |n| % Len，最坏为 O(Len)。
func (d *Deque[T]) Rotate(n int) {
	if d == nil || d.n <= 1 {
		return
	}
	size := d.n
	n %= size
	if n < 0 {
		n += size
	}
	if n == 0 {
		return
	}
	// 右旋 n 位等价于左旋 size-n 位；选择步数较少的方向执行。
	if n <= size-n {
		for i := 0; i < n; i++ {
			item, _ := d.Pop()
			d.AppendLeft(item)
		}
		return
	}
	for i := 0; i < size-n; i++ {
		item, _ := d.PopLeft()
		d.Append(item)
	}
}

// Clear 删除所有元素并释放底层存储，之后 Len 为 0。
// d 为 nil 时无操作。复杂度 O(1)。
func (d *Deque[T]) Clear() {
	if d == nil {
		return
	}
	d.buf = nil
	d.head = 0
	d.n = 0
}

// Copy 返回一个包含相同元素的独立 Deque，后续对两者的修改互不影响。
// d 为 nil 时返回一个空 Deque。复杂度 O(n)。
func (d *Deque[T]) Copy() *Deque[T] {
	if d == nil {
		return &Deque[T]{}
	}
	c := &Deque[T]{
		buf: make([]T, d.n),
		n:   d.n,
	}
	if d.n == 0 {
		return c
	}
	if d.head+d.n <= len(d.buf) {
		copy(c.buf, d.buf[d.head:d.head+d.n])
		return c
	}
	k := copy(c.buf, d.buf[d.head:])
	copy(c.buf[k:], d.buf[:d.n-k])
	return c
}

// grow 把环形缓冲区的容量翻倍，并把逻辑元素按顺序搬到新缓冲区的开头。
// 仅在缓冲区已满时调用。
func (d *Deque[T]) grow() {
	old := d.buf
	newCap := len(old) * 2
	if newCap < 8 {
		newCap = 8
	}
	d.buf = make([]T, newCap)
	if d.n > 0 {
		if d.head+d.n <= len(old) {
			copy(d.buf, old[d.head:d.head+d.n])
		} else {
			k := copy(d.buf, old[d.head:])
			copy(d.buf[k:], old[:d.n-k])
		}
	}
	d.head = 0
}
