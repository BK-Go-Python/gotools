package dict

import (
	"reflect"
	"testing"
)

// assertDict compares a Dict with the expected Python-semantics mapping.
func assertDict(t *testing.T, got Dict, want map[string]interface{}) {
	t.Helper()
	if !reflect.DeepEqual(map[string]interface{}(got), want) {
		t.Fatalf("got %s (%#v), want %#v", got.String(), map[string]interface{}(got), want)
	}
}

func TestNewFrom(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		// Python: {} == {} and len({}) == 0
		d := New()
		if d == nil {
			t.Fatal("New() returned nil")
		}
		if d.Len() != 0 {
			t.Fatalf("Len() = %d, want 0", d.Len())
		}
		if got := d.String(); got != "{}" {
			t.Fatalf("String() = %q, want %q", got, "{}")
		}
		d.Set("a", 1)
		if !d.Has("a") {
			t.Fatal("New() dictionary is not writable")
		}
	})

	t.Run("From", func(t *testing.T) {
		cases := []struct {
			name string
			in   map[string]interface{}
			want map[string]interface{}
		}{
			{"empty", map[string]interface{}{}, map[string]interface{}{}},
			{"nil", nil, map[string]interface{}{}},
			{"values", map[string]interface{}{"a": 1, "b": "x"}, map[string]interface{}{"a": 1, "b": "x"}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := From(tc.in)
				if got == nil {
					t.Fatal("From returned nil")
				}
				assertDict(t, got, tc.want)
			})
		}

		t.Run("shallow copy", func(t *testing.T) {
			src := map[string]interface{}{"a": 1}
			got := From(src)
			src["a"] = 2
			if got.Get("a") != 1 {
				t.Fatalf("From did not copy: Get(a) = %#v, want 1", got.Get("a"))
			}
			got.Set("b", 3)
			if _, ok := src["b"]; ok {
				t.Fatal("From result shares top-level storage with the source map")
			}
		})
	})
}

func TestGetSetSetDefaultHas(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		cases := []struct {
			name string
			dict Dict
			key  string
			def  []interface{}
			want interface{}
		}{
			// Python: {'a': 1}.get('a') == 1
			{"present", Dict{"a": 1}, "a", nil, 1},
			// Python: {}.get('a') is None
			{"missing no default", Dict{}, "a", nil, nil},
			// Python: {}.get('a', 42) == 42
			{"missing default", Dict{}, "a", []interface{}{42}, 42},
			// Python: {'a': 1}.get('a', 42) == 1
			{"present ignores default", Dict{"a": 1}, "a", []interface{}{42}, 1},
			// Python: {'a': None}.get('a', 42) is None
			{"present nil ignores default", Dict{"a": nil}, "a", []interface{}{42}, nil},
			{"missing nil default", Dict{}, "a", []interface{}{nil}, nil},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.dict.Get(tc.key, tc.def...); !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("Get = %#v, want %#v", got, tc.want)
				}
			})
		}
	})

	t.Run("Set", func(t *testing.T) {
		// Python: d['a'] = 1; d['b'] = 2
		d := New()
		if ret := d.Set("a", 1); ret == nil {
			t.Fatal("Set returned nil")
		}
		d.Set("b", 2).Set("a", 3)
		assertDict(t, d, map[string]interface{}{"a": 3, "b": 2})
	})

	t.Run("SetDefault", func(t *testing.T) {
		cases := []struct {
			name  string
			start Dict
			key   string
			def   interface{}
			want  interface{}
			after map[string]interface{}
		}{
			// Python: {}.setdefault('a', 1) == 1 and inserts a=1.
			{"missing inserts", Dict{}, "a", 1, 1, map[string]interface{}{"a": 1}},
			// Python: {'a': 1}.setdefault('a', 9) == 1, value unchanged.
			{"present keeps", Dict{"a": 1}, "a", 9, 1, map[string]interface{}{"a": 1}},
			// Python: {}.setdefault('a') == None and inserts a=None.
			{"missing nil", Dict{}, "a", nil, nil, map[string]interface{}{"a": nil}},
			// Python: {'a': None}.setdefault('a', 9) is None, value unchanged.
			{"present nil keeps", Dict{"a": nil}, "a", 9, nil, map[string]interface{}{"a": nil}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.start.SetDefault(tc.key, tc.def); !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("SetDefault = %#v, want %#v", got, tc.want)
				}
				assertDict(t, tc.start, tc.after)
			})
		}
	})

	t.Run("Has", func(t *testing.T) {
		d := Dict{"a": nil}
		if !d.Has("a") {
			t.Fatal("Has(a) = false, want true for nil value")
		}
		if d.Has("b") {
			t.Fatal("Has(b) = true, want false")
		}
		if New().Has("a") {
			t.Fatal("Has(a) = true on empty dict")
		}
	})
}

func TestDeletePop(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		cases := []struct {
			name  string
			start Dict
			key   string
			want  bool
			after map[string]interface{}
		}{
			// Python: del d['a'] on {'a': 1} succeeds.
			{"present", Dict{"a": 1, "b": 2}, "a", true, map[string]interface{}{"b": 2}},
			// Python: del d['a'] on {} raises KeyError.
			{"missing", Dict{"a": 1}, "x", false, map[string]interface{}{"a": 1}},
			{"empty", Dict{}, "a", false, map[string]interface{}{}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.start.Delete(tc.key); got != tc.want {
					t.Fatalf("Delete = %v, want %v", got, tc.want)
				}
				assertDict(t, tc.start, tc.after)
			})
		}
	})

	t.Run("Pop", func(t *testing.T) {
		cases := []struct {
			name  string
			start Dict
			key   string
			def   []interface{}
			want  interface{}
			after map[string]interface{}
		}{
			// Python: {'a': 1}.pop('a') == 1 and removes the key.
			{"present", Dict{"a": 1, "b": 2}, "a", nil, 1, map[string]interface{}{"b": 2}},
			// Python: {'a': None}.pop('a') is None and removes the key.
			{"present nil", Dict{"a": nil}, "a", nil, nil, map[string]interface{}{}},
			// Python: {}.pop('a', 42) == 42, dictionary unchanged.
			{"missing default", Dict{"b": 2}, "a", []interface{}{42}, 42, map[string]interface{}{"b": 2}},
			// Python raises KeyError without a default; gotools documents a nil
			// return instead.
			{"missing no default", Dict{"b": 2}, "a", nil, nil, map[string]interface{}{"b": 2}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.start.Pop(tc.key, tc.def...); !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("Pop = %#v, want %#v", got, tc.want)
				}
				assertDict(t, tc.start, tc.after)
			})
		}
	})
}

func TestPopItem(t *testing.T) {
	cases := []struct {
		name    string
		start   Dict
		want    Item
		rest    map[string]interface{}
		wantErr bool
	}{
		// Python 3.13 popitem() is LIFO; gotools picks the first key in sorted
		// order for deterministic Go map behaviour (documented).
		{"sorted first", Dict{"b": 2, "a": 1, "c": 3}, Item{Key: "a", Value: 1}, map[string]interface{}{"b": 2, "c": 3}, false},
		{"single", Dict{"x": 9}, Item{Key: "x", Value: 9}, map[string]interface{}{}, false},
		// Python: {}.popitem() raises KeyError.
		{"empty", Dict{}, Item{}, map[string]interface{}{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.start.PopItem()
			if tc.wantErr {
				if err == nil {
					t.Fatal("PopItem error = nil, want error")
				}
			} else if err != nil {
				t.Fatalf("PopItem unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("PopItem = %#v, want %#v", got, tc.want)
			}
			assertDict(t, tc.start, tc.rest)
		})
	}
}

func TestKeysValuesItemsClear(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		cases := []struct {
			name   string
			start  Dict
			keys   []string
			values []interface{}
			items  []Item
		}{
			// Python dicts keep insertion order; gotools sorts keys (documented).
			{
				"three",
				Dict{"b": 2, "a": 1, "c": 3},
				[]string{"a", "b", "c"},
				[]interface{}{1, 2, 3},
				[]Item{{"a", 1}, {"b", 2}, {"c", 3}},
			},
			{
				"empty",
				Dict{},
				[]string{},
				[]interface{}{},
				[]Item{},
			},
			{
				"nil value",
				Dict{"a": nil},
				[]string{"a"},
				[]interface{}{nil},
				[]Item{{"a", nil}},
			},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				keys := tc.start.Keys()
				values := tc.start.Values()
				items := tc.start.Items()
				if !reflect.DeepEqual(keys, tc.keys) {
					t.Fatalf("Keys = %#v, want %#v", keys, tc.keys)
				}
				if !reflect.DeepEqual(values, tc.values) {
					t.Fatalf("Values = %#v, want %#v", values, tc.values)
				}
				if !reflect.DeepEqual(items, tc.items) {
					t.Fatalf("Items = %#v, want %#v", items, tc.items)
				}
				// Keys, Values and Items must describe the same order.
				for i, key := range keys {
					if items[i].Key != key || !reflect.DeepEqual(values[i], items[i].Value) {
						t.Fatalf("order mismatch at %d: key=%q item=%#v value=%#v", i, key, items[i], values[i])
					}
				}
			})
		}
	})

	t.Run("Clear", func(t *testing.T) {
		d := Dict{"b": 2, "a": 1}
		ret := d.Clear()
		if ret == nil {
			t.Fatal("Clear returned nil")
		}
		if d.Len() != 0 || len(d.Keys()) != 0 {
			t.Fatalf("after Clear: Len = %d, Keys = %#v", d.Len(), d.Keys())
		}
		if got := d.String(); got != "{}" {
			t.Fatalf("after Clear String = %q, want {}", got)
		}
		// Clearing an empty dictionary is a no-op.
		Dict{}.Clear()
	})
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name    string
		start   Dict
		other   interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		// Python: {'a': 1}.update({'b': 2}) -> {'a': 1, 'b': 2}
		{"Dict", Dict{"a": 1}, Dict{"b": 2}, map[string]interface{}{"a": 1, "b": 2}, false},
		{"Dict overwrite", Dict{"a": 1}, Dict{"a": 2}, map[string]interface{}{"a": 2}, false},
		{
			"map string interface",
			Dict{"a": 1},
			map[string]interface{}{"a": 9, "b": true},
			map[string]interface{}{"a": 9, "b": true},
			false,
		},
		// Python: d.update({'n': 7}) accepts any mapping with string keys.
		{"map string int", Dict{"a": 1}, map[string]int{"a": 2, "n": 7}, map[string]interface{}{"a": 2, "n": 7}, false},
		{"pointer to map", Dict{}, &map[string]int{"p": 9}, map[string]interface{}{"p": 9}, false},
		// Python: d.update(None) raises TypeError; gotools treats nil as a no-op.
		{"nil no-op", Dict{"a": 1}, nil, map[string]interface{}{"a": 1}, false},
		// Python: d.update('abc') raises ValueError; d.update(1) raises TypeError.
		{"string invalid", Dict{"a": 1}, "abc", map[string]interface{}{"a": 1}, true},
		{"int invalid", Dict{"a": 1}, 42, map[string]interface{}{"a": 1}, true},
		{"slice invalid", Dict{"a": 1}, []string{"x"}, map[string]interface{}{"a": 1}, true},
		{"non string keys invalid", Dict{"a": 1}, map[int]string{1: "x"}, map[string]interface{}{"a": 1}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.start.Update(tc.other)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Update(%#v) error = nil, want error", tc.other)
				}
			} else if err != nil {
				t.Fatalf("Update(%#v) unexpected error: %v", tc.other, err)
			}
			assertDict(t, tc.start, tc.want)
		})
	}
}

func TestMergeCopyLen(t *testing.T) {
	t.Run("Merge", func(t *testing.T) {
		// Python: {'a': 1, 'b': 2} | {'b': 20, 'c': 3} -> {'a': 1, 'b': 20, 'c': 3}
		orig := Dict{"a": 1, "b": 2}
		other1 := Dict{"b": 20, "c": 3}
		other2 := Dict{"d": 4}
		got := orig.Merge(other1, other2)
		assertDict(t, got, map[string]interface{}{"a": 1, "b": 20, "c": 3, "d": 4})
		assertDict(t, orig, map[string]interface{}{"a": 1, "b": 2})
		assertDict(t, other1, map[string]interface{}{"b": 20, "c": 3})
		assertDict(t, other2, map[string]interface{}{"d": 4})

		got.Set("a", 99)
		if orig.Get("a") != 1 {
			t.Fatal("Merge result aliases the original dictionary")
		}

		// Python: d | {} is a new copy of d.
		copyOnly := orig.Merge()
		copyOnly.Set("z", 9)
		if orig.Has("z") {
			t.Fatal("Merge with no arguments aliases the original dictionary")
		}
	})

	t.Run("Copy", func(t *testing.T) {
		orig := Dict{"a": 1, "b": []interface{}{1, 2}}
		cp := orig.Copy()
		cp.Set("a", 9)
		cp.Set("c", 3)
		if orig.Get("a") != 1 || orig.Has("c") {
			t.Fatal("Copy is not independent at the top level")
		}
		assertDict(t, cp, map[string]interface{}{"a": 9, "b": []interface{}{1, 2}, "c": 3})

		// Python: copy() is shallow; nested mutable values are shared.
		cp.Get("b").([]interface{})[0] = 99
		if got := orig.Get("b").([]interface{})[0]; got != 99 {
			t.Fatalf("Copy is deep: nested value = %#v, want shared mutation", got)
		}
	})

	t.Run("Len", func(t *testing.T) {
		cases := []struct {
			name string
			dict Dict
			want int
		}{
			{"empty", Dict{}, 0},
			{"one", Dict{"a": nil}, 1},
			{"three", Dict{"a": 1, "b": 2, "c": 3}, 3},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.dict.Len(); got != tc.want {
					t.Fatalf("Len = %d, want %d", got, tc.want)
				}
			})
		}
	})
}

func TestGetPath(t *testing.T) {
	base := func() Dict {
		return Dict{
			"user": map[string]interface{}{
				"name":    "Ada",
				"address": Dict{"city": "Paris"},
			},
			"tags":   []interface{}{"go", "python"},
			"matrix": []interface{}{[]interface{}{1, 2}, []interface{}{3, 4}},
			"none":   nil,
		}
	}
	cases := []struct {
		name string
		path string
		def  []interface{}
		want interface{}
	}{
		{"nested map", "user.name", nil, "Ada"},
		{"nested Dict", "user.address.city", nil, "Paris"},
		{"slice index", "tags.1", nil, "python"},
		{"nested slice index", "matrix.1.0", nil, 3},
		{"nil value", "none", []interface{}{"fallback"}, nil},
		// GetPath is a gotools extension; Python dicts have no dotted lookup.
		{"missing with default", "user.age", []interface{}{"unknown"}, "unknown"},
		{"missing without default", "user.age", nil, nil},
		{"missing root", "missing.x", []interface{}{0}, 0},
		{"slice out of range", "tags.2", []interface{}{"none"}, "none"},
		{"slice non numeric", "tags.x", []interface{}{"none"}, "none"},
		{"through scalar", "user.name.length", []interface{}{"none"}, "none"},
		{"empty path", "", []interface{}{"default"}, "default"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := base().GetPath(tc.path, tc.def...); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("GetPath(%q) = %#v, want %#v", tc.path, got, tc.want)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	cases := []struct {
		name    string
		dict    Dict
		want    string
		wantErr bool
	}{
		// Python: json.dumps({}) == '{}'
		{"empty", Dict{}, "{}", false},
		// Python: json.dumps({"a": 1, "b": 2}) uses spaces; gotools emits compact
		// JSON with sorted keys: {"a":1,"b":2}.
		{"simple", Dict{"b": 2, "a": 1}, `{"a":1,"b":2}`, false},
		{
			"nested",
			Dict{
				"user": map[string]interface{}{"name": "Ada", "age": 30},
				"tags": []interface{}{"go", "py"},
			},
			`{"tags":["go","py"],"user":{"age":30,"name":"Ada"}}`,
			false,
		},
		// Python: json.dumps({"html": "<b>&"}) does not escape HTML characters.
		{"html not escaped", Dict{"html": "<b>&"}, `{"html":"<b>&"}`, false},
		{"null value", Dict{"x": nil}, `{"x":null}`, false},
		// Python: json.dumps({"bad": set()}) raises TypeError.
		{"unsupported value", Dict{"bad": make(chan int)}, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.dict.ToJSON()
			if tc.wantErr {
				if err == nil {
					t.Fatal("ToJSON error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ToJSON unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("ToJSON = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestEqualString(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		cases := []struct {
			name string
			a    Dict
			b    Dict
			want bool
		}{
			// Python: {} == {} is True
			{"empty", Dict{}, Dict{}, true},
			{"same", Dict{"a": 1, "b": "x"}, Dict{"a": 1, "b": "x"}, true},
			// Python: {'a': 1} == {'a': 1.0} is True
			{"cross numeric", Dict{"a": 1}, Dict{"a": 1.0}, true},
			// Python: {'a': [1, 2]} == {'a': [1.0, 2.0]} is True
			{"nested cross numeric", Dict{"a": []interface{}{1, 2}}, Dict{"a": []interface{}{1.0, 2.0}}, true},
			{"different value", Dict{"a": 1}, Dict{"a": 2}, false},
			{"missing key", Dict{"a": 1, "b": 2}, Dict{"a": 1}, false},
			{"extra key", Dict{"a": 1}, Dict{"a": 1, "b": 2}, false},
			{"nil vs empty", Dict(nil), Dict{}, true},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.a.Equal(tc.b); got != tc.want {
					t.Fatalf("Equal = %v, want %v", got, tc.want)
				}
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		cases := []struct {
			name string
			dict Dict
			want string
		}{
			// Python: repr({}) == '{}'
			{"empty", Dict{}, "{}"},
			// Python: repr({'a': 1, 'b': 'x', 'c': None}) ==
			// "{'a': 1, 'b': 'x', 'c': None}"
			{"mixed", Dict{"a": 1, "b": "x", "c": nil}, `{'a': 1, 'b': 'x', 'c': None}`},
			// Python: repr({'f': 2.0, 'flag': True}) == "{'f': 2.0, 'flag': True}"
			{"bool float", Dict{"f": 2.0, "flag": true}, `{'f': 2.0, 'flag': True}`},
			// Python: repr({'d': {'z': 1}, 'l': [1, 'a']}) ==
			// "{'d': {'z': 1}, 'l': [1, 'a']}"
			{"nested", Dict{"d": Dict{"z": 1}, "l": []interface{}{1, "a"}}, `{'d': {'z': 1}, 'l': [1, 'a']}`},
			// Python: repr({'q': "it's"}) == '{"q": "it\'s"}'
			{"quotes", Dict{"q": "it's"}, `{'q': "it's"}`},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.dict.String(); got != tc.want {
					t.Fatalf("String = %q, want %q", got, tc.want)
				}
			})
		}
	})
}
