package maximast

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Station struct {
	coord  Coord
	coordM CoordM
	height float64
	rangeS float64
}

func ParseStationFile(filename string, seperator string) []Station {

	stationsFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening stations file %v \n", err)
		os.Exit(-1)
	}
	defer stationsFile.Close()

	r := bufio.NewReader(stationsFile)

	var coords []Coord
	var dists []float64

	coord, dist, err := readStation(r, seperator)
	for err == nil {
		coords = append(coords, coord)
		dists = append(dists, dist)
		//fmt.Println(s)
		coord, dist, err = readStation(r, seperator)
	}

	coordMs := CalculateToMeter(coords, FindSmallestCoord(coords))

	var stations []Station
	for index, coord := range coords {
		stations = append(stations, Station{
			coord,
			coordMs[index],
			dists[index],
			0	,
		})
	}

	return stations
}

// https://golang.org/pkg/bufio/#Reader.ReadLine
// x/y/dist => / is seperator for this exmpl.
// Returns Coord, float64(distance), error
func readStation(r *bufio.Reader, sep string) (Coord, float64, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	s := strings.Split(string(ln), sep)
	if len(s) != 3 {
		return Coord{}, 0, fmt.Errorf("Error: Length of line does not exist of three parts (%v)!", len(s))
	}

	x, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return Coord{}, 0, err
	}

	y, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		return Coord{}, 0, err
	}

	dist, err := strconv.ParseFloat(s[2], 64)
	if err != nil {
		return Coord{}, 0, err
	}

	// Create the new station!
	return Coord{x, y}, dist, nil
}

func (s Station) String() string {
	return fmt.Sprintf("Station: \n \tX: %v \n \tY: %v \n \tXm: %v \n \tYm: %v \n\tDist: %v \n", s.coord.x, s.coord.y, s.coordM.xM, s.coordM.yM, s.height)
}

func FindSmallestCoord(coords []Coord) Coord {
	if len(coords) == 0 {
		return Coord{}
	}

	var smallestCoord Coord = coords[0]

	for _, coord := range coords {
		if smallestCoord.x < coord.x {
			smallestCoord.x = coord.x
		}
		if smallestCoord.y < coord.y {
			smallestCoord.y = coord.y
		}
	}
	return smallestCoord
}

func FindSmallestStationCoords(stations []Station) Coord {
	if len(stations) == 0 {
		return Coord{}
	}

	var smallestCoord Coord = stations[0].coord

	for _, station := range stations {
		if station.coord.x < smallestCoord.x {
			smallestCoord.x = station.coord.x
		}
		if station.coord.y < smallestCoord.y {
			smallestCoord.y = station.coord.y
		}
	}
	return smallestCoord
}

func FindBiggestStationCoords(stations []Station) Coord {
	if len(stations) == 0 {
		return Coord{}
	}

	var smallestCoord Coord = stations[0].coord

	for _, station := range stations {
		if station.coord.x > smallestCoord.x {
			smallestCoord.x = station.coord.x
		}
		if station.coord.y > smallestCoord.y {
			smallestCoord.y = station.coord.y
		}
	}
	return smallestCoord
}
