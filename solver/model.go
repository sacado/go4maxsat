package solver

import "math/rand"

type Model []bool

// Returns 1 if the given clause is satisfied, 0 else
func (m Model) satisfied(clause Clause) int {
	for _, lit := range clause {
		v := lit.Var()

		if lit.Positive() == m[v] {
			return 1
		}
	}

	return 0
}

func (m Model) nbSatisfied(pb Problem) int {
	nb := 0

	for _, clause := range pb.Clauses {
		nb += m.satisfied(clause)
	}

	return nb
}

// Returns 1 if the given clause would be satisfied with vFlip flipped, 0 else
func (m Model) satisfiedIfFlipped(clause Clause, vFlip Var) int {
	for _, lit := range clause {
		v := lit.Var()
		assign := m[v]

		if v == vFlip {
			assign = !assign
		}

		if lit.Positive() == assign {
			return 1
		}
	}

	return 0
}

// Returns the list of all unsat clauses
func (m Model) unsatClauses(pb Problem) []Clause {
	unsat := make([]Clause, 0, len(pb.Clauses))

	for _, clause := range pb.Clauses {
		sat := false

		for _, lit := range clause {
			if lit.Positive() == m[lit.Var()] {
				sat = true
				break
			}
		}

		if !sat {
			unsat = append(unsat, clause)
		}
	}

	return unsat
}

func randomModel(pb Problem) Model {
	m := make(Model, pb.NbVars)

	for i := range m {
		m[i] = rand.Intn(2) == 0
	}

	return m
}

// Performs nbMutations random flips in m
// Mutations can happen twice on the same item
func (m Model) mutate(nbMutations int) {
	for i := 0; i < nbMutations; i++ {
		idx := rand.Intn(len(m))
		m[idx] = !m[idx]
	}
}

func (m Model) copy() Model {
	m2 := make(Model, len(m))
	copy(m2, m)

	return m2
}
