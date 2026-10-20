package py

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

type reflectUser struct {
	Name   string `json:"name"`
	Age    int    `json:"age,omitempty"`
	hidden string
}

func (u reflectUser) Greet() string { return "hi" }

func TestGetAttr(t *testing.T) {
	user := reflectUser{Name: "Tom", Age: 18}
	if got, err := GetAttr(user, "name"); err != nil || got != "Tom" {
		t.Fatalf("GetAttr name = (%v,%v)", got, err)
	}
	if got, err := GetAttr(&user, "Age"); err != nil || got != 18 {
		t.Fatalf("GetAttr Age = (%v,%v)", got, err)
	}
	if got, err := GetAttr(user, "Greet"); err != nil || !Callable(got) {
		t.Fatalf("GetAttr Greet = (%v,%v)", got, err)
	}
	if got, err := GetAttr(user, "missing", "fallback"); err != nil || got != "fallback" {
		t.Fatalf("GetAttr default = (%v,%v)", got, err)
	}
	if _, err := GetAttr(user, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetAttr missing err = %v", err)
	}
	if !HasAttr(user, "name") || HasAttr(user, "missing") {
		t.Fatal("HasAttr mismatch")
	}
}

func TestSetAttr(t *testing.T) {
	user := reflectUser{Name: "Tom", Age: 18}
	if err := SetAttr(&user, "name", "Bob"); err != nil || user.Name != "Bob" {
		t.Fatalf("SetAttr name = (%v,%v)", user.Name, err)
	}
	if err := SetAttr(&user, "Age", 20); err != nil || user.Age != 20 {
		t.Fatalf("SetAttr Age = (%v,%v)", user.Age, err)
	}
	if err := SetAttr(&user, "Age", "x"); !errors.Is(err, ErrType) {
		t.Fatalf("SetAttr wrong type err = %v", err)
	}
	if err := SetAttr(user, "Age", 30); !errors.Is(err, ErrType) {
		t.Fatalf("SetAttr non-pointer err = %v", err)
	}

	m := map[string]any{"a": 1}
	if err := SetAttr(m, "b", 2); err != nil || m["b"] != 2 {
		t.Fatalf("SetAttr map = (%v,%v)", m, err)
	}
	if err := SetAttr(m, "b", nil); err != nil {
		t.Fatalf("SetAttr map nil = %v", err)
	}
}

func TestDelAttr(t *testing.T) {
	m := map[string]int{"a": 1}
	if err := DelAttr(m, "a"); err != nil || len(m) != 0 {
		t.Fatalf("DelAttr = (%v,%v)", m, err)
	}
	if err := DelAttr(m, "a"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("DelAttr missing err = %v", err)
	}
	if err := DelAttr(reflectUser{}, "Name"); !errors.Is(err, ErrType) {
		t.Fatalf("DelAttr struct err = %v", err)
	}
}

func TestDir(t *testing.T) {
	user := reflectUser{Name: "Tom"}
	got := Dir(user)
	want := []string{"Age", "Greet", "Name"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Dir = %v, want %v", got, want)
	}
	if got := Dir(map[string]int{"b": 2, "a": 1}); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("Dir map = %v", got)
	}
	if got := Dir(nil); len(got) != 0 {
		t.Fatalf("Dir nil = %v", got)
	}
}

func TestCallableAndTypes(t *testing.T) {
	if !Callable(func() {}) || Callable(1) {
		t.Fatal("Callable mismatch")
	}
	if !IsInstance(1, 0) || IsInstance("x", 0) {
		t.Fatal("IsInstance mismatch")
	}
	if !IsInstance(1, reflect.TypeOf(0)) {
		t.Fatal("IsInstance reflect.Type failed")
	}
	if IsInstance(nil, 0) {
		t.Fatal("IsInstance(nil, int) should be false")
	}

	if !IsSubclass((*bytes.Buffer)(nil), reflect.TypeOf((*io.Writer)(nil)).Elem()) {
		t.Fatal("IsSubclass interface failed")
	}
	if !IsSubclass((*bytes.Buffer)(nil), (*bytes.Buffer)(nil)) {
		t.Fatal("IsSubclass same type failed")
	}
	if IsSubclass((*bytes.Buffer)(nil), (*bytes.Reader)(nil)) {
		t.Fatal("IsSubclass unrelated types returned true")
	}
}

func TestTypeOfAndVars(t *testing.T) {
	if TypeOf(nil) != nil {
		t.Fatal("TypeOf(nil) != nil")
	}
	if TypeOf(1).Kind() != reflect.Int {
		t.Fatal("TypeOf(1) != int")
	}

	user := reflectUser{Name: "Tom", Age: 18}
	vars, err := Vars(&user)
	if err != nil {
		t.Fatal(err)
	}
	if vars["name"] != "Tom" || vars["age"] != 18 {
		t.Fatalf("Vars = %v", vars)
	}
	if _, ok := vars["hidden"]; ok {
		t.Fatal("Vars included unexported field")
	}

	m := map[string]any{"a": 1}
	copied, err := Vars(m)
	if err != nil || copied["a"] != 1 {
		t.Fatalf("Vars map = (%v,%v)", copied, err)
	}
	copied["a"] = 2
	if m["a"] != 1 {
		t.Fatal("Vars map did not copy")
	}
}
