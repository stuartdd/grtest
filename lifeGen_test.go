package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

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
//
func TestLifeTiming(t *testing.T) {
	rle := &RLE{}
	rle.Load("testdata/1234_synth.rle")
	tims := make([]int64, 0)
	timsPos := 0
	var timTot int64 = 0
	tim := time.Now().UnixMilli()
	lg := NewLifeGen(func(l *LifeGen) {
		tim = time.Now().UnixMilli() - tim
		timTot = timTot + tim
		tims = append(tims, tim)
		tim = time.Now().UnixMilli()
		timsPos++
	})
	lg.AddCellsAtOffset(0, 0, rle.coords, lg.CurrentGenId())
	lg.AddCellsAtOffset(100, 100, rle.coords, lg.CurrentGenId())
	lg.AddCellsAtOffset(200, 200, rle.coords, lg.CurrentGenId())
	tim2 := time.Now().UnixMilli()
	for i := 0; i < 4; i++ {
		lg.NextGen()
	}
	fmt.Printf("Time:%d Total:%d TestTime:%d\n", tims, timTot, time.Now().UnixMilli()-tim2)
}

func TestLifeNextGen(t *testing.T) {
	rle := &RLE{}
	rle.Load("testdata/ibeacon.rle")
	if len(rle.coords) != 36 {
		t.Errorf("ibeacon: Expected len(coords):%d actual len(coords):%d", 36, len(rle.coords))
	}
	lg := NewLifeGen(func(l *LifeGen) {
		fmt.Println(l)
	})
	lg.AddCellsAtOffset(0, 0, rle.coords, lg.CurrentGenId())
	lg.NextGen()
}

func TestLifeRLE(t *testing.T) {
	rle := &RLE{}
	rle.Load("testdata/rats.rle")
	assertStr(t, "$rats", rle.name)
	assertStr(t, "David Buckingham", rle.owner)
	assertStr(t, "testdata/rats.rle", rle.fileName)
	assertStr(t, "www.conwaylife.com/wiki/index.php?title=$rats", rle.comment)
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
	lg := NewLifeGen(nil)
	lg.AddCellsAtOffset(0, 0, []int{2, 2}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 0)
	lg.AddCellsAtOffset(0, 0, []int{1, 1}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 1)
	lg.AddCellsAtOffset(0, 0, []int{1, 2}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 2)
	lg.AddCellsAtOffset(0, 0, []int{1, 3}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 3)
	lg.AddCellsAtOffset(0, 0, []int{2, 1}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 4)
	lg.AddCellsAtOffset(0, 0, []int{2, 3}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 5)
	lg.AddCellsAtOffset(0, 0, []int{3, 1}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 6)
	lg.AddCellsAtOffset(0, 0, []int{3, 2}, lg.currentGenId)
	testCountNear(t, lg, 2, 2, 7)
	lg.AddCellsAtOffset(0, 0, []int{3, 3}, lg.currentGenId)
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
	lg := NewLifeGen(nil)
	lg.AddCellsAtOffset(0, 0, []int{1, 1, 2, 2}, lg.currentGenId)
	testGen(t, lg, "Add Cells:", "1,1 2,2")
	lg.AddCellsAtOffset(0, 0, []int{0, 0, 1, 1, 2, 2}, lg.currentGenId)
	testGen(t, lg, "Add Cells:", "0,0 1,1 2,2")
}

func TestLifeGen(t *testing.T) {
	lg := NewLifeGen(nil)
	testGen(t, lg, "Empty gen:", "None")
	testGet(t, lg, 0, 0, 0)
	lg.AddCell(100, 100, lg.currentGenId)
	testGen(t, lg, "Single cell:", "100,100")
	testGet(t, lg, 0, 0, 0)
	testGet(t, lg, 100, 100, 1)
	lg.AddCell(100, 100, lg.currentGenId)
	testGen(t, lg, "Add cell:", "100,100")
	lg.AddCell(100, 50, lg.currentGenId)
	testGen(t, lg, "Add cell:", "100,50 100,100")
	testGet(t, lg, 100, 50, 1)
	testGet(t, lg, 100, 100, 1)
	lg.AddCell(50, 50, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 100,50 100,100")
	lg.AddCell(50, 75, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 50,75 100,50 100,100")
	lg.AddCell(100, 75, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 50,75 100,50 100,75 100,100")
	lg.AddCell(100, 200, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 50,75 100,50 100,75 100,100 100,200")
	lg.AddCell(200, 200, lg.currentGenId)
	testGen(t, lg, "Add cell:", "50,50 50,75 100,50 100,75 100,100 100,200 200,200")
	lg.AddCell(0, 0, lg.currentGenId)
	testGen(t, lg, "Add cell:", "0,0 50,50 50,75 100,50 100,75 100,100 100,200 200,200")
	testGet(t, lg, 100, 50, 1)
	testGet(t, lg, 100, 100, 1)
	testGet(t, lg, 50, 74, 0)
}

func testGet(t *testing.T, lg *LifeGen, x, y, exp int) {
	b := lg.GetCell(x, y, nil)
	if b != exp {
		t.Errorf("TestGet: Expected '%d' actual '%d'", exp, b)
	}
}

func testCountNear(t *testing.T, lg *LifeGen, x, y, exp int) {
	b := lg.CountNear(x, y, nil)
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
