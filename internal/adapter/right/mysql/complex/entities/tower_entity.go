package complexentities

type TowerEntity struct {
	ID            int64
	ComplexID     int64
	Tower         string
	Floors        int
	TotalUnits    int
	UnitsPerFloor int
}
