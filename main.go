// go4maxsat project main.go
package main

import (
	"fmt"
	"go4maxsat/solver"
	"math/rand"
	"os"
	"runtime"
)

const nbSolvers = 100

var nbRoutines = runtime.NumCPU()

type solution struct {
	score int
	model solver.Model
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) != 3 || os.Args[1] != "tabu" && os.Args[1] != "sa" {
		fmt.Fprintf(os.Stderr, "Syntax : %s tabu|sa file.cnf\n", os.Args[0])
		os.Exit(1)
	}

	strategy := os.Args[1]
	pb, err := solver.LoadDimacs(os.Args[2])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("c parsed %d vars %d clauses\n", pb.NbVars, len(pb.Clauses))
		chBest := make(chan solution)
		chEnd := make(chan bool)

		for i := 0; i < nbRoutines; i++ {
			go func() {
				best := len(pb.Clauses)

				for j := 0; j < nbSolvers; j++ {
					var s solver.Solver

					if strategy == "tabu" {
						s = solver.NewTabu(pb, rand.Intn(10))
					} else {
						s = solver.NewAnnealingSolver(pb)
					}

					model := s.Solve(100000)
					score := s.Score()

					if score < best {
						best = score
						chBest <- solution{score, model}

						if score == 0 {
							break
						}
					}
				}

				chEnd <- true
			}()
		}

		best := len(pb.Clauses)
		var model solver.Model
		ended := 0

		for ended < nbRoutines && best > 0 {
			select {
			case sol := <-chBest:
				if sol.score < best {
					best = sol.score
					fmt.Printf("o %d\n", best)
					model = sol.model
				}

			case _ = <-chEnd:
				ended++
			}
		}

		fmt.Printf("%v\n", model)
	}
}
