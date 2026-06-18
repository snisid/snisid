package domain

type RecoveryContext string

const (
	PoliceOperation    RecoveryContext = "POLICE_OPERATION"
	Checkpoint         RecoveryContext = "CHECKPOINT"
	PortSeizure        RecoveryContext = "PORT_SEIZURE"
	AirportSeizure     RecoveryContext = "AIRPORT_SEIZURE"
	CommunitySurrender RecoveryContext = "COMMUNITY_SURRENDER"
	CrimeScene         RecoveryContext = "CRIME_SCENE"
	Raid               RecoveryContext = "RAID"
	BorderSeizure      RecoveryContext = "BORDER_SEIZURE"
	Other              RecoveryContext = "OTHER"
)

type WeaponDisposition string

const (
	HeldAsEvidence      WeaponDisposition = "HELD_AS_EVIDENCE"
	Destroyed           WeaponDisposition = "DESTROYED"
	ReturnedToOwner     WeaponDisposition = "RETURNED_TO_OWNER"
	TransferredToPolice WeaponDisposition = "TRANSFERRED_TO_POLICE"
	SentToInterpol      WeaponDisposition = "SENT_TO_INTERPOL"
	Pending             WeaponDisposition = "PENDING"
)
