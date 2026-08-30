import pytest

from services.didcomm_mediator import DIDCommMediator, ForwardedMessage


class TestDIDCommMediator:
    @pytest.fixture
    def mediator(self):
        return DIDCommMediator()

    def test_forward_packed_message(self, mediator):
        msg = mediator.forward("did:key:bob", {"ciphertext": "encrypted"}, sender_did="did:key:alice")
        assert msg.recipient_did == "did:key:bob"
        assert msg.sender_did == "did:key:alice"
        assert msg.packed_message == {"ciphertext": "encrypted"}
        assert not msg.is_delivered

    def test_forward_request_valid(self, mediator):
        forward = {
            "type": "https://didcomm.org/routing/2.0/forward",
            "id": "msg-1",
            "to": ["did:key:bob"],
            "from": "did:key:alice",
            "body": {"next": "did:key:bob"},
            "attachments": [{"data": {"json": {"ciphertext": "encrypted"}}}],
        }
        msg = mediator.forward_request(forward)
        assert msg is not None
        assert msg.recipient_did == "did:key:bob"
        assert msg.sender_did == "did:key:alice"

    def test_forward_request_invalid_type(self, mediator):
        msg = mediator.forward_request({"type": "https://didcomm.org/basicmessage/2.0/message"})
        assert msg is None

    def test_forward_request_no_recipient(self, mediator):
        msg = mediator.forward_request({"type": "https://didcomm.org/routing/2.0/forward"})
        assert msg is None

    def test_fetch_messages(self, mediator):
        mediator.forward("did:key:bob", {"data": "msg1"}, "did:key:alice")
        mediator.forward("did:key:bob", {"data": "msg2"}, "did:key:charlie")
        mediator.forward("did:key:someone", {"data": "msg3"}, "did:key:dave")

        pending = mediator.fetch_messages("did:key:bob")
        assert len(pending) == 2

    def test_deliver_message(self, mediator):
        mediator.forward("did:key:bob", {"data": "test"}, "did:key:alice")
        msg_id = list(mediator._messages.keys())[0]

        delivered = mediator.deliver(msg_id)
        assert delivered is not None
        assert delivered.is_delivered
        assert delivered.delivered_at is not None

        # Verify it's no longer pending
        pending = mediator.fetch_messages("did:key:bob")
        assert len(pending) == 0

    def test_deliver_nonexistent(self, mediator):
        assert mediator.deliver("nonexistent") is None

    def test_get_inbox(self, mediator):
        mediator.forward("did:key:bob", {"data": "first"}, "did:key:alice")
        mediator.forward("did:key:bob", {"data": "second"}, "did:key:charlie")

        inbox = mediator.get_inbox("did:key:bob")
        assert len(inbox) == 2
        # Newest first
        assert inbox[0].packed_message["data"] == "second"

    def test_get_pending_count(self, mediator):
        mediator.forward("did:key:bob", {"data": "1"}, "did:key:alice")
        mediator.forward("did:key:bob", {"data": "2"}, "did:key:charlie")
        assert mediator.get_pending_count("did:key:bob") == 2

        # Deliver one
        msg_id = list(mediator._messages.keys())[0]
        mediator.deliver(msg_id)
        assert mediator.get_pending_count("did:key:bob") == 1

    def test_delete_message(self, mediator):
        mediator.forward("did:key:bob", {"data": "test"}, "did:key:alice")
        msg_id = list(mediator._messages.keys())[0]

        assert mediator.delete_message(msg_id) is True
        assert mediator.delete_message("nonexistent") is False

    def test_clear(self, mediator):
        mediator.forward("did:key:bob", {"data": "1"}, "did:key:alice")
        mediator.forward("did:key:charlie", {"data": "2"}, "did:key:bob")
        assert len(mediator._messages) == 2
        mediator.clear()
        assert len(mediator._messages) == 0
