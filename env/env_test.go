package env

import (
	"slices"
	"strings"
	"testing"
)

func TestGetSetUnset(t *testing.T) {
	const key = "GOTOOLS_ENV_TEST_KEY"
	Unset(key)
	if got := Get(key); got != "" {
		t.Fatalf("Get missing = %q", got)
	}
	if got := Get(key, "fallback"); got != "fallback" {
		t.Fatalf("Get default = %q", got)
	}
	if Has(key) {
		t.Fatal("Has should be false")
	}
	if err := Set(key, "value"); err != nil {
		t.Fatal(err)
	}
	if got := Get(key); got != "value" {
		t.Fatalf("Get = %q", got)
	}
	if value, ok := Lookup(key); !ok || value != "value" {
		t.Fatalf("Lookup = (%q,%v)", value, ok)
	}
	if !Has(key) {
		t.Fatal("Has should be true")
	}
	if err := Unset(key); err != nil {
		t.Fatal(err)
	}
	if Has(key) {
		t.Fatal("Unset failed")
	}
}

func TestAllKeysExpand(t *testing.T) {
	const key = "GOTOOLS_ENV_TEST_ALL"
	if err := Set(key, "v"); err != nil {
		t.Fatal(err)
	}
	defer Unset(key)

	if !slices.Contains(Keys(), key) {
		t.Fatal("Keys missing test key")
	}
	all := All()
	if !slices.Contains(all, key+"=v") {
		t.Fatalf("All missing %s=v", key)
	}
	if !sortStrings(all) {
		t.Fatal("All should be sorted")
	}
	if got := Expand("x=${" + key + "}"); got != "x=v" {
		t.Fatalf("Expand = %q", got)
	}
	if got := Expand("$" + key); got != "v" {
		t.Fatalf("Expand short = %q", got)
	}
}

func sortStrings(values []string) bool {
	for i := 1; i < len(values); i++ {
		if strings.Compare(values[i-1], values[i]) > 0 {
			return false
		}
	}
	return true
}
