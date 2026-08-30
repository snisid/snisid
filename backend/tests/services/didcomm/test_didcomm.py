import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from services.didcomm import DIDCommMessage, DIDCommMessenger, DIDCommPacker
from services.didcomm.api import router as didcomm_router


class TestDIDCommMessage:
    def test_to_dict_full(self):
        msg = DIDCommMessage(
            id="msg-1",
            type="https://example.com/protocol/1.0/message",
            body={"hello": "world"},
            from_did="did:key:alice",
            to_did="did:key:bob",
            thid="thread-1",
        )
        d = msg.to_dict()
        assert d["id"] == "msg-1"
        assert d["from"] == "did:key:alice"
        assert d["to"] == ["did:key:bob"]
        assert d["thid"] == "thread-1"

    def test_to_dict_minimal(self):
        msg = DIDCommMessage(id="msg-1", type="test", body={})
        d = msg.to_dict()
        assert "from" not in d
        assert "to" not in d

    def test_from_dict_roundtrip(self):
        original = DIDCommMessage(
            id="msg-1",
            type="test",
            body={"key": "value"},
            from_did="did:key:alice",
            to_did="did:key:bob",
            thid="t-1",
            pthid="pt-1",
        )
        d = original.to_dict()
        restored = DIDCommMessage.from_dict(d)
        assert restored.id == "msg-1"
        assert restored.from_did == "did:key:alice"
        assert restored.to_did == "did:key:bob"

    def test_attachments(self):
        msg = DIDCommMessage(
            id="msg-1",
            type="test",
            body={},
            attachments=[{"id": "attach-1", "data": {"base64": "dGVzdA=="}}],
        )
        assert len(msg.attachments) == 1


class TestDIDCommPacker:
    @pytest.fixture
    def packer(self):
        return DIDCommPacker()

    @pytest.fixture
    def sample_message(self):
        return DIDCommMessage(
            id="msg-1",
            type="https://didcomm.org/trust-ping/2.0/ping",
            body={"response_requested": True},
            from_did="did:key:alice",
            to_did="did:key:bob",
        )

    def test_pack_signed(self, packer, sample_message):
        packed = packer.pack(sample_message, sender_did="did:key:alice")
        assert "payload" in packed
        assert "signatures" in packed
        assert len(packed["signatures"]) == 1

    def test_pack_encrypted(self, packer, sample_message):
        packed = packer.pack(
            sample_message,
            sender_did="did:key:alice",
            recipient_did="did:key:bob",
        )
        assert "ciphertext" in packed
        assert "recipients" in packed

    def test_unpack_signed(self, packer, sample_message):
        packed = packer.pack(sample_message, sender_did="did:key:alice")
        unpacked = packer.unpack(packed)
        assert unpacked.id == "msg-1"
        assert unpacked.body["response_requested"] is True

    def test_unpack_encrypted(self, packer, sample_message):
        packed = packer.pack(
            sample_message,
            sender_did="did:key:alice",
            recipient_did="did:key:bob",
        )
        unpacked = packer.unpack(packed)
        assert unpacked.id == "msg-1"

    def test_roundtrip_with_attachments(self, packer):
        msg = DIDCommMessage(
            id="msg-2",
            type="test",
            body={"data": "test"},
            from_did="did:key:alice",
            to_did="did:key:bob",
            attachments=[{"id": "a1", "data": {"json": {"foo": "bar"}}}],
        )
        packed = packer.pack(msg, sender_did="did:key:alice", recipient_did="did:key:bob")
        unpacked = packer.unpack(packed)
        assert len(unpacked.attachments) == 1


class TestDIDCommMessenger:
    @pytest.fixture
    def messenger(self):
        return DIDCommMessenger()

    def test_send_receive(self, messenger):
        msg = DIDCommMessage(
            id="m1",
            type="test",
            body={"hello": "world"},
            from_did="did:key:alice",
            to_did="did:key:bob",
        )
        packed = messenger.send(msg, "did:key:alice", "did:key:bob")
        received = messenger.receive(packed)
        assert received.id == "m1"
        assert received.body["hello"] == "world"
        assert len(messenger._sent) == 1
        assert len(messenger._inbox) == 1

    def test_trust_ping_flow(self, messenger):
        ping = messenger.create_trust_ping(
            from_did="did:key:alice", to_did="did:key:bob"
        )
        assert ping.type == "https://didcomm.org/trust-ping/2.0/ping"
        assert ping.body["response_requested"] is True

        packed_ping = messenger.send(
            ping, "did:key:alice", "did:key:bob"
        )
        received_ping = messenger.receive(packed_ping)

        response = messenger.create_trust_ping_response(received_ping)
        assert response.type == "https://didcomm.org/trust-ping/2.0/ping_response"
        assert response.thid == ping.id

    def test_create_trust_ping_no_response(self, messenger):
        ping = messenger.create_trust_ping(
            from_did="did:key:alice",
            to_did="did:key:bob",
            response_requested=False,
        )
        assert ping.body["response_requested"] is False


class TestDIDCommApi:
    @pytest.fixture
    def client(self):
        app = FastAPI()
        app.include_router(didcomm_router)
        return TestClient(app)

    def test_pack(self, client):
        resp = client.post(
            "/didcomm/pack",
            json={
                "type": "https://didcomm.org/trust-ping/2.0/ping",
                "body": {"response_requested": True},
                "from_did": "did:key:alice",
                "to_did": "did:key:bob",
            },
        )
        assert resp.status_code == 200
        data = resp.json()
        assert "signatures" in data or "ciphertext" in data

    def test_unpack(self, client):
        pack_resp = client.post(
            "/didcomm/pack",
            json={
                "type": "https://didcomm.org/trust-ping/2.0/ping",
                "body": {"hello": "world"},
                "from_did": "did:key:alice",
                "to_did": "did:key:bob",
            },
        )
        packed = pack_resp.json()
        resp = client.post("/didcomm/unpack", json={"packed": packed})
        assert resp.status_code == 200
        data = resp.json()
        assert data["body"]["hello"] == "world"

    def test_trust_ping_endpoint(self, client):
        resp = client.post(
            "/didcomm/trust-ping",
            params={"from_did": "did:key:alice", "to_did": "did:key:bob"},
        )
        assert resp.status_code == 200
