package solver

type Solver interface {
	Solve(nbTries int) Model
	Score() int // Current score, i.e. total weight of unsat clauses
	// Small is better, objective is 0
}
