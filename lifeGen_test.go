package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLifeRemoveCells(t *testing.T) {
	rle, err := NewRleFile("testdata/ibeacon.rle")
	if err != nil {
		t.Errorf("RLE File load failed. %e", err)
	}
	if len(rle.coords) != 36 {
		t.Errorf("ibeacon: Expected len(coords):%d actual len(coords):%d", 36, len(rle.coords))
	}
	lg := NewLifeGen(func(l *LifeGen) {
		fmt.Println(l)
	}, RUN_FOR_EVER)
	lg.AddCellsAtOffset(0, 0, 0b01, rle.coords, LIFE_GEN_1)
	n := lg.CountCells()
	if n != 18 {
		t.Errorf("ibeacon: Expected count:%d actual count:%d", 18, n)
	}
	lg.RemoveCells(0b01)
	n = lg.CountCells()
	if n != 0 {
		t.Errorf("ibeacon:Remove All Expected count:%d actual count:%d", 0, n)
	}
	lg.AddCellsAtOffset(0, 0, 0b01, rle.coords, LIFE_GEN_1)
	count := 0
	lg.VisitAllCells(func(lc *LifeCell) bool {
		if lc.x == 3 {
			lc.mode = 0
			count++
		}
		return true
	})

	lg.RemoveCells(0b01)
	n = lg.CountCells()
	if n != count {
		t.Errorf("ibeacon: Remove all except x=3 Expected count:%d actual count:%d", count, n)
	}

	lg.AddCellsAtOffset(0, 0, 0b01, rle.coords, LIFE_GEN_1)
	lg.VisitAllCells(func(lc *LifeCell) bool {
		fmt.Printf("%d ", lc.mode)
		return true
	})
	fmt.Println("  Y")
	count = 3
	lg.VisitAllCells(func(lc *LifeCell) bool {
		if count > 0 {
			lc.mode = 0
			count--
		} else {
			return false
		}
		return true
	})
	lg.VisitAllCells(func(lc *LifeCell) bool {
		fmt.Printf("%d ", lc.mode)
		return true
	})
	lg.RemoveCells(0b01)
	n = lg.CountCells()
	if n != 3 {
		t.Errorf("ibeacon: Remove all except first 3 Expected count:%d actual count:%d", 3, n)
	}
}
func TestLifeVisitAllCells(t *testing.T) {
	rle, err := NewRleFile("testdata/ibeacon.rle")
	if err != nil {
		t.Errorf("TestLifeVisitAllCells: RLE File load failed. %e", err)
	}
	if len(rle.coords) != 36 {
		t.Errorf("TestLifeVisitAllCells: Expected len(coords):%d actual len(coords):%d", 36, len(rle.coords))
	}
	lg := NewLifeGen(func(l *LifeGen) {
		fmt.Println(l)
	}, RUN_FOR_EVER)
	count := 0
	lg.VisitAllCells(func(lc *LifeCell) bool {
		count++
		return true
	})
	if count != 0 { // No cells added.
		t.Errorf("TestLifeVisitAllCells: Expected %d visits actual visits %d", 0, count)
	}
	count = 0
	lg.AddCellsAtOffset(0, 0, 0b01, rle.coords, LIFE_GEN_1)
	lg.VisitAllCells(func(lc *LifeCell) bool {
		count++
		return true
	})
	if count != 18 { // 36 coords 18 cell.
		t.Errorf("TestLifeVisitAllCells: Expected %d visits actual visits %d", 18, count)
	}
}

//    -2,-1, 0, 1, 2, 3, 4 X
// -2  .  .  .  .  .  .  .
// -1  .  y  y  y  y  y  .
//  0  .  y  x  x  x  y  .
//  1  .  y  y  y  y  y  .
//  2  .  .  .  .  .  .  .
//  3  .  .  .  .  .  .  .
//  4  .  .  .  .  .  .  .
//  Y
// DeadCells
// x:-1, y:-1 x:-1, y:0 x:-1, y:1 x:0, y:-1 x:0, y:1 x:1, y:-1 x:1, y:1 x:2, y:-1 x:2, y:1 x:3, y:-1 x:3, y:0 x:3, y:1
//
// Base Times:
// Time:[970 924 887 868] Total:3649 TestTime:3634
// Time:[955 908 867 857] Total:3587 TestTime:3572
// Time:[990 904 878 856] Total:3628 TestTime:3611
// Time:[974 898 871 860] Total:3603 TestTime:3588
// Time:[966 911 885 875] Total:3637 TestTime:3623
//
// Refactor 1 CountNearFast
// Time:[900 934 889 830] Total:3553 TestTime:3539
// Time:[903 885 854 836] Total:3478 TestTime:3464
//
// Inline LifeCell.id(). Order of magnitude!
// Time:[508 561 542 533] Total:2144
// Time:[521 556 549 524] Total:2150
//
// All int64 x and y (Slight regression)
// Time:[527 563 561 527] Total:2178
// Time:[542 553 566 534] Total:2195
//
// Call to index mid point
// Time:[520 558 563 579] Total:2220
// Time:[494 637 612 567] Total:2310
// Time:[526 579 604 585] Total:2294
//
// GetCell GetCellFast using index mid point
// Time:[383 349 355 331] Total:1418
// Time:[351 363 345 334] Total:1393
// Time:[326 350 344 332] Total:1352
//

var deadCells = &LifeDeadCells{count: 0, root: nil}

func TestLifeTiming(t *testing.T) {
	LifeTiming(t)
	LifeTiming(t)
	LifeTiming(t)
}

func LifeTiming(t *testing.T) {
	rle, err := NewRleFile("testdata/1234_synth.rle")
	if err != nil {
		t.Errorf("RLE File load failed. %e", err)
	}
	tims := make([]int64, 0)
	var timTot int64 = 0
	tim := time.Now().UnixMilli()
	lg := NewLifeGen(func(l *LifeGen) {
		//
		// Count overall times. Note this is in a separate thread to the NextGen process
		// So we delay at the end (see below)
		//
		tim = time.Now().UnixMilli() - tim
		timTot = timTot + tim
		tims = append(tims, tim)
		tim = time.Now().UnixMilli()
	}, RUN_FOR_EVER)
	lg.AddCellsAtOffset(0, 0, 0, rle.coords, lg.currentGenId)
	lg.AddCellsAtOffset(100, 100, 0, rle.coords, lg.currentGenId)
	lg.AddCellsAtOffset(200, 200, 0, rle.coords, lg.currentGenId)
	for i := 0; i < 4; i++ {
		lg.NextGen()
	}
	for len(tims) < 4 {
		time.Sleep(time.Millisecond * 10) // Wait as for all the NextGen callbacks to finish
	}
	fmt.Printf("// Time:%d Total:%d\n", tims, timTot)
}

func TestLifeNextGen(t *testing.T) {
	rle, err := NewRleFile("testdata/ibeacon.rle")
	if err != nil {
		t.Errorf("RLE File load failed. %e", err)
	}
	if len(rle.coords) != 36 {
		t.Errorf("ibeacon: Expected len(coords):%d actual len(coords):%d", 36, len(rle.coords))
	}
	lg := NewLifeGen(func(l *LifeGen) {
		fmt.Println(l)
	}, RUN_FOR_EVER)
	lg.AddCellsAtOffset(0, 0, 0, rle.coords, lg.currentGenId)
	lg.NextGen()
}

func TestLifeRLE(t *testing.T) {
	rle, err := NewRleFile("testdata/rats.rle")
	if err != nil {
		t.Errorf("RLE File load failed. %e", err)
	}
	assertStr(t, "$rats", rle.name)
	assertStr(t, "David Buckingham", rle.owner)
	assertStr(t, "testdata/rats.rle", rle.fileName)
	assertStr(t, "A period 6 oscillator found in 1972.", rle.comment)
	if len(rle.decoded) != 286 {
		t.Errorf("TestRle: Expected len(decoded):%d actual len(decoded):%d", 64, len(rle.decoded))
	}
	if len(rle.coords) != 64 {
		t.Errorf("TestRle: Expected len(coords):%d actual len(coords):%d", 64, len(rle.coords))
	}
	fmt.Println(rle)
}

func assertStr(t *testing.T, exp, act string) {
	if exp != act {
		t.Errorf("TestRle: Expected '%s' actual '%s'", exp, act)
	}
}

func TestLifeGenCountCells(t *testing.T) {
	lg := NewLifeGen(nil, RUN_FOR_EVER)
	lg.AddCellsAtOffset(0, 0, 0, []int64{2, 2}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 0)
	lg.AddCellsAtOffset(0, 0, 0, []int64{1, 1}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 1)
	lg.AddCellsAtOffset(0, 0, 0, []int64{1, 2}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 2)
	lg.AddCellsAtOffset(0, 0, 0, []int64{1, 3}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 3)
	lg.AddCellsAtOffset(0, 0, 0, []int64{2, 1}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 4)
	lg.AddCellsAtOffset(0, 0, 0, []int64{2, 3}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 5)
	lg.AddCellsAtOffset(0, 0, 0, []int64{3, 1}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 6)
	lg.AddCellsAtOffset(0, 0, 0, []int64{3, 2}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 7)
	lg.AddCellsAtOffset(0, 0, 0, []int64{3, 3}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 8)

	testCountNear(t, lg, 1, 1, 3)
	testCountNear(t, lg, 1, 2, 5)
	testCountNear(t, lg, 1, 3, 3)

	testCountNear(t, lg, 2, 1, 5)
	testCountNear(t, lg, 2, 3, 5)

	testCountNear(t, lg, 3, 1, 3)
	testCountNear(t, lg, 3, 2, 5)
	testCountNear(t, lg, 3, 3, 3)

	testCountNear(t, lg, 0, 0, 1)
	testCountNear(t, lg, 0, 1, 2)
	testCountNear(t, lg, 0, 2, 3)
	testCountNear(t, lg, 0, 3, 2)
	testCountNear(t, lg, 0, 4, 1)

	testCountNear(t, lg, 4, 0, 1)
	testCountNear(t, lg, 4, 1, 2)
	testCountNear(t, lg, 4, 2, 3)
	testCountNear(t, lg, 4, 3, 2)
	testCountNear(t, lg, 4, 4, 1)

}

func TestLifeGenAddCells(t *testing.T) {
	lg := NewLifeGen(nil, RUN_FOR_EVER)
	lg.AddCellsAtOffset(0, 0, 0, []int64{1, 1, 2, 2}, lg.currentGenId)
	testGen(t, lg, "Add Cells:", "1,1 2,2")
	lg.AddCellsAtOffset(0, 0, 0, []int64{0, 0, 1, 1, 2, 2}, lg.currentGenId)
	testGen(t, lg, "Add Cells:", "0,0 1,1 2,2")
}

func TestLifeGen(t *testing.T) {
	lg := NewLifeGen(nil, RUN_FOR_EVER)
	testGen(t, lg, "Empty gen:", "None")
	testGet(t, lg, 0, 0, 0)
	lg.AddCell(100, 100, 0, lg.currentGenId)
	testGen(t, lg, "Single cell:", "100,100")
	testGet(t, lg, 0, 0, 0)
	testGet(t, lg, 100, 100, 1)
	lg.AddCell(100, 100, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "100,100")
	lg.AddCell(100, 50, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "100,50 100,100")
	testGet(t, lg, 100, 50, 1)
	testGet(t, lg, 100, 100, 1)
	lg.AddCell(50, 50, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 100,50 100,100")
	lg.AddCell(50, 75, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 50,75 100,50 100,100")
	lg.AddCell(100, 75, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 50,75 100,50 100,75 100,100")
	lg.AddCell(100, 200, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 50,75 100,50 100,75 100,100 100,200")
	lg.AddCell(200, 200, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 50,75 100,50 100,75 100,100 100,200 200,200")
	lg.AddCell(0, 0, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "0,0 50,50 50,75 100,50 100,75 100,100 100,200 200,200")
	testGet(t, lg, 100, 50, 1)
	testGet(t, lg, 100, 100, 1)
	testGet(t, lg, 50, 74, 0)
}

func testGet(t *testing.T, lg *LifeGen, x, y int64, exp int) {
	b := lg.GetCell(x, y, deadCells)
	if b != exp {
		t.Errorf("TestGet: Expected '%d' actual '%d'", exp, b)
	}
}

func testCountNear(t *testing.T, lg *LifeGen, x, y int64, exp int) {
	b := lg.CountNear(x, y, deadCells)
	if b != exp {
		t.Errorf("CountNear: Expected '%d' actual '%d'", exp, b)
	}
}

func testGen(t *testing.T, lg *LifeGen, id, exp string) {
	s := lg.Short()
	s = strings.TrimSpace(s)
	if s != exp {
		t.Errorf("%s: Expected '%s' actual '%s'", id, exp, s)
	}
}
