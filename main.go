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
const nbRoutines = 4

type solution struct {
	score int
	model solver.Model
}

func main() {
	runtime.GOMAXPROCS(4)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Syntax : %s file.cnf\n", os.Args[0])
		os.Exit(1)
	}

	pb, err := solver.LoadDimacs(os.Args[1])

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
					s := solver.NewTabu(pb, rand.Intn(10))
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
