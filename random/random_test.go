package random

import (
	"slices"
	"testing"
)

func TestSeedDeterministic(t *testing.T) {
	Seed(42)
	a, b := Float(), Float()
	Seed(42)
	c, d := Float(), Float()
	if a != c || b != d {
		t.Fatalf("Seed not deterministic: %v %v vs %v %v", a, b, c, d)
	}
	if a < 0 || a >= 1 || b < 0 || b >= 1 {
		t.Fatalf("Float out of range: %v %v", a, b)
	}
}

func TestRandrangeAndRandint(t *testing.T) {
	Seed(1)
	for i := 0; i < 100; i++ {
		value, err := Randrange(1, 10, 2)
		if err != nil || value < 1 || value >= 10 || (value-1)%2 != 0 {
			t.Fatalf("Randrange = (%d,%v)", value, err)
		}
	}
	if _, err := Randrange(5, 5); err == nil {
		t.Fatal("empty range should error")
	}
	if _, err := Randrange(0, 10, 0); err == nil {
		t.Fatal("zero step should error")
	}
	if got := Randint(3, 3); got != 3 {
		t.Fatalf("Randint(3,3) = %d", got)
	}
}

func TestChoiceShuffleSample(t *testing.T) {
	Seed(2)
	items := []int{1, 2, 3, 4}
	if value, ok := Choice(items); !ok || !slices.Contains(items, value) {
		t.Fatalf("Choice = (%d,%v)", value, ok)
	}
	if _, ok := Choice([]int{}); ok {
		t.Fatal("Choice empty should be false")
	}

	shuffled := append([]int(nil), items...)
	Shuffle(shuffled)
	if len(shuffled) != len(items) {
		t.Fatalf("Shuffle len = %d", len(shuffled))
	}
	slices.Sort(shuffled)
	if !slices.Equal(shuffled, []int{1, 2, 3, 4}) {
		t.Fatalf("Shuffle changed elements: %v", shuffled)
	}

	sample, err := Sample(items, 2)
	if err != nil || len(sample) != 2 {
		t.Fatalf("Sample = (%v,%v)", sample, err)
	}
	if _, err := Sample(items, 5); err == nil {
		t.Fatal("oversized sample should error")
	}
	choices, err := Choices(items, 3)
	if err != nil || len(choices) != 3 {
		t.Fatalf("Choices = (%v,%v)", choices, err)
	}
}
