package listingmodel

import globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"

// PropertyTypeOption aggregates the available property type metadata
// for listings. Each option binds the numeric code returned to clients
// with the internal bit flag used by complex services.
type PropertyTypeOption struct {
	Code        int64
	Label       string
	PropertyBit globalmodel.PropertyType
}

// PropertyTypeCatalog holds every supported property type option for listings.
var PropertyTypeCatalog = []PropertyTypeOption{
	{Code: 1, Label: "Apartamento", PropertyBit: globalmodel.Apartment},
	{Code: 2, Label: "Loja", PropertyBit: globalmodel.CommercialStore},
	{Code: 4, Label: "Laje", PropertyBit: globalmodel.CommercialFloor},
	{Code: 8, Label: "Sala", PropertyBit: globalmodel.Suite},
	{Code: 16, Label: "Casa", PropertyBit: globalmodel.House},
	{Code: 32, Label: "Casa na Planta", PropertyBit: globalmodel.OffPlanHouse},
	{Code: 64, Label: "Terreno Residencial", PropertyBit: globalmodel.ResidencialLand},
	{Code: 128, Label: "Terreno Comercial", PropertyBit: globalmodel.CommercialLand},
	{Code: 256, Label: "Prédio", PropertyBit: globalmodel.Building},
	{Code: 512, Label: "Galpão", PropertyBit: globalmodel.Warehouse},
}

var propertyTypeByBit map[globalmodel.PropertyType]PropertyTypeOption

var propertyTypeByCode map[int64]PropertyTypeOption

func init() {
	propertyTypeByBit = make(map[globalmodel.PropertyType]PropertyTypeOption, len(PropertyTypeCatalog))
	propertyTypeByCode = make(map[int64]PropertyTypeOption, len(PropertyTypeCatalog))

	for _, option := range PropertyTypeCatalog {
		propertyTypeByBit[option.PropertyBit] = option
		propertyTypeByCode[option.Code] = option
	}
}

// PropertyTypeOptionFromBit returns the catalog option for the provided bit mask.
func PropertyTypeOptionFromBit(bit globalmodel.PropertyType) (PropertyTypeOption, bool) {
	option, ok := propertyTypeByBit[bit]
	return option, ok
}

// PropertyTypeOptionFromCode returns the catalog option for the provided decimal code.
func PropertyTypeOptionFromCode(code int64) (PropertyTypeOption, bool) {
	option, ok := propertyTypeByCode[code]
	return option, ok
}
