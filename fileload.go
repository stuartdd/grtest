package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

type RLE struct {
	fileName string
	decoded  string
	coords   []int64
	name     string
	owner    string
	comment  string
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

func (rle *RLE) rleDecodeString(rleStr string) (string, []int64) {
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
	coords := make([]int64, 0)
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

type MyFileDialog struct {
	file string
	path string
	wait bool
	err  error
}

func runMyFileDialog(parent fyne.Window, path string, callback func(string, error)) *MyFileDialog {
	dil := &MyFileDialog{wait: true, path: path}
	if path == "" {
		dil.GetCurrentPath()
	}
	fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if uc == nil {
			dil.err = fmt.Errorf("no file selected")
			dil.file = ""
		} else {
			dil.file = uc.URI().Path()
			dil.err = err
		}
		dil.wait = false
		if callback != nil {
			callback(dil.file, dil.err)
		}
	}, parent)
	l, err := dil.GetLastValueAsListableURI()
	if err != nil {
		dil.err = err
		return dil
	}
	fd.SetLocation(l)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".rle", ".RLE", ".Rle"}))
	fd.Show()
	dil.wait = true
	for dil.wait {
		time.Sleep(200 * time.Millisecond)
	}
	return dil
}

func (d *MyFileDialog) GetLastValueAsListableURI() (fyne.ListableURI, error) {
	u, err := storage.ParseURI("file://" + d.path)
	if err != nil {
		return nil, err
	}
	l, err := storage.ListerForURI(u)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (d *MyFileDialog) GetCurrentPath() {
	d.path = ""
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			if pair[0] == "PWD" {
				d.path = pair[1]
			} else {
				if pair[0] == "HOME" && d.path == "" {
					d.path = pair[1]
				}
			}
		}
	}
}
