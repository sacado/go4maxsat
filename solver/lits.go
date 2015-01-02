package solver

// Lit is a literal, i.e. a var with a sign indicating negation
type Lit int

// Var is a variable, i.e. a positive int
type Var int

// Clause is a disjunction of lits
type Clause []Lit

// Problem is a problem, i.e. a # of vars,
// a weight for each clause and clauses
type Problem struct {
	NbVars  int
	Weights []int // Weight of each clause, 1 by default
	Clauses []Clause
}

// Returns the var associated with l
func (l Lit) Var() Var {
	if l < 0 {
		return Var(-l - 1)
	} else {
		return Var(l - 1)
	}
}

// True if l is positive
func (l Lit) Positive() bool {
	return l > 0
}

// Returns the negation of l
func (l Lit) Negation() Lit {
	return -l
}

// Returns a lit associated with the var
// if 'signed', the lit will be negative
// else it will be positive
func (v Var) Lit(signed bool) Lit {
	if signed {
		return -Lit(v) - 1
	} else {
		return Lit(v) + 1
	}
}

func (clause Clause) contains(v Var) bool {
	for _, lit := range clause {
		vi := lit.Var()

		if vi == v {
			return true
		}
	}

	return false
}
