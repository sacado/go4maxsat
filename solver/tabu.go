package solver

// A local-search SAT solver based on tabu list

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Tabu struct {
	problem      Problem
	model        Model
	nbUnsat      int // Current # of UNSAT clauses
	weightUnsat  int // Total weight of unsat clauses
	nbFlips      int // Total # of vars flipped
	tabuLength   int
	changed      []int   // For each var, last time it was flipped
	unsatClauses []int   // Indexes of currently unsat clauses
	whereFalse   []int   // Where each clause is in unsatClauses
	appears      appears // Where each literal appears
	nbTrue       []int   // Nb of true lits in each clause
	breakCount   []int   // For each var, total weight of broken clauses
}

func NewTabu(pb Problem, tabuLength int) *Tabu {
	model := randomModel(pb)
	appears := newAppears(pb)
	nbUnsat := 0
	weightUnsat := 0
	nbFlips := pb.NbVars // All vars are flipped on first random assignment
	changed := make([]int, pb.NbVars)
	unsatClauses := make([]int, len(pb.Clauses))
	whereFalse := make([]int, len(pb.Clauses))
	nbTrue := make([]int, len(pb.Clauses))
	breakCount := make([]int, pb.NbVars)

	for i, clause := range pb.Clauses {
		var unit Var // A var that could be the only true lit in the clause
		for _, lit := range clause {
			v := lit.Var()

			if lit.Positive() == model[v] {
				unit = v
				nbTrue[i]++
			}
		}

		if nbTrue[i] == 0 {
			unsatClauses[nbUnsat] = i
			whereFalse[i] = nbUnsat
			nbUnsat++
			weightUnsat += pb.Weights[i]
		} else if nbTrue[i] == 1 {
			// Flipping that var would make the clause unsat
			breakCount[unit] += pb.Weights[i]
		}
	}

	return &Tabu{
		pb,
		model,
		nbUnsat,
		weightUnsat,
		nbFlips,
		tabuLength,
		changed,
		unsatClauses,
		whereFalse,
		appears,
		nbTrue,
		breakCount,
	}
}

func (s *Tabu) Score() int {
	return s.weightUnsat
}

// Picks and returns a variable to flip
// Vars are chosen from a random unsat clause
// If all lits from the clause are tabu,
// No var is return and ok is false.
// Else, ok is true and a var is returned
func (s *Tabu) pick() (res Var, ok bool) {
	idx := s.unsatClauses[rand.Intn(s.nbUnsat)]
	cl := s.problem.Clauses[idx]
	top := make([]Var, len(cl))
	nbTop := 0
	var best int

	for i := 0; i < len(cl); i++ {
		v := cl[i].Var()
		breakCount := s.breakCount[v]

		if breakCount == 0 {
			if best > 0 {
				best = 0
				nbTop = 1
				top[0] = v
			} else {
				top[nbTop] = v
				nbTop++
			}
		} else if s.tabuLength < s.nbFlips-s.changed[cl[i].Var()] {
			if nbTop == 0 || breakCount < best {
				best = breakCount
				nbTop = 1
				top[0] = v
			} else if breakCount == best {
				top[nbTop] = v
				nbTop++
			}
		}
	}

	if nbTop == 0 {
		return res, false
	} else if nbTop == 1 {
		return top[0], true
	} else {
		return top[rand.Intn(nbTop)], true
	}
}

func (s *Tabu) flip(v Var) {
	s.nbFlips++
	s.changed[v] = s.nbFlips
	s.model[v] = !s.model[v]
	newLit := v.Lit(!s.model[v])
	oldLit := newLit.Negation()
	newAppears := s.appears.get(newLit)
	oldAppears := s.appears.get(oldLit)

	for _, idx := range newAppears {
		s.nbTrue[idx]++

		if s.nbTrue[idx] == 1 {
			// Clause was broken, it is now sat but will be broken by a flip of v
			s.breakCount[v] += s.problem.Weights[idx]
			s.nbUnsat--
			s.weightUnsat -= s.problem.Weights[idx]
			s.unsatClauses[s.whereFalse[idx]] = s.unsatClauses[s.nbUnsat]
			s.whereFalse[s.unsatClauses[s.nbUnsat]] = s.whereFalse[idx]
		} else if s.nbTrue[idx] == 2 {
			// Another var would have made the clause broken
			// Find it, decrement its breakcount
			clause := s.problem.Clauses[idx]

			for _, lit := range clause {
				vi := lit.Var()

				if vi != v && lit.Positive() == s.model[vi] {
					s.breakCount[vi] -= s.problem.Weights[idx]
					break
				}
			}
		}
	}

	for _, idx := range oldAppears {
		s.nbTrue[idx]--

		if s.nbTrue[idx] == 0 {
			s.unsatClauses[s.nbUnsat] = idx
			s.whereFalse[idx] = s.nbUnsat
			s.nbUnsat++
			s.weightUnsat += s.problem.Weights[idx]
			// That clause cannot be broken : it *is* broken
			s.breakCount[v] -= s.problem.Weights[idx]
		} else if s.nbTrue[idx] == 1 {
			// Clause will be broken if last lit is flipped
			// Find it and update breakCount
			clause := s.problem.Clauses[idx]

			for _, lit := range clause {
				vi := lit.Var()

				if s.model[vi] == lit.Positive() {
					// It is the only positive lit in the clause
					s.breakCount[vi] += s.problem.Weights[idx]
					break
				}
			}
		}
	}
}

func (s *Tabu) Solve(nbTries int) Model {
	for i := 0; i < nbTries && s.nbUnsat > 0; i++ {
		v, ok := s.pick()

		if ok {
			s.flip(v)
		}
	}

	return s.model
}
