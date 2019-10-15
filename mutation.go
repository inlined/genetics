package genetics

import (
	"fmt"

	"github.com/inlined/rand"
)

const (
	randomResettingMutation = "RandomResettingMutation"
	swapMutation            = "SwapMutation"
	scrambleMutation        = "ScrambleMutation"
	inversionMutation       = "InversionMutation"
)

// Mutator introduces randomness to the population.
// While mutations should be rare to avoid turning the algorithm into a random
// walk, some mutations are necessary to enforce convergence.
// Mutators work on unpacked Chromosomes because species' bit length is
// important to some algorithms.
type Mutator interface {
	fmt.Stringer
	Mutate(r rand.Rand, c *Chromosome)
}

// RandomResettingMutation (equivalent to Bit Flip Mutation for Species with a bitwidth of 1)
// Will randomly set an allele to one of the acceptable values. Assumes Species' Genomes
// accept values of any bit size. This is most useful for chromasomes where genes affect
// independent behavior (e.g. not permutation-based algorithms).
type RandomResettingMutation struct{}

func (RandomResettingMutation) String() string {
	return randomResettingMutation
}

// Mutate implements the Mutator interface
func (m RandomResettingMutation) Mutate(r rand.Rand, c *Chromosome) {
	n := r.Int31n(int32(len(c.Genes)))
	v := r.Int31n(int32(c.Species.MaxAllele))
	c.Genes[n] = Gene(v)
}

// SwapMutation mutations swap the value of two genomes.
// SwapMutation is a mutation most appropriate for permutation genes
// (e.g. graph algorithms)
type SwapMutation struct{}

func (SwapMutation) String() string {
	return swapMutation
}

// Mutate implements the mutator interface
func (m SwapMutation) Mutate(r rand.Rand, c *Chromosome) {
	// To avoid worrying about a collision with the same index, we'll
	// instead calculate both an index and an offset from that index
	// (wrapping around as a cyclical buffer)
	len := int32(len(c.Genes))
	i0 := r.Int31n(len - 1)
	d := r.Int31n(len-i0-1) + 1
	i1 := i0 + d
	v0 := c.Genes[i0]
	c.Genes[i0] = c.Genes[i1]
	c.Genes[i1] = v0
}

// ScrambleMutation picks two crossover points and scrambles the alleles
// in the middle segment. This is most appropraite for permutation-encoded
// Genes, such as graph algorithms.
type ScrambleMutation struct{}

func (ScrambleMutation) String() string {
	return scrambleMutation
}

// Mutate implements Mutator
func (m ScrambleMutation) Mutate(r rand.Rand, c *Chromosome) {
	s := c.Species
	l := r.Int31n(int32(s.NumGenes) - 1)
	d := r.Int31n(int32(s.NumGenes)-l-1) + 1
	u := d + l
	for i := l; i < u; i++ {
		d2 := r.Int31n(d + 1)
		c.Genes[i], c.Genes[l+d2] = c.Genes[l+d2], c.Genes[i]
	}
}

// InversionMutation picks two crossover points and then flipps the alleles
// in the middle segment. This is most appropraite for permutation-encoded
// Genes, such as graph algorithms.
type InversionMutation struct{}

func (InversionMutation) String() string {
	return inversionMutation
}

// Mutate implements Mutator
func (m InversionMutation) Mutate(r rand.Rand, c *Chromosome) {
	s := c.Species
	l := r.Int31n(int32(s.NumGenes) - 1)
	d := r.Int31n(int32(s.NumGenes)-l-1) + 1
	u := d + l
	for ; l < u; l, u = l+1, u-1 {
		c.Genes[l], c.Genes[u] = c.Genes[u], c.Genes[l]
	}
}
