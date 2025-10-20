package complexmodel

type complexSize struct {
	id          int64
	complexID   int64
	size        float64
	description string
}

func (cs *complexSize) ID() int64 {
	return cs.id
}

func (cs *complexSize) SetID(id int64) {
	cs.id = id
}

func (cs *complexSize) ComplexID() int64 {
	return cs.complexID
}

func (cs *complexSize) SetComplexID(complexID int64) {
	cs.complexID = complexID
}

func (cs *complexSize) Size() float64 {
	return cs.size
}

func (cs *complexSize) SetSize(size float64) {
	cs.size = size
}

func (cs *complexSize) Description() string {
	return cs.description
}

func (cs *complexSize) SetDescription(description string) {
	cs.description = description
}
