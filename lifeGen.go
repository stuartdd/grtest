package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type LifeGenId int

type LifeCell struct {
	x, y int
	next *LifeCell
}

type LifeDeadCells struct {
	root  *LifeCell
	count int
}

type LifeGen struct {
	generations     []*LifeCell
	cellCount       []int
	currentGenId    LifeGenId
	countGen        int
	startTimeMillis int64
	timeMillis      int64
	genDone         func(l *LifeGen)
}

const (
	LIFE_GEN_1 LifeGenId = 0
	LIFE_GEN_2 LifeGenId = 1
)

func (ldc *LifeDeadCells) Add(x, y int) {
	t := ldc.root
	if t == nil {
		ldc.root = &LifeCell{x: x, y: y, next: nil}
		ldc.count = 1
		return
	}
	l := t
	for t != nil {
		if t.x == x && t.y == y {
			return
		}
		l = t
		t = t.next
	}
	l.next = &LifeCell{x: x, y: y, next: t}
	ldc.count++
}

func (lc *LifeCell) id() int64 {
	return int64(lc.x)*100000000 + int64(lc.y)
}

func (lc *LifeCell) String() string {
	return fmt.Sprintf("%016d", lc.id())
}

func NewLifeGen(genDone func(*LifeGen)) *LifeGen {
	generations := make([]*LifeCell, 2)
	generations[LIFE_GEN_1] = nil
	generations[LIFE_GEN_2] = nil
	cellCount := make([]int, 2)
	cellCount[LIFE_GEN_1] = 0
	cellCount[LIFE_GEN_2] = 0
	return &LifeGen{generations: generations, currentGenId: LIFE_GEN_1, startTimeMillis: 0, timeMillis: 0, countGen: 0, cellCount: cellCount, genDone: genDone}
}

func (lg *LifeGen) NextGen() {
	if lg.startTimeMillis != 0 {
		return
	}

	lg.startTimeMillis = time.Now().UnixMilli()
	deadCells := &LifeDeadCells{count: 0, root: nil}
	count := 0
	gen1 := lg.CurrentGenId()
	gen2 := lg.NextGenId()
	current := lg.CurrentGenRoot()
	cn := 0
	xc := 0
	yc := 0

	for current != nil {
		xc = current.x
		yc = current.y
		cn = lg.CountNear(xc, yc, deadCells)
		if cn == 2 || cn == 3 {
			count = count + lg.AddCell(xc, yc, gen2)
		}
		current = current.next
	}

	dc := deadCells.root
	for dc != nil {
		xc = dc.x
		yc = dc.y
		cn = lg.CountNear(xc, yc, nil)
		if cn == 3 {
			count = count + lg.AddCell(xc, yc, gen2)
		}
		dc = dc.next
	}

	//

	lg.countGen = lg.countGen + 1
	lg.timeMillis = time.Now().UnixMilli() - lg.startTimeMillis
	lg.cellCount[gen2] = count
	lg.startTimeMillis = 0
	lg.currentGenId = gen2
	lg.generations[gen1] = nil
	lg.cellCount[gen1] = 0
	if lg.genDone != nil {
		lg.genDone(lg)
	}
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
	c := lg.generations[lg.currentGenId]
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

func (lg *LifeGen) GetBounds() (int, int, int, int) {
	maxx := math.MinInt32
	maxy := math.MinInt32
	minx := math.MaxInt32
	miny := math.MaxInt32
	cell := lg.CurrentGenRoot()
	for cell != nil {
		if cell.x > maxx {
			maxx = cell.x
		}
		if cell.x < minx {
			minx = cell.x
		}
		if cell.y > maxy {
			maxy = cell.y
		}
		if cell.y < miny {
			miny = cell.y
		}
		cell = cell.next
	}
	return minx, miny, maxx, maxy
}

func (lg *LifeGen) AddCells(c []int, gen LifeGenId) int {
	n := 0
	for i := 0; i < len(c); i = i + 2 {
		n = n + lg.AddCell(c[i], c[i+1], gen)
	}
	lg.cellCount[gen] = lg.cellCount[gen] + n
	return n
}

func (lg *LifeGen) AddCell(x, y int, genId LifeGenId) int {
	toAdd := &LifeCell{x: x, y: y, next: nil}
	toAddid := toAdd.id()
	if lg.generations[genId] == nil { // First cell (root)
		lg.generations[genId] = toAdd
		return 1
	}
	var current *LifeCell
	var p *LifeCell

	current = lg.generations[genId]
	p = nil
	if current.id() == toAddid {
		return 0 // Already exists
	}
	if current.id() > toAddid {
		lg.generations[genId] = toAdd
		toAdd.next = current
		return 1
	}

	for current != nil {
		if current.id() == toAddid {
			return 0 // Already exists
		}
		if current.id() > toAddid {
			t := p.next
			p.next = toAdd
			toAdd.next = t
			return 1
		}
		p = current
		current = current.next
	}
	p.next = toAdd // Add to the last cell
	return 1
}

func (lg *LifeGen) String() string {
	c := lg.generations[lg.currentGenId]
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
	c := lg.generations[lg.currentGenId]
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

func (ldc *LifeDeadCells) String() string {
	t := ldc.root
	var sb strings.Builder
	sb.WriteString("DeadCells ")
	for t != nil {
		sb.WriteString(fmt.Sprintf("x:%d, y:%d ", t.x, t.y))
		t = t.next
	}
	return sb.String()
}
