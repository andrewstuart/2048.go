package game

import (
  "io"
  "net/http"
  "fmt"
  "math/rand"
  "time"

  "code.google.com/p/go.net/websocket"
)


type Socket struct {
}
type Player struct {
  ID string
  Name string
  Socket Socket
  Grid Grid
}

type Game struct {
  ID string
  Players []Player
}
