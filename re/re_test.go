package re

import (
	"slices"
	"testing"
)

func TestMatchSearch(t *testing.T) {
	ok, err := Match(`\d+`, "123abc")
	if err != nil || !ok {
		t.Fatalf("Match = (%v,%v)", ok, err)
	}
	ok, err = Match(`\d+`, "abc123")
	if err != nil || ok {
		t.Fatalf("Match should be anchored = (%v,%v)", ok, err)
	}
	ok, err = FullMatch(`\d+`, "123")
	if err != nil || !ok {
		t.Fatalf("FullMatch = (%v,%v)", ok, err)
	}
	ok, err = FullMatch(`\d+`, "123a")
	if err != nil || ok {
		t.Fatalf("FullMatch extra = (%v,%v)", ok, err)
	}
	ok, err = Search(`\d+`, "abc123")
	if err != nil || !ok {
		t.Fatalf("Search = (%v,%v)", ok, err)
	}
	if _, err := Match(`(`, "x"); err == nil {
		t.Fatal("invalid pattern should error")
	}
}

func TestFindAllSubSplit(t *testing.T) {
	all, err := FindAll(`\d+`, "a1b22c333")
	if err != nil || !slices.Equal(all, []string{"1", "22", "333"}) {
		t.Fatalf("FindAll = (%v,%v)", all, err)
	}
	match, ok, err := FindString(`\d+`, "a1b22")
	if err != nil || !ok || match != "1" {
		t.Fatalf("FindString = (%q,%v,%v)", match, ok, err)
	}
	sub, err := Sub(`(\w+)\s+(\w+)`, `${2} ${1}`, "hello world", 1)
	if err != nil || sub != "world hello" {
		t.Fatalf("Sub = (%q,%v)", sub, err)
	}
	parts, err := Split(`\s+`, "a b  c", 1)
	if err != nil || !slices.Equal(parts, []string{"a", "b  c"}) {
		t.Fatalf("Split = (%v,%v)", parts, err)
	}
}

func TestCompiled(t *testing.T) {
	compiled, err := Compile(`[a-z]+`)
	if err != nil {
		t.Fatal(err)
	}
	if !compiled.SearchString("123abc") {
		t.Fatal("SearchString failed")
	}
	if got := compiled.FindAllString("ab12cd", -1); !slices.Equal(got, []string{"ab", "cd"}) {
		t.Fatalf("FindAllString = %v", got)
	}
	if got := compiled.Sub("X", "ab12cd", 1); got != "X12cd" {
		t.Fatalf("Sub = %q", got)
	}
	if got := compiled.Split("ab12cd", -1); !slices.Equal(got, []string{"", "12", ""}) {
		t.Fatalf("Split = %v", got)
	}
}

func TestQuoteMeta(t *testing.T) {
	quoted := QuoteMeta("a+b")
	if quoted != `a\+b` {
		t.Fatalf("QuoteMeta = %q", quoted)
	}
}
