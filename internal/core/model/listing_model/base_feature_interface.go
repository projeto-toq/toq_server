package listingmodel

type BaseFeatureInterface interface {
	ID() int64
	SetID(id int64)
	Feature() string
	SetFeature(feature string)
	Description() string
	SetDescription(description string)
	Priority() uint8
	SetPriority(priority uint8)
}

func NewBaseFeature() BaseFeatureInterface {
	return &baseFeature{}
}
