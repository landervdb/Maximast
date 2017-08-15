package maximast

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Area struct {
	coords []Coord
	coordMs []CoordM
}



func ParseAreaFile(filename string, seperator string, stations []Station) Area {
	areaFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening area file %v \n", err)
		os.Exit(-1)
	}
	defer areaFile.Close()

	r := bufio.NewReader(areaFile)

	var coords []Coord

	c, err := readArea(r, seperator)
	for err == nil {
		coords = append(coords, c)
		//fmt.Println(c)
		c, err = readArea(r, seperator)
	}

	coordMs := CalculateToMeter(coords, FindSmallestStationCoords(stations))

	return Area{coords, coordMs}
}

// https://golang.org/pkg/bufio/#Reader.ReadLine
// x/y => / = seperator for this exmpl.
func readArea(r *bufio.Reader, sep string) (Coord, error) {
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
	if len(s) != 2 {
		return Coord{}, fmt.Errorf("Error: Length of line does not exist of two parts (%v)!", len(s))
	}

	x, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return Coord{}, err
	}

	y, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		return Coord{}, err
	}

	// Create the new coord!
	newCoord := Coord{x, y}

	return newCoord, err
}

func (a Area) String() string {
	var areaString []string
	areaString =  append(areaString, fmt.Sprint("Area: \n "))
	if a.coords != nil {
		for _, coord := range a.coords {
			areaString = append(areaString, fmt.Sprintf("\n \t x: %v \n \t y: %v \n", coord.x, coord.y))
		}
	}

	if a.coordMs != nil {
		for _,coordM := range a.coordMs{
			areaString = append(areaString, fmt.Sprintf("\n \t x: %v \n \t y: %v \n", coordM.xM, coordM.yM))
		}
	}

	return strings.Join(areaString, " ")
}
