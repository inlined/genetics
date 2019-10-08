package genetics_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/inlined/rand"
	"github.com/inlined/xkcd"

	"github.com/inlined/genetics"
)

func TestParentSelection(t *testing.T) {
	for _, test := range []struct {
		tag             string
		strategy        genetics.NaturalSelection
		numSelected     int
		fitness         []genetics.Fitness
		rand            rand.Rand
		expectedParents []int
	}{
		{
			tag:             "SUS pick every other (even)",
			strategy:        genetics.StochasticUniversalSampling{},
			numSelected:     3,
			fitness:         []genetics.Fitness{2, 2, 2, 2, 2, 2},
			rand:            xkcd.Rand(1),
			expectedParents: []int{0, 2, 4},
		}, {
			tag:             "SUS pick every other (odd)",
			strategy:        genetics.StochasticUniversalSampling{},
			numSelected:     3,
			fitness:         []genetics.Fitness{2, 2, 2, 2, 2, 2},
			rand:            xkcd.Rand(3),
			expectedParents: []int{1, 3, 5},
		}, {
			// This is an edge case and a major sign to switch the selection mechanism to ranked scoring
			// TODO: should this also return a diversity score to help automate the scoring algorithm?
			tag:             "SUS top-exclusively",
			strategy:        genetics.StochasticUniversalSampling{},
			numSelected:     3,
			fitness:         []genetics.Fitness{10, 1, 1},
			rand:            xkcd.Rand(1),
			expectedParents: []int{0, 0, 0},
		}, {
			tag:             "SUS redundant picks",
			strategy:        genetics.StochasticUniversalSampling{},
			numSelected:     3,
			fitness:         []genetics.Fitness{10, 1, 1},
			rand:            xkcd.Rand(2),
			expectedParents: []int{0, 0, 1},
		}, {
			tag:             "Ranked wheel begin",
			strategy:        genetics.RankedSelection{},
			numSelected:     3,
			fitness:         []genetics.Fitness{10, 5, 1}, // Ranked weights: 3 2 1; d=6/3=2
			rand:            xkcd.Rand(0),
			expectedParents: []int{0, 0, 1},
		}, {
			tag:             "Ranked wheel end",
			strategy:        genetics.RankedSelection{},
			numSelected:     3,                            // d=6/3=2
			fitness:         []genetics.Fitness{10, 5, 1}, // Ranked weights: 3 2 1
			rand:            xkcd.Rand(1),
			expectedParents: []int{0, 1, 2},
		}, {
			tag:             "Ranked wheel scrambled",
			strategy:        genetics.RankedSelection{},
			numSelected:     2,                                // d = 10 / 2 = 5
			fitness:         []genetics.Fitness{4, 20, 16, 3}, // Ranked weights: 2, 4, 3, 1
			rand:            xkcd.Rand(4),
			expectedParents: []int{2, 3},
		}, {
			tag:             "Tournament of 1", // To some extent, this verifies I understand the cryptic rules for deal()
			strategy:        genetics.TournamentSelection{Size: 1},
			numSelected:     2,                                // d = 10 / 2 = 5
			fitness:         []genetics.Fitness{4, 20, 16, 3}, // Ranked weights: 2, 4, 3, 1
			rand:            xkcd.Rand(3, 1),                  // deal {3}, {1}
			expectedParents: []int{3, 1},
		}, {
			tag:             "Tournament of 2",
			strategy:        genetics.TournamentSelection{Size: 2},
			numSelected:     2,                                // d = 10 / 2 = 5
			fitness:         []genetics.Fitness{4, 20, 16, 3}, // Ranked weights: 2, 4, 3, 1
			rand:            xkcd.Rand(3, 2, 1, 2),            // deal {3, 2}, {1, 2}
			expectedParents: []int{2 /* winner of 3 vs 2 */, 1 /* winner of 1 vs 2 */},
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			got := test.strategy.SelectParents(test.rand, test.numSelected, test.fitness)
			if diff := cmp.Diff(got, test.expectedParents); diff != "" {
				t.Fatalf("Got wrong indexes; got=%v; want=%v; diff=%v", got, test.expectedParents, diff)
			}
		})
	}
}
