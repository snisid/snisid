package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/chef-svc/internal/domain"
)

type MemberRepo struct {
	db *sql.DB
}

func NewMemberRepo(db *sql.DB) *MemberRepo {
	return &MemberRepo{db: db}
}

func (r *MemberRepo) Create(m *domain.CriminalMember) error {
	m.MemberID = uuid.New().String()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	aliases, _ := json.Marshal(m.Aliases)
	languages, _ := json.Marshal(m.KnownLanguages)
	photos, _ := json.Marshal(m.PhotoRefs)
	communes, _ := json.Marshal(m.TerritoryCommunes)
	weapons, _ := json.Marshal(m.WeaponTypes)

	_, err := r.db.Exec(`
		INSERT INTO criminal_members (
			member_id, national_chef_id, snisid_person_id, fir_record_id,
			afis_subject_id, rdep_deportee_id, primary_gang_id, role_in_gang,
			role_description, joined_date, rank_level, aliases, known_languages,
			tattoo_description, physical_description, photo_refs, territory_dept,
			territory_communes, known_armed, weapon_types, trained_combatant,
			status, un_designated, ofac_designated, ofac_sdn_ref,
			interpol_notice_ref, last_known_address, last_seen_at,
			intel_confidence, created_by, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,
			$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32
		)`,
		m.MemberID, m.NationalChefID, m.SNISIDPersonID, m.FIRRecordID,
		m.AFISSubjectID, m.RDEPDeporteeID, m.PrimaryGangID, m.RoleInGang,
		m.RoleDescription, m.JoinedDate, m.RankLevel, aliases, languages,
		m.TattooDescription, m.PhysicalDescription, photos, m.TerritoryDept,
		communes, m.KnownArmed, weapons, m.TrainedCombatant,
		m.Status, m.UNDesignated, m.OFACDesignated, m.OFACSDNRef,
		m.InterpolNoticeRef, m.LastKnownAddress, m.LastSeenAt,
		m.IntelConfidence, m.CreatedBy, m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *MemberRepo) GetByID(id string) (*domain.CriminalMember, error) {
	m := &domain.CriminalMember{}
	var aliases, languages, photos, communes, weapons []byte

	err := r.db.QueryRow(`
		SELECT member_id, national_chef_id, snisid_person_id, fir_record_id,
			afis_subject_id, rdep_deportee_id, primary_gang_id, role_in_gang,
			role_description, joined_date, rank_level, aliases, known_languages,
			tattoo_description, physical_description, photo_refs, territory_dept,
			territory_communes, known_armed, weapon_types, trained_combatant,
			status, un_designated, ofac_designated, ofac_sdn_ref,
			interpol_notice_ref, last_known_address, last_seen_at,
			intel_confidence, created_by, created_at, updated_at
		FROM criminal_members WHERE member_id = $1`, id,
	).Scan(
		&m.MemberID, &m.NationalChefID, &m.SNISIDPersonID, &m.FIRRecordID,
		&m.AFISSubjectID, &m.RDEPDeporteeID, &m.PrimaryGangID, &m.RoleInGang,
		&m.RoleDescription, &m.JoinedDate, &m.RankLevel, &aliases, &languages,
		&m.TattooDescription, &m.PhysicalDescription, &photos, &m.TerritoryDept,
		&communes, &m.KnownArmed, &weapons, &m.TrainedCombatant,
		&m.Status, &m.UNDesignated, &m.OFACDesignated, &m.OFACSDNRef,
		&m.InterpolNoticeRef, &m.LastKnownAddress, &m.LastSeenAt,
		&m.IntelConfidence, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(aliases, &m.Aliases)
	json.Unmarshal(languages, &m.KnownLanguages)
	json.Unmarshal(photos, &m.PhotoRefs)
	json.Unmarshal(communes, &m.TerritoryCommunes)
	json.Unmarshal(weapons, &m.WeaponTypes)

	return m, nil
}

func (r *MemberRepo) GetByGang(gangID string) ([]domain.CriminalMember, error) {
	rows, err := r.db.Query(`
		SELECT member_id, role_in_gang, status, rank_level, territory_dept, last_seen_at
		FROM criminal_members WHERE primary_gang_id = $1 ORDER BY rank_level ASC`, gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.CriminalMember
	for rows.Next() {
		var m domain.CriminalMember
		if err := rows.Scan(&m.MemberID, &m.RoleInGang, &m.Status, &m.RankLevel, &m.TerritoryDept, &m.LastSeenAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *MemberRepo) GetSanctioned() ([]domain.CriminalMember, error) {
	rows, err := r.db.Query(`
		SELECT member_id, primary_gang_id, role_in_gang, status,
			ofac_designated, un_designated, ofac_sdn_ref, interpol_notice_ref
		FROM criminal_members
		WHERE ofac_designated = true OR un_designated = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.CriminalMember
	for rows.Next() {
		var m domain.CriminalMember
		if err := rows.Scan(&m.MemberID, &m.PrimaryGangID, &m.RoleInGang, &m.Status,
			&m.OFACDesignated, &m.UNDesignated, &m.OFACSDNRef, &m.InterpolNoticeRef); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *MemberRepo) GetLeaders() ([]domain.CriminalMember, error) {
	rows, err := r.db.Query(`
		SELECT member_id, primary_gang_id, role_in_gang, status,
			rank_level, territory_dept, last_seen_at
		FROM criminal_members
		WHERE role_in_gang IN ('SUPREME_LEADER', 'ZONE_COMMANDER', 'LIEUTENANT')
		ORDER BY rank_level ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.CriminalMember
	for rows.Next() {
		var m domain.CriminalMember
		if err := rows.Scan(&m.MemberID, &m.PrimaryGangID, &m.RoleInGang, &m.Status,
			&m.RankLevel, &m.TerritoryDept, &m.LastSeenAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *MemberRepo) UpdateStatus(id string, status domain.MemberStatus) error {
	_, err := r.db.Exec(`
		UPDATE criminal_members SET status = $1, updated_at = $2
		WHERE member_id = $3`, status, time.Now(), id)
	return err
}

func (r *MemberRepo) Update(m *domain.CriminalMember) error {
	m.UpdatedAt = time.Now()
	aliases, _ := json.Marshal(m.Aliases)
	languages, _ := json.Marshal(m.KnownLanguages)
	photos, _ := json.Marshal(m.PhotoRefs)
	communes, _ := json.Marshal(m.TerritoryCommunes)
	weapons, _ := json.Marshal(m.WeaponTypes)

	_, err := r.db.Exec(`
		UPDATE criminal_members SET
			national_chef_id=$2, snisid_person_id=$3, fir_record_id=$4,
			afis_subject_id=$5, rdep_deportee_id=$6, primary_gang_id=$7,
			role_in_gang=$8, role_description=$9, joined_date=$10,
			rank_level=$11, aliases=$12, known_languages=$13,
			tattoo_description=$14, physical_description=$15, photo_refs=$16,
			territory_dept=$17, territory_communes=$18, known_armed=$19,
			weapon_types=$20, trained_combatant=$21, status=$22,
			un_designated=$23, ofac_designated=$24, ofac_sdn_ref=$25,
			interpol_notice_ref=$26, last_known_address=$27, last_seen_at=$28,
			intel_confidence=$29, updated_at=$30
		WHERE member_id = $1`,
		m.MemberID, m.NationalChefID, m.SNISIDPersonID, m.FIRRecordID,
		m.AFISSubjectID, m.RDEPDeporteeID, m.PrimaryGangID, m.RoleInGang,
		m.RoleDescription, m.JoinedDate, m.RankLevel, aliases, languages,
		m.TattooDescription, m.PhysicalDescription, photos, m.TerritoryDept,
		communes, m.KnownArmed, weapons, m.TrainedCombatant,
		m.Status, m.UNDesignated, m.OFACDesignated, m.OFACSDNRef,
		m.InterpolNoticeRef, m.LastKnownAddress, m.LastSeenAt,
		m.IntelConfidence, m.UpdatedAt,
	)
	return err
}

type IntelNoteRepo struct {
	db *sql.DB
}

func NewIntelNoteRepo(db *sql.DB) *IntelNoteRepo {
	return &IntelNoteRepo{db: db}
}

func (r *IntelNoteRepo) Create(n *domain.IntelligenceNote) error {
	n.NoteID = uuid.New().String()
	n.CreatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO intelligence_notes (note_id, member_id, source_id, note_type, content, confidence, collected_at, created_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		n.NoteID, n.MemberID, n.SourceID, n.NoteType, n.Content, n.Confidence, n.CollectedAt, n.CreatedBy, n.CreatedAt)
	return err
}

func (r *IntelNoteRepo) GetByMemberID(memberID string) ([]domain.IntelligenceNote, error) {
	rows, err := r.db.Query(`
		SELECT note_id, member_id, source_id, note_type, content, confidence, collected_at, created_by, created_at
		FROM intelligence_notes WHERE member_id = $1 ORDER BY created_at DESC`, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []domain.IntelligenceNote
	for rows.Next() {
		var n domain.IntelligenceNote
		if err := rows.Scan(&n.NoteID, &n.MemberID, &n.SourceID, &n.NoteType, &n.Content, &n.Confidence, &n.CollectedAt, &n.CreatedBy, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

type SightingRepo struct {
	db *sql.DB
}

func NewSightingRepo(db *sql.DB) *SightingRepo {
	return &SightingRepo{db: db}
}

func (r *SightingRepo) Create(s *domain.Sighting) error {
	s.SightingID = uuid.New().String()
	s.CreatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO sightings (sighting_id, member_id, source_id, dept, commune, latitude, longitude, spotted_at, confidence, notes, created_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		s.SightingID, s.MemberID, s.SourceID, s.Dept, s.Commune, s.Latitude, s.Longitude, s.SpottedAt, s.Confidence, s.Notes, s.CreatedBy, s.CreatedAt)
	return err
}

func (r *SightingRepo) GetByMemberID(memberID string) ([]domain.Sighting, error) {
	rows, err := r.db.Query(`
		SELECT sighting_id, member_id, source_id, dept, commune, latitude, longitude, spotted_at, confidence, notes, created_by, created_at
		FROM sightings WHERE member_id = $1 ORDER BY spotted_at DESC`, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sightings []domain.Sighting
	for rows.Next() {
		var s domain.Sighting
		if err := rows.Scan(&s.SightingID, &s.MemberID, &s.SourceID, &s.Dept, &s.Commune, &s.Latitude, &s.Longitude, &s.SpottedAt, &s.Confidence, &s.Notes, &s.CreatedBy, &s.CreatedAt); err != nil {
			return nil, err
		}
		sightings = append(sightings, s)
	}
	return sightings, nil
}
