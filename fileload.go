package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type RLE struct {
	fileName string
	decoded  string
	coords   []int64
	name     string
	owner    string
	comment  string
	minX     int64
	minY     int64
	maxX     int64
	maxY     int64
}

func NewRleFile(fileName string) (*RLE, error) {
	rle := &RLE{fileName: fileName}
	file, err := os.Open(rle.fileName)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	ln := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#N") {
			rle.name = strings.TrimSpace(line[2:])
		} else {
			if strings.HasPrefix(line, "#C") {
				if rle.comment == "" {
					rle.comment = strings.TrimSpace(line[2:])
				}
			} else {
				if strings.HasPrefix(line, "#O") {
					rle.owner = strings.TrimSpace(line[2:])
				} else {
					if !strings.HasPrefix(line, "#") {
						if ln > 0 {
							sb.WriteString(line)
						}
						ln++
					}
				}
			}
		}
	}
	rle.decoded, rle.coords = rle.rleDecodeString(sb.String())
	if len(rle.coords) == 0 {
		rle.minX = 0
		rle.minY = 0
		rle.maxX = 0
		rle.maxY = 0
	} else {
		rle.minX = math.MaxInt64
		rle.minY = math.MaxInt64
		rle.maxX = math.MinInt64
		rle.maxY = math.MinInt64
	}
	for i := 0; i < len(rle.coords); i = i + 2 {
		if rle.coords[i] < rle.minX {
			rle.minX = rle.coords[i]
		}
		if rle.coords[i] > rle.maxX {
			rle.maxX = rle.coords[i]
		}
		if rle.coords[i+1] < rle.minY {
			rle.minY = rle.coords[i+1]
		}
		if rle.coords[i+1] > rle.maxY {
			rle.maxY = rle.coords[i+1]
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return rle, nil
}

func (rle *RLE) RleCenter() (int64, int64) {
	if len(rle.coords) == 0 {
		return 0, 0
	}
	return (rle.maxX - rle.minX) / 2, (rle.maxY - rle.minY) / 2
}

func (rle *RLE) rleDecodeString(rleStr string) (string, []int64) {
	var result strings.Builder
	coords := make([]int64, 0)
	for len(rleStr) > 0 {
		letterIndex := strings.IndexFunc(rleStr, func(r rune) bool { return !unicode.IsDigit(r) })
		multiply := 1
		if letterIndex != 0 {
			multiply, _ = strconv.Atoi(rleStr[:letterIndex])
		}
		result.WriteString(strings.Repeat(string(rleStr[letterIndex]), multiply))
		rleStr = rleStr[letterIndex+1:]
	}
	out := result.String()

	var sb strings.Builder
	count := 0
	width := 0
	var y int64 = 0
	var x int64 = 0
	for _, c := range out {
		switch c {
		case '$':
			y++
			x = 0
			if count == 0 {
				for i := 0; i <= width; i++ {
					sb.WriteString("| ")
				}
				sb.WriteString("\n")
			} else {
				sb.WriteString("|\n")
				width = count
			}
			count = 0
		case 'b':
			sb.WriteString("| ")
			count++
			x++
		case 'o':
			coords = append(coords, x)
			coords = append(coords, y)
			sb.WriteString("|O")
			count++
			x++
		case '!':
			for i := 0; i <= (width - count); i++ {
				sb.WriteString("| ")
			}
		}
	}
	return sb.String(), coords
}

func (rle *RLE) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Owner  :%s\n", rle.owner))
	sb.WriteString(fmt.Sprintf("Name   :%s\n", rle.name))
	sb.WriteString(fmt.Sprintf("File   :%s\n", rle.fileName))
	sb.WriteString(fmt.Sprintf("Comment:%s\n", rle.comment))
	sb.WriteString(fmt.Sprintf("Cells  :%d ", len(rle.coords)/2))

	for i := 0; i < len(rle.coords); i = i + 2 {
		sb.WriteString(fmt.Sprintf("%3d, %3d ", rle.coords[i], rle.coords[i+1]))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%s\n", rle.decoded))
	return sb.String()
}

func PathToParentPath(p string) (string, error) {
	fd, err := os.Stat(p)
	if os.IsNotExist(err) {
		return "", err
	}
	if !fd.IsDir() {
		return "", fmt.Errorf("path %s is not a dir", p)
	}
	fpp := filepath.Dir(p)
	fd, err = os.Stat(fpp)
	if os.IsNotExist(err) {
		return "", err
	}
	if !fd.IsDir() {
		return "", fmt.Errorf("path %s is not a dir", p)
	}
	return fpp, nil
}
