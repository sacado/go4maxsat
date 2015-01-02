package solver

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

// Parses a clause, provided line is a DIMACS CNF clause
func parseClause(line string) (cl Clause, err error) {
	tokens := strings.Fields(line)
	cl = make(Clause, len(tokens)-1)

	for j := 0; j < len(tokens)-1; j++ {
		lit, err := strconv.ParseInt(tokens[j], 10, 32)

		if err != nil {
			return cl, err
		}

		cl[j] = Lit(lit)
	}

	return cl, nil
}

// Parses CNF clauses from scanner and updates pb accordingly
// Might return an error
func parseCNFClauses(scanner *bufio.Scanner, pb *Problem) error {
	for i := range pb.Clauses {
		if !scanner.Scan() {
			return errors.New("Invalid # of clauses")
		}

		pb.Weights[i] = 1
		cl, err := parseClause(scanner.Text())

		if err != nil {
			return err
		}

		pb.Clauses[i] = cl
	}

	return nil
}

// Parses a DIMACS WCNF clause, i.e a clause whose first literal is a weight
func parseWeightedClause(line string) (weight int, cl Clause, err error) {
	tokens := strings.Fields(line)
	cl = make(Clause, len(tokens)-2)

	w, err := strconv.ParseInt(tokens[0], 10, 32)

	if err != nil {
		return weight, cl, err
	}

	for j := 1; j < len(tokens)-1; j++ {
		lit, err := strconv.ParseInt(tokens[j], 10, 32)

		if err != nil {
			return weight, cl, err
		}

		cl[j-1] = Lit(lit)
	}

	return int(w), cl, nil
}

// Parses WCNF clauses from scanner and updates pb accordingly
// Might return an error
func parseWCNFClauses(scanner *bufio.Scanner, pb *Problem) error {
	for i := range pb.Clauses {
		if !scanner.Scan() {
			return errors.New("Invalid # of clauses")
		}

		weight, cl, err := parseWeightedClause(scanner.Text())

		if err != nil {
			return err
		}

		pb.Weights[i] = weight
		pb.Clauses[i] = cl
	}

	return nil
}

func parseDimacs(f *os.File) (pb Problem, err error) {
	scanner := bufio.NewScanner(f)

	if !scanner.Scan() {
		return pb, errors.New("Could not read header")
	}

	line := scanner.Text()
	toks := strings.Fields(line)

	if len(toks) != 4 || toks[0] != "p" || (toks[1] != "cnf" && toks[1] != "wcnf") {
		return pb, errors.New("Invalid header syntax")
	}

	nbVars, err := strconv.ParseInt(toks[2], 10, 32)

	if err != nil {
		return pb, errors.New("Invalid # of vars")
	}

	pb.NbVars = int(nbVars)
	nbClauses, err := strconv.ParseInt(toks[3], 10, 32)

	if err != nil {
		return pb, errors.New("Invalid # of clauses")
	}

	pb.Weights = make([]int, nbClauses)
	pb.Clauses = make([]Clause, nbClauses)

	if toks[1] == "cnf" {
		return pb, parseCNFClauses(scanner, &pb)
	} else {
		return pb, parseWCNFClauses(scanner, &pb)
	}
}

// Loads a CNF file whose path is 'filename'
// and parses the underlying Problem
// An error might be returned during parsing or file reading
func LoadDimacs(filename string) (pb Problem, err error) {
	f, err := os.Open(filename)

	if err != nil {
		return pb, err
	}

	return parseDimacs(f)
}
