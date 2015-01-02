package solver

// An 'appears' is the list, for each lit, of indexes of
// the clauses it appears in.
// For a var v, the position of literal v is 2*v and of -v is 2*v + 1
type appears [][]int

func newAppears(pb Problem) appears {
	res := make(appears, 2*pb.NbVars)

	for idx, clause := range pb.Clauses {
		for _, lit := range clause {
			index := res.indexFor(lit)

			if res[index] == nil {
				res[index] = []int{idx}
			} else {
				res[index] = append(res[index], idx)
			}
		}
	}

	return res
}

func (appears appears) get(lit Lit) []int {
	return appears[appears.indexFor(lit)]
}

func (appears appears) indexFor(lit Lit) int {
	v := lit.Var()
	offset := 0

	if !lit.Positive() {
		offset = 1
	}

	return 2*int(v) + offset
}
