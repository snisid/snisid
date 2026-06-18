package domain

type Source string

const (
	UN_2653         Source = "UN_2653"
	OFAC_SDN        Source = "OFAC_SDN"
	EU_CONSOLIDATED Source = "EU_CONSOLIDATED"
	INTERPOL        Source = "INTERPOL"
	CANADA_OSFI     Source = "CANADA_OSFI"
	UK_OFSI         Source = "UK_OFSI"
	OTHER           Source = "OTHER"
)

type Measure string

const (
	ASSETS_FREEZE Measure = "ASSETS_FREEZE"
	TRAVEL_BAN    Measure = "TRAVEL_BAN"
	ARMS_EMBARGO  Measure = "ARMS_EMBARGO"
	ALL_MEASURES  Measure = "ALL_MEASURES"
)

type EntityType string

const (
	INDIVIDUAL   EntityType = "INDIVIDUAL"
	ORGANIZATION EntityType = "ORGANIZATION"
	VESSEL       EntityType = "VESSEL"
	AIRCRAFT     EntityType = "AIRCRAFT"
)
