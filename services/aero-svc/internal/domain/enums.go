package domain

type AircraftType string

const (
	COMMERCIAL_JET AircraftType = "COMMERCIAL_JET"
	TURBOPROP      AircraftType = "TURBOPROP"
	PISTON_SINGLE  AircraftType = "PISTON_SINGLE"
	PISTON_TWIN    AircraftType = "PISTON_TWIN"
	HELICOPTER     AircraftType = "HELICOPTER"
	ULTRALIGHT     AircraftType = "ULTRALIGHT"
	DRONE_LARGE    AircraftType = "DRONE_LARGE"
	UNKNOWN        AircraftType = "UNKNOWN"
)

type StripStatus string

const (
	ACTIVE             StripStatus = "ACTIVE"
	INACTIVE           StripStatus = "INACTIVE"
	DESTROYED          StripStatus = "DESTROYED"
	LEGALIZED          StripStatus = "LEGALIZED"
	UNDER_SURVEILLANCE StripStatus = "UNDER_SURVEILLANCE"
)
