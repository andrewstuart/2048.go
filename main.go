package main

import (
  "os"
  "fmt"
  "twentyfortyeight/grid"
  "encoding/json"
)

func printGrid (g grid.Grid) {
  fmt.Println()
  for y:= g.Size - 1; y >= 0; y-- {
    for x:= 0; x < g.Size; x++ {
      if(g.Cells[x][y].Tile != nil) {
        new := g.Cells[x][y].Tile.New

        if(new) {
          fmt.Printf("\033[1m")
        }
        fmt.Printf("%d ", g.Cells[x][y].Tile.Value)
        if(new) {
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

    for i:=0; i < 10000; i++ {
      grid.Shift(i % 4 + 1)
    }
}
