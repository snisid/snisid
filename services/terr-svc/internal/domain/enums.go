package domain

type ControlLevel string

const (
	FULL_CONTROL     ControlLevel = "FULL_CONTROL"
	STRONG_INFLUENCE ControlLevel = "STRONG_INFLUENCE"
	CONTESTED        ControlLevel = "CONTESTED"
	WEAK_INFLUENCE   ControlLevel = "WEAK_INFLUENCE"
	STATE_CONTROLLED ControlLevel = "STATE_CONTROLLED"
	NO_MAN_LAND      ControlLevel = "NO_MAN_LAND"
)

type Source string

const (
	PNH_FIELD_REPORT   Source = "PNH_FIELD_REPORT"
	SATELLITE_ANALYSIS Source = "SATELLITE_ANALYSIS"
	INFORMANT          Source = "INFORMANT"
	NGO_REPORT         Source = "NGO_REPORT"
	ACLED              Source = "ACLED"
	MEDIA_CROSS_CHECK  Source = "MEDIA_CROSS_CHECK"
	LAPI_ANALYSIS      Source = "LAPI_ANALYSIS"
)
