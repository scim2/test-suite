package spec

// All contains every cataloged requirement across all three RFCs.
var All []Requirement

func init() {
	All = make([]Requirement, 0, len(RFC7642)+len(RFC7643)+len(RFC7644))
	All = append(All, RFC7642...)
	All = append(All, RFC7643...)
	All = append(All, RFC7644...)
}

// ByFeature returns all requirements that belong to the given feature.
func ByFeature(f Feature) []Requirement {
	var out []Requirement
	for _, r := range All {
		if r.Feature == f {
			out = append(out, r)
		}
	}
	return out
}

// ByID returns the requirement with the given ID, or nil if not found.
func ByID(id string) *Requirement {
	for i := range All {
		if All[i].ID == id {
			return &All[i]
		}
	}
	return nil
}

// ByLevel returns all requirements at the given compliance level.
func ByLevel(l Level) []Requirement {
	var out []Requirement
	for _, r := range All {
		if r.Level == l {
			out = append(out, r)
		}
	}
	return out
}

// concat joins multiple requirement slices into one.
func concat(slices ...[]Requirement) []Requirement {
	n := 0
	for _, s := range slices {
		n += len(s)
	}
	out := make([]Requirement, 0, n)
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}
