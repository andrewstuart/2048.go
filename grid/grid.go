package grid

import (
  "fmt"
  "math/rand"
  "time"
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
  MergeHistory []*tile
  Current pos
  Prev pos
  New bool
}

func (t *tile) Move (dest *Cell) {
  t.New = false;
  fmt.Printf("Moving tile at %d to Cell at %d", t.Current, dest.Pos)
  fmt.Println()
  t.Prev = t.Current
  t.Current = dest.Pos
}

func (t *tile) Merge (tn *tile) {
  fmt.Printf("Merging tile %d with tile %d", t.Current, tn.Current)
  t.Value += tn.Value //Future proofs for other merge rules. Fibonacci game??? Yes please. TODO
  fmt.Printf("New value is %d", t.Value)
  fmt.Println()
  t.MergeHistory = append(t.MergeHistory, tn)
}

func randVal() *rand.Rand {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))

  return r
}

func randTileVal() int {
  if randVal().Float64() > 0.9 {
    return 4
  } else {
    return 2
  }
}

//Grid tracks the tiles
type Grid struct {
  Size int
  StartCells int
  Tiles []*tile
  Cells [][]*Cell
}

type Cell struct {
  Tile *tile
  Pos pos
}

//Create the grid
func (g *Grid) Build () {
  g.Cells = make([][]*Cell, g.Size)
  g.Tiles = make([]*tile, 0)

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
    g.NewTile()
  }
}

func (g *Grid) Reset() {
  for i := range g.Cells {
    for j := range g.Cells[i] {
      g.Cells[i][j].Tile = nil
    }
  }
  g.Tiles = make([]*tile, 0)
}

func (g *Grid) NewTile() tile {

  avail := g.EmptyCells()
  //Two random values between 0 and Grid.Size
  i := randVal().Int() % len(avail)
  cell := avail[i]

  id := idSeed //Will concurrency screw with this?
  idSeed += 1

  newTile := tile{
    ID: id,
    Value: randTileVal(),
    MergeHistory: make([]*tile, 0),
    Current: cell.Pos,
    New: true,
  }

  g.Cells[cell.Pos.X][cell.Pos.Y].Tile = &newTile

  g.Tiles = append(g.Tiles, &newTile)

  return newTile
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
func (g *Grid) Shift(d int) *Grid {
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

      if(cell.Tile != nil && dest.Pos != cell.Pos) { //If there's something here and somewhere to move it. Otherwise do nothing this iteration
        if(dest.Tile != nil && dest.Tile.Value == cell.Tile.Value) { //If the value at the second finger matches the value at the current finger, merge.
          //Do a merge.
          dest.Tile.Merge(cell.Tile)
          //Always increment finger 2 after a merge
          f2 += delta
        } else if dest.Tile != nil {
          f2 += delta //TODO I wish this was prettier.

          //Now reevaluate the cells changing.
          if(v.X == 0) {
            dest = g.Cells[i][f2]
          } else {
            dest = g.Cells[f2][i]
          }

          cell.Tile.Move(dest)
          dest.Tile = cell.Tile
        } else {
          cell.Tile.Move(dest)
          dest.Tile = cell.Tile
        }
        cell.Tile = nil
      }
    }
  }
  g.NewTile()
  return g
}

func (g *Grid) EmptyCells() []*Cell {
  var ret []*Cell
  for i:=0; i < g.Size; i++ {
    for j:=0; j < g.Size; j++ {
      if(g.Cells[i][j].Tile == nil) {
        ret = append(ret, g.Cells[i][j])
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
