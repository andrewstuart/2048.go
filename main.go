package main

import (
	"fmt"
	"log"
	"math/rand"

	"astuart.co/2048.go/grid"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func printGrid(g grid.Grid) {
	fmt.Println()
	for y := g.Size - 1; y >= 0; y-- {
		for x := 0; x < g.Size; x++ {
			if g.Cells[x][y].Tile != nil {
				new := g.Cells[x][y].Tile.New

				if new {
					fmt.Printf("\033[1m")
				}
				fmt.Printf("%d ", g.Cells[x][y].Tile.Value)
				if new {
					fmt.Printf("\033[0m")
				}
			} else {
				fmt.Print("X ")
			}
		}
		fmt.Println()
	}
}

func main() {
	gr := grid.NewGrid(4, 2, 512)

	t := tview.NewTable().SetBorders(true).SetFixed(4, 4)
	txt := tview.NewTextView()

	draw := func() {
		for i := range gr.Cells {
			for j, c := range gr.Cells[i] {
				if c.Tile != nil {
					color := tcell.ColorRed
					if rand.Int63n(2)%2 == 1 {
						color = tcell.ColorOrange
					}

					t.SetCell(j, i, tview.NewTableCell(fmt.Sprint(c.Tile.Value)).SetBackgroundColor(color).SetTextColor(tcell.ColorBlack))
					continue
				}
				t.SetCell(j, i, tview.NewTableCell("  "))
			}
		}
		txt.SetText(fmt.Sprintf("Score: %d", gr.Score))
	}

	draw()
	fl := tview.NewFlex().AddItem(t, 50, 9, true).AddItem(txt, 20, 1, false).SetFullScreen(true).SetDirection(tview.FlexRow)
	a := tview.NewApplication().SetRoot(fl, true)
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var dir int
		switch event.Key() {
		case tcell.KeyUp:
			dir = grid.DirDown
		case tcell.KeyDown:
			dir = grid.DirUp
		case tcell.KeyLeft:
			dir = grid.DirLeft
		case tcell.KeyRight:
			dir = grid.DirRight
		case tcell.KeyCtrlSpace:
			gr.Reset()
			draw()
			return event
		default:
			return event
		}
		if err := gr.Shift(dir); err != nil {
			txt.SetText(fmt.Sprintf("Error: %s; Press ctrl-space to reset", err))
			return event
		}
		draw()
		return event
	})
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
