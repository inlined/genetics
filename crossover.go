package genetics

import (
	"math"
	"sort"

	"github.com/inlined/rand"
)

// Crossover is a strategy for generating two children based
// on two parents.
type Crossover interface {
	Crossover(r rand.Rand, a, b Chromosome) (x, y Chromosome)
}

// MultiPointCrossover is a generalization of the Crossover method.
// N crossover points are selected and children are made of parens'
// chromosomes alternating sources at the crossover points.
// Multi-point crossovers are appropriate for numeric chromosomes
// NOTE: It might be more appropriate to allow crossovers mid-allele,
// which  might require different int encodings.
type MultiPointCrossover struct {
	Points int
}

// Crossover imnplements Crossover.
// Inefficiency: This algorithm makes n^2 data copies because it assumes
// N is ~1-3
func (c MultiPointCrossover) Crossover(r rand.Rand, a, b Chromosome) (x, y Chromosome) {
	s := a.Species
	x = s.New()
	y = s.New()
	temp := s.New()
	copy(x.Genes[:], a.Genes[:])
	copy(y.Genes[:], b.Genes[:])
	indexes := rand.Deal(r, s.NumGenes, c.Points)
	sort.Ints(indexes)
	for _, n := range indexes {
		copy(temp.Genes[n:], x.Genes[n:])
		copy(x.Genes[n:], y.Genes[n:])
		copy(y.Genes[n:], temp.Genes[n:])
	}
	return x, y
}

// WholeArithmeticRecombination picks a random float weight from 0-1. The children are
// a weighted average of the parents with inverse weights.
// Whole arithmetic recombinatinos are appropriate for numeric chromosomes and will
// trend towards the average value of the population.
type WholeArithmeticRecombination struct{}

// Crossover implements Crossover
func (c WholeArithmeticRecombination) Crossover(r rand.Rand, a, b Chromosome) (x, y Chromosome) {
	f := r.Float64()
	s := a.Species
	x = s.New()
	y = s.New()
	for i := 0; i < s.NumGenes; i++ {
		// Because we're dealing with integers, a strict linear interpolation
		// will floor twice.
		// To avoid the edge case where 0.5 rounds up twice, we'll only do float
		// math once and then calculate an int delta that's applied twice.
		f1 := f*float64(a.Genes[i]) + (1-f)*float64(b.Genes[i])
		x.Genes[i] = Gene(math.Round(f1))
		d := x.Genes[i] - a.Genes[i]
		y.Genes[i] = b.Genes[i] - d
	}
	return x, y
}

// DavisOrderCrossover aka OX1 picks two crossover points, dividing the genomes
// into three segments. The middle segment is preserved whereas the right and
// left are rotationally filled with the left and right of the other chromosome.
// OX1 is appropraite for permutative genes, such as graph algorithms.
type DavisOrderCrossover struct{}

// Crossover implements Crossover
func (c DavisOrderCrossover) Crossover(r rand.Rand, a, b Chromosome) (x, y Chromosome) {

	indexes := rand.Deal(r, len(b.Genes)+1, 2)
	if indexes[0] > indexes[1] {
		indexes[0], indexes[1] = indexes[1], indexes[0]
	}
	return davisCrossoverOne(a, b, indexes[0], indexes[1]), davisCrossoverOne(b, a, indexes[0], indexes[1])
}

func davisCrossoverOne(p1, p2 Chromosome, lower, upper int) Chromosome {
	s := p1.Species
	child := s.New()
	seen := make([]bool, s.NumGenes)

	// 1. Preserve the range [lower, upper) of p1
	for i := lower; i < upper; i++ {
		seen[p1.Genes[i]] = true
		child.Genes[i] = p1.Genes[i]
	}

	// 2. Fill in unseen portions of [0, len) p2 into p1 starting
	// at upper
	insert := upper % s.NumGenes
	for read := 0; read < s.NumGenes; read++ {
		if seen[p2.Genes[read]] {
			continue
		}
		child.Genes[insert] = p2.Genes[read]
		insert = (insert + 1) % s.NumGenes
	}
	return child
}
