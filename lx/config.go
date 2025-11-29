package lx

// Options holds CLI-level configuration before effective values are derived.
type Options struct {
	Files []string

	Head  int
	Tail  int
	NBoth int

	HeadSet bool
	TailSet bool
	NSet    bool

	PrefixDelimiter  string
	PostfixDelimiter string
	LineNumbers      bool
}

// Effective derives a fully configured Runner from the options, applying
// -n / --head / --tail override rules.
func (o Options) Effective() Runner {
	effHead := o.Head
	effTail := o.Tail

	// -n sets both head and tail, unless overridden by explicit --head/-h or --tail/-t.
	if o.NSet {
		if !o.HeadSet {
			effHead = o.NBoth
		}
		if !o.TailSet {
			effTail = o.NBoth
		}
	}

	return NewRunner(
		effHead,
		effTail,
		o.PrefixDelimiter,
		o.PostfixDelimiter,
		o.LineNumbers,
	)
}

