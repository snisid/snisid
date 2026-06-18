package domain

type ExplType string

const (
	ExplTypeIED              ExplType = "IED"
	ExplTypeGrenade          ExplType = "GRENADE"
	ExplTypeRPG              ExplType = "RPG"
	ExplTypeMortar           ExplType = "MORTAR"
	ExplTypeLandmine         ExplType = "LANDMINE"
	ExplTypeDynamite         ExplType = "DYNAMITE"
	ExplTypeBlastingCap      ExplType = "BLASTING_CAP"
	ExplTypeAmmunitionBulk   ExplType = "AMMUNITION_BULK"
	ExplTypeMilitaryOrdnance ExplType = "MILITARY_ORDNANCE"
	ExplTypeUnknown          ExplType = "UNKNOWN"
)

type ExplStatus string

const (
	ExplStatusRecovered      ExplStatus = "RECOVERED"
	ExplStatusDestroyed      ExplStatus = "DESTROYED"
	ExplStatusDetonated      ExplStatus = "DETONATED"
	ExplStatusStoredEvidence ExplStatus = "STORED_EVIDENCE"
	ExplStatusTransferred    ExplStatus = "TRANSFERRED"
)
