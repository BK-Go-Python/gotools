package list

import (
	"reflect"
	"testing"
)

// assertList compares a List with the expected Python-semantics elements.
func assertList(t *testing.T, got List, want []interface{}) {
	t.Helper()
	if !reflect.DeepEqual([]interface{}(got), want) {
		t.Fatalf("got %s (%#v), want %#v", got.String(), []interface{}(got), want)
	}
}

func TestNewFromMustFrom(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		cases := []struct {
			name string
			args []interface{}
			want []interface{}
		}{
			// Python: list() == []
			{"empty", nil, []interface{}{}},
			// Python: list((1, 'a', None)) == [1, 'a', None]
			{"values", []interface{}{1, "a", nil}, []interface{}{1, "a", nil}},
			{"nested", []interface{}{New(1, 2), nil}, []interface{}{List{1, 2}, nil}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := New(tc.args...)
				assertList(t, got, tc.want)
				if got.Len() != len(tc.want) {
					t.Fatalf("Len() = %d, want %d", got.Len(), len(tc.want))
				}
			})
		}
	})

	t.Run("From", func(t *testing.T) {
		cases := []struct {
			name    string
			in      interface{}
			want    []interface{}
			wantErr bool
		}{
			// Python: list([1, 2, 3]) == [1, 2, 3]
			{"slice", []int{1, 2, 3}, []interface{}{1, 2, 3}, false},
			{"interface slice", []string{"x", "y"}, []interface{}{"x", "y"}, false},
			// Python: list('aé') == ['a', 'é'] (iterates characters, not bytes)
			{"string runes", "aé", []interface{}{"a", "é"}, false},
			{"array", [2]int{4, 5}, []interface{}{4, 5}, false},
			// Python: list({'b': 2, 'a': 1}) follows insertion order; gotools sorts
			// map keys for determinism: ['a', 'b'].
			{"map keys sorted", map[string]int{"b": 2, "a": 1}, []interface{}{"a", "b"}, false},
			// Python: list(None) raises TypeError.
			{"nil", nil, nil, true},
			{"not iterable", 42, nil, true},
			{"struct", struct{}{}, nil, true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := From(tc.in)
				if tc.wantErr {
					if err == nil {
						t.Fatalf("From(%#v) error = nil, want error", tc.in)
					}
					if got != nil {
						t.Fatalf("From(%#v) = %#v, want nil on error", tc.in, got)
					}
					return
				}
				if err != nil {
					t.Fatalf("From(%#v) unexpected error: %v", tc.in, err)
				}
				assertList(t, got, tc.want)
			})
		}
	})

	t.Run("MustFrom", func(t *testing.T) {
		got := MustFrom([]int{1, 2})
		assertList(t, got, []interface{}{1, 2})

		defer func() {
			if recover() == nil {
				t.Fatal("MustFrom(42) did not panic")
			}
		}()
		MustFrom(42)
	})
}

func TestAppendExtend(t *testing.T) {
	t.Run("Append", func(t *testing.T) {
		cases := []struct {
			name  string
			start List
			items []interface{}
			want  []interface{}
		}{
			{"empty", List{}, []interface{}{1}, []interface{}{1}},
			{"one", List{1, 2}, []interface{}{3}, []interface{}{1, 2, 3}},
			{"variadic", List{1}, []interface{}{2, 3}, []interface{}{1, 2, 3}},
			{"nil element", List{1}, []interface{}{nil}, []interface{}{1, nil}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.start
				if ret := got.Append(tc.items...); ret != &got {
					t.Fatal("Append did not return the receiver")
				}
				assertList(t, got, tc.want)
			})
		}
	})

	t.Run("Extend", func(t *testing.T) {
		cases := []struct {
			name    string
			start   List
			in      interface{}
			want    []interface{}
			wantErr bool
		}{
			// Python: [1].extend([2, 3]) -> [1, 2, 3]
			{"int slice", List{1}, []int{2, 3}, []interface{}{1, 2, 3}, false},
			{"interface slice", List{"a"}, []interface{}{"b", "c"}, []interface{}{"a", "b", "c"}, false},
			// Python: ['x'].extend('ab') -> ['x', 'a', 'b']
			{"string", List{"x"}, "ab", []interface{}{"x", "a", "b"}, false},
			{"array", List{1}, [2]int{2, 3}, []interface{}{1, 2, 3}, false},
			{"empty", List{1}, []interface{}{}, []interface{}{1}, false},
			// Python: [].extend(None) raises TypeError.
			{"nil", List{1}, nil, []interface{}{1}, true},
			{"not iterable", List{1}, 42, []interface{}{1}, true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.start
				err := got.Extend(tc.in)
				if tc.wantErr {
					if err == nil {
						t.Fatalf("Extend(%#v) error = nil, want error", tc.in)
					}
				} else if err != nil {
					t.Fatalf("Extend(%#v) unexpected error: %v", tc.in, err)
				}
				assertList(t, got, tc.want)
			})
		}
	})
}

func TestInsert(t *testing.T) {
	cases := []struct {
		name  string
		start List
		index int
		value interface{}
		want  []interface{}
	}{
		// Python: [1, 2, 3].insert(1, 9) -> [1, 9, 2, 3]
		{"middle", List{1, 2, 3}, 1, 9, []interface{}{1, 9, 2, 3}},
		{"head", List{1, 2, 3}, 0, 0, []interface{}{0, 1, 2, 3}},
		{"at len", List{1, 2, 3}, 3, 4, []interface{}{1, 2, 3, 4}},
		// Python: insert(10, 9) clamps to len.
		{"past len clamps", List{1, 2, 3}, 10, 9, []interface{}{1, 2, 3, 9}},
		// Python: [1, 2, 3].insert(-1, 9) -> [1, 2, 9, 3]
		{"negative one", List{1, 2, 3}, -1, 9, []interface{}{1, 2, 9, 3}},
		{"negative len", List{1, 2, 3}, -3, 0, []interface{}{0, 1, 2, 3}},
		// Python: insert(-10, 9) clamps to index 0.
		{"negative past head clamps", List{1, 2, 3}, -10, 9, []interface{}{9, 1, 2, 3}},
		{"empty positive", List{}, 5, "x", []interface{}{"x"}},
		{"empty negative", List{}, -5, "x", []interface{}{"x"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.start
			if ret := got.Insert(tc.index, tc.value); ret != &got {
				t.Fatal("Insert did not return the receiver")
			}
			assertList(t, got, tc.want)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name    string
		start   List
		value   interface{}
		want    []interface{}
		wantErr bool
	}{
		// Python: [1, 2, 2, 3].remove(2) removes only the first match.
		{"first duplicate", List{1, 2, 2, 3}, 2, []interface{}{1, 2, 3}, false},
		{"head", List{1, 2, 3}, 1, []interface{}{2, 3}, false},
		{"tail", List{1, 2, 3}, 3, []interface{}{1, 2}, false},
		// Python: [1, 2].remove(1.0) -> [2] because 1 == 1.0.
		{"cross numeric", List{1, 2}, 1.0, []interface{}{2}, false},
		{"nil", List{nil, 1}, nil, []interface{}{1}, false},
		// Python: raises ValueError when the value is absent.
		{"missing", List{1, 2}, 9, []interface{}{1, 2}, true},
		{"empty", List{}, 1, []interface{}{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.start
			err := got.Remove(tc.value)
			if tc.wantErr {
				if err == nil {
					t.Fatal("Remove error = nil, want error")
				}
			} else if err != nil {
				t.Fatalf("Remove unexpected error: %v", err)
			}
			assertList(t, got, tc.want)
		})
	}
}

func TestPop(t *testing.T) {
	cases := []struct {
		name    string
		start   List
		want    interface{}
		rest    []interface{}
		wantErr bool
	}{
		// Python: [1, 2, 3].pop() == 3, list becomes [1, 2].
		{"last", List{1, 2, 3}, 3, []interface{}{1, 2}, false},
		{"single", List{1}, 1, []interface{}{}, false},
		// Python: [].pop() raises IndexError.
		{"empty", List{}, nil, []interface{}{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.start
			value, err := got.Pop()
			if tc.wantErr {
				if err == nil {
					t.Fatal("Pop error = nil, want error")
				}
			} else if err != nil {
				t.Fatalf("Pop unexpected error: %v", err)
			}
			if !reflect.DeepEqual(value, tc.want) {
				t.Fatalf("Pop value = %#v, want %#v", value, tc.want)
			}
			assertList(t, got, tc.rest)
		})
	}
}

func TestPopAt(t *testing.T) {
	cases := []struct {
		name    string
		start   List
		index   int
		want    interface{}
		rest    []interface{}
		wantErr bool
	}{
		{"head", List{1, 2, 3}, 0, 1, []interface{}{2, 3}, false},
		{"middle", List{1, 2, 3}, 1, 2, []interface{}{1, 3}, false},
		{"tail", List{1, 2, 3}, 2, 3, []interface{}{1, 2}, false},
		// Python: [1, 2, 3].pop(-1) == 3
		{"negative tail", List{1, 2, 3}, -1, 3, []interface{}{1, 2}, false},
		{"negative head", List{1, 2, 3}, -3, 1, []interface{}{2, 3}, false},
		// Python: pop(3) / pop(-4) raise IndexError.
		{"past end", List{1, 2, 3}, 3, nil, []interface{}{1, 2, 3}, true},
		{"past head", List{1, 2, 3}, -4, nil, []interface{}{1, 2, 3}, true},
		{"empty", List{}, 0, nil, []interface{}{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.start
			value, err := got.PopAt(tc.index)
			if tc.wantErr {
				if err == nil {
					t.Fatal("PopAt error = nil, want error")
				}
			} else if err != nil {
				t.Fatalf("PopAt unexpected error: %v", err)
			}
			if !reflect.DeepEqual(value, tc.want) {
				t.Fatalf("PopAt value = %#v, want %#v", value, tc.want)
			}
			assertList(t, got, tc.rest)
		})
	}
}

func TestClearCopy(t *testing.T) {
	t.Run("Clear", func(t *testing.T) {
		cases := []struct {
			name  string
			start List
		}{
			{"non empty", List{1, "a", nil}},
			{"empty", List{}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.start
				if ret := got.Clear(); ret != &got {
					t.Fatal("Clear did not return the receiver")
				}
				assertList(t, got, []interface{}{})
				if got.Len() != 0 {
					t.Fatalf("Len() = %d, want 0", got.Len())
				}
			})
		}
	})

	t.Run("Copy", func(t *testing.T) {
		orig := New(1, "a", nil)
		cp := orig.Copy()
		cp[0] = 99
		assertList(t, orig, []interface{}{1, "a", nil})
		assertList(t, cp, []interface{}{99, "a", nil})

		// Python: copy() is shallow; nested mutable values are shared.
		inner := New(1, 2)
		origNested := List{inner}
		cpNested := origNested.Copy()
		cpNested[0].(List)[0] = 9
		assertList(t, origNested[0].(List), []interface{}{9, 2})
	})
}

func TestGetSet(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		cases := []struct {
			name    string
			list    List
			index   int
			want    interface{}
			wantErr bool
		}{
			{"first", List{10, 20, 30}, 0, 10, false},
			{"last", List{10, 20, 30}, 2, 30, false},
			// Python: l[-1] == 30, l[-3] == 10
			{"negative last", List{10, 20, 30}, -1, 30, false},
			{"negative first", List{10, 20, 30}, -3, 10, false},
			{"past end", List{10, 20, 30}, 3, nil, true},
			{"past head", List{10, 20, 30}, -4, nil, true},
			{"empty", List{}, 0, nil, true},
			{"empty negative", List{}, -1, nil, true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := tc.list.Get(tc.index)
				if tc.wantErr {
					if err == nil {
						t.Fatal("Get error = nil, want error")
					}
					return
				}
				if err != nil {
					t.Fatalf("Get unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("Get = %#v, want %#v", got, tc.want)
				}
			})
		}
	})

	t.Run("Set", func(t *testing.T) {
		cases := []struct {
			name    string
			start   List
			index   int
			value   interface{}
			want    []interface{}
			wantErr bool
		}{
			{"head", List{1, 2, 3}, 0, 9, []interface{}{9, 2, 3}, false},
			{"tail", List{1, 2, 3}, 2, 9, []interface{}{1, 2, 9}, false},
			// Python: l[-1] = 9
			{"negative", List{1, 2, 3}, -1, 9, []interface{}{1, 2, 9}, false},
			{"past end", List{1, 2, 3}, 3, 9, []interface{}{1, 2, 3}, true},
			{"past head", List{1, 2, 3}, -4, 9, []interface{}{1, 2, 3}, true},
			{"empty", List{}, 0, 9, []interface{}{}, true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.start
				err := got.Set(tc.index, tc.value)
				if tc.wantErr {
					if err == nil {
						t.Fatal("Set error = nil, want error")
					}
				} else if err != nil {
					t.Fatalf("Set unexpected error: %v", err)
				}
				assertList(t, got, tc.want)
			})
		}
	})
}

func TestIndexCountContains(t *testing.T) {
	t.Run("Index", func(t *testing.T) {
		cases := []struct {
			name    string
			list    List
			value   interface{}
			want    int
			wantErr bool
		}{
			{"first", List{1, 2, 1, 3}, 1, 0, false},
			{"last", List{1, 2, 1, 3}, 3, 3, false},
			// Python: [1, 2].index(2.0) == 1
			{"cross numeric", List{1, 2}, 2.0, 1, false},
			{"global string", List{"x", 1}, "x", 0, false},
			{"missing", List{1, 2, 1, 3}, 9, -1, true},
			{"empty", List{}, 1, -1, true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := tc.list.Index(tc.value)
				if tc.wantErr {
					if err == nil {
						t.Fatal("Index error = nil, want error")
					}
				} else if err != nil {
					t.Fatalf("Index unexpected error: %v", err)
				}
				if got != tc.want {
					t.Fatalf("Index = %d, want %d", got, tc.want)
				}
			})
		}
	})

	t.Run("Count", func(t *testing.T) {
		cases := []struct {
			name  string
			list  List
			value interface{}
			want  int
		}{
			{"duplicates", List{1, 2, 1, 3}, 1, 2},
			// Python: [1, 2, 1].count(1.0) == 2
			{"cross numeric", List{1, 2, 1}, 1.0, 2},
			{"missing", List{1, 2}, 9, 0},
			{"empty", List{}, 1, 0},
			{"nil", List{nil, 1, nil}, nil, 2},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.list.Count(tc.value); got != tc.want {
					t.Fatalf("Count = %d, want %d", got, tc.want)
				}
			})
		}
	})

	t.Run("Contains", func(t *testing.T) {
		cases := []struct {
			name  string
			list  List
			value interface{}
			want  bool
		}{
			{"present", List{1, 2, 3}, 2, true},
			// Python: 1.0 in [1, 2] is True
			{"cross numeric", List{1, 2}, 1.0, true},
			{"absent", List{1, 2}, 4, false},
			{"empty", List{}, 1, false},
			{"nil", List{nil}, nil, true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.list.Contains(tc.value); got != tc.want {
					t.Fatalf("Contains = %v, want %v", got, tc.want)
				}
			})
		}
	})
}

func TestReverse(t *testing.T) {
	cases := []struct {
		name  string
		start List
		want  []interface{}
	}{
		// Python: [1, 2, 3, 4].reverse() -> [4, 3, 2, 1]
		{"even", List{1, 2, 3, 4}, []interface{}{4, 3, 2, 1}},
		{"odd", List{1, 2, 3}, []interface{}{3, 2, 1}},
		{"single", List{1}, []interface{}{1}},
		{"empty", List{}, []interface{}{}},
		{"palindrome", List{"a", "b", "a"}, []interface{}{"a", "b", "a"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.start
			if ret := got.Reverse(); ret != &got {
				t.Fatal("Reverse did not return the receiver")
			}
			assertList(t, got, tc.want)
		})
	}
}

func TestSort(t *testing.T) {
	cases := []struct {
		name  string
		start List
		want  []interface{}
	}{
		// Python: sorted([3, 1.5, -2, 0, 2.5, 10]) == [-2, 0, 1.5, 2.5, 3, 10]
		{"mixed numbers", List{3, 1.5, -2, 0, 2.5, 10}, []interface{}{-2, 0, 1.5, 2.5, 3, 10}},
		// Python: sorted([2, 1.0, 1, 3]) == [1.0, 1, 2, 3] (stable, 1.0 == 1).
		{"stable equal numbers", List{2, 1.0, 1, 3}, []interface{}{1.0, 1, 2, 3}},
		{"strings", List{"banana", "apple", "cherry"}, []interface{}{"apple", "banana", "cherry"}},
		{"already sorted", List{1, 2, 3}, []interface{}{1, 2, 3}},
		{"reverse sorted", List{3, 2, 1}, []interface{}{1, 2, 3}},
		{"single", List{1}, []interface{}{1}},
		{"empty", List{}, []interface{}{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.start
			if ret := got.Sort(); ret != &got {
				t.Fatal("Sort did not return the receiver")
			}
			assertList(t, got, tc.want)
		})
	}

	t.Run("custom less", func(t *testing.T) {
		got := List{1, 3, 2}
		// Python: sorted([1, 3, 2], reverse=True) == [3, 2, 1]
		got.Sort(func(a, b interface{}) bool { return a.(int) > b.(int) })
		assertList(t, got, []interface{}{3, 2, 1})
	})
}

type sortPerson struct {
	name string
	age  int
}

func TestSortBy(t *testing.T) {
	people := func() List {
		return List{
			sortPerson{"a", 30},
			sortPerson{"b", 20},
			sortPerson{"c", 30},
			sortPerson{"d", 10},
		}
	}
	age := func(v interface{}) interface{} { return v.(sortPerson).age }
	length := func(v interface{}) interface{} { return len(v.(string)) }

	cases := []struct {
		name    string
		start   List
		key     func(interface{}) interface{}
		reverse []bool
		want    []interface{}
	}{
		// Python: sorted(people, key=lambda p: p.age) keeps a before c on ties.
		{
			"age ascending stable",
			people(), age, nil,
			[]interface{}{
				sortPerson{"d", 10}, sortPerson{"b", 20},
				sortPerson{"a", 30}, sortPerson{"c", 30},
			},
		},
		// Python: sorted(people, key=lambda p: p.age, reverse=True) is stable:
		// equal ages keep their original relative order (a before c).
		{
			"age descending stable",
			people(), age, []bool{true},
			[]interface{}{
				sortPerson{"a", 30}, sortPerson{"c", 30},
				sortPerson{"b", 20}, sortPerson{"d", 10},
			},
		},
		// Python: sorted(['bb', 'a', 'cc', 'd'], key=len) == ['a', 'd', 'bb', 'cc']
		{
			"length ascending stable",
			List{"bb", "a", "cc", "d"}, length, nil,
			[]interface{}{"a", "d", "bb", "cc"},
		},
		// Python: sorted(..., key=len, reverse=True) == ['bb', 'cc', 'a', 'd']
		{
			"length descending stable",
			List{"bb", "a", "cc", "d"}, length, []bool{true},
			[]interface{}{"bb", "cc", "a", "d"},
		},
		{
			"explicit false is ascending",
			List{"bb", "a", "cc", "d"}, length, []bool{false},
			[]interface{}{"a", "d", "bb", "cc"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.start
			if ret := got.SortBy(tc.key, tc.reverse...); ret != &got {
				t.Fatal("SortBy did not return the receiver")
			}
			assertList(t, got, tc.want)
		})
	}
}

func TestSlice(t *testing.T) {
	base := func() List { return List{0, 1, 2, 3, 4} }
	cases := []struct {
		name     string
		start    List
		startIdx int
		endIdx   int
		want     []interface{}
	}{
		// Python: [0, 1, 2, 3, 4][1:4] == [1, 2, 3]
		{"middle", base(), 1, 4, []interface{}{1, 2, 3}},
		{"full", base(), 0, 5, []interface{}{0, 1, 2, 3, 4}},
		// Python: l[-2:5] == [3, 4]
		{"negative start", base(), -2, 5, []interface{}{3, 4}},
		// Python: l[-4:-1] == [1, 2, 3]
		{"negative both", base(), -4, -1, []interface{}{1, 2, 3}},
		// Python: l[-100:100] == l
		{"clamped both", base(), -100, 100, []interface{}{0, 1, 2, 3, 4}},
		// Python: l[5:2] == []
		{"start past end", base(), 5, 2, []interface{}{}},
		{"equal", base(), 2, 2, []interface{}{}},
		// Python: l[4:2] == []
		{"reversed bounds", base(), 4, 2, []interface{}{}},
		// Python: l[0:-100] == []
		{"negative end clamps to zero", base(), 0, -100, []interface{}{}},
		{"past end", base(), 10, 20, []interface{}{}},
		{"empty list", List{}, 0, 0, []interface{}{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertList(t, tc.start.Slice(tc.startIdx, tc.endIdx), tc.want)
		})
	}

	t.Run("returns a new list", func(t *testing.T) {
		src := New(0, 1, 2, 3, 4)
		part := src.Slice(1, 3)
		part[0] = 99
		assertList(t, src, []interface{}{0, 1, 2, 3, 4})
		assertList(t, part, []interface{}{99, 2})
	})
}

func TestToStringsJoin(t *testing.T) {
	t.Run("ToStrings", func(t *testing.T) {
		cases := []struct {
			name string
			list List
			want []string
		}{
			// Python: list(map(str, [1, 'a', None, 2.5, True])) ==
			// ['1', 'a', 'None', '2.5', 'True']
			{"mixed", List{1, "a", nil, 2.5, true}, []string{"1", "a", "None", "2.5", "True"}},
			// Python: str([1, 2]) == '[1, 2]'
			{"nested", List{List{1, 2}}, []string{"[1, 2]"}},
			{"empty", List{}, []string{}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.list.ToStrings(); !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("ToStrings = %#v, want %#v", got, tc.want)
				}
			})
		}
	})

	t.Run("Join", func(t *testing.T) {
		cases := []struct {
			name string
			list List
			sep  string
			want string
		}{
			// Python: ','.join(['a', 'b', 'c']) == 'a,b,c'
			{"comma", List{"a", "b", "c"}, ",", "a,b,c"},
			{"empty sep", List{"a", "b"}, "", "ab"},
			{"single", List{"x"}, "--", "x"},
			// Python: ','.join([]) == ''
			{"empty", List{}, ",", ""},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.list.Join(tc.sep); got != tc.want {
					t.Fatalf("Join = %q, want %q", got, tc.want)
				}
			})
		}
	})
}

func TestMapFilterEach(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		cases := []struct {
			name string
			list List
			fn   func(interface{}) interface{}
			want []interface{}
		}{
			// Python: [x * 2 for x in [1, 2, 3]] == [2, 4, 6]
			{"double", List{1, 2, 3}, func(v interface{}) interface{} { return v.(int) * 2 }, []interface{}{2, 4, 6}},
			{"empty", List{}, func(v interface{}) interface{} { return v }, []interface{}{}},
			{"stringify", List{1, 2}, func(v interface{}) interface{} { return v.(int) + 1 }, []interface{}{2, 3}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				orig := tc.list.Copy()
				got := tc.list.Map(tc.fn)
				assertList(t, got, tc.want)
				assertList(t, tc.list, []interface{}(orig))
			})
		}
	})

	t.Run("Filter", func(t *testing.T) {
		even := func(v interface{}) bool { return v.(int)%2 == 0 }
		cases := []struct {
			name string
			list List
			fn   func(interface{}) bool
			want []interface{}
		}{
			// Python: [x for x in [1, 2, 3, 4] if x % 2 == 0] == [2, 4]
			{"even", List{1, 2, 3, 4}, even, []interface{}{2, 4}},
			{"none", List{1, 3}, even, []interface{}{}},
			{"all", List{2, 4}, even, []interface{}{2, 4}},
			{"empty", List{}, even, []interface{}{}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				assertList(t, tc.list.Filter(tc.fn), tc.want)
			})
		}
	})

	t.Run("Each", func(t *testing.T) {
		var indexes []int
		var values []interface{}
		// Python: for i, x in enumerate(['a', 'b', 'c'])
		List{"a", "b", "c"}.Each(func(i int, v interface{}) {
			indexes = append(indexes, i)
			values = append(values, v)
		})
		if !reflect.DeepEqual(indexes, []int{0, 1, 2}) {
			t.Fatalf("indexes = %#v, want [0 1 2]", indexes)
		}
		if !reflect.DeepEqual(values, []interface{}{"a", "b", "c"}) {
			t.Fatalf("values = %#v, want ['a' 'b' 'c']", values)
		}

		calls := 0
		List{}.Each(func(int, interface{}) { calls++ })
		if calls != 0 {
			t.Fatalf("Each called fn %d times on an empty list", calls)
		}
	})
}

func TestEqualLenString(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		cases := []struct {
			name string
			a    List
			b    List
			want bool
		}{
			// Python: [] == [] is True
			{"empty", List{}, List{}, true},
			{"same", List{1, "a", nil}, List{1, "a", nil}, true},
			// Python: [1, 2] == [1.0, 2.0] is True
			{"cross numeric", List{1, 2}, List{1.0, 2.0}, true},
			{"nested cross numeric", List{List{1}}, List{List{1.0}}, true},
			{"different order", List{1, 2}, List{2, 1}, false},
			{"different length", List{1, 2}, List{1}, false},
			{"different value", List{1, 2}, List{1, 3}, false},
			{"nil vs empty", List(nil), List{}, true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.a.Equal(tc.b); got != tc.want {
					t.Fatalf("Equal = %v, want %v", got, tc.want)
				}
			})
		}
	})

	t.Run("Len", func(t *testing.T) {
		cases := []struct {
			name string
			list List
			want int
		}{
			{"empty", List{}, 0},
			{"one", List{nil}, 1},
			{"many", List{1, 2, 3}, 3},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.list.Len(); got != tc.want {
					t.Fatalf("Len = %d, want %d", got, tc.want)
				}
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		cases := []struct {
			name string
			list List
			want string
		}{
			// Python: repr([]) == '[]'
			{"empty", List{}, "[]"},
			// Python: repr([1, 'a', None, True, 2.0]) == "[1, 'a', None, True, 2.0]"
			{"mixed", List{1, "a", nil, true, 2.0}, "[1, 'a', None, True, 2.0]"},
			// Python: repr([[1, 2], 'x']) == "[[1, 2], 'x']"
			{"nested", List{List{1, 2}, "x"}, "[[1, 2], 'x']"},
			// Python: repr(["it's"]) == '["it\'s"]'
			{"quote choice", List{"it's"}, `["it's"]`},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.list.String(); got != tc.want {
					t.Fatalf("String = %q, want %q", got, tc.want)
				}
			})
		}
	})
}

func TestErrors(t *testing.T) {
	cases := []struct {
		name string
		run  func() error
	}{
		// Python: [].pop() raises IndexError.
		{"pop empty", func() error {
			l := List{}
			_, err := l.Pop()
			return err
		}},
		// Python: [1].remove(2) raises ValueError.
		{"remove missing", func() error {
			l := List{1}
			return l.Remove(2)
		}},
		// Python: [1].index(2) raises ValueError.
		{"index missing", func() error {
			_, err := (List{1}).Index(2)
			return err
		}},
		// Python: [].extend(1) raises TypeError.
		{"extend not iterable", func() error {
			l := List{}
			return l.Extend(1)
		}},
		{"extend nil", func() error {
			l := List{}
			return l.Extend(nil)
		}},
		// Python: [].pop(0) raises IndexError.
		{"pop at empty", func() error {
			l := List{}
			_, err := l.PopAt(0)
			return err
		}},
		// Python: [1].pop(2) raises IndexError.
		{"pop at out of range", func() error {
			l := List{1}
			_, err := l.PopAt(2)
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.run(); err == nil {
				t.Fatal("error = nil, want error")
			}
		})
	}
}
