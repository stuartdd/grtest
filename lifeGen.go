package main

import (
	"fmt"
	"strings"
)

type LifeCell struct {
	x, y int
	next *LifeCell
}

type LifeGen struct {
	generations []*LifeCell
	currentGen  int
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
	return &LifeGen{generations: generations, currentGen: 0}
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
