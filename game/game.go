package game

import (
  "net/http"
  "encoding/json"
  "time"
  "grid/grid"

  "code.google.com/p/go.net/websocket"
)

//Constants for game status
const (
  pending = iota
  active = iota
  over = iota
)

type Player struct {
  ID string
  Name string
  Socket *websocket.Conn
  Grid Grid
}

func (p *Player)

type Game struct {
  ID string
  NumPlayers int
  Status int //Should use contstants above
  MaxScore int
  Players []*Player
}

func (g *Game) AddPlayer(p *Player) {
  g.Players = append(g.Players, p);

  if(len(g.Players) == g.NumPlayers) {
  }
}

//MergeRule type is meant to return an integer value for the second tile based on two tiles being merged. Semantically, a is merged onto b.
type MergeRule func (a, b *grid.Tile) int

//Score is meant to return an integer score and possibly update a grid based on a scoring rule.
type ScoreRule func (t *grid.Tile, g *grid.Grid) int

//Rule is a struct to pass logic into a game.
type Rule struct {
  Score ScoreRule
  Merge MergeRule
}

func (g *Game) NewPlayer(p *Player) {
  g.Players = append(g.Players, p)
}
