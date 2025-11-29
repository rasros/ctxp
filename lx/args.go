package lx

// NormalizeArgs rewrites "-n2" / "-t10" / "-h5" into ["-n","2"] / ["-t","10"] / ["-h","5"]
// so that urfave/cli/v3 parses them as int flags.
func NormalizeArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	out := make([]string, 0, len(args)+4)
	for _, a := range args {
		if len(a) > 2 && a[0] == '-' && (a[1] == 'n' || a[1] == 't' || a[1] == 'h') {
			digits := a[2:]
			isDigits := true
			for _, ch := range digits {
				if ch < '0' || ch > '9' {
					isDigits = false
					break
				}
			}
			if isDigits {
				out = append(out, a[:2], digits)
				continue
			}
		}
		out = append(out, a)
	}

	return out
}

