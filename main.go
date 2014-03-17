package main

import (
  "fmt"
  "os"
  "twentyfortyeight/grid"
  "encoding/json"
)

func main() {
  grid := grid.Grid{
    Size: 4,
  }

  grid.Build()

  enc := json.NewEncoder(os.Stdout)

  enc.Encode(grid)

  fmt.Println(grid)
}
