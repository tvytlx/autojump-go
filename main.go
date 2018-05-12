package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/sahilm/fuzzy"
)

var help = `Automatically jump to directory passed as an argument.

usage: autojump-go [--help] [--add DIRECTORY] [DIRECTORY]

positional arguments:
  DIRECTORY             directory to jump to

optional arguments:
  --help            show this help message and exit
  --add DIRECTORY
`
var dataPath string

// Data stores data file
type Data struct {
	value map[string]float64
	fp    *os.File
}

// Load value from file
func (d *Data) Load(dataPath string) {
	var _, err = os.Stat(dataPath)

	if os.IsNotExist(err) {
		var f, err = os.Create(dataPath)
		if err != nil {
			return
		}
		d.fp = f
	} else {
		f, err := os.OpenFile(dataPath, os.O_RDWR, 0644)
		if err != nil {
			return
		}
		d.fp = f
	}

	var b = new(bytes.Buffer)
	b.ReadFrom(d.fp)
	lines := strings.Split(b.String(), "\n")
	for _, line := range lines[:len(lines)-1] {
		s := strings.SplitN(line, ",", 2)
		weight, _ := strconv.ParseFloat(s[0], 64)
		d.value[s[1]] = weight
	}
}

// Close and save the data file
func (d *Data) Close() {
	defer d.fp.Close()

	// write
	d.fp.Truncate(0)
	for k, v := range d.value {
		d.fp.WriteString(fmt.Sprintf("%.3f,%s\n", v, k))
	}

	if err := d.fp.Sync(); err != nil {
		panic(err)
	}
}

// Add will add a new path or increase the weight if the path is
// already existed.
func (d *Data) Add(path string) {
	// if not exists
	if _, ok := d.value[path]; !ok {
		d.value[path] = 0
		return
	}

	d.value[path] = math.Sqrt(math.Pow(d.value[path], 2) + 100)
}

// GetPaths return all the paths in data file with the same order
func (d *Data) GetPaths() (paths []string) {
	for path := range d.value {
		paths = append(paths, path)
	}

	return paths
}

// GetWeights return all the weights in data file with the same order
func (d *Data) GetWeights() (weights []float64) {
	for _, weight := range d.value {
		weights = append(weights, weight)
	}

	return weights
}

// Match the most possible path from input
func Match(input string, data *Data) string {
	paths := []string{}
	for path := range data.value {
		paths = append(paths, path)
	}

	matches := fuzzy.Find(input, paths)
	bestChoice := "."
	maxWeights := float64(0)

	// no more than 10 candidate matches
	var maxNum int
	if maxNum = len(matches); maxNum > 10 {
		maxNum = 10
	}
	for _, match := range matches[:maxNum] {
		weight := data.value[match.Str]
		if weight > maxWeights {
			maxWeights = weight
			bestChoice = match.Str
		}
	}

	return bestChoice
}

func handle(flag, path string) string {
	d := Data{value: make(map[string]float64)}
	d.Load(dataPath)

	defer d.Close()

	switch flag {

	case "--add":
		d.Add(path)

	case "":
		return Match(path, &d)

	default:
		panic("args flag error")
	}

	return ""
}

func init() {
	if home := os.Getenv("HOME"); home != "" {
		dataPath = home + "/.autojump-go.txt"
	} else {
		panic("$HOME is empty")
	}
}

func main() {
	args := os.Args

	switch l := len(args); l {

	case 2: // autojump-go dir|--help
		switch args[1] {

		case "--help":
			fmt.Print(help)

		default:
			fmt.Print(handle("", args[1]))
		}

	case 3: // autojump-go --add dir
		fmt.Print(handle(args[1], args[2]))

	default:
		panic("args error")
	}
}