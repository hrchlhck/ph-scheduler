package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func CheckError(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func ToInt(s string) int {
	n, e := strconv.Atoi(s)

	CheckError(e)

	return n
}

func ToFloat(s string, bitSize int) float64 {
	n, e := strconv.ParseFloat(s, bitSize)
	CheckError(e)
	return n
}

func GetFields(file string, singleLine bool) [][]string {
	data, err := os.ReadFile(file)

	CheckError(err)

	var sepData []string = strings.Split(string(data), "\n")

	var ret [][]string = make([][]string, len(sepData))

	for i, d := range sepData {
		var fields []string = strings.Fields(d)
		ret[i] = make([]string, len(fields))
		ret[i] = fields

		if singleLine {
			return ret
		}
	}

	return ret
}
