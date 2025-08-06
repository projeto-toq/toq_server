package listingmodel

type baseFeature struct {
	id          int64
	feature     string
	description string
	priority    uint8
}

func (b *baseFeature) ID() int64 {
	return b.id
}

func (b *baseFeature) SetID(id int64) {
	b.id = id
}

func (b *baseFeature) Feature() string {
	return b.feature
}

func (b *baseFeature) SetFeature(feature string) {
	b.feature = feature
}

func (b *baseFeature) Description() string {
	return b.description
}

func (b *baseFeature) SetDescription(description string) {
	b.description = description
}

func (b *baseFeature) Priority() uint8 {
	return b.priority
}

func (b *baseFeature) SetPriority(priority uint8) {
	b.priority = priority
}
