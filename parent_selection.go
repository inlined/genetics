package genetics

import (
	"sort"

	"github.com/inlined/rand"
)

// NaturalSelection is an interface to pick the selection method.
// A NaturalSelection MAY NOT BE GOROUTINE SAFE. It may only be used in one Evolve function at a time.
// This helps avoid the lock incurred by the top-level rand functions.
// TODO: consider nested interfaces (NaturalSelection has a Seed() function to return a Selector
// that implements SelectParents). This would avoid re-generating the roulette wheel in
// StochasticUniversalSampling
type NaturalSelection interface {
	SelectParents(rand rand.Rand, numParents int, fitness []Fitness) (indexes []int)
}

// StochasticUniversalSampling creates a "roulette" wheel where each parent
// gets a slice in proportion to their fitness. We then spin the wheel with
// two fixed points to select which parents win.
// If src is nill, a new source is created with the current time.
type StochasticUniversalSampling struct{}

// SelectParents implements the NaturalSelection interface.
func (s StochasticUniversalSampling) SelectParents(rand rand.Rand, numParents int, fitness []Fitness) (indexes []int) {
	totalFitness := Fitness(0)
	for _, f := range fitness {
		totalFitness += f
	}

	// Use a fixed distance (uniform distribution) across the wheel.
	// Note: we choose here to use integer arithmetic instead of a float distribution.
	// This uses faster ALUs but introduces the possibility of error when totalFitness !>> numParents
	distance := totalFitness / Fitness(numParents)
	// Spin the wheel up to distance (equivalent to spinning the wheel randomly and then taking the modulo
	// of the size)
	pos := Fitness(rand.Int63n(int64(distance)))

	// Iterate through the fitness scores as if it were a roulete wheel (e.g. incrementing f by
	// fitness[n] rather than one) and remember the indexes which contain any pointers P.
	// In edge cases, a position may hit the same parent multiple times; in this case, the parent
	// is selected repeatedly.
	// TODO: Should this be instead selected with a weight to avoid a parent mating with itself?
	indexes = make([]int, 0, numParents)
	accumFitness := Fitness(0)
	for n := 0; len(indexes) < numParents; n++ {
		accumFitness += fitness[n]
		for ; pos < accumFitness; pos += distance {
			indexes = append(indexes, n)
		}
	}

	return indexes
}

// RankedSelection gives each chromosome odds of reproduction not based on its proportional
// fitness, but its rank in overall fitness. This ensures that populations trend towards
// optimal solutions still as the problem is converging.
type RankedSelection struct{}

// SelectParents selects parents in proportion to their fitness' rank.
func (s RankedSelection) SelectParents(rand rand.Rand, numParents int, fitness []Fitness) (indexes []int) {
	// Lazy version of a Schwartzian transform; may be memory wasteful
	zipped := make([]tie, len(fitness))
	for n, f := range fitness {
		zipped[n] = tie{index: n, fitness: f}
	}
	sort.Slice(zipped, func(i, j int) bool {
		return zipped[i].fitness > zipped[j].fitness
	})
	rankedIndexes := make([]int, len(fitness))
	for n, t := range zipped {
		rankedIndexes[n] = t.index
	}

	// Edited version of SUS. Should we waste the cycles trying to use a universal internal
	// datastructure?
	totalRank := len(fitness) * (len(fitness) + 1) / 2

	// Use a fixed distance (uniform distribution) across the wheel.
	// Note: we choose here to use integer arithmetic instead of a float distribution.
	// This uses faster ALUs but introduces the possibility of error when totalFitness !>> numParents
	distance := totalRank / numParents
	// Spin the wheel up to distance (equivalent to spinning the wheel randomly and then taking the modulo
	// of the size)
	pos := int(rand.Int31n(int32(distance)))

	// Iterate through the fitness scores as if it were a roulete wheel (e.g. incrementing f by
	// fitness[n] rather than one) and remember the indexes which contain any pointers P.
	// In edge cases, a position may hit the same parent multiple times; in this case, the parent
	// is selected repeatedly.
	// TODO: Should this be instead selected with a weight to avoid a parent mating with itself?
	indexes = make([]int, 0, numParents)
	accumRank := 0
	for n := 0; len(indexes) < numParents; n++ {
		accumRank += len(rankedIndexes) - n
		for ; pos < accumRank; pos += distance {
			indexes = append(indexes, rankedIndexes[n])
		}
	}

	return indexes
}

// TournamentSelection picks each parent by picking Size candidates from a fitness list
// at random and selecting the parent with the greatest fitness.
type TournamentSelection struct {
	Size int
}

func (s TournamentSelection) selectOneParent(r rand.Rand, fitness []Fitness) int {
	indexes := rand.Deal(r, len(fitness), s.Size)
	maxFitness := fitness[indexes[0]]
	maxIndex := indexes[0]
	for n := 1; n < s.Size; n++ {
		if fitness[indexes[n]] >= maxFitness {
			maxFitness = fitness[indexes[n]]
			maxIndex = indexes[n]
		}
	}
	return maxIndex
}

// SelectParents selects the len(indexes) parents who win a s.Size-way tournament
func (s TournamentSelection) SelectParents(rand rand.Rand, numParents int, fitness []Fitness) (indexes []int) {
	indexes = make([]int, numParents)
	for n := 0; n < numParents; n++ {
		indexes[n] = s.selectOneParent(rand, fitness)
	}
	return indexes
}
