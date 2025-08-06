package complexmodel

type complexZipCode struct {
	id        int64
	complexID int64
	zipCode   string
}

func (c *complexZipCode) ID() int64 {
	return c.id
}

func (c *complexZipCode) SetID(id int64) {
	c.id = id
}

func (c *complexZipCode) ComplexID() int64 {
	return c.complexID
}

func (c *complexZipCode) SetComplexID(complexID int64) {
	c.complexID = complexID
}

func (c *complexZipCode) ZipCode() string {
	return c.zipCode
}

func (c *complexZipCode) SetZipCode(zipCode string) {
	c.zipCode = zipCode
}
