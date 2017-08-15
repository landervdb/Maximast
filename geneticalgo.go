package maximast

import (
	"flag"
	"math"
	"fmt"
	"math/rand"
)

var (
	maxIFlag = flag.Float64(
		"maxI",
		43.0,
		"Maximum amount of power usage per station! default on 43!")
	maxNewParents = flag.Int(
		"maxParents",
		1000,
		"Maximum amount of parents in populations! Default 1000!")
	bandwidthFlag = flag.Bool(
		"lowBandwidth",
		false,
		"Low bandwidth or high bandwidth(default high)")
	sizeRandomPoints = flag.Int(
		"numbPointsCoverage",
		40,
		"Number of points to cover the area with random points size x size! (default 40 or 40 x 40)")

)

type Population struct {
	area           Area
	stations       []Station
	constants      []float64
	multiple       []float64
	coveragePoints []Coord
	fits           []Fit
	mainIndiviual  []Parent
}

type Fit struct {
	parents      []Parent
	ID           int
	fitness      float64
	oldFitness   float64
	energy       float64
	coveragePerc float64
}

type Coord struct {
	x float64
	y float64
}

type CoordM struct {
	xM float64
	yM float64
}

type Parent struct {
	power float64
}

func (coord Coord) String() string {
	return fmt.Sprintf("COORD: \n \tX: %v \n \tY: %v", coord.x, coord.y)
}

func (coord CoordM) String() string {
	return fmt.Sprintf("COORD in METERS: \n \tX: %v \n \tY: %v", coord.xM, coord.yM)
}

func NewPopulation(stations []Station, area Area) Population {

	constants, multiple := precalc_values(area)

	var fits []Fit
	for i := 0; i < *maxNewParents; i++ {
		var parents []Parent
		for j := 0; j < len(stations); j++ {
			parents = append(parents, Parent{
				power: math.Mod(rand.Float64(), *maxIFlag),
			})
		}
		fits = append(fits, Fit{
			parents: parents,
			ID:      i,
		})
	}

	var mainIndividual []Parent
	for i := 0; i < len(stations); i++ {
		mainIndividual = append(mainIndividual, Parent{*maxIFlag})
	}

	return Population{
		area,
		stations,
		constants,
		multiple,
		[]Coord{},
		fits,
		mainIndividual}
}

func precalc_values(area Area) ([]float64, []float64) {
	polyCorners := len(area.coords)
	var j int = polyCorners - 1

	constant := make([]float64, polyCorners)
	multiple := make([]float64, polyCorners)
	for i := 0; i < polyCorners; i++ {
		if area.coordMs[j].xM == area.coordMs[i].yM {
			constant[i] = area.coordMs[i].xM
			multiple[i] = 0
		} else {
			constant[i] = area.coordMs[i].xM -
				(area.coordMs[i].yM*area.coordMs[j].xM)/
					(area.coordMs[j].yM-area.coordMs[i].yM) +
				(area.coordMs[i].yM*area.coordMs[i].xM)/
					(area.coordMs[j].yM-area.coordMs[i].yM)
			multiple[i] = (area.coordMs[j].xM - area.coordMs[i].xM) / (area.coordMs[j].yM - area.coordMs[i].yM)
		}
		j = i
	}
	return constant, multiple
}

func (p *Population) CalculateRadiusAndEnergy() float64 {
	var I float64
	var sum_energy float64 = 0
	var tmpenergy float64 = 0

	stations := (*p).stations

	mainIndividual := (*p).mainIndiviual

	for index, station := range stations {
		//Power of every station!
		I = mainIndividual[index].power

		//Calculate energy from stations.
		tmpenergy = p.calculateEnergy(I)
		sum_energy += tmpenergy

		//Calculate range van de masten.
		if tmpenergy == 0 {
			station.rangeS = 0
		} else {
			station.rangeS = p.calculateRange(station.height, I)
		}
	}

	return sum_energy
}

func (p *Population) calculateRange(height float64, I float64) float64 {
	var PL, Lbsh, ka, kd float64

	PL = 0.0
	if *bandwidthFlag {
		PL = I + 98.8
	} else {
		PL = I + 78.3
	}

	//Lbsh,ka,kd: groter of kleiner dan 12 de hoogte!
	if height > 12 {
		ka = 54.0
		kd = 18.0
		Lbsh = -18.0 * math.Log10(1.0+height-12.0)
	} else {
		Lbsh = 0.0
		ka = 54.0 - 0.8*(height-12.0)
		kd = 18.0 - 15*((height-12.0)/12.0)
	}

	return 1000.0 * math.Pow(float64(10.0), float64((PL-Lbsh-ka-121.34)/(20.0+kd)))
}

func (p *Population) calculateEnergy(I float64) float64 {
	if I == 0 {
		return 0
	} else {
		return ((300.0 * (math.Pow(float64(10.0), float64(I/10.0)))) / 12800.0) + 1205.0
	}
}



func CalculateToMeter(coords []Coord, smallestCoord Coord) []CoordM {
	var mb, ml float64
	var x float64

	x = smallestCoord.x * math.Pi / 180

	mb = 111132.92 - 559.82*math.Cos(2*x) + 1.175*math.Cos(4*x) - 0.0023*math.Cos(6*x)
	ml = 111412.84*math.Cos(x) - 93.5*math.Cos(3*x) - 0.118*math.Cos(5*x)

	var coordMs []CoordM
	for _, coord := range coords {
		coordMs = append(coordMs, CoordM{
			(coord.x - smallestCoord.x) * mb,
			(coord.y - smallestCoord.y) * ml,
		})
	}
	return coordMs
}



func (p *Population) RandomCoveragePoints() {
	stations := (*p).stations
	if len(stations) == 0 {
		return
	}

	minCoord := FindSmallestStationCoords(stations)
	maxCoord := FindBiggestStationCoords(stations)

	var coverageCoords []Coord

	tmp := Coord{}
	var i, j float64
	for i = 0; int(i) < *sizeRandomPoints; i++ {
		for j = 0; int(j) < *sizeRandomPoints; j++ {
			tmp.x = minCoord.x + ((maxCoord.x-minCoord.x)/float64(*sizeRandomPoints))*j
			tmp.y = maxCoord.y - ((maxCoord.y)/float64(*sizeRandomPoints))*i
			if p.checkCoordInArea(tmp) {
				coverageCoords = append(coverageCoords, Coord{
					x: minCoord.x + ((maxCoord.x-minCoord.x)/float64(*sizeRandomPoints))*j,
					y: maxCoord.y - ((maxCoord.y)/float64(*sizeRandomPoints))*i,
				})
			}
		}
	}
	(*p).coveragePoints = coverageCoords

}

func (p *Population)checkCoordInArea(coord Coord) bool {
	multiple := (*p).multiple
	constant := (*p).constants
	area := (*p).area
	x := coord.x
	y := coord.y
	polyCorners := len(area.coords)
	var j int = polyCorners - 1
	var oddNodes bool = false

	for i := 0; i < polyCorners; i++ {
		if ((area.coordMs[i].yM < y) &&
			(area.coordMs[j].yM >= y)) ||
			((area.coordMs[j].yM < y) &&
				(area.coordMs[i].yM >= y)) {
			oddNodes = ((y*multiple[i]+constant[i] < x) != oddNodes) // XOR
		}
		j = i
	}
	return oddNodes
}

