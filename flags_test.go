package genetics_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/inlined/genetics"
)

func TestNaturalSelectionFlag(t *testing.T) {
	for _, test := range []struct {
		tag  string
		flag string
		err  error
		val  genetics.NaturalSelection
	}{
		{
			tag:  "StochasticUniversalSampling",
			flag: "StochasticUniversalSampling",
			val:  genetics.StochasticUniversalSampling{},
		}, {
			tag:  "RankedSelection",
			flag: "RankedSelection",
			val:  genetics.RankedSelection{},
		}, {
			tag:  "TournamentSelection",
			flag: "TournamentSelection(2)",
			val:  genetics.TournamentSelection{Size: 2},
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			var flag genetics.NaturalSelectionFlag
			err := flag.Set(test.flag)
			if err == nil && test.err != nil {
				t.Errorf("expected error %s", err)
				return
			}
			if err != nil && test.err == nil {
				t.Errorf("failed with err %s", err)
				return
			}
			if err != nil && test.err != nil && err.Error() != test.err.Error() {
				t.Errorf("expected error %s got error %s", test.err, err)
				return
			}

			if diff := cmp.Diff(test.val, flag.Get()); diff != "" {
				t.Errorf("failed to parse %s; got=%s want=%s diff=%s", test.flag, flag.Get(), test.val, diff)
			}
		})
	}
}
