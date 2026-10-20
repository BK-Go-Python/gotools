package collections

import "github.com/BK-Go-Python/gotools/internal/pyval"

// OrderedItem 是 OrderedDict 中的一个键值对。
type OrderedItem[K comparable, V any] struct {
	Key   K
	Value V
}

// odNode 是 OrderedDict 内部双向链表的节点。
// head 指向最早插入的节点，tail 指向最近插入的节点。
type odNode[K comparable, V any] struct {
	key        K
	value      V
	prev, next *odNode[K, V]
}

// OrderedDict 是保持插入顺序的映射，对应 Python 的 collections.OrderedDict。
//
// 内部使用 map 加双向链表实现，Set/Get/Has/Delete/Pop/PopItem/MoveToEnd
// 均为 O(1) 时间复杂度。零值可以直接使用；nil 指针上的方法都不会 panic，
// 只读方法返回零值或空结果，写方法是无操作。
type OrderedDict[K comparable, V any] struct {
	items map[K]*odNode[K, V]
	head  *odNode[K, V]
	tail  *odNode[K, V]
}

// NewOrderedDict 创建一个空的 OrderedDict。
func NewOrderedDict[K comparable, V any]() *OrderedDict[K, V] {
	return &OrderedDict[K, V]{items: make(map[K]*odNode[K, V])}
}

// Set 等价于 Python 的 od[key] = value。键已存在时只更新值，不改变位置；
// 新键追加到末尾。它返回自身以支持链式调用，nil 接收者上的 Set 无操作。
func (d *OrderedDict[K, V]) Set(key K, value V) *OrderedDict[K, V] {
	if d == nil {
		return nil
	}
	if d.items == nil {
		d.items = make(map[K]*odNode[K, V])
	}
	if n, ok := d.items[key]; ok {
		n.value = value
		return d
	}
	n := &odNode[K, V]{key: key, value: value, prev: d.tail}
	if d.tail != nil {
		d.tail.next = n
	} else {
		d.head = n
	}
	d.tail = n
	d.items[key] = n
	return d
}

// Get 等价于 Python 的 od.get(key, default)。键存在时返回对应的值；
// 否则返回第一个 def，未提供 def 时返回 V 的零值。
func (d *OrderedDict[K, V]) Get(key K, def ...V) V {
	if n, ok := d.lookup(key); ok {
		return n.value
	}
	if len(def) > 0 {
		return def[0]
	}
	var zero V
	return zero
}

// Has 等价于 Python 的 key in od。
func (d *OrderedDict[K, V]) Has(key K) bool {
	_, ok := d.lookup(key)
	return ok
}

// Delete 删除指定键。键存在时删除并返回 true，否则返回 false。
func (d *OrderedDict[K, V]) Delete(key K) bool {
	n, ok := d.lookup(key)
	if !ok {
		return false
	}
	d.unlink(n)
	delete(d.items, key)
	return true
}

// Pop 等价于 Python 的 od.pop(key, default)。键存在时删除并返回其值；
// 否则返回第一个 def，未提供 def 时返回 V 的零值（Python 会抛出 KeyError）。
func (d *OrderedDict[K, V]) Pop(key K, def ...V) V {
	n, ok := d.lookup(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		var zero V
		return zero
	}
	value := n.value
	d.unlink(n)
	delete(d.items, key)
	return value
}

// PopItem 等价于 Python 的 od.popitem(last=...)。默认 last 为 true，
// 移除并返回最近插入的项；last 为 false 时移除最早插入的项。
// 字典为空时返回 (OrderedItem{}, false)。
func (d *OrderedDict[K, V]) PopItem(last ...bool) (OrderedItem[K, V], bool) {
	popLast := true
	if len(last) > 0 {
		popLast = last[0]
	}
	if d == nil {
		return OrderedItem[K, V]{}, false
	}
	n := d.tail
	if !popLast {
		n = d.head
	}
	if n == nil {
		return OrderedItem[K, V]{}, false
	}
	item := OrderedItem[K, V]{Key: n.key, Value: n.value}
	d.unlink(n)
	delete(d.items, n.key)
	return item, true
}

// Keys 按插入顺序返回所有键。
func (d *OrderedDict[K, V]) Keys() []K {
	if d == nil {
		return []K{}
	}
	keys := make([]K, 0, len(d.items))
	for n := d.head; n != nil; n = n.next {
		keys = append(keys, n.key)
	}
	return keys
}

// Values 按插入顺序返回所有值。
func (d *OrderedDict[K, V]) Values() []V {
	if d == nil {
		return []V{}
	}
	values := make([]V, 0, len(d.items))
	for n := d.head; n != nil; n = n.next {
		values = append(values, n.value)
	}
	return values
}

// Items 按插入顺序返回所有键值对。
func (d *OrderedDict[K, V]) Items() []OrderedItem[K, V] {
	if d == nil {
		return []OrderedItem[K, V]{}
	}
	items := make([]OrderedItem[K, V], 0, len(d.items))
	for n := d.head; n != nil; n = n.next {
		items = append(items, OrderedItem[K, V]{Key: n.key, Value: n.value})
	}
	return items
}

// Len 等价于 Python 的 len(od)。
func (d *OrderedDict[K, V]) Len() int {
	if d == nil {
		return 0
	}
	return len(d.items)
}

// Clear 等价于 Python 的 od.clear()，清空字典并保持可用。
func (d *OrderedDict[K, V]) Clear() {
	if d == nil {
		return
	}
	d.items = nil
	d.head = nil
	d.tail = nil
}

// Copy 返回浅拷贝，等价于 Python 的 od.copy()，副本与原始字典互不影响。
func (d *OrderedDict[K, V]) Copy() *OrderedDict[K, V] {
	out := NewOrderedDict[K, V]()
	if d == nil {
		return out
	}
	for n := d.head; n != nil; n = n.next {
		out.Set(n.key, n.value)
	}
	return out
}

// MoveToEnd 等价于 Python 的 od.move_to_end(key)，把已存在的键移到末尾。
// 键存在时返回 true，否则返回 false 且不改变顺序。
func (d *OrderedDict[K, V]) MoveToEnd(key K) bool {
	n, ok := d.lookup(key)
	if !ok {
		return false
	}
	if d.tail == n {
		return true
	}
	d.unlink(n)
	n.prev = d.tail
	if d.tail != nil {
		d.tail.next = n
	}
	d.tail = n
	if d.head == nil {
		d.head = n
	}
	return true
}

// Equal 比较两个 OrderedDict 是否键值对相同且顺序相同，即 Python 的
// OrderedDict 有序相等比较。值的比较使用 pyval.Equal，支持 int/float
// 等跨类型相等，也支持切片、映射等不可比较类型，不会 panic。
func (d *OrderedDict[K, V]) Equal(other *OrderedDict[K, V]) bool {
	if d == other {
		return true
	}
	if d == nil {
		return other.Len() == 0
	}
	if other == nil {
		return d.Len() == 0
	}
	if len(d.items) != len(other.items) {
		return false
	}
	dn, on := d.head, other.head
	for dn != nil && on != nil {
		if dn.key != on.key || !pyval.Equal(dn.value, on.value) {
			return false
		}
		dn, on = dn.next, on.next
	}
	return dn == nil && on == nil
}

// lookup 返回 key 对应的链表节点。
func (d *OrderedDict[K, V]) lookup(key K) (*odNode[K, V], bool) {
	if d == nil || d.items == nil {
		return nil, false
	}
	n, ok := d.items[key]
	return n, ok
}

// unlink 从双向链表中摘除节点 n。
func (d *OrderedDict[K, V]) unlink(n *odNode[K, V]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		d.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		d.tail = n.prev
	}
	n.prev = nil
	n.next = nil
}
