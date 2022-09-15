package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type LifeCell struct {
	x, y int
	next *LifeCell
}

type LifeGen struct {
	generations     []*LifeCell
	currentGen      int
	countGen        int
	startTimeMillis int64
	timeMillis      int64
}

type RLE struct {
	fileName string
	decoded  string
	coords   []int
	name     string
	owner    string
	comment  string
}

var (
	generations []*LifeCell
)

func (lc *LifeCell) id() int64 {
	return int64(lc.x)*100000000 + int64(lc.y)
}

func (lc *LifeCell) String() string {
	return fmt.Sprintf("%016d", lc.id())
}

func NewLifeGen() *LifeGen {
	generations = make([]*LifeCell, 2)
	generations[0] = nil
	generations[1] = nil
	return &LifeGen{generations: generations, currentGen: 0, startTimeMillis: 0, timeMillis: 0, countGen: 0}
}

func (lg *LifeGen) NextGen() {
	if lg.startTimeMillis != 0 {
		return
	}
	lg.startTimeMillis = time.Now().UnixMilli()
	fmt.Printf("Life Gen Next")
	time.Sleep(time.Millisecond * 100)

	//
	//
	//

	lg.countGen = lg.countGen + 1
	lg.timeMillis = time.Now().UnixMilli() - lg.startTimeMillis
	lg.startTimeMillis = 0
}

func (lg *LifeGen) CurrentGen() *LifeCell {
	return lg.generations[lg.currentGen]
}

func (lg *LifeGen) CountNear(x, y int) int {
	count := lg.GetCell(x-1, y-1)
	count = count + lg.GetCell(x-1, y)
	count = count + lg.GetCell(x-1, y+1)
	count = count + lg.GetCell(x, y-1)
	count = count + lg.GetCell(x, y+1)
	count = count + lg.GetCell(x+1, y-1)
	count = count + lg.GetCell(x+1, y)
	count = count + lg.GetCell(x+1, y+1)
	return count
}

func (lg *LifeGen) GetCell(x, y int) int {
	f := &LifeCell{x: x, y: y}
	c := generations[lg.currentGen]
	for c != nil {
		if c.id() == f.id() {
			return 1
		}
		if c.id() > f.id() {
			return 0
		}
		c = c.next
	}
	return 0
}

func (lg *LifeGen) AddCells(c []int) {
	for i := 0; i < len(c); i = i + 2 {
		lg.AddCell(c[i], c[i+1])
	}
}

func (lg *LifeGen) AddCell(x, y int) {
	c := &LifeCell{x: x, y: y, next: nil}
	cid := c.id()
	gen := lg.currentGen
	if generations[gen] == nil {
		generations[gen] = c
	} else {
		var n *LifeCell
		var p *LifeCell

		n = generations[gen]
		p = nil
		if n.id() == cid {
			return
		}
		if n.id() > cid {
			t := generations[gen]
			generations[0] = c
			c.next = t
			return
		}

		for n != nil {
			if n.id() == cid {
				return
			}
			if n.id() > cid {
				t := p.next
				p.next = c
				c.next = t
				return
			}
			p = n
			n = n.next
		}
		p.next = c
	}
}

func (lg *LifeGen) String() string {
	c := generations[lg.currentGen]
	if c == nil {
		return "None"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Gen:%d\n", lg.currentGen))
	for c != nil {
		sb.WriteString(fmt.Sprintf("X:%d Y:%d id:%s\n", c.x, c.y, c))
		c = c.next
	}
	return sb.String()
}

func (lg *LifeGen) Short() string {
	c := generations[lg.currentGen]
	if c == nil {
		return "None"
	}
	var sb strings.Builder
	for c != nil {
		sb.WriteString(fmt.Sprintf("%d,%d ", c.x, c.y))
		c = c.next
	}
	return sb.String()
}

func (rle *RLE) Load(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	rle.fileName = fileName
	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	ln := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#N") {
			rle.name = strings.TrimSpace(line[2:])
		} else {
			if strings.HasPrefix(line, "#C") {
				rle.comment = strings.TrimSpace(line[2:])
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
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func (rle *RLE) rleDecodeString(rleStr string) (string, []int) {
	var result strings.Builder
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
	coords := make([]int, 0)
	count := 0
	width := 0
	y := 0
	x := 0
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
	for i := 0; i < len(rle.coords); i = i + 2 {
		sb.WriteString(fmt.Sprintf("%3d, %3d ", rle.coords[i], rle.coords[i+1]))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%s\n", rle.decoded))
	return sb.String()
}
