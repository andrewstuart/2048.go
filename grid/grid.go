package grid

import (
  "math/rand"
  "time"
)

var idSeed uint= 0;

//Tiles track value and position on the grid
type tile struct {
  ID uint
  Value int
  MergeHistory []*tile
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
  Tiles []*tile
  Cells [][]*Cell
}

type Cell struct {
  Tile *tile
}

//Create the grid
func (g *Grid) Build () {
  g.Cells = make([][]*Cell, g.Size)
  g.Tiles = make([]*tile, g.Size * g.Size)

  for i := range g.Cells {
    g.Cells[i] = make([]*Cell, g.Size)
    for j := range g.Cells[i] {
      g.Cells[i][j] = new(Cell)
    }
  }
}

func (g *Grid) Reset() {
  for i := range g.Cells {
    for j := range g.Cells[i] {
      g.Cells[i][j].Tile = nil
    }
  }
  g.Tiles = make([]*tile, g.Size * g.Size)
}

func (g *Grid) NewTile() tile {

  //Two random values between 0 and Grid.Size
  x := randVal().Int() % g.Size
  y := randVal().Int() % g.Size

  id := idSeed
  idSeed += 1

  newTile := tile{
    ID: id,
    Value: randTileVal(),
    MergeHistory: make([]*tile, 11),
  }

  g.Cells[x][y].Tile = &newTile

  g.Tiles = append(g.Tiles, &newTile)

  return newTile
}

type vector struct {
  X int
  Y int
}

var dMap = map[int]vector{
  1: vector{ //Up
    X: 0,
    Y: 1,
  },
  2: vector{ //Right
    X: 1,
    Y: 0,
  },
  3: vector{ //Down
    X: 0,
    Y: -1,
  },
  4: vector{ //Left
    X: -1,
    Y: 0,
  },
}

func (g *Grid) getEdge (v *vector) (max, min, delta int) {
  if(v.X == 1 || v.Y == 1) {
    max = g.Size
    delta = 1
  } else {
    max = 0
    delta = -1
  }

  return max, g.Size - max - delta, delta
}

//Get direction
//For table row (vector perpendicular to direction)
  //For row cell (vector opposite to direction)
    //If cell has tile
      //Add to 'new row' slice
      //If tile value matches prev tile value
        //Merge tile to prev 
        //Remove tile ([2, 2, 4] .. [2 <-- 2, 4] .. [4, 4])
func (g *Grid) Move(d int) *Grid {
  v := dMap[d]

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

        dest.Tile = cell.Tile
        cell.Tile = nil
        loc += delta
      }
    }
  }

  return g
}
