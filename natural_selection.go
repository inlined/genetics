package genetics

import (
	"github.com/inlined/rand"
)

// NaturalSelection is an interface to pick the selection method.
// A NaturalSelection MAY NOT BE GOROUTINE SAFE. It may only be used in one Evolve function at a time.
// This helps avoid the lock incurred by the top-level rand functions.
// TODO: consider nested interfaces (NaturalSelection has a Seed() function to return a Selector
// that implements SelectParents). This would avoid re-generating the roulette wheel in
// StochasticUniversalSampling
type NaturalSelection interface {
	SelectParents(numParents int, fitness []Fitness) (indexes []int)
}

// StochasticUniversalSampling creates a "roulette" wheel where each parent
// gets a slice in proportion to their fitness. We then spin the wheel with
// two fixed points to select which parents win.
// If src is nill, a new source is created with the current time.
type StochasticUniversalSampling struct {
	Rand rand.Rand
}

// SelectParents implements the NaturalSelection interface.
func (s *StochasticUniversalSampling) SelectParents(numParents int, fitness []Fitness) (indexes []int) {
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
	pos := Fitness(s.Rand.Int63n(int64(distance)))

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

// RankedSelection weighs 

// Mutator introduces randomness to the population.
// While mutations should be rare to avoid turning the algorithm into a random
// walk, some mutations are necessary to enforce convergence.
// Mutators work on unpacked Chromosomes because species' bit length is
// important to some algorithms.
type Mutator interface {
	Mutate(c *Chromosome)
}

// RandomResetting (equivalent to Bit Flip Mutation for Species with a bitwidth of 1)
// Will randomly set an allele to one of the acceptable values. Assumes Species' Genomes
// accept values of any bit size. This is most useful for chromasomes where genes affect
// independent behavior (e.g. not permutation-based algorithms).
type RandomResetting struct {
	Rand rand.Rand
	Freq float32
}

// Mutate implements the Mutator interface
func (r *RandomResetting) Mutate(c *Chromosome) {
	f := r.Rand.Float32()
	if f > r.Freq {
		return
	}

	n := r.Rand.Int31n(int32(len(c.Genes)))
	v := r.Rand.Int31n(1 << c.Species.BitsPerGene)
	c.Genes[n] = Gene(v)
}

// Swapping ...
type Swapping struct {
	Rand rand.Rand
	Freq float32
}

// Mutate implements the mutator interface
func (s *Swapping) Mutate(c *Chromosome) {
	f := s.Rand.Float32()
	if f > s.Freq {
		return
	}

	// To avoid worrying about a collision with the same index, we'll
	// instead calculate both an index and an offset from that index
	// (wrapping around as a cyclical buffer)
	len := int32(len(c.Genes))
	i0 := s.Rand.Int31n(len)
	d := s.Rand.Int31n(len - 1)
	i1 := (i0 + d) % len
	v0 := c.Genes[i0]
	c.Genes[i0] = c.Genes[i1]
	c.Genes[i1] = v0
}
