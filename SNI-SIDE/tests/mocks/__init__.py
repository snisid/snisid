"""
Mock classes for testing SNI-SIDE components
"""


class MockNeo4jSession:
    def __init__(self):
        self.nodes = {}
        self.relationships = []

    async def run(self, query: str, params: dict = None):
        return MockNeo4jResult()

    async def close(self):
        pass

    async def __aenter__(self):
        return self

    async def __aexit__(self, *args):
        pass


class MockNeo4jResult:
    async def data(self):
        return []

    async def single(self):
        return None


class MockNeo4jDriver:
    def __init__(self):
        self.sessions = []

    def session(self):
        return MockNeo4jSession()

    async def close(self):
        pass


class MockRedis:
    def __init__(self):
        self.data = {}

    async def get(self, key):
        return self.data.get(key)

    async def setex(self, key, ttl, value):
        self.data[key] = value

    async def delete(self, key):
        self.data.pop(key, None)

    async def close(self):
        pass


class MockKafkaConsumer:
    def __init__(self):
        self.topics = []
        self.messages = []

    async def start(self):
        pass

    async def stop(self):
        pass

    async def getmany(self, timeout_ms=1000, max_records=500):
        return {}

    async def commit(self):
        pass


class MockKafkaProducer:
    def __init__(self):
        self.sent_messages = []

    async def start(self):
        pass

    async def stop(self):
        pass

    async def send_and_wait(self, topic, key=None, value=None):
        self.sent_messages.append({
            "topic": topic,
            "key": key,
            "value": value,
            "timestamp": None,
        })
        return MockRecordMetadata()

    async def send(self, topic, key=None, value=None):
        self.sent_messages.append({
            "topic": topic,
            "key": key,
            "value": value,
        })
        return MockFuture()


class MockRecordMetadata:
    def __init__(self):
        self.offset = 0


class MockFuture:
    def __init__(self):
        pass

    async def __aenter__(self):
        return self

    async def __aexit__(self, *args):
        pass


class MockClickHouseClient:
    def __init__(self):
        self.executed = []

    def execute(self, query, params=None):
        self.executed.append({"query": query, "params": params})

    def disconnect(self):
        pass
