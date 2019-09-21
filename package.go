// Package genetics implements swappable components designed for a variety of genetics
// algorithms. They were designed as a learning exercise and are based on the instructions
// in tutorialspoint.com/genetic_algorithms.
package genetics

// TODO: Incremental vs generational offspring

// private types for interface implementations
/*
type geometric struct{}
type exponential struct{}

var (
	// Geometric growth
	Geometric GeneScorer = geometric{}

	// Exponential growth
//	Exponential GeneScorer = exponential{}
)

// GeneScorer ...
type GeneScorer interface {
	ScoreFitness(trait int64, gene Gene) Fitness
}

func (geometric) ScoreFitness(trait int64, gene Gene) Fitness {
	if gene > maxGene {
		gene = maxGene
	}
	return Fitness(trait * int64(gene))
}

func (exponential) FitnessTrait(trait int64, genome Gene) Fitness {
	traitNegative := false
	if trait < 0 {
		traitNegative = true
		trait = -trait
	}
	genomeNegative := false
	if genome < 0 {
		genomeNegative = true
		genome = -genome
	}
	if genome > maxGene {
		genome = maxGene
	}
	var total Fitness = 1
	for genome != 0 {
		total *= Fitness(trait)
		genome--
	}
	if traitNegative != genomeNegative {
		total = -total
	}
	return total
}
*/
