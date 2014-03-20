package grid

import (
  "math/rand"
  "time"
  "fmt"
)

var idSeed uint= 0;

//Position structure
type pos struct {
  X int
  Y int
}

//Tiles track value and position on the grid
type tile struct {
  ID uint
  Value int
  MergeHistory tileList
  Current pos
  Prev pos
  New bool
  Score chan int
}

func (t *tile) Move (dest *Cell) {
  t.Prev = t.Current
  t.Current = dest.Pos
}

func (t *tile) Merge (tn *tile) {
  t.Value += tn.Value //Future proofs for other merge rules. Fibonacci game??? Yes please. TODO
  t.MergeHistory = append(t.MergeHistory, tn)
}

func randVal() *rand.Rand {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))

  return r
}

func randTileVal() int {
  if randVal().Float64() > 0.9 {
    return 4 } else {
    return 2
  }
}

type tileList []*tile

func (t tileList) remove(tr *tile) tileList {
  for i, tl := range t {
    if(tl == tr) {
      return append(t[:i], t[i+1:]...)
    }
  }
  return t
}

//Grid tracks the tiles
type Grid struct {
  Size int
  StartCells int
  Tiles tileList
  Cells [][]*Cell
  totalScore int
  Score chan int
}

type Cell struct {
  Tile *tile
  Pos pos
}

//Create the grid
func (g *Grid) Build () {
  g.Cells = make([][]*Cell, g.Size)
  g.Tiles = make(tileList, 0)

  for i := range g.Cells {
    g.Cells[i] = make([]*Cell, g.Size)
    for j := range g.Cells[i] {
      cell := new(Cell)

      cell.Pos = pos{
        X: i,
        Y: j,
      }

      g.Cells[i][j] = cell
    }
  }

  start := g.StartCells

  for start != 0 {
    start = start - 1
    g.newTile()
  }
}

func (g *Grid) Reset() {
  for i := range g.Cells {
    for j := range g.Cells[i] {
      g.Cells[i][j].Tile = nil
    }
  }
  g.Tiles = make(tileList, 0)
}

func (g *Grid) newTile() {

  if(len(g.Tiles) == g.Size * g.Size) {
    if(!g.matchesRemaining()) {
      fmt.Println("YOU LOSE")
      g.Reset()
    }
  }

  avail := g.EmptyCells()

  if(len(avail) == 0) {
    return
  }
  //Two random values between 0 and Grid.Size
  i := randVal().Int() % len(avail)
  cell := avail[i]

  id := idSeed //Will concurrency screw with this?
  idSeed += 1

  newTile := tile{
    ID: id,
    Value: randTileVal(),
    MergeHistory: make(tileList, 0),
    Current: cell.Pos,
    New: true,
  }

  g.Tiles = append(g.Tiles, &newTile)
  cell.Tile = &newTile

  return
}

type vector struct {
  X int
  Y int
  Label string
}

var dMap = map[int]vector{
  1: vector{ //Up
    X: 0,
    Y: 1,
    Label: "Up",
  },
  2: vector{ //Right
    X: 1,
    Y: 0,
    Label: "Right",
  },
  3: vector{ //Down
    X: 0,
    Y: -1,
    Label: "Down",
  },
  4: vector{ //Left
    X: -1,
    Y: 0,
    Label: "Left",
  },
}

func (g *Grid) getEdge (v *vector) (start, end, delta int) {
  if(v.X == 1 || v.Y == 1) {
    end = - 1
    start = g.Size - 1
    delta = -1
  } else {
    end = g.Size
    start = 0
    delta = 1
  }

  return start, end, delta
}

//Get direction
//For table slice (vector perpendicular to direction)
//For row cell (vector opposite to direction)
//If cell has tile
//Add to 'new row' slice
//If tile value matches prev tile value
//Merge tile to prev 
//Remove tile ([2, 2, 4] .. [2 <-- 2, 4] .. [4, 4])
func (g *Grid) Shift(d int) (*Grid) {
  v := dMap[d]

  start, end, delta := g.getEdge(&v)

  for i:= start; i != end; i += delta {
    f2 := start
    for j:= start; j != end; j += delta {

      var cell, dest *Cell

      if(v.X == 0) { //Horizontal motion is zero. Direction is in Y
        cell = g.Cells[i][j]
        dest = g.Cells[i][f2]
      } else {
        cell = g.Cells[j][i]
        dest = g.Cells[f2][i]
      }

      if(cell.Tile != nil) { //If there's something here and somewhere to move it. Otherwise do nothing this iteration
        cell.Tile.New = false
        if dest.Pos != cell.Pos { //If they're not the same cells
          if dest.Tile != nil {
            if dest.Tile.Value == cell.Tile.Value { //If the value at the second finger matches the value at the current finger, merge.
              //Do a merge.
              dest.Tile.Merge(cell.Tile)
              //Now remove the old 
              g.Tiles = g.Tiles.remove(cell.Tile)
              //Always increment finger2 after a merge
              f2 += delta
              cell.Tile = nil
            } else {
              f2 += delta
              if(v.X == 0) {
                //Now reevaluate the cells changing.
                dest = g.Cells[i][f2]
              } else {
                dest = g.Cells[f2][i]
              }

              if dest.Pos != cell.Pos {
                cell.Tile.Move(dest)
                dest.Tile = cell.Tile
                cell.Tile = nil
              }
            }
          } else {
            cell.Tile.Move(dest)
            dest.Tile = cell.Tile
            cell.Tile = nil
          }
        }
      }
    }
  }
  g.newTile()
  return g
}

func (g *Grid) matchesRemaining() bool {
  for i := 0; i < g.Size; i++ {
    for j := 0 + i % 2; j < g.Size; j = j + 2 {
      for k:=1; k <= 4; k++ {
        v := dMap[k]
        cell := g.Cells[i][j]

        y := i + v.Y
        x := j + v.X

        if(y >= 0 && y < g.Size && x >= 0 && x < g.Size) {
          cmp := g.Cells[x][y]

          fmt.Println(cell.Tile.Value, cmp.Tile.Value)

          if(cell.Tile.Value == cmp.Tile.Value) {
            return true
          }
        }
      }
    }
  }

  return false
}

func (g *Grid) EmptyCells() (ret []*Cell) {
  for i:=0; i < g.Size; i++ {
    for j:=0; j < g.Size; j++ {
      c := g.Cells[i][j]
      if(c.Tile == nil) {
        ret = append(ret, c)
      }
    }
  }

  return ret
}

func NewGrid(s, c int) Grid {
  g := Grid{
    Size: s,
    StartCells: c,
  }

  g.Build()

  return g
}
