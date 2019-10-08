package genetics_test

import (
	"encoding/base64"
	"encoding/binary"
	"testing"

	"github.com/inlined/genetics"
	"github.com/inlined/rand"
	"github.com/inlined/xkcd"
)

func TestMutations(t *testing.T) {
	uint64ToChromosome := func(x uint32) string {
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, x)
		return base64.StdEncoding.EncodeToString(bs)
	}
	s := genetics.NewSpecies(4, 0xFF)

	for _, test := range []struct {
		tag      string
		mutator  genetics.Mutator
		rand     rand.Rand
		expected uint32
	}{
		{
			tag:      "reset first gene",
			mutator:  genetics.RandomResettingMutation{},
			rand:     xkcd.Rand(0, 0xDA),
			expected: 0xDAADF00D,
		}, {
			tag:      "reset middle gene",
			mutator:  genetics.RandomResettingMutation{},
			rand:     xkcd.Rand(2, 0xD0), // Index 2, value 0xD0
			expected: 0xBAADD00D,
		}, {
			tag:      "reset last gene",
			mutator:  genetics.RandomResettingMutation{},
			rand:     xkcd.Rand(3, 0x01),
			expected: 0xBAADF001,
		}, {
			tag:      "swap with first gene",
			mutator:  genetics.SwapMutation{},
			rand:     xkcd.Rand(0, 0), // Index 0, offset 0 + 1
			expected: 0xADBAF00D,
		}, {
			tag:      "swap with last gene",
			mutator:  genetics.SwapMutation{},
			rand:     xkcd.Rand(2, 0),
			expected: 0xBAAD0DF0,
		}, {
			tag:      "swap first and last gene",
			mutator:  genetics.SwapMutation{},
			rand:     xkcd.Rand(0, 2),
			expected: 0x0DADF0BA,
		}, {
			tag:      "swap middle genes",
			mutator:  genetics.SwapMutation{},
			rand:     xkcd.Rand(1, 0),
			expected: 0xBAF0AD0D,
		}, {
			tag:      "scramble first genes",
			mutator:  genetics.ScrambleMutation{},
			rand:     xkcd.Rand(0, 0, 1), // Index 0, offset 0+1, swap 0 with 1
			expected: 0xADBAF00D,
		}, {
			tag:      "scramble last genes",
			mutator:  genetics.ScrambleMutation{},
			rand:     xkcd.Rand(2, 0, 1),
			expected: 0xBAAD0DF0,
		}, {
			tag:      "scramble middle genes",
			mutator:  genetics.ScrambleMutation{},
			rand:     xkcd.Rand(1, 0, 1),
			expected: 0xBAF0AD0D,
		}, {
			tag:      "scramble many genes",
			mutator:  genetics.ScrambleMutation{},
			rand:     xkcd.Rand(1, 1, 1, 2), // Index 1, offset 1 + 1, swap 0 with 1, 1 with 2
			expected: 0xBAF00DAD,
		}, {
			tag:      "invert first genes",
			mutator:  genetics.InversionMutation{},
			rand:     xkcd.Rand(0, 0), // Index 0, offset 0 + 1
			expected: 0xADBAF00D,
		}, {
			tag:      "invert last genes",
			mutator:  genetics.InversionMutation{},
			rand:     xkcd.Rand(2, 0),
			expected: 0xBAAD0DF0,
		}, {
			tag:      "invert middle genes",
			mutator:  genetics.InversionMutation{},
			rand:     xkcd.Rand(1, 0),
			expected: 0xBAF0AD0D,
		}, {
			tag:      "invert all genes",
			mutator:  genetics.InversionMutation{},
			rand:     xkcd.Rand(0, 2),
			expected: 0x0DF0ADBA,
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			start := uint64ToChromosome(0xBAADF00D)
			want := uint64ToChromosome(test.expected)

			chromosome, err := s.ParseChromosome(start)
			if err != nil {
				t.Errorf("genetics.DeserializeChromosome(%s) failed: %s", start, err)
				return
			}
			test.mutator.Mutate(test.rand, &chromosome)
			if chromosome.String() != want {
				w, _ := s.ParseChromosome(want)
				t.Errorf("Mutate(): got=%x want=%x", chromosome.Genes, w.Genes)
			}
		})
	}
}
