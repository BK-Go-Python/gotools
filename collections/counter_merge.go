package collections

// AddCounter 等价于 Python 的 Counter + other：逐项相加，只保留正计数。
func (c Counter[T]) AddCounter(other Counter[T]) Counter[T] {
	out := make(Counter[T], len(c)+len(other))
	for value, count := range c {
		if total := count + other[value]; total > 0 {
			out[value] = total
		}
	}
	for value, count := range other {
		if _, exists := c[value]; !exists && count > 0 {
			out[value] = count
		}
	}
	return out
}

// SubCounter 等价于 Python 的 Counter - other：逐项相减，只保留正计数。
func (c Counter[T]) SubCounter(other Counter[T]) Counter[T] {
	out := make(Counter[T], len(c))
	for value, count := range c {
		if total := count - other[value]; total > 0 {
			out[value] = total
		}
	}
	return out
}

// Union 等价于 Python 的 Counter | other：逐项取较大计数，只保留正计数。
func (c Counter[T]) Union(other Counter[T]) Counter[T] {
	out := make(Counter[T], len(c)+len(other))
	for value, count := range c {
		if otherCount := other[value]; otherCount > count {
			count = otherCount
		}
		if count > 0 {
			out[value] = count
		}
	}
	for value, count := range other {
		if _, exists := c[value]; !exists && count > 0 {
			out[value] = count
		}
	}
	return out
}

// Intersection 等价于 Python 的 Counter & other：逐项取较小计数，只保留正计数。
func (c Counter[T]) Intersection(other Counter[T]) Counter[T] {
	out := make(Counter[T])
	for value, count := range c {
		if otherCount := other[value]; otherCount < count {
			count = otherCount
		}
		if count > 0 {
			out[value] = count
		}
	}
	return out
}

// Negate 等价于 Python 的一元 -Counter：计数取负，只保留正计数。
func (c Counter[T]) Negate() Counter[T] {
	out := make(Counter[T], len(c))
	for value, count := range c {
		if count = -count; count > 0 {
			out[value] = count
		}
	}
	return out
}

// UpdateCounter 等价于 Python 的 c.update(other_counter)：计数相加，保留零和负计数。
func (c Counter[T]) UpdateCounter(other Counter[T]) Counter[T] {
	if c == nil {
		c = make(Counter[T], len(other))
	}
	for value, count := range other {
		c[value] += count
	}
	return c
}

// SubtractCounter 等价于 Python 的 c.subtract(other_counter)：计数相减，允许负计数。
func (c Counter[T]) SubtractCounter(other Counter[T]) Counter[T] {
	if c == nil {
		c = make(Counter[T], len(other))
	}
	for value, count := range other {
		c[value] -= count
	}
	return c
}
