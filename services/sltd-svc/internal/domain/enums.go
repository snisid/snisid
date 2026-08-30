package domain

type DocType string

const (
	DocTypePassport        DocType = "PASSPORT"
	DocTypeNationalID      DocType = "NATIONAL_ID"
	DocTypeTravelDocument  DocType = "TRAVEL_DOCUMENT"
	DocTypeVisa            DocType = "VISA"
	DocTypeResidencePermit DocType = "RESIDENCE_PERMIT"
	DocTypeRefugeeDocument DocType = "REFUGEE_DOCUMENT"
	DocTypeLaissezPasser   DocType = "LAISSEZ_PASSER"
)

type DocStatus string

const (
	DocStatusLost      DocStatus = "LOST"
	DocStatusStolen    DocStatus = "STOLEN"
	DocStatusRevoked   DocStatus = "REVOKED"
	DocStatusExpired   DocStatus = "EXPIRED"
	DocStatusFound     DocStatus = "FOUND"
	DocStatusRecovered DocStatus = "RECOVERED"
	DocStatusCancelled DocStatus = "CANCELLED"
)
