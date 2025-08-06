package complexmodel

type complexTower struct {
	id            int64
	complexID     int64
	tower         string
	floors        int
	totalUnits    int
	unitsPerFloor int
}

func (ct *complexTower) ID() int64 {
	return ct.id
}

func (ct *complexTower) SetID(id int64) {
	ct.id = id
}

func (ct *complexTower) ComplexID() int64 {
	return ct.complexID
}

func (ct *complexTower) SetComplexID(complexID int64) {
	ct.complexID = complexID
}

func (ct *complexTower) Tower() string {
	return ct.tower
}

func (ct *complexTower) SetTower(tower string) {
	ct.tower = tower
}

func (ct *complexTower) Floors() int {
	return ct.floors
}

func (ct *complexTower) SetFloors(floors int) {
	ct.floors = floors
}

func (ct *complexTower) TotalUnits() int {
	return ct.totalUnits
}

func (ct *complexTower) SetTotalUnits(totalUnits int) {
	ct.totalUnits = totalUnits
}

func (ct *complexTower) UnitsPerFloor() int {
	return ct.unitsPerFloor
}

func (ct *complexTower) SetUnitsPerFloor(unitsPerFloor int) {
	ct.unitsPerFloor = unitsPerFloor
}
