package complexmodel

type ComplexTowerInterface interface {
	ID() int64
	SetID(int64)
	ComplexID() int64
	SetComplexID(int64)
	Tower() string
	SetTower(string)
	Floors() *int
	SetFloors(*int)
	TotalUnits() *int
	SetTotalUnits(*int)
	UnitsPerFloor() *int
	SetUnitsPerFloor(*int)
}

func NewComplexTower() ComplexTowerInterface {
	return &complexTower{}
}
