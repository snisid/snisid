import pytest

from services.revocation import (
    RevocationEvent,
    RevocationEventType,
    RevocationNotifier,
    WalletRevocationHook,
)


class TestRevocationEvent:
    def test_create_event(self):
        event = RevocationEvent(
            credential_id="vc-1",
            event_type=RevocationEventType.CREDENTIAL_REVOKED,
            subject_id="did:key:alice",
            reason="Lost ID card",
        )
        assert event.id is not None
        assert event.credential_id == "vc-1"
        assert event.subject_id == "did:key:alice"

    def test_to_dict(self):
        event = RevocationEvent(
            credential_id="vc-2",
            event_type=RevocationEventType.CREDENTIAL_SUSPENDED,
            subject_id="did:key:bob",
        )
        d = event.to_dict()
        assert d["credential_id"] == "vc-2"
        assert d["event_type"] == "credential.suspended"


class TestRevocationNotifier:
    @pytest.fixture
    def notifier(self):
        return RevocationNotifier()

    def test_notify_revocation(self, notifier):
        event = notifier.notify_revocation("vc-1", "did:key:alice", "Fraud")
        assert event.event_type == RevocationEventType.CREDENTIAL_REVOKED
        assert len(notifier._history) == 1

    def test_notify_suspension(self, notifier):
        event = notifier.notify_suspension("vc-2", "did:key:bob", "Pending review")
        assert event.event_type == RevocationEventType.CREDENTIAL_SUSPENDED

    def test_notify_reinstatement(self, notifier):
        event = notifier.notify_reinstatement("vc-3", "did:key:charlie")
        assert event.event_type == RevocationEventType.CREDENTIAL_REINSTATED

    def test_subscribe_callback(self, notifier):
        received = []

        def callback(event):
            received.append(event)

        notifier.subscribe(RevocationEventType.CREDENTIAL_REVOKED, callback)
        notifier.notify_revocation("vc-4", "did:key:dave")
        assert len(received) == 1
        assert received[0].credential_id == "vc-4"

    def test_get_history_by_credential(self, notifier):
        notifier.notify_revocation("vc-a", "did:key:alice")
        notifier.notify_revocation("vc-b", "did:key:bob")
        notifier.notify_revocation("vc-a", "did:key:alice", "Duplicate")

        history = notifier.get_history(credential_id="vc-a")
        assert len(history) == 2

    def test_get_history_by_subject(self, notifier):
        notifier.notify_revocation("vc-1", "did:key:alice")
        notifier.notify_revocation("vc-2", "did:key:bob")
        history = notifier.get_history(subject_id="did:key:alice")
        assert len(history) == 1

    def test_get_history_limit(self, notifier):
        for i in range(10):
            notifier.notify_revocation(f"vc-{i}", f"did:key:user{i}")
        history = notifier.get_history(limit=3)
        assert len(history) == 3


class TestWalletRevocationHook:
    @pytest.fixture
    def notifier(self):
        return RevocationNotifier()

    @pytest.fixture
    def hook(self, notifier):
        return WalletRevocationHook(notifier, "did:key:alice")

    def test_track_credential(self, notifier, hook):
        hook.track_credential("vc-1")
        assert "vc-1" in notifier.get_wallet_credentials("did:key:alice")

    def test_untrack_credential(self, notifier, hook):
        hook.track_credential("vc-1")
        hook.untrack_credential("vc-1")
        assert "vc-1" not in notifier.get_wallet_credentials("did:key:alice")

    def test_receive_notification_for_tracked(self, notifier, hook):
        hook.track_credential("vc-tracked")
        hook.track_credential("vc-other")

        notifier.notify_revocation("vc-tracked", "did:key:alice", "Revoked")
        notifier.notify_revocation("vc-other", "did:key:bob", "Other")

        assert len(hook.get_notifications()) == 2

    def test_ignores_untracked_credentials(self, notifier, hook):
        notifier.notify_revocation("vc-unknown", "did:key:alice")
        assert len(hook.get_notifications()) == 0

    def test_check_status(self, notifier, hook):
        hook.track_credential("vc-status")
        notifier.notify_suspension("vc-status", "did:key:alice", "Under review")
        status = hook.check_status("vc-status")
        assert status == RevocationEventType.CREDENTIAL_SUSPENDED

    def test_check_status_none(self, notifier, hook):
        assert hook.check_status("vc-never") is None

    def test_clear_notifications(self, notifier, hook):
        hook.track_credential("vc-clear")
        notifier.notify_revocation("vc-clear", "did:key:alice")
        assert len(hook.get_notifications()) == 1
        hook.clear_notifications()
        assert len(hook.get_notifications()) == 0

    def test_receive_reinstatement(self, notifier, hook):
        hook.track_credential("vc-reinstate")
        notifier.notify_reinstatement("vc-reinstate", "did:key:alice", "Error corrected")
        assert len(hook.get_notifications()) == 1
        assert hook.get_notifications()[0].event_type == RevocationEventType.CREDENTIAL_REINSTATED

    def test_multiple_wallets(self, notifier):
        alice_hook = WalletRevocationHook(notifier, "did:key:alice")
        bob_hook = WalletRevocationHook(notifier, "did:key:bob")

        alice_hook.track_credential("vc-shared")
        bob_hook.track_credential("vc-shared")

        notifier.notify_revocation("vc-shared", "did:key:alice")

        assert len(alice_hook.get_notifications()) == 1
        assert len(bob_hook.get_notifications()) == 1
