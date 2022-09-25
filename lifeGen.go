package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const indexMult = 100000000

type LifeGenId int

type LifeCell struct {
	x, y, ind int64
	next      *LifeCell
}

type LifeDeadCells struct {
	root  *LifeCell
	count int
}

type LifeGen struct {
	generations     []*LifeCell
	cellIndex       []*LifeCell
	cellCount       []int
	currentGenId    LifeGenId
	countGen        int
	startTimeMillis int64
	timeMillis      int64
	onGenDone       func(l *LifeGen)
	onGenStopped    func(l *LifeGen)
	runFor          int
}

const (
	LIFE_GEN_1 LifeGenId = 0
	LIFE_GEN_2 LifeGenId = 1
)

// Add a cell position to the dead cell list.
// Dont add duplicates
// Order is NOT important.
func (ldc *LifeDeadCells) addDeadCell(x, y int64) {
	t := ldc.root
	if t == nil {
		ldc.root = &LifeCell{x: x, y: y, next: nil, ind: x*indexMult + y}
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
	l.next = &LifeCell{x: x, y: y, next: t, ind: x*indexMult + y}
	ldc.count++
}

func NewLifeGen(genDone func(*LifeGen)) *LifeGen {
	lg := &LifeGen{generations: make([]*LifeCell, 2), cellIndex: make([]*LifeCell, 2), cellCount: make([]int, 2)}
	lg.Clear()
	return lg
}

func (lg *LifeGen) SetRunFor(n int, f func(*LifeGen)) {
	lg.onGenStopped = nil
	lg.runFor = n
	lg.onGenStopped = f
}

func (lg *LifeGen) GetRunFor() int {
	return lg.runFor
}

func (lg *LifeGen) IsRunning() bool {
	return lg.runFor > 0
}

func (lg *LifeGen) Clear() {
	lg.cellIndex[LIFE_GEN_1] = nil
	lg.cellIndex[LIFE_GEN_2] = nil
	lg.generations[LIFE_GEN_1] = nil
	lg.generations[LIFE_GEN_2] = nil
	lg.cellCount[LIFE_GEN_1] = 0
	lg.cellCount[LIFE_GEN_2] = 0
	lg.currentGenId = LIFE_GEN_1
	lg.countGen = 0
	lg.runFor = 0
	lg.startTimeMillis = 0
	lg.timeMillis = 0
}

// Ensure correct generation switching
func (lg *LifeGen) NextGenId() LifeGenId {
	if lg.currentGenId == LIFE_GEN_1 {
		return LIFE_GEN_2
	}
	return LIFE_GEN_1
}

// Scan current gen for the mid point cell (by ind).
// This is so we dont have to scan ALL the cells to find the given cell.
// Note this can only be done if the cells are ordered by their ind value.
// Used by GetCell and GetCellFast
func (lg *LifeGen) index(gen LifeGenId) *LifeCell {
	g := lg.generations[lg.currentGenId]
	if g == nil || g.next == nil {
		return nil // None or 1 cell. Do not index!
	}
	low := g.ind
	hi := low
	for g != nil {
		hi = g.ind
		g = g.next
	}
	mid := low + ((hi - low) / 2)
	g = lg.generations[lg.currentGenId]
	for g != nil {
		if g.ind >= mid {
			return g
		}
		g = g.next
	}
	return nil
}

// Scan the current generation and produce the next generation.
// Then swap generations so the next gen becomes the current gen
func (lg *LifeGen) NextGen() {

	lg.runFor = lg.runFor - 1
	if lg.runFor < 0 {
		if lg.onGenStopped != nil {
			f := lg.onGenStopped
			lg.onGenStopped = nil
			f(lg)
		}
		return
	}
	// If startTimeMillis is not 0 then we a concurrently calling NextGen before it is finished!
	//
	// Record start time
	// Set up the index
	// Clear the dead cell list
	// Get current and next generation ids.
	//
	lg.startTimeMillis = time.Now().UnixMilli()
	lg.cellIndex[lg.currentGenId] = lg.index(lg.currentGenId)
	deadCells := &LifeDeadCells{count: 0, root: nil}
	count := 0
	gen1 := lg.currentGenId
	gen2 := lg.NextGenId()

	//
	// scan current gen adding cells to next gen keeping track of any surrounding dead cells
	// for later processing.
	//
	cn := 0
	var xc int64 = 0
	var yc int64 = 0
	current := lg.generations[lg.currentGenId]
	for current != nil {
		xc = current.x
		yc = current.y
		cn = lg.CountNear(xc, yc, deadCells)
		//
		// Number of surrounding live cells
		// 		2 Means the cell continues in next gen
		// 		3 Means the cell is new in next gen
		//
		if cn == 2 || cn == 3 {
			count = count + lg.AddCell(xc, yc, gen2)
		}
		current = current.next
	}
	//
	// Now we have a list of all the surrounding dead cells we need to see if they are alive in next gen
	//
	dc := deadCells.root
	for dc != nil {
		xc = dc.x
		yc = dc.y
		cn = lg.CountNearFast(xc, yc)
		//
		// If a dead cell position has 3 live surrounding cells it is alive in the nex generation
		//
		if cn == 3 {
			count = count + lg.AddCell(xc, yc, gen2)
		}
		dc = dc.next
	}

	// Count the generation
	lg.countGen = lg.countGen + 1
	// Set the cell count
	lg.cellCount[gen2] = count

	// Swap generations and clear the next gen and next gen cell count
	lg.currentGenId = gen2
	lg.generations[gen1] = nil
	lg.cellCount[gen1] = 0

	// time the process and clear the start time
	lg.timeMillis = time.Now().UnixMilli() - lg.startTimeMillis
	lg.startTimeMillis = 0
	//
	// Call the function requested at the end of the Generation process
	// This is NOT included in the timing as it may involve GUI stuff
	// If is run as a separate thread so it will not block the generation processing
	//
	if lg.onGenDone != nil {
		go lg.onGenDone(lg)
	}
}

// Count cells around a dead cell to see if it will be live in the next gen
// Only need to count up to 4 so finish early if count > 3
func (lg *LifeGen) CountNearFast(x, y int64) int {
	count := lg.GetCellFast(x-1, y-1)
	count = count + lg.GetCellFast(x-1, y)
	count = count + lg.GetCellFast(x-1, y+1)
	count = count + lg.GetCellFast(x, y-1)
	if count > 3 {
		return count
	}
	count = count + lg.GetCellFast(x, y+1)
	if count > 3 {
		return count
	}
	count = count + lg.GetCellFast(x+1, y-1)
	if count > 3 {
		return count
	}
	count = count + lg.GetCellFast(x+1, y)
	if count > 3 {
		return count
	}
	count = count + lg.GetCellFast(x+1, y+1)
	return count
}

// Count cells around a live cell to see if it will be alive in the next generation
// Because we dont store dead cells we need to remember any dead cell positions
// surrounding the current cell so we can check them later.
func (lg *LifeGen) CountNear(x, y int64, deadCells *LifeDeadCells) int {
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

// Get cell returns a cell if it is in the current (sorted) live cell list.
// If it not then it is recorded as a dead cell for CountNear.
// No check is made on deadCellList parameter as it WILL never be nil.
// If not counting dead cells use GetCellFast.
// If the index has been found then start at the index if the cell is above or equal to it.
// Return 0 if not found, 1 if found.
func (lg *LifeGen) GetCell(x, y int64, deadCellList *LifeDeadCells) int {
	f := &LifeCell{x: x, y: y, ind: x*indexMult + y}
	var c *LifeCell
	if lg.cellIndex[lg.currentGenId] != nil && f.ind >= lg.cellIndex[lg.currentGenId].ind {
		// Index found and cell ind is >= to the indexed cell ind. Start from indexed cell
		c = lg.cellIndex[lg.currentGenId]
	} else {
		// Index not found or ind is < the indexed cell. Start from root cell.
		c = lg.generations[lg.currentGenId]
	}
	for c != nil {
		if c.ind == f.ind {
			// Cell found
			return 1
		}
		if c.ind > f.ind {
			// We passed where the cell should be. Dont waste time searching further
			// The cell is not found (assumed dead!) so add it to the dead cell list
			deadCellList.addDeadCell(x, y)
			return 0
		}
		c = c.next
	}
	// We got the end. The cell is not found (assumed dead!) so add it to the dead cell list
	deadCellList.addDeadCell(x, y)
	return 0
}

// Get cell returns a cell if it is in the current (sorted) live cell list.
// This is faster that GetCell as it does not cound surrounding dead cells.
// If the index has been found then start at the index if the cell is above or equal to it.
// Return 0 if not found, 1 if found.
func (lg *LifeGen) GetCellFast(x, y int64) int {
	f := &LifeCell{x: x, y: y, ind: x*indexMult + y}
	var c *LifeCell
	if lg.cellIndex[lg.currentGenId] != nil && f.ind >= lg.cellIndex[lg.currentGenId].ind {
		// Index found and cell ind is >= to the indexed cell ind. Start from indexed cell
		c = lg.cellIndex[lg.currentGenId]
	} else {
		// Index not found or ind is < the indexed cell. Start from root cell.
		c = lg.generations[lg.currentGenId]
	}
	for c != nil {
		if c.ind == f.ind {
			// Cell found
			return 1
		}
		if c.ind > f.ind {
			// We passed where the cell should be. Dont waste time searching further
			return 0
		}
		c = c.next
	}
	// We got the end. The cell is not found.
	return 0
}

// Get the minimum and maximum cell x,y positions
// Used to scale the GUI if needed.
func (lg *LifeGen) GetBounds() (int64, int64, int64, int64) {
	var maxx int64 = math.MinInt64
	var maxy int64 = math.MinInt64
	var minx int64 = math.MaxInt64
	var miny int64 = math.MaxInt64
	cell := lg.generations[lg.currentGenId]
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

// Add a list of cells to a specific generation.
// The cells can be added at a specific offset to allow new cells to be added in different places.
func (lg *LifeGen) AddCellsAtOffset(x, y int64, c []int64, gen LifeGenId) int {
	n := 0
	for i := 0; i < len(c); i = i + 2 {
		n = n + lg.AddCell(x+c[i], y+c[i+1], gen)
	}
	lg.cellCount[gen] = lg.cellCount[gen] + n
	return n
}

func (lg *LifeGen) RemoveCell(x, y int64, genId LifeGenId) {
	c := lg.generations[genId]
	if c == nil {
		return
	}
	if c.x == x && c.y == y {
		lg.generations[genId] = c.next
		return
	}
	p := c
	c = c.next
	for c != nil {
		if c.x == x && c.y == y {
			p.next = c.next
			return
		}
		p = c
		c = c.next
	}
}

// Add a cell to a specific generation defined by it's x,y value.
// No duplicates are added.
// Order is maintained.
//
//	   	The root cell has the lowest ind value
//			The last cess has the highes ind value
func (lg *LifeGen) AddCell(x, y int64, genId LifeGenId) int {
	lg.cellIndex[genId] = nil
	toAdd := &LifeCell{x: x, y: y, next: nil, ind: x*indexMult + y}
	toAddid := toAdd.ind
	if lg.generations[genId] == nil { // Generation has NO cells so cell becomes the root cell
		lg.generations[genId] = toAdd
		return 1
	}
	var current *LifeCell

	current = lg.generations[genId]
	if current.ind == toAddid {
		return 0 // First cell has the same id as teh cell to add so dont add it
	}

	if current.ind > toAddid {
		lg.generations[genId] = toAdd // New cell has ind less that the root so add it in front as the new root cell
		toAdd.next = current
		return 1
	}

	// Scan up the list and insert in order of ind.
	// Prev keeps track of the previous cell to make insertion easy.
	var prev *LifeCell
	for current != nil {
		if current.ind == toAddid {
			return 0 // Already exists so dont add it
		}
		if current.ind > toAddid {
			t := prev.next // Found the first cell with ind > new cell ind. So insert in front of it
			prev.next = toAdd
			toAdd.next = t
			return 1
		}
		// Next cell!
		prev = current
		current = current.next
	}
	// Not found so add it to the end.
	prev.next = toAdd
	return 1
}

// Debugging string utils
//
//	Return the cells ind as a 16 digit string
func (lc *LifeCell) String() string {
	return fmt.Sprintf("%016d", lc.ind)
}

// Debugging string utils
//
//	List a generation verbose
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

// Debugging string utils
//
//	List a generation (just x,y) values
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

// Debugging string utils
//
//	String the list of dead cells
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
