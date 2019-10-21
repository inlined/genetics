package genetics

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	errAlreadySet      = "%sFlag.Set(%s): already set to %s"
	errUnexpectedFn    = "%sFlag.Set(%s): unknown function name %s"
	errUnexpectedParam = "%sFlag.Set(%s): function %s does not accept parameters"
	errInvalidParam    = "%sFlag.Set(%s): param %s should %s"
)

var (
	flagFmt = regexp.MustCompile(`^(\w+)(\((\w*)\))?$`)
)

// NaturalSelectionFlag allows developers to pick a NaturalSelection
// strategy using flag.Value. Vallid values include:
// --flag=StochasticUniversalSampling
// --flag=RankedSelection
// --flag=TournamentSelection(3)
type NaturalSelectionFlag struct {
	selection NaturalSelection
}

func (f NaturalSelectionFlag) String() string {
	if f.selection == nil {
		return stochasticUniversalSampling
	}
	return f.selection.String()
}

// Set implements flag.Value
func (f *NaturalSelectionFlag) Set(s string) error {
	if f.selection != nil {
		return fmt.Errorf(errAlreadySet, "NaturalSelection", s, f)
	}

	match := flagFmt.FindStringSubmatch(s)
	fn, arg := match[1], match[3]

	switch fn {
	case stochasticUniversalSampling:
		f.selection = StochasticUniversalSampling{}
	case rankedSelection:
		f.selection = RankedSelection{}
	case tournamentSelection:
		n, err := strconv.Atoi(arg)
		if err != nil || n < 2 {
			return fmt.Errorf(errInvalidParam, "NaturalSelection", s, arg, "a whole number >= 2")
		}
		f.selection = TournamentSelection{Size: n}
	default:
		return fmt.Errorf(errUnexpectedFn, "NaturalSelection", s, fn)
	}

	if fn != tournamentSelection && arg != "" {
		return fmt.Errorf(errUnexpectedParam, "NaturalSelection", fn, arg)
	}

	return nil
}

// Get returns a parsed NaturalSelection value
func (f *NaturalSelectionFlag) Get() NaturalSelection {
	if f.selection == nil {
		return StochasticUniversalSampling{}
	}
	return f.selection
}

// CrossoverFlag allows users to set Crossover strategies. Can only
// be set once. Values include:
// --flag=MultiPointCrossover(2)
// --flag=WholeArithmeticRecombination
// --flag=DavisOrderCrossover
type CrossoverFlag struct {
	crossover Crossover
}

func (f CrossoverFlag) String() string {
	if f.crossover == nil {
		return fmt.Sprintf("%s(1)", multiPointCrossover)
	}
	return f.crossover.String()
}

// Set implements flag.Value
func (f *CrossoverFlag) Set(s string) error {
	if f.crossover != nil {
		return fmt.Errorf(errAlreadySet, "Crossover", s, f)
	}

	match := flagFmt.FindStringSubmatch(s)
	fn, arg := match[1], match[3]

	switch fn {
	case wholeArithmeticRecombination:
		f.crossover = WholeArithmeticRecombination{}
	case davisOrderCrossover:
		f.crossover = DavisOrderCrossover{}
	case multiPointCrossover:
		n, err := strconv.Atoi(arg)
		if err != nil || n < 2 {
			return fmt.Errorf(errInvalidParam, "Crossover", s, arg, "a whole number >= 2")
		}
		f.crossover = MultiPointCrossover{Points: n}
	default:
		return fmt.Errorf(errUnexpectedFn, "Crossover", s, fn)
	}

	if fn != multiPointCrossover && arg != "" {
		return fmt.Errorf(errUnexpectedParam, "Crossover", fn, arg)
	}

	return nil
}

// Get returns the parsed Crossover
func (f CrossoverFlag) Get() Crossover {
	if f.crossover == nil {
		return MultiPointCrossover{1}
	}
	return f.crossover
}

// MutationFlag allows developers to specify a Mutator strategy
// using flag.Value. Valid values include:
// --flag=RandomResettingMutation
// --flag=SwapMutation
// --flag=ScrambleMutation
// --flag=InversionMutation
type MutationFlag struct {
	mutator Mutator
}

func (f MutationFlag) String() string {
	if f.mutator == nil {
		return scrambleMutation
	}
	return f.mutator.String()
}

// Set implements flag.Value
func (f *MutationFlag) Set(s string) error {
	if f.mutator != nil {
		return fmt.Errorf(errAlreadySet, "Mutation", s, f)
	}

	match := flagFmt.FindStringSubmatch(s)
	fn, arg := match[1], match[3]

	switch fn {
	case randomResettingMutation:
		f.mutator = RandomResettingMutation{}
	case swapMutation:
		f.mutator = SwapMutation{}
	case scrambleMutation:
		f.mutator = ScrambleMutation{}
	case inversionMutation:
		f.mutator = InversionMutation{}
	default:
		return fmt.Errorf(errUnexpectedFn, "Mutation", s, fn)
	}

	if arg != "" {
		return fmt.Errorf(errUnexpectedParam, "Mutation", fn, arg)
	}

	return nil
}

// Get returns the parsed Mutator
func (f MutationFlag) Get() Mutator {
	if f.mutator == nil {
		return ScrambleMutation{}
	}
	return f.mutator
}
