import logging, os, json
from typing import Optional
from neo4j import GraphDatabase, AsyncGraphDatabase

logger = logging.getLogger("sniside.graph-updater")

NEO4J_URI = os.getenv("NEO4J_URI", "bolt://neo4j:7687")
NEO4J_USER = os.getenv("NEO4J_USER", "neo4j")
NEO4J_PASSWORD = os.getenv("NEO4J_PASSWORD", "sniside-neo4j")


class Neo4jGraphUpdater:
    def __init__(self):
        self.driver = None

    def start(self):
        self.driver = AsyncGraphDatabase.driver(NEO4J_URI, auth=(NEO4J_USER, NEO4J_PASSWORD))
        logger.info("Neo4j driver initialized")

    def stop(self):
        if self.driver:
            self.driver.close()

    async def _run(self, query: str, params: dict = None):
        async with self.driver.session() as session:
            result = await session.run(query, params or {})
            return await result.data()

    async def upsert_citizen(self, niu: str, data: dict):
        await self._run("""
            MERGE (c:Citizen {niu: $niu})
            SET c.full_name = $full_name,
                c.status = $status,
                c.risk_level = $risk_level,
                c.updated_at = timestamp()
        """, {
            "niu": niu,
            "full_name": data.get("full_name", "UNKNOWN"),
            "status": data.get("status", "ACTIVE"),
            "risk_level": data.get("risk_level", "LOW"),
        })
        logger.debug(f"Citizen {niu} upserted")

    async def upsert_alias(self, niu: str, alias: str):
        await self._run("""
            MERGE (a:Alias {name: $alias, citizen_niu: $niu})
        """, {"alias": alias, "niu": niu})

    async def upsert_vehicle(self, plate: str, data: dict):
        await self._run("""
            MERGE (v:Vehicle {plate: $plate})
            SET v.make = $make, v.model = $model, v.year = $year,
                v.color = $color, v.status = $status
        """, {
            "plate": plate,
            "make": data.get("make", ""),
            "model": data.get("model", ""),
            "year": data.get("year", 0),
            "color": data.get("color", ""),
            "status": data.get("status", "ACTIVE"),
        })

    async def create_case(self, case_id: str, data: dict):
        await self._run("""
            MERGE (c:Case {case_id: $case_id})
            SET c.type = $type, c.status = $status,
                c.title = $title, c.opened_at = $opened_at
        """, {
            "case_id": case_id,
            "type": data.get("case_type", ""),
            "status": data.get("status", "OPEN"),
            "title": data.get("title", ""),
            "opened_at": data.get("timestamp", 0),
        })

    async def create_gang(self, gang_name: str, data: dict):
        await self._run("""
            MERGE (g:Gang {name: $name})
            SET g.territory = $territory, g.activities = $activities,
                g.risk_level = $risk_level
        """, {
            "name": gang_name,
            "territory": data.get("territory", ""),
            "activities": json.dumps(data.get("criminal_activities", [])),
            "risk_level": data.get("risk_level", "MEDIUM"),
        })

    async def create_border_crossing(self, niu: str, crossing_id: str, data: dict):
        await self._run("""
            MERGE (bc:BorderCrossing {crossing_id: $crossing_id})
            SET bc.port = $port, bc.direction = $direction,
                bc.origin = $origin, bc.destination = $destination,
                bc.timestamp = $timestamp
        """, {
            "crossing_id": crossing_id,
            "port": data.get("port_of_entry", ""),
            "direction": data.get("direction", "EXIT"),
            "origin": data.get("origin_country", ""),
            "destination": data.get("destination_country", ""),
            "timestamp": data.get("timestamp", 0),
        })
        if niu:
            await self._run("""
                MATCH (c:Citizen {niu: $niu})
                MATCH (bc:BorderCrossing {crossing_id: $crossing_id})
                MERGE (c)-[r:TRAVELLED_TO]->(bc)
                SET r.timestamp = $timestamp
            """, {"niu": niu, "crossing_id": crossing_id, "timestamp": data.get("timestamp", 0)})

    async def create_bank_account(self, acct_id: str, account_number: str, data: dict):
        await self._run("""
            MERGE (a:BankAccount {account_id: $acct_id})
            SET a.account_number = $acct_num, a.bank = $bank,
                a.type = $type, a.country = $country
        """, {
            "acct_id": acct_id,
            "acct_num": account_number or acct_id,
            "bank": data.get("bank", ""),
            "type": data.get("account_type", "CHECKING"),
            "country": data.get("country", ""),
        })

    async def create_ip(self, ip: str, data: dict):
        await self._run("""
            MERGE (n:IP {ip: $ip})
            SET n.country = $country, n.asn = $asn,
                n.isp = $isp, n.threat_score = $threat_score
        """, {
            "ip": ip,
            "country": data.get("country", ""),
            "asn": data.get("asn", ""),
            "isp": data.get("isp", ""),
            "threat_score": data.get("threat_score", 0),
        })

    async def create_domain(self, domain: str, data: dict):
        await self._run("""
            MERGE (d:Domain {domain: $domain})
            SET d.registrant = $registrant, d.threat_score = $threat_score,
                d.category = $category
        """, {
            "domain": domain,
            "registrant": data.get("registrant", ""),
            "threat_score": data.get("threat_score", 0),
            "category": data.get("category", ""),
        })

    async def create_wallet(self, wallet: str, data: dict):
        await self._run("""
            MERGE (w:Wallet {address: $wallet})
            SET w.transaction_count = $tx_count,
                w.total_value_btc = $total_btc,
                w.exchange = $exchange
        """, {
            "wallet": wallet,
            "tx_count": data.get("transaction_count", 0),
            "total_btc": data.get("total_value_btc", 0),
            "exchange": data.get("exchange", ""),
        })

    async def create_evidence(self, evidence_id: str, data: dict):
        await self._run("""
            MERGE (e:DigitalEvidence {evidence_id: $evidence_id})
            SET e.type = $type, e.status = $status,
                e.source = $source, e.collected_at = $collected_at
        """, {
            "evidence_id": evidence_id,
            "type": data.get("evidence_type", ""),
            "status": data.get("status", "COLLECTED"),
            "source": data.get("source", ""),
            "collected_at": data.get("timestamp", 0),
        })

    async def create_relationship(self, from_id: str, rel_type: str, to_id: str, props: dict):
        rel_type_clean = rel_type.upper().replace("-", "_").replace(" ", "_")
        query = f"""
            MATCH (a)
            WHERE a.niu = $from_id OR a.plate = $from_id
               OR a.case_id = $from_id OR a.name = $from_id
               OR a.crossing_id = $from_id OR a.account_id = $from_id
               OR a.ip = $from_id OR a.domain = $from_id
               OR a.address = $from_id OR a.evidence_id = $from_id
               OR a.firearm_id = $from_id OR a.route_id = $from_id
               OR a.seizure_id = $from_id
            MATCH (b)
            WHERE b.niu = $to_id OR b.plate = $to_id
               OR b.case_id = $to_id OR b.name = $to_id
               OR b.crossing_id = $to_id OR b.account_id = $to_id
               OR b.ip = $to_id OR b.domain = $to_id
               OR b.address = $to_id OR b.evidence_id = $to_id
               OR b.firearm_id = $to_id OR b.route_id = $to_id
               OR b.seizure_id = $to_id
            MERGE (a)-[r:{rel_type_clean}]->(b)
            SET r.updated_at = timestamp()
        """
        set_clauses = []
        for k, v in props.items():
            if isinstance(v, str):
                set_clauses.append(f"r.{k} = '{v.replace(chr(39), chr(39)+chr(39))}'")
            else:
                set_clauses.append(f"r.{k} = ${k}")
        if set_clauses:
            query += ", " + ", ".join(set_clauses) + ";"
        else:
            query += ";"

        string_params = {k: v for k, v in props.items() if not isinstance(v, str)}
        string_params["from_id"] = from_id
        string_params["to_id"] = to_id

        try:
            await self._run(query, string_params)
            logger.debug(f"Relationship {rel_type_clean}: {from_id} -> {to_id}")
        except Exception as e:
            logger.warning(f"Failed to create relationship {rel_type_clean} {from_id}->{to_id}: {e}")

    async def remove_relationship(self, from_id: str, rel_type: str, to_id: str):
        rel_type_clean = rel_type.upper().replace("-", "_")
        query = f"""
            MATCH (a)-[r:{rel_type_clean}]->(b)
            WHERE (a.niu = $from_id OR a.plate = $from_id OR a.account_id = $from_id)
              AND (b.niu = $to_id OR b.plate = $to_id OR b.account_id = $to_id)
            DELETE r
        """
        await self._run(query, {"from_id": from_id, "to_id": to_id})
