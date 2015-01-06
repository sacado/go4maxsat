package solver

// A local-search MAX-SAT solver based on simulated annealing

import (
	//"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type AnnealingSolver struct {
	problem      Problem
	model        Model
	nbUnsat      int     // Current # of UNSAT clauses
	weightUnsat  int     // Total weight of unsat clauses
	unsatClauses []int   // Indexes of currently unsat clauses
	whereFalse   []int   // Where each clause is in unsatClauses
	appears      appears // Where each literal appears
	nbTrue       []int   // Nb of true lits in each clause
	breakCount   []int   // For each var, total weight of broken clauses
	makeCount    []int   // For each var, total weight of clauses that would be true if flipped
}

func NewAnnealingSolver(pb Problem) *AnnealingSolver {
	model := randomModel(pb)
	appears := newAppears(pb)
	nbUnsat := 0
	weightUnsat := 0
	unsatClauses := make([]int, len(pb.Clauses))
	whereFalse := make([]int, len(pb.Clauses))
	nbTrue := make([]int, len(pb.Clauses))
	breakCount := make([]int, pb.NbVars)
	makeCount := make([]int, pb.NbVars)

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

			for _, lit := range clause {
				makeCount[lit.Var()] += pb.Weights[i]
			}
		} else if nbTrue[i] == 1 {
			// Flipping that var would make the clause unsat
			breakCount[unit] += pb.Weights[i]
		}
	}

	return &AnnealingSolver{
		pb,
		model,
		nbUnsat,
		weightUnsat,
		unsatClauses,
		whereFalse,
		appears,
		nbTrue,
		breakCount,
		makeCount,
	}
}

func (s *AnnealingSolver) Score() int {
	return s.weightUnsat
}

func (s *AnnealingSolver) flip(v Var) {
	s.model[v] = !s.model[v]
	newLit := v.Lit(!s.model[v])
	oldLit := newLit.Negation()

	for _, idx := range s.appears.get(newLit) {
		s.nbTrue[idx]++
		weight := s.problem.Weights[idx]

		if s.nbTrue[idx] == 1 {
			// Clause was broken, it is now sat but will be broken by a flip of v
			s.breakCount[v] += weight
			s.nbUnsat--
			s.weightUnsat -= weight
			s.unsatClauses[s.whereFalse[idx]] = s.unsatClauses[s.nbUnsat]
			s.whereFalse[s.unsatClauses[s.nbUnsat]] = s.whereFalse[idx]

			// Flipping lits won't make the clause true anymore
			for _, lit := range s.problem.Clauses[idx] {
				s.makeCount[lit.Var()] -= weight
			}
		} else if s.nbTrue[idx] == 2 {
			// Another var would have made the clause broken
			// Find it, decrement its breakcount
			for _, lit := range s.problem.Clauses[idx] {
				vi := lit.Var()

				if vi != v && lit.Positive() == s.model[vi] {
					s.breakCount[vi] -= weight
					break
				}
			}
		}
	}

	for _, idx := range s.appears.get(oldLit) {
		s.nbTrue[idx]--
		weight := s.problem.Weights[idx]

		if s.nbTrue[idx] == 0 {
			s.unsatClauses[s.nbUnsat] = idx
			s.whereFalse[idx] = s.nbUnsat
			s.nbUnsat++
			s.weightUnsat += weight
			// That clause cannot be broken : it *is* broken
			s.breakCount[v] -= weight

			// Any lit can make this clause true
			for _, lit := range s.problem.Clauses[idx] {
				s.makeCount[lit.Var()] += weight
			}
		} else if s.nbTrue[idx] == 1 {
			// Clause will be broken if last lit is flipped
			// Find it and update breakCount
			for _, lit := range s.problem.Clauses[idx] {
				vi := lit.Var()

				if s.model[vi] == lit.Positive() {
					// It is the only positive lit in the clause
					s.breakCount[vi] += weight
					break
				}
			}
		}
	}
}

func (s *AnnealingSolver) Solve(nbTries int) Model {
	t := 1.0
	//fmt.Printf("Score=%d, model=%v\n", s.weightUnsat, s.model)

	for i := 1; i <= nbTries; i++ {
		v := Var(rand.Intn(s.problem.NbVars))
		//fmt.Printf("Chose var %d that would break %d and make %d\n", v, s.breakCount[v], s.makeCount[v])
		delta := float64(s.makeCount[v]-s.breakCount[v]) / float64(s.problem.NbVars)
		t *= 0.99
		//fmt.Printf("t=%g, delta=%g, p=%g\n", t, delta, proba(delta, t))

		if proba(delta, t) > rand.Float64() {
			s.flip(v)
			//fmt.Printf("Flipped, score=%d model=%v\n", s.weightUnsat, s.model)
		}
	}

	return s.model
}

func proba(delta, temp float64) float64 {
	if delta >= 0 {
		return 1.0
	} else {
		num := 1 - delta
		return math.Exp(-num / temp)
	}
}
