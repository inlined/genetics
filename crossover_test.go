package genetics_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/inlined/rand"
	"github.com/inlined/xkcd"

	"github.com/inlined/genetics"
)

func TestCrossover(t *testing.T) {
	// Permutative crossovers
	for _, test := range []struct {
		tag      string
		p1       []genetics.Gene // if unspecified: [1, 6)
		p2       []genetics.Gene // if unspecified: [6, 11)
		strategy genetics.Crossover
		rand     rand.Rand
		c1       []genetics.Gene
		c2       []genetics.Gene
	}{
		{
			tag:      "crossover once midpoint",
			strategy: genetics.MultiPointCrossover{Points: 1},
			rand:     xkcd.Rand(2),
			c1:       []genetics.Gene{1, 2, 8, 9, 10},
			c2:       []genetics.Gene{6, 7, 3, 4, 5},
		}, {
			tag:      "crossover once beginning",
			strategy: genetics.MultiPointCrossover{Points: 1},
			rand:     xkcd.Rand(0),
			c1:       []genetics.Gene{6, 7, 8, 9, 10},
			c2:       []genetics.Gene{1, 2, 3, 4, 5},
		}, {
			tag:      "crossover once end",
			strategy: genetics.MultiPointCrossover{Points: 1},
			rand:     xkcd.Rand(4),
			c1:       []genetics.Gene{1, 2, 3, 4, 10},
			c2:       []genetics.Gene{6, 7, 8, 9, 5},
		}, {
			tag:      "crossover twice, adjacent",
			strategy: genetics.MultiPointCrossover{Points: 2},
			rand:     xkcd.Rand(1, 2),
			c1:       []genetics.Gene{1, 7, 3, 4, 5},
			c2:       []genetics.Gene{6, 2, 8, 9, 10},
		}, {
			tag:      "crossover twice, span",
			strategy: genetics.MultiPointCrossover{Points: 2},
			rand:     xkcd.Rand(1, 3),
			c1:       []genetics.Gene{1, 7, 8, 4, 5},
			c2:       []genetics.Gene{6, 2, 3, 9, 10},
		}, {
			tag:      "crossover twice, out of order",
			strategy: genetics.MultiPointCrossover{Points: 2},
			rand:     xkcd.Rand(3, 1),
			c1:       []genetics.Gene{1, 7, 8, 4, 5},
			c2:       []genetics.Gene{6, 2, 3, 9, 10},
		}, {
			tag:      "crossover thrice",
			strategy: genetics.MultiPointCrossover{Points: 3},
			rand:     xkcd.Rand(3, 1, 4),
			c1:       []genetics.Gene{1, 7, 8, 4, 10},
			c2:       []genetics.Gene{6, 2, 3, 9, 5},
		}, {
			tag:      "recombination, flip",
			strategy: genetics.WholeArithmeticRecombination{},
			rand:     xkcd.Rand(0.0),
			c1:       []genetics.Gene{6, 7, 8, 9, 10},
			c2:       []genetics.Gene{1, 2, 3, 4, 5},
		}, {
			tag:      "recombination, center",
			strategy: genetics.WholeArithmeticRecombination{},
			rand:     xkcd.Rand(0.5),
			c1:       []genetics.Gene{4, 5, 6, 7, 8},
			c2:       []genetics.Gene{3, 4, 5, 6, 7},
		}, {
			tag:      "recombination, biased",
			strategy: genetics.WholeArithmeticRecombination{},
			rand:     xkcd.Rand(0.2),
			c1:       []genetics.Gene{5, 6, 7, 8, 9},
			c2:       []genetics.Gene{2, 3, 4, 5, 6},
		}, {
			tag:      "OX1",
			strategy: genetics.DavisOrderCrossover{},
			rand:     xkcd.Rand(1, 3),
			p1:       []genetics.Gene{0, 1, 2, 3, 4},
			p2:       []genetics.Gene{0, 1, 2, 3, 4},
			c1:       []genetics.Gene{4, 1, 2, 0, 3},
			c2:       []genetics.Gene{4, 1, 2, 0, 3},
		}, {
			tag:      "OX1 dealing different orders",
			strategy: genetics.DavisOrderCrossover{},
			rand:     xkcd.Rand(3, 1),
			p1:       []genetics.Gene{0, 1, 2, 3, 4},
			p2:       []genetics.Gene{0, 1, 2, 3, 4},
			c1:       []genetics.Gene{4, 1, 2, 0, 3},
			c2:       []genetics.Gene{4, 1, 2, 0, 3},
		}, {
			tag:      "OX1 uses crossover",
			strategy: genetics.DavisOrderCrossover{},
			rand:     xkcd.Rand(1, 3),
			p1:       []genetics.Gene{0, 1, 2, 3, 4},
			p2:       []genetics.Gene{4, 3, 2, 1, 0},
			c1:       []genetics.Gene{0, 1, 2, 4, 3},
			c2:       []genetics.Gene{4, 3, 2, 0, 1},
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			s := genetics.NewSpecies(5, 20)
			p1 := s.New(test.p1...)
			if test.p1 == nil {
				p1 = s.New(1, 2, 3, 4, 5)
			}
			p2 := s.New(test.p2...)
			if test.p2 == nil {
				p2 = s.New(6, 7, 8, 9, 10)
			}

			got1, got2 := test.strategy.Crossover(test.rand, p1, p2)

			if diff := cmp.Diff(got1.Genes, test.c1); diff != "" {
				t.Errorf("Crossover() returned unexpected gene 1. Got=%v; want=%v; diff=%s", got1.Genes, test.c1, diff)
			}
			if diff := cmp.Diff(got2.Genes, test.c2); diff != "" {
				t.Errorf("Crossover() returned unexpected gene 2. Got=%v; want=%v; diff=%s", got2.Genes, test.c2, diff)
			}
		})
	}
}
