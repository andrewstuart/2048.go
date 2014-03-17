package grid

import (
  "fmt"
  "math/rand"
  "time"
)

var idSeed uint= 0;

//Tiles track value and position on the grid
type tile struct {
  ID uint
  Value int
  MergeHistory []*tile
  X int
  Y int
}

func (t *tile) Move (dest *Cell) {
  fmt.Printf("Moving tile at X %d,Y %d to Cell at X %d, Y %d", t.X, t.Y, dest.X, dest.Y)
  fmt.Println()
  t.X = dest.X
  t.Y = dest.Y
}

func (t *tile) Merge (tn *tile) {
  fmt.Printf("Tile %d %d merging with tile %d %d", t.X, t.Y, tn.X, tn.Y)
  fmt.Println()
  t.Value += tn.Value //Future proofs for other merge rules. Fibonacci game??? Yes please. TODO
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

func randPos(max int) int {
  return randVal().Int() % max }

//Grid tracks the tiles
type Grid struct {
  Size int
  StartCells int
  Tiles []*tile
  Cells [][]*Cell
}

type Cell struct {
  Tile *tile
  X int
  Y int
}

//Create the grid
func (g *Grid) Build () {
  g.Cells = make([][]*Cell, g.Size)
  g.Tiles = make([]*tile, 0)

  for i := range g.Cells {
    g.Cells[i] = make([]*Cell, g.Size)
    for j := range g.Cells[i] {
      cell := new(Cell)
      cell.X = i
      cell.Y = j
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

  id := idSeed
  idSeed += 1

  newTile := tile{
    ID: id,
    Value: randTileVal(),
    MergeHistory: make([]*tile, 0),
    X: cell.X,
    Y: cell.Y,
  }

  g.Cells[cell.X][cell.Y].Tile = &newTile

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

func (g *Grid) getEdge (v *vector) (max, min, delta int) {
  if(v.X == 1 || v.Y == 1) {
    max = g.Size - 1
    min = - 1
    delta = -1
  } else {
    min = g.Size
    max = 0
    delta = 1
  }

  fmt.Printf("Max: %d, Min: %d, Delta: %d", max, min, delta)
  fmt.Println()

  return max, min, delta
}

//Get direction
//For table row (vector perpendicular to direction)
  //For row cell (vector opposite to direction)
    //If cell has tile
      //Add to 'new row' slice
      //If tile value matches prev tile value
        //Merge tile to prev 
        //Remove tile ([2, 2, 4] .. [2 <-- 2, 4] .. [4, 4])
func (g *Grid) Shift(d int) *Grid {
  v := dMap[d]
  fmt.Println(v)



  start, end, delta := g.getEdge(&v)

  for i:= start; i != end; i += delta {
    loc := start
    merged := false
    for j:= start; j != end; j += delta {
      var col int
      var cell, dest *Cell
      if(v.X == 0) { //Horizontal motion is zero. Direction is in Y
        col = i
        cell = g.Cells[col][j]
        dest = g.Cells[col][loc]
      } else {
        col = j
        cell = g.Cells[col][i]
        dest = g.Cells[col][loc]
      }

      if(cell.Tile != nil) { //If there is a tile for the cell we're checking
        if(loc != start) { //And if the pointer is not left at the beginning, which would mean we've moved no cells,
          potentialMerge := g.Cells[col][loc - delta].Tile; //Peek backwards to the last tile we moved

          if(!merged && potentialMerge.Value == cell.Tile.Value) { //If we didn't merge last time around and If the tiles match
              //Do a merge.
              potentialMerge.Merge(cell.Tile)
              cell.Tile = nil
              merged = true
              continue //Then continue
          } else {
            merged = false
          }
        }

        cell.Tile.Move(dest)
        dest.Tile = cell.Tile

        if(cell.X != dest.X && cell.Y != dest.Y) {
          cell.Tile = nil
        }

        loc += delta
      }
    }
    fmt.Println("LOC", loc)
  }

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
