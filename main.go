package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var isOneAndZero = regexp.MustCompile(`^[10]+$`).MatchString

func parseFlags() (size int, uniformRate, mutationRate float64, target []byte, elitist bool) {
	flag.IntVar(&size, "size", 50, "sets the population size")
	flag.Float64Var(&uniformRate, "uniform-rate", 0.5, "sets the crossover uniform rate")
	flag.Float64Var(&mutationRate, "mutation-rate", 0.015, "sets the mutation rate")
	flag.BoolVar(&elitist, "elitist", true, "sets elitism")
	targetString := flag.String("target", "110011", "establishes the target gene sequence, 1s and 0s sequence")
	flag.Parse()

	if *targetString == "" || !isOneAndZero(*targetString) {
		log.Fatal("invalid target")
	}

	target = stringToGenes(*targetString)

	return
}

func stringToGenes(target string) []byte {
	genes := make([]byte, len(target))
	for i, c := range target {
		if c == '1' {
			genes[i] = 1
		} else {
			genes[i] = 0
		}
	}
	return genes
}

func main() {
	size, uniformRate, mutationRate, target, elitism := parseFlags()
	pop := newInitializedPopulation(
		size,
		uniformRate,
		mutationRate,
		target,
		elitism,
	)

	for pop.fittest().fitness() < len(pop.target) {
		fmt.Printf("Generation: %d Fittest: %d\n%s\n\n", pop.generation, pop.fittest().fitness(), pop.fittest())
		pop.evolve()
	}

	fmt.Println("Solution found!")
	fmt.Println("Generations:", pop.generation)
	fmt.Println("Genes:", pop.fittest())
}

type population struct {
	individuals  []*individual
	target       []byte
	generation   int
	elitist      bool
	uniformRate  float64
	mutationRate float64
}

func newPopulation(size int, uniformRate, mutationRate float64, target []byte, elitist bool) *population {
	return &population{
		individuals: make([]*individual, size),
		target:      target,
		elitist:     elitist,
	}
}

func newInitializedPopulation(size int, uniformRate, mutationRate float64, target []byte, elitist bool) *population {
	pop := newPopulation(size, uniformRate, mutationRate, target, elitist)
	for i := 0; i < size; i++ {
		pop.individuals[i] = newIndividual(len(target), pop.fitness)
	}
	return pop
}

func (pop *population) evolve() {
	evolved := newPopulation(len(pop.individuals), pop.uniformRate, pop.mutationRate, pop.target, pop.elitist)

	startAt := 0
	if pop.elitist {
		evolved.individuals[0] = pop.fittest()
		startAt = 1
	}

	for i := startAt; i < len(pop.individuals); i++ {
		a := pop.individualByTournament(5)
		b := pop.individualByTournament(5)
		c := crossover(a, b, 0.5)
		c.mutate(0.015)
		evolved.individuals[i] = c
	}

	pop.generation++

	pop.individuals = evolved.individuals
}

func (pop *population) individualByTournament(size int) *individual {
	tournament := newPopulation(size, pop.uniformRate, pop.mutationRate, pop.target, pop.elitist)

	for i := 0; i < size; i++ {
		tournament.individuals[i] = pop.individuals[rand.Intn(len(pop.individuals))]
	}

	return tournament.fittest()
}

func (pop *population) fittest() *individual {
	fittest := pop.individuals[0]
	for _, ind := range pop.individuals[1:] {
		if fittest.fitness() < ind.fitness() {
			fittest = ind
		}
	}
	return fittest
}

func (pop *population) fitness(ind *individual) int {
	if len(ind.genes) != len(pop.target) {
		return 0
	}

	fitness := 0
	for i := 0; i < len(pop.target); i++ {
		if pop.target[i] == ind.genes[i] {
			fitness++
		}
	}

	return fitness
}

type individual struct {
	genes       []byte
	fitnessFunc func(*individual) int
}

func newIndividual(size int, fitness func(*individual) int) *individual {
	ind := &individual{
		genes:       make([]byte, size),
		fitnessFunc: fitness,
	}

	for i := 0; i < size; i++ {
		ind.genes[i] = byte(rand.Intn(2))
	}

	return ind
}

func crossover(a, b *individual, uniformRate float64) *individual {
	if len(a.genes) != len(b.genes) {
		return nil
	}

	c := newIndividual(len(a.genes), a.fitnessFunc)

	for i := 0; i < len(a.genes); i++ {
		if rand.Float64() <= uniformRate {
			c.genes[i] = a.genes[i]
		} else {
			c.genes[i] = b.genes[i]
		}
	}

	return c
}

func (ind *individual) fitness() int {
	return ind.fitnessFunc(ind)
}

func (ind *individual) mutate(mutationRate float64) {
	for i := 0; i < len(ind.genes); i++ {
		if rand.Float64() <= mutationRate {
			ind.genes[i] = byte(rand.Intn(2))
		}
	}
}

func (ind *individual) String() string {
	s := []byte{}
	for i := 0; i < len(ind.genes); i++ {
		if ind.genes[i] == 0 {
			s = append(s, []byte("0")...)
		} else {
			s = append(s, []byte("1")...)
		}
	}
	return string(s)
}
