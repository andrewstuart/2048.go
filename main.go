package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"astuart.co/2048.go/grid"
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
	grid := grid.NewGrid(4, 2, 512)

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(grid.Tiles)

	t := tview.NewTable().SetBorder(true).SetTitle("2048")

	if err := tview.NewApplication().SetRoot(b, true).Run(); err != nil {
		log.Fatal(err)
	}
}
