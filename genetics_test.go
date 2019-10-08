package genetics_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/inlined/genetics"
	"github.com/inlined/rand"
)

func TestSerialization(t *testing.T) {
	// Can't test the top bit due to a lack of rand.uint63n
	for nGenes := 1; nGenes < 64; nGenes++ {
		s := genetics.NewSpecies(nGenes, 0xFF)
		t.Run(fmt.Sprintf("Length_%d", nGenes), func(t *testing.T) {
			for run := 0; run < 100; run++ {
				want, err := s.NewRand(rand.New())
				if err != nil {
					t.Error(err)
				}
				str := want.String()

				got, err := s.ParseChromosome(str)
				if err != nil {
					t.Errorf("Species.ParseChromosome(%s); err = %s", str, err)
				}
				if d := cmp.Diff(want, got); d != "" {
					t.Errorf("Round trip serialization failed; got=%+v; want=%+v; diff=%s", got, want, d)
				}
			}
		})
	}
}

func TestNewPerm(t *testing.T) {
	numGenes := 20
	maxAllele := genetics.Gene(19)
	s := genetics.NewSpecies(numGenes, maxAllele)
	c, err := s.NewPerm(rand.New())
	if err != nil {
		t.Errorf("NewPerm(); err=%s", err)
	}

	seen := make([]bool, numGenes)
	for _, v := range c.Genes {
		seen[int(v)] = true
	}

	for i, v := range seen {
		if !v {
			t.Errorf("NewPerm() did not create allele for %d", i)
		}
	}
}

func TestNewPermFailure(t *testing.T) {
	numGenes := 20
	maxAllele := genetics.Gene(18)
	_, err := genetics.NewSpecies(numGenes, maxAllele).NewPerm(rand.New())
	if err == nil {
		t.Error("NewSpecies(20, 18).NewPerm() should fail")
	}
}
