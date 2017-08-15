package main

import (
	"Maximast"
	"flag"
	"fmt"
	"os"
)

var (
	areaFlag = flag.String(
		"area",
		"../data/area.txt",
		"Location of input data about the area")
	stationFlag = flag.String(
		"stations",
		"../data/stations.txt",
		"Location of the stations in the area")
	seperatorFlag = flag.String(
		"seperator",
		"/",
		"The argument the files are seperated with!")
	minimumCoverage = flag.Float64(
		"minCoverage",
		0.9,
		"The minium coverage before we can accept the result (default 0.9), below this value we only look for coverage. Starting from this value we also look for energy on top of coverage optimization. Needs to be between 0 and 1.")
	maxCoverageLook = flag.Float64(
		"maxCoverage",
		0.95,
		"The maximum coverage, at this value we will not look for coverage anymore and only look for the best energy! Needs to be between 0 and 1.")
)

//Genetic Algorithm to find the optimal energy/coverage transmission masts.
func main() {

	parseFlags()

	stations := maximast.ParseStationFile(*stationFlag, *seperatorFlag)
	area := maximast.ParseAreaFile(*areaFlag, *seperatorFlag, stations)
	population := maximast.NewPopulation(stations, area)

	population.RandomCoveragePoints()
	fmt.Printf("Sum Energy: %v\n", population.CalculateRadiusAndEnergy())

}

func parseFlags() {

	if 0 > *minimumCoverage || 1 < *minimumCoverage {
		fmt.Printf("Error: Minium coverage has to be between 0 and 1 (have %v) ", *minimumCoverage)
		os.Exit(1)
	}

	if 0 > *maxCoverageLook || 1 < *maxCoverageLook {
		fmt.Printf("Error: Maximum coverage has to be between 0 and 1 (have %v) ", *maxCoverageLook)
		os.Exit(1)
	}
}
