package main

import (
  "fmt"
  "twentyfortyeight/grid"
  "encoding/json"
  "code.google.com/p/go.net/websocket"
  "net/http"
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

type Move struct {
  Direction int
}

type Message struct {
  Name string
  Data Move
}

type GameSock struct {
  enc json.Encoder
  dec json.Decoder
  G grid.Grid
}

func NewGameSock (ws *websocket.Conn) {
  enc := json.NewEncoder(ws)
  dec := json.NewDecoder(ws)

  grid, move := grid.NewGrid(4, 2, 2048)

  for {
    select {
    case g := <-grid:
      fmt.Println("got something back")
      enc.Encode(g);
    default:
      var V Message
      err := dec.Decode(&V);

      if err == nil {
        move <- V.Data.Direction + 1
        enc.Encode(<-grid)
      } else {
        panic(err.Error())
      }
    }
  }
}


func main() {

  http.Handle("/game", websocket.Handler(NewGameSock))
  err := http.ListenAndServe(":12345", nil)
  if err != nil {
    panic ("ListenAndServe: " + err.Error())
  }
}
