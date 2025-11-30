package lx

import "testing"

func TestOptionsEffective_NoN_UsesHeadTail(t *testing.T) {
	opts := Options{
		Head:  2,
		Tail:  3,
		NBoth: 5,

		HeadSet: true,
		TailSet: true,
		NSet:    false,
	}

	r := opts.Effective()

	if r.Head != 2 || r.Tail != 3 {
		t.Fatalf("Effective() Head/Tail = (%d,%d), want (2,3)", r.Head, r.Tail)
	}
}

func TestOptionsEffective_NOnly_SplitsTotalNBetweenHeadAndTail(t *testing.T) {
	opts := Options{
		Head:  0,
		Tail:  0,
		NBoth: 5,

		HeadSet: false,
		TailSet: false,
		NSet:    true,
	}

	r := opts.Effective()

	if r.Head != 3 || r.Tail != 2 {
		t.Fatalf("Effective() Head/Tail = (%d,%d), want (3,2)", r.Head, r.Tail)
	}
}

func TestOptionsEffective_NWithHeadOverride(t *testing.T) {
	opts := Options{
		Head:  2,
		Tail:  0,
		NBoth: 5,

		HeadSet: true,  // explicit head
		TailSet: false, // tail not explicitly set
		NSet:    true,
	}

	r := opts.Effective()

	// Total N = 5, head explicitly 2 => tail should be 3.
	if r.Head != 2 || r.Tail != 3 {
		t.Fatalf("Effective() Head/Tail = (%d,%d), want (2,3)", r.Head, r.Tail)
	}
}

func TestOptionsEffective_NWithTailOverride(t *testing.T) {
	opts := Options{
		Head:  0,
		Tail:  7,
		NBoth: 5,

		HeadSet: false, // head not explicitly set
		TailSet: true,  // explicit tail
		NSet:    true,
	}

	r := opts.Effective()

	// Total N = 5, tail explicitly 7 => clamped to 5, head = 0.
	if r.Head != 0 || r.Tail != 5 {
		t.Fatalf("Effective() Head/Tail = (%d,%d), want (0,5)", r.Head, r.Tail)
	}
}

func TestOptionsEffective_NWithBothOverrides(t *testing.T) {
	opts := Options{
		Head:  2,
		Tail:  7,
		NBoth: 5,

		HeadSet: true,
		TailSet: true,
		NSet:    true,
	}

	r := opts.Effective()

	// Both head and tail explicitly set; -n should not override them.
	if r.Head != 2 || r.Tail != 7 {
		t.Fatalf("Effective() Head/Tail = (%d,%d), want (2,7)", r.Head, r.Tail)
	}
}

func TestOptionsEffective_PrefixPostfixAndLineNumbers(t *testing.T) {
	opts := Options{
		Head: 1,
		Tail: 0,

		PrefixDelimiter:  "PRE\n",
		PostfixDelimiter: "POST\n",
		LineNumbers:      true,

		HeadSet: true,
	}

	r := opts.Effective()

	if r.Head != 1 || r.Tail != 0 {
		t.Fatalf("Effective() Head/Tail = (%d,%d), want (1,0)", r.Head, r.Tail)
	}
	if r.PrefixDelimiter != "PRE\n" {
		t.Fatalf("Effective() PrefixDelimiter = %q, want %q", r.PrefixDelimiter, "PRE\n")
	}
	if r.PostfixDelimiter != "POST\n" {
		t.Fatalf("Effective() PostfixDelimiter = %q, want %q", r.PostfixDelimiter, "POST\n")
	}
	if !r.LineNumbers {
		t.Fatalf("Effective() LineNumbers = false, want true")
	}
}
