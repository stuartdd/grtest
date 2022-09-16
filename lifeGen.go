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

type LifeGenId int

type LifeCell struct {
	x, y int
	next *LifeCell
}

type LifeDeadCells struct {
	root *LifeCell
}

type LifeGen struct {
	generations     []*LifeCell
	currentGenId    LifeGenId
	countGen        int
	startTimeMillis int64
	timeMillis      int64
	cellCount       int
}

type RLE struct {
	fileName string
	decoded  string
	coords   []int
	name     string
	owner    string
	comment  string
}

const (
	LIFE_GEN_1 LifeGenId = 0
	LIFE_GEN_2 LifeGenId = 1
)

var (
	generations []*LifeCell
)

func (ldc *LifeDeadCells) Add(x, y int) {
	t := ldc.root
	ldc.root = &LifeCell{x: x, y: y, next: t}
}

func (ldc *LifeDeadCells) Root() *LifeCell {
	return ldc.root
}

func (lc *LifeCell) id() int64 {
	return int64(lc.x)*100000000 + int64(lc.y)
}

func (lc *LifeCell) String() string {
	return fmt.Sprintf("%016d", lc.id())
}

func NewLifeGen() *LifeGen {
	generations = make([]*LifeCell, 2)
	generations[LIFE_GEN_1] = nil
	generations[LIFE_GEN_2] = nil
	return &LifeGen{generations: generations, currentGenId: LIFE_GEN_1, startTimeMillis: 0, timeMillis: 0, countGen: 0}
}

func (lg *LifeGen) NextGen() {
	if lg.startTimeMillis != 0 {
		return
	}
	lg.startTimeMillis = time.Now().UnixMilli()
	fmt.Printf("Life Gen Next")
	time.Sleep(time.Millisecond * 100)

	//
	deadCells := &LifeDeadCells{}
	count := 0
	gen1 := lg.CurrentGenId()
	gen2 := lg.NextGenId()
	c := lg.CurrentGenRoot()
	for c != nil {
		x := c.x
		y := c.y
		i := lg.CountNear(x, y, deadCells)
		if i == 2 || i == 3 {
			lg.AddCell(x, y, gen2)
			count++
		}
		c = c.next
	}
	//

	lg.countGen = lg.countGen + 1
	lg.timeMillis = time.Now().UnixMilli() - lg.startTimeMillis
	lg.cellCount = count
	lg.startTimeMillis = 0
	lg.currentGenId = gen2
	generations[gen1] = nil
}

func (lg *LifeGen) NextGenId() LifeGenId {
	if lg.currentGenId == LIFE_GEN_1 {
		return LIFE_GEN_2
	}
	return LIFE_GEN_1
}

func (lg *LifeGen) CurrentGenId() LifeGenId {
	return lg.currentGenId
}

func (lg *LifeGen) CurrentGenRoot() *LifeCell {
	return lg.generations[lg.currentGenId]
}

func (lg *LifeGen) CountNear(x, y int, deadCells *LifeDeadCells) int {
	count := lg.GetCell(x-1, y-1, deadCells)
	count = count + lg.GetCell(x-1, y, deadCells)
	count = count + lg.GetCell(x-1, y+1, deadCells)
	count = count + lg.GetCell(x, y-1, deadCells)
	count = count + lg.GetCell(x, y+1, deadCells)
	count = count + lg.GetCell(x+1, y-1, deadCells)
	count = count + lg.GetCell(x+1, y, deadCells)
	count = count + lg.GetCell(x+1, y+1, deadCells)
	return count
}

func (lg *LifeGen) GetCell(x, y int, deadCells *LifeDeadCells) int {
	f := &LifeCell{x: x, y: y}
	c := generations[lg.currentGenId]
	for c != nil {
		if c.id() == f.id() {
			return 1
		}
		if c.id() > f.id() {
			if deadCells != nil {
				deadCells.Add(x, y)
			}
			return 0
		}
		c = c.next
	}
	if deadCells != nil {
		deadCells.Add(x, y)
	}
	return 0
}

func (lg *LifeGen) AddCells(c []int, gen LifeGenId) {
	for i := 0; i < len(c); i = i + 2 {
		lg.AddCell(c[i], c[i+1], gen)
	}
}

func (lg *LifeGen) AddCell(x, y int, genId LifeGenId) {
	c := &LifeCell{x: x, y: y, next: nil}
	cid := c.id()
	if generations[genId] == nil {
		generations[genId] = c
	} else {
		var n *LifeCell
		var p *LifeCell

		n = generations[genId]
		p = nil
		if n.id() == cid {
			return
		}
		if n.id() > cid {
			t := generations[genId]
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
	c := generations[lg.currentGenId]
	if c == nil {
		return "None"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Gen:%d\n", lg.currentGenId))
	for c != nil {
		sb.WriteString(fmt.Sprintf("X:%d Y:%d id:%s\n", c.x, c.y, c))
		c = c.next
	}
	return sb.String()
}

func (lg *LifeGen) Short() string {
	c := generations[lg.currentGenId]
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
