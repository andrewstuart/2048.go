package main

import (
  "os"
  "fmt"
  "twentyfortyeight/grid"
  "encoding/json"
)

func printGrid (g grid.Grid) {
  for y:= g.Size - 1; y >= 0; y-- {
    for x:= 0; x < g.Size; x++ {
      if(g.Cells[x][y].Tile != nil) {
        fmt.Printf("%d ", g.Cells[x][y].Tile.Value)
      } else {
        fmt.Print("0 ")
      }
    }
    fmt.Println()
  }
}


func main() {

  grid := grid.NewGrid(4, 2)

  enc := json.NewEncoder(os.Stdout)

  enc.Encode(grid)

  printGrid(grid)

  grid.Shift(1)
  printGrid(grid)

  grid.Shift(2)
  printGrid(grid)

  grid.Shift(3)
  printGrid(grid)

  grid.Shift(4)
  printGrid(grid)
}
