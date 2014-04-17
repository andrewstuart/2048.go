package main

import (
  "fmt"
  "twentyfortyeight/grid"
  "encoding/json"
  "code.google.com/p/go.net/websocket"
  "net/http"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
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

type move int

type Message struct {
  Name string `json:"name"`
  Move move   `json:"move"`
}

type Player struct {
  Socket *websocket.Conn
  Grid grid.Grid
  History []*move //Persist later as a bulk SQL transaction rather than each time
}

type Event struct {
  Name string
  Data *grid.Grid
}

var db, _ = sql.Open("mysql", "twentyfortyeight:Pipeline97@tcp(mysql:3306)/twentyfortyeight?parseTime=true")

func newSqlGame() int64 {
  result, err := db.Exec(fmt.Sprintf("insert into games (size, maxScore) values (%d, %d)", 4, 2048))

  if err != nil {
    panic(err.Error())
  }

  id, err := result.LastInsertId()

  if err != nil {
    panic(err.Error())
  }

  return id
}

func NewGameSock (ws *websocket.Conn) {
  enc := json.NewEncoder(ws)
  dec := json.NewDecoder(ws)

  id := newSqlGame()

  grid, gridch, move := grid.NewGrid(4, 2, 2048)

  evt := Event{
    Name: "grid",
    Data: grid,
  }

  enc.Encode(evt)

  for {
    var V Message
    err := dec.Decode(&V);

    if err == nil {

      d := V.Data.Direction + 1
      _, err := db.Exec("INSERT INTO moves (direction, gameId) values ( ? , ? )", d, id)

      if err != nil {
        panic(err.Error())
      }

      move <- d
      enc.Encode(<-gridch)
    } else {
      panic(err.Error())
      return
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
