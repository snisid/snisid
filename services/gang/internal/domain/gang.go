package domain

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	GangID                     uuid.UUID           `json:"gang_id"`
	NationalGangID             string              `json:"national_gang_id"`
	Name                       string              `json:"name"`
	Aliases                    []string            `json:"aliases"`
	StructureType              *GangStructureType  `json:"structure_type,omitempty"`
	PrimaryActivity            GangPrimaryActivity `json:"primary_activity"`
	ActivityLevel              GangActivityLevel   `json:"activity_level"`
	EstimatedMembers           *int                `json:"estimated_members,omitempty"`
	ArmedMembersPct            *int16              `json:"armed_members_pct,omitempty"`
	HeavyWeapons               bool                `json:"heavy_weapons"`
	PrimaryDeptCode            string              `json:"primary_dept_code"`
	TerritoryCommunes          []string            `json:"territory_communes"`
	TerritoryGeoJSON           *map[string]any     `json:"territory_geojson,omitempty"`
	EstimatedRevenueUSDMonthly *float64            `json:"estimated_revenue_usd_monthly,omitempty"`
	PrimaryIncomeSources       []string            `json:"primary_income_sources"`
	UNDesignationDate          *time.Time          `json:"un_designation_date,omitempty"`
	OFACDesignation            bool                `json:"ofac_designation"`
	OFACSDNRef                 *string             `json:"ofac_sdn_ref,omitempty"`
	AlliedGangIDs              []uuid.UUID         `json:"allied_gang_ids"`
	RivalGangIDs               []uuid.UUID         `json:"rival_gang_ids"`
	EstablishedDate            *time.Time          `json:"established_date,omitempty"`
	CurrentLeaderID            *uuid.UUID          `json:"current_leader_id,omitempty"`
	IntelConfidence            *int16              `json:"intel_confidence,omitempty"`
	LastIntelUpdate            *time.Time          `json:"last_intel_update,omitempty"`
	IsActive                   bool                `json:"is_active"`
	CreatedBy                  uuid.UUID           `json:"created_by"`
	CreatedAt                  time.Time           `json:"created_at"`
	UpdatedAt                  time.Time           `json:"updated_at"`
}

type CreateOrganizationRequest struct {
	Name              string              `json:"name" validate:"required"`
	Aliases           []string            `json:"aliases"`
	StructureType     *GangStructureType  `json:"structure_type"`
	PrimaryActivity   GangPrimaryActivity `json:"primary_activity" validate:"required"`
	ActivityLevel     GangActivityLevel   `json:"activity_level"`
	EstimatedMembers  *int                `json:"estimated_members"`
	ArmedMembersPct   *int16              `json:"armed_members_pct"`
	HeavyWeapons      bool                `json:"heavy_weapons"`
	PrimaryDeptCode   string              `json:"primary_dept_code" validate:"required,len=2"`
	TerritoryCommunes []string            `json:"territory_communes"`
	OFACDesignation   bool                `json:"ofac_designation"`
	OFACSDNRef        *string             `json:"ofac_sdn_ref"`
	EstablishedDate   *time.Time          `json:"established_date"`
	IntelConfidence   *int16              `json:"intel_confidence"`
}
