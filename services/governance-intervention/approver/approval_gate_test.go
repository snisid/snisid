package approver

import (
	"testing"

	"github.com/snisid/platform/services/governance-intervention/proposer"
	"github.com/stretchr/testify/assert"
)

func makeIntervention(action string) *proposer.Intervention {
	return &proposer.Intervention{
		ID:     "INT-TEST",
		Action: action,
		Status: proposer.StatusProposed,
	}
}

func TestAuthorize_NoJustification(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("MONITOR")

	err := g.Authorize(i, "admin", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JUSTIFICATION_REQUIRED")
	assert.Equal(t, proposer.StatusProposed, i.Status)
}

func TestAuthorize_FreezeAccount_InvestigatorAllowed(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("FREEZE_ACCOUNT")

	err := g.Authorize(i, "investigator", "high risk transaction")
	assert.NoError(t, err)
	assert.Equal(t, proposer.StatusApproved, i.Status)
}

func TestAuthorize_FreezeAccount_AdminAllowed(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("FREEZE_ACCOUNT")

	err := g.Authorize(i, "admin", "emergency freeze")
	assert.NoError(t, err)
	assert.Equal(t, proposer.StatusApproved, i.Status)
}

func TestAuthorize_FreezeAccount_AnalystDenied(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("FREEZE_ACCOUNT")

	err := g.Authorize(i, "analyst", "looks suspicious")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "INSUFFICIENT_PRIVILEGES_FOR_ACTION")
	assert.Equal(t, proposer.StatusProposed, i.Status)
}

func TestAuthorize_Monitor_AnyoneAllowed(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("MONITOR")

	err := g.Authorize(i, "intern", "routine check")
	assert.NoError(t, err)
	assert.Equal(t, proposer.StatusApproved, i.Status)
}

func TestAuthorize_Investigate_AnyRole(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("INVESTIGATE")

	err := g.Authorize(i, "analyst", "flagged by system")
	assert.NoError(t, err)
	assert.Equal(t, proposer.StatusApproved, i.Status)
}

func TestReject_UpdatesStatus(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("FREEZE_ACCOUNT")

	g.Reject(i, "insufficient evidence")
	assert.Equal(t, proposer.StatusRejected, i.Status)
}

func TestReject_NilIntervention(t *testing.T) {
	g := &ApprovalGate{}
	g.Reject(nil, "reason")
}

func TestAuthorize_EmptyJustificationEdgeCase(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("INVESTIGATE")

	err := g.Authorize(i, "admin", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JUSTIFICATION_REQUIRED")
}

func TestAuthorize_JustificationWithSpaces(t *testing.T) {
	g := &ApprovalGate{}
	i := makeIntervention("MONITOR")

	err := g.Authorize(i, "analyst", "   ")
	assert.NoError(t, err)
	assert.Equal(t, proposer.StatusApproved, i.Status)
}
