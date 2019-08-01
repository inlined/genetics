package genetics_test

import (
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/inlined/xkcdrand"

	"github.com/inlined/goldmine/pkg/genetics"
)

func TestStochasticUniversalSampling(t *testing.T) {
	for _, test := range []struct {
		tag             string
		numParents      int
		source          rand.Source
		fitness         []genetics.Fitness
		expectedIndexes []int
	}{
		{
			tag:             "pick every other (even)",
			numParents:      3,
			source:          xkcdrand.Sequence(1),
			fitness:         []genetics.Fitness{2, 2, 2, 2, 2, 2},
			expectedIndexes: []int{0, 2, 4},
		},
		{
			tag:             "pick every other (odd)",
			numParents:      3,
			source:          xkcdrand.Sequence(3),
			fitness:         []genetics.Fitness{2, 2, 2, 2, 2, 2},
			expectedIndexes: []int{1, 3, 5},
		},
		{
			// This is an edge case and a major sign to switch the selection mechanism to ranked scoring
			// TODO: should this also return a diversity score to help automate the scoring algorithm?
			tag:             "top-exclusively",
			numParents:      3,
			source:          xkcdrand.Sequence(1),
			fitness:         []genetics.Fitness{10, 1, 1},
			expectedIndexes: []int{0, 0, 0},
		},
		{
			tag:             "redundant picks",
			numParents:      3,
			source:          xkcdrand.Sequence(2),
			fitness:         []genetics.Fitness{10, 1, 1},
			expectedIndexes: []int{0, 0, 1},
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			s := genetics.StochasticUniversalSampling(test.source)
			got := s.SelectParents(test.numParents, test.fitness)
			if diff := cmp.Diff(got, test.expectedIndexes); diff != "" {
				t.Fatalf("Got wrong indexes; got=%v; want=%v; diff=%v", got, test.expectedIndexes, diff)
			}
		})
	}
}

func TestRandomResetting(t *testing.T) {
	const mutationRate = 0.01
	for _, test := range []struct {
		tag                string
		bitsPerGenome      uint8
		numGenomes         uint8
		mutationRate       float32
		randSource         rand.Source
		chromosome         uint64
		expectedChromosome uint64
	}{
		{
			tag:                "no mutation",
			bitsPerGenome:      8,
			numGenomes:         4,
			randSource:         xkcdrand.Sequence(xkcdrand.Float64(0.5)),
			chromosome:         0xBAADF00D,
			expectedChromosome: 0xBAADF00D,
		}, {
			tag:                "nibble-wide; first chromosome",
			bitsPerGenome:      4,
			numGenomes:         4,
			randSource:         xkcdrand.SequenceInt32(0, 0, 0xD),
			chromosome:         0xF00D,
			expectedChromosome: 0xD00D,
		}, {
			tag:                "nibble-wide; last chromosome",
			bitsPerGenome:      4,
			numGenomes:         4,
			randSource:         xkcdrand.SequenceInt32(0, 3, 0xF),
			chromosome:         0xF00D,
			expectedChromosome: 0xF00F,
		}, {
			tag:                "byte-wide; last chromosome",
			bitsPerGenome:      8,
			numGenomes:         2,
			randSource:         xkcdrand.SequenceInt32(0, 1, 0xF),
			chromosome:         0xF00D,
			expectedChromosome: 0xF00F,
		}, {
			tag:                "binary",
			bitsPerGenome:      1,
			numGenomes:         16,
			randSource:         xkcdrand.SequenceInt32(0, 2, 1),
			chromosome:         0xD00D,
			expectedChromosome: 0xF00D,
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			s, err := genetics.NewSpecies(test.bitsPerGenome, test.numGenomes)
			if err != nil {
				t.Errorf("genetics.NewSpecies(%d, %d) failed: %s", test.bitsPerGenome, test.numGenomes, err)
				return
			}
			chromosome, err := s.DeserializeChromosome(test.chromosome)
			if err != nil {
				t.Errorf("genetics.DeserializeChromosome() failed: %s", err)
				return
			}
			mutator := genetics.RandomResetting{
				RNG:  rand.New(test.randSource),
				Freq: mutationRate,
			}
			mutator.Mutate(&chromosome)
			got, err := s.SerializeChromosome(chromosome)
			if err != nil {
				t.Errorf("genetics.SerializeChromosome() failed: %s", err)
				return
			}
			if got != test.expectedChromosome {
				t.Errorf("RandomResetting.Mutate(%x): got=%x want=%x", test.chromosome, got, test.expectedChromosome)
			}
		})
	}
}
