package genetics_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/inlined/genetics"

	"github.com/inlined/rand"
)

type searchParams struct {
	sampleRate     int
	generationSize int
	numGenerations int
}

func (p searchParams) newSoln() solution {
	return solution{
		chromosome: genetics.Chromosome{},
		samples:    make([]genetics.Fitness, p.numGenerations/p.sampleRate),
	}
}

type solution struct {
	chromosome genetics.Chromosome
	score      genetics.Fitness
	samples    []genetics.Fitness
}

func (s *solution) consider(chromosome genetics.Chromosome, fitness genetics.Fitness) {
	if s.chromosome.Genes == nil || fitness > s.score {
		s.score = fitness
		s.chromosome = chromosome
	}
}

func (s *solution) maybeSample(generation, sampleRate int) {
	if (generation+1)%sampleRate == 0 {
		s.samples[(generation+1)/sampleRate-1] = s.score
	}
}

type knapsackProblem struct {
	maxWeight int
	weights   []int
	values    []int
}

func scoreKnapsack(strategy []genetics.Gene, knapsack knapsackProblem) genetics.Fitness {
	value := 0
	weight := 0
	for i := 0; i < len(strategy); i++ {
		if strategy[i] == 0 {
			continue
		}
		if weight+knapsack.weights[i] > knapsack.maxWeight {
			continue
		}

		weight += knapsack.weights[i]
		value += knapsack.values[i]
	}

	return genetics.Fitness(value)
}

func solveKnapsackRandomly(params searchParams, knapsack knapsackProblem, rng rand.Rand) solution {
	soln := params.newSoln()
	s := genetics.NewSpecies(len(knapsack.weights), 1)
	for generation := 0; generation < params.numGenerations; generation++ {
		for i := 0; i < params.generationSize; i++ {
			chromosome, _ := s.NewRand(rng)
			fitness := scoreKnapsack(chromosome.Genes, knapsack)
			soln.consider(chromosome, fitness)
		}
		soln.maybeSample(generation, params.sampleRate)
	}
	return soln
}

func solveKnapsackGenetically(params searchParams, knapsack knapsackProblem, rng rand.Rand) solution {
	soln := params.newSoln()
	s := genetics.NewSpecies(len(knapsack.weights), 1)

	evolver := genetics.Evolver{
		ReplacementCount: params.generationSize / 2,
		MutationRate:     0.03,
		Selector:         genetics.StochasticUniversalSampling{},
		Crossover:        genetics.MultiPointCrossover{Points: 2},
		Mutator:          genetics.InversionMutation{},
	}

	if evolver.ReplacementCount%2 == 1 {
		evolver.ReplacementCount += 1
	}

	// Generation 0: random
	pop := make([]genetics.Chromosome, params.generationSize)
	fitness := make([]genetics.Fitness, params.generationSize)
	for i := range pop {
		pop[i], _ = s.NewRand(rng)
	}
	for generation := 0; generation < params.numGenerations; generation++ {
		for i := 0; i < params.generationSize; i++ {
			fitness[i] = scoreKnapsack(pop[i].Genes, knapsack)
			soln.consider(pop[i], fitness[i])
		}
		soln.maybeSample(generation, params.sampleRate)

		evolver.Evolve(rng, pop, fitness)
	}
	return soln
}

// Knapsack solutins are modeled as a binary genome.
// If Genen] is 0, it is skipped
// If Gene[n] is 1 and there is still available weight, it is taken
func TestKnapsackProblem(t *testing.T) {
	params := searchParams{
		sampleRate:     10,
		generationSize: 50,
		numGenerations: 100,
	}

	rng := rand.New()
	rng.Seed(time.Now().Unix())
	numItems := 50
	knapsack := knapsackProblem{
		maxWeight: 5000,
		weights:   make([]int, numItems),
		values:    make([]int, numItems),
	}

	for n := 0; n < numItems; n++ {
		knapsack.weights[n] = int(rng.Int31n(int32(knapsack.maxWeight * 10 / numItems)))
		knapsack.values[n] = int(rng.Int31n(int32(knapsack.maxWeight * 10 / numItems)))
	}

	randSolution := solveKnapsackRandomly(params, knapsack, rng)
	geneticSolution := solveKnapsackGenetically(params, knapsack, rng)

	fmt.Printf("Random growth: %v\n", randSolution.samples)
	fmt.Printf("Genetic growth: %v\n", geneticSolution.samples)

	if randSolution.score >= geneticSolution.score {
		t.Errorf("Evolution did not benefit over randomness: %d vs %d", randSolution.score, geneticSolution.score)
	}
}

// Scores are negative because there is no well-known maximum traversal
// This requires parent selection to support negative numbers, such as ranked selection or tournament selection
func scoreTravellingSalesperson(chromosome genetics.Chromosome, weights [][]int) genetics.Fitness {
	f := genetics.Fitness(0)
	for i := 1; i < chromosome.Species.NumGenes; i++ {
		from := chromosome.Genes[i-1]
		to := chromosome.Genes[i]
		// Graph is stored as a symmetric upper/left matrix
		if from > to {
			to, from = from, to
		}
		f -= genetics.Fitness(weights[to][from])
	}
	return f
}

func solveTravellingSalespersonRandomly(params searchParams, weights [][]int, rng rand.Rand) solution {
	soln := params.newSoln()
	s := genetics.NewSpecies(len(weights), genetics.Gene(len(weights)-1))
	for generation := 0; generation < params.numGenerations; generation++ {
		for i := 0; i < params.generationSize; i++ {
			chromosome, _ := s.NewPerm(rng)
			fitness := scoreTravellingSalesperson(chromosome, weights)
			soln.consider(chromosome, fitness)
		}
		soln.maybeSample(generation, params.sampleRate)
	}
	return soln
}

func solveTravellingSalespersonGenetically(params searchParams, weights [][]int, rng rand.Rand) solution {
	soln := solution{
		samples: make([]genetics.Fitness, params.numGenerations/params.sampleRate),
		score:   genetics.Fitness(math.MinInt64),
	}
	s := genetics.NewSpecies(len(weights), genetics.Gene(len(weights)-1))

	evolver := genetics.Evolver{
		ReplacementCount: params.generationSize / 2,
		MutationRate:     0.03,
		Selector:         genetics.TournamentSelection{Size: 4},
		Crossover:        genetics.DavisOrderCrossover{},
		Mutator:          genetics.ScrambleMutation{},
	}

	if evolver.ReplacementCount%2 == 1 {
		evolver.ReplacementCount += 1
	}

	// Generation 0: random
	pop := make([]genetics.Chromosome, params.generationSize)
	fitness := make([]genetics.Fitness, params.generationSize)
	for i := range pop {
		pop[i], _ = s.NewPerm(rng)
	}
	for generation := 0; generation < params.numGenerations; generation++ {
		for i := 0; i < params.generationSize; i++ {
			fitness[i] = scoreTravellingSalesperson(pop[i], weights)
			soln.consider(pop[i], fitness[i])
		}
		soln.maybeSample(generation, params.sampleRate)

		evolver.Evolve(rng, pop, fitness)
	}
	return soln
}

// Knapsack solutins are modeled as a binary genome.
// If Genen] is 0, it is skipped
// If Gene[n] is 1 and there is still available weight, it is taken
func TestTravellingSalesperson(t *testing.T) {
	params := searchParams{
		sampleRate:     10,
		generationSize: 50,
		numGenerations: 100,
	}

	rng := rand.New()
	rng.Seed(time.Now().Unix())
	const numCities = 50
	weights := make([][]int, numCities)
	for i := 0; i < numCities; i++ {
		weights[i] = make([]int, i)
		for j := 0; j < i; j++ {
			weights[i][j] = int(rng.Int31n(100))
		}
	}

	randSolution := solveTravellingSalespersonRandomly(params, weights, rng)
	geneticSolution := solveTravellingSalespersonGenetically(params, weights, rng)

	fmt.Printf("Random growth: %v\n", randSolution.samples)
	fmt.Printf("Genetic growth: %v\n", geneticSolution.samples)

	if randSolution.score >= geneticSolution.score {
		t.Errorf("Evolution did not benefit over randomness: %d vs %d", randSolution.score, geneticSolution.score)
	}
}
