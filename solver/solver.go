package solver

// A Solver is a type that can solve a problem and return a Model
// that satisfies as many clauses as possible.
type Solver interface {
	Solve(nbTries int) Model
	Score() int // Current score, i.e. total weight of unsat clauses
	// Small is better, objective is 0
}
