package genetics

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/inlined/rand"
)

// Should Species become Population, which has an array of array of Genes?
// Then Population could create a new generation, which wouldn't be too bad via
// shallow copies of slices. Though that would possibly pin memory and cause leaks
// Still, it might help encapsulate concepts.

// Gene is a single trait to control behavior.
// IF THIS TYPE IS CHANGED FROM BYTE, Species.NewRand() MUST CHANGE
type Gene = int

// Fitness is an arbitrary fitness number based on genomes and their matching traits.
type Fitness int64

// Chromosome represents a single genetic strategy for a Species.
type Chromosome struct {
	Species *Species
	Genes   []Gene
}

// String prints Gene list of a Chromosome but does not preserve the name of the Species.
func (c Chromosome) String() string {
	//return base64.StdEncoding.EncodeToString(c.Genes)
	return "DEPRECATED"
}

// Species is a factory for all Genes in a repeated evolutionary experiment.
// Separating this from the actual Chromosome allows easier reuse of genetic algorithms
// in multiple circumstances as well as experimentation with the ordering of Chromosomes
// which influences the rate at which they may be separated by crossovers.
type Species struct {
	NumGenes  int
	MaxAllele Gene
}

// NewSpecies initializes a Species
func NewSpecies(numGenes int, maxAllele Gene) *Species {
	return &Species{
		NumGenes:  numGenes,
		MaxAllele: maxAllele,
	}
}

// New creates a Chromosome of the species. Any passed Genes
// are initialized starting at index 0. Any surpluss Genes
// are ignored and any missing Genes are 0-initialized.
func (s *Species) New(g ...Gene) Chromosome {
	c := Chromosome{
		Genes:   make([]Gene, s.NumGenes),
		Species: s,
	}
	for i := 0; i < len(g) && i < s.NumGenes; i++ {
		c.Genes[i] = g[i]
	}
	return c
}

// NewRand creates a random-initialized Chromosome of the species;
// each allele is independently randomized
func (s *Species) NewRand(rng rand.Rand) (Chromosome, error) {
	child := s.New()
	b := make([]byte, s.NumGenes)
	if n, err := rng.Read(b); n != s.NumGenes || err != nil {
		return Chromosome{}, fmt.Errorf("rand.Read(); wanted %d bytes; got %d bytes; err=%s", s.NumGenes, n, err)
	}
	for n, v := range b {
		// Type jumping to avoid overflow when MaxAllele is the maximum Gene value
		child.Genes[n] = Gene(int(v) % (int(s.MaxAllele) + 1))
	}
	return child, nil
}

// NewPerm creates a random
func (s *Species) NewPerm(rng rand.Rand) (Chromosome, error) {
	// a permutation of [0, NumGenes) must fit in MaxAllele
	if int(s.MaxAllele) < s.NumGenes-1 {
		return Chromosome{}, fmt.Errorf("NewPerm() cannot generate %d elements with max %d", s.NumGenes, s.MaxAllele)
	}
	child := s.New()
	r := rng.Perm(s.NumGenes)
	for i, v := range r {
		child.Genes[i] = Gene(v)
	}
	return child, nil
}

// ParseChromosome creates an in-memory representation for Chromosomes encoded with SerializeChromosome
func (s *Species) ParseChromosome(encoded string) (Chromosome, error) {
	/*
		g, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return Chromosome{}, err
		}

		if len(g) != int(s.NumGenes) {
			return Chromosome{}, fmt.Errorf("Species.ParseChromosome(%s); expected %d alleles, got %d", encoded, s.NumGenes, len(g))
		}
		for _, a := range g {
			if a > s.MaxAllele {
				return Chromosome{}, fmt.Errorf("Chromosome%v has allele %d greater than maximum %d", g, a, s.MaxAllele)
			}
		}
		return Chromosome{
			Genes:   g,
			Species: s,
		}, nil
	*/
	return Chromosome{}, errors.New("DEPRECATED")
}

// Evolver replaces one generation of genes with another
type Evolver struct {
	ReplacementCount int
	MutationRate     float32
	Selector         NaturalSelection
	Crossover        Crossover
	Mutator          Mutator
}

// Evolve replaces a handful of the population with the next generation
func (e Evolver) Evolve(rand rand.Rand, pop []Chromosome, scores []Fitness) {
	indexes := e.Selector.SelectParents(rand, e.ReplacementCount, scores)
	rand.Shuffle(len(indexes), func(i, j int) {
		indexes[i], indexes[j] = indexes[j], indexes[i]
	})
	children := make([]Chromosome, e.ReplacementCount)
	for i := 0; i < e.ReplacementCount; i += 2 {
		children[i], children[i+1] = e.Crossover.Crossover(rand, pop[indexes[i]], pop[indexes[i+1]])
		if rand.Float32() < e.MutationRate {
			e.Mutator.Mutate(rand, &children[i])
		}
		if rand.Float32() < e.MutationRate {
			e.Mutator.Mutate(rand, &children[i+1])
		}
	}

	minIndexes := kMinIndexes(scores, e.ReplacementCount)
	for child, parent := range minIndexes {
		pop[parent] = children[child]
	}
}

func kMinIndexes(f []Fitness, k int) []int {
	h := make(maxTieHeap, k)
	for i := 0; i < k; i++ {
		h[i] = tie{
			index:   i,
			fitness: f[i],
		}
	}
	heap.Init(h)

	for i := k; i < len(f); i++ {
		if f[i] < h[0].fitness {
			h[0].index = i
			h[0].fitness = f[i]
			heap.Fix(h, 0)
		}
	}

	res := make([]int, k)
	for i, v := range h {
		res[i] = v.index
	}
	return res
}
