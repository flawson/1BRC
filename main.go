package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type StationEntry struct {
	TempTotal float64
	TempMin   float64
	TempMax   float64
	Entries   int64
}

var TempMap map[string]*StationEntry

func parseBytes(b []byte) {
	station, temp, ok := bytes.Cut(b, []byte(";"))
	if !ok {
		fmt.Printf("No semi-colon in line")
		return
	}

	stStr := string(station)
	stTemp, err := strconv.ParseFloat(strings.TrimSpace(string(temp)), 64)
	if err != nil {
		fmt.Printf("Unparsable float value: %v\n", err.Error())
		return
	}

	if se, ok := TempMap[stStr]; !ok {
		TempMap[stStr] = &StationEntry{
			TempTotal: stTemp,
			TempMin:   stTemp,
			TempMax:   stTemp,
			Entries:   1,
		}
	} else {
		if stTemp > se.TempMax {
			TempMap[stStr].TempMax = stTemp
		}
		if stTemp < se.TempMin {
			TempMap[stStr].TempMin = stTemp
		}
		TempMap[stStr].TempTotal += stTemp
		TempMap[stStr].Entries++
	}
}

func ReadFile(fh io.Reader) {
	reader := bufio.NewReader(fh)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			switch err {
			case io.EOF:
				return
			}
		}
		parseBytes(line)
	}
}

func main() {
	TempMap = make(map[string]*StationEntry)
	inputFile := os.Args[1]
	fh, err := os.OpenFile(inputFile, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	ReadFile(fh)
	for stationName, entry := range TempMap {
		fmt.Printf("%s=%.3f/%.3f/%.3f\n", stationName, entry.TempMin, entry.TempTotal/float64(entry.Entries), entry.TempMax)
	}
}
