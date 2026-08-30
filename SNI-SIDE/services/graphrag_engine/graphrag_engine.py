import asyncio, json, logging, os, uuid, time, re
from datetime import datetime
from typing import Optional
from enum import Enum

import httpx
from neo4j import AsyncGraphDatabase
from prometheus_client import Counter, Histogram, Gauge, start_http_server

logger = logging.getLogger("sniside.graphrag")

NEO4J_URI = os.getenv("NEO4J_URI", "bolt://neo4j:7687")
NEO4J_USER = os.getenv("NEO4J_USER", "neo4j")
NEO4J_PASSWORD = os.getenv("NEO4J_PASSWORD", "sniside-neo4j")
MILVUS_URI = os.getenv("MILVUS_URI", "http://milvus:19530")
LLM_ENDPOINT = os.getenv("LLM_ENDPOINT", "http://ollama:11434")
LLM_MODEL = os.getenv("LLM_MODEL", "mistral:7b-instruct-v0.3-q4_K_M")
LLM_TIMEOUT = int(os.getenv("LLM_TIMEOUT", "30"))
MAX_CONTEXT_NODES = int(os.getenv("MAX_CONTEXT_NODES", "500"))
METRICS_PORT = int(os.getenv("METRICS_PORT", "9103"))
EMBEDDING_DIM = int(os.getenv("EMBEDDING_DIM", "768"))

reports_generated = Counter("sniside_graphrag_reports_total", "Intelligence reports generated", ["type"])
report_latency = Histogram("sniside_graphrag_latency_seconds", "Report generation latency", ["type"],
                           buckets=[.1, .25, .5, 1, 2.5, 5, 10, 30, 60])
context_nodes = Histogram("sniside_graphrag_context_nodes", "Context graph size", buckets=[1, 5, 10, 25, 50, 100, 250, 500])
llm_errors = Counter("sniside_graphrag_llm_errors_total", "LLM call errors", ["stage"])


class ReportType(str, Enum):
    ENTITY_PROFILE = "ENTITY_PROFILE"
    LINK_ANALYSIS = "LINK_ANALYSIS"
    NETWORK_MAP = "NETWORK_MAP"
    FINANCIAL_FLOW = "FINANCIAL_FLOW"
    TEMPORAL_ANALYSIS = "TEMPORAL_ANALYSIS"
    CROSS_ENTITY = "CROSS_ENTITY"


CYPHER_TEMPLATES = {
    "ENTITY_PROFILE": """
        MATCH (c {niu: $id})
        OPTIONAL MATCH (c)-[r]->(connected)
        RETURN c as source, type(r) as rel_type, properties(r) as rel_props,
               labels(connected) as target_labels, connected as target
        LIMIT $limit
    """,
    "LINK_ANALYSIS": """
        MATCH path = (c {niu: $id})-[*1..2]-(connected)
        UNWIND nodes(path) AS n
        UNWIND relationships(path) AS r
        RETURN DISTINCT n, type(r) as rel_type, properties(r) as rel_props
        LIMIT $limit
    """,
    "NETWORK_MAP": """
        MATCH path = (g {name: $id})-[*1..3]-(connected)
        UNWIND nodes(path) AS n
        UNWIND relationships(path) AS r
        RETURN DISTINCT n, type(r) as rel_type, properties(r) as rel_props
        LIMIT $limit
    """,
    "FINANCIAL_FLOW": """
        MATCH path = (a:BankAccount {account_id: $id})-[:TRANSFERRED_TO*1..5]->(target)
        UNWIND nodes(path) AS n
        UNWIND relationships(path) AS r
        RETURN DISTINCT n, type(r) as rel_type, properties(r) as rel_props
        LIMIT $limit
    """,
    "TEMPORAL_ANALYSIS": """
        MATCH (c {niu: $id})-[r]-(connected)
        WHERE r.timestamp >= $since
        RETURN c as source, type(r) as rel_type, properties(r) as rel_props,
               labels(connected) as target_labels, connected as target,
               r.timestamp as event_time
        ORDER BY event_time DESC
        LIMIT $limit
    """,
    "CROSS_ENTITY": """
        MATCH (a {niu: $id1}), (b {niu: $id2})
        OPTIONAL MATCH path = shortestPath((a)-[*..6]-(b))
        UNWIND nodes(path) AS n
        UNWIND relationships(path) AS r
        RETURN DISTINCT n, type(r) as rel_type, properties(r) as rel_props
        LIMIT $limit
    """,
}

PROMPT_TEMPLATES = {
    "ENTITY_PROFILE": """Tu es un analyste de renseignement criminel senior pour le SNI-SIDE, système national de renseignement.

Analyse le profil criminel complet de l'entité suivante à partir des données du graphe d'intelligence national.

ENTITÉ CIBLE:
- Type: {entity_type}
- Identifiant: {entity_id}
- Nom: {entity_label}

CONTEXTE DU GRAPHE (nœuds et relations):
{context}

RAPPORTS ANTÉRIEURS:
{history}

Réponds UNIQUEMENT en JSON valide avec cette structure:
{{
  "executive_summary": "résumé opérationnel (3-5 phrases)",
  "key_findings": ["liste des constats clés"],
  "risk_assessment": {{
    "overall_risk": "CRITICAL|HIGH|MEDIUM|LOW",
    "risk_factors": ["facteurs de risque"],
    "confidence": 0.95
  }},
  "connections_analysis": "analyse des connexions suspectes",
  "recommendations": ["recommandations opérationnelles"],
  "indicators": {{
    "nius": ["listes des NIU liés"],
    "vehicles": ["plaques liées"],
    "phones": ["numéros liés"],
    "addresses": ["adresses liées"],
    "accounts": ["comptes liés"]
  }}
}}""",

    "LINK_ANALYSIS": """Tu es un analyste de renseignement criminel senior.

Analyse les connexions et le réseau autour de l'entité suivante. Identifie les chemins suspects, les relations inhabituelles, et les patterns criminels.

ENTITÉ CIBLE:
- Type: {entity_type}
- Identifiant: {entity_id}
- Nom: {entity_label}

CONTEXTE DU GRAPHE (connexions jusqu'à 2 hops):
{context}

Réponds UNIQUEMENT en JSON valide avec:
{{
  "executive_summary": "résumé",
  "key_findings": ["constats"],
  "network_metrics": {{
    "node_count": N,
    "relationship_count": N,
    "density": 0.0,
    "central_nodes": ["nœuds centraux"]
  }},
  "suspicious_patterns": [
    {{"pattern": "description", "severity": "HIGH|MEDIUM|LOW", "entities": ["entités impliquées"]}}
  ],
  "recommendations": ["recommandations"]
}}""",

    "FINANCIAL_FLOW": """Tu es un expert AML (Anti-Money Laundering).

Analyse les flux financiers suspects suivants et identifie les patterns de blanchiment, les réseaux de sociétés-écrans, et les transferts anormaux.

CONTEXTE FINANCIER:
{context}

Réponds UNIQUEMENT en JSON valide avec:
{{
  "executive_summary": "résumé AML",
  "key_findings": ["constats"],
  "flow_analysis": "analyse détaillée des flux",
  "risk_score": 0.0,
  "aml_indicators": ["indicateurs AML détectés"],
  "recommendations": ["recommandations"]
}}""",

    "TEMPORAL_ANALYSIS": """Tu es un analyste spécialisé en chronologie criminelle.

Analyse l'évolution temporelle des activités de l'entité suivante et identifie les tendances, les changements de comportement, et les événements significatifs.

ENTITÉ CIBLE:
- Type: {entity_type}
- Identifiant: {entity_id}
- Nom: {entity_label}
- Période: depuis {since}

CONTEXTE CHRONOLOGIQUE:
{context}

Réponds UNIQUEMENT en JSON valide avec:
{{
  "executive_summary": "résumé chronologique",
  "key_findings": ["constats"],
  "timeline": [
    {{"date": "...", "event": "...", "significance": "HIGH|MEDIUM|LOW"}}
  ],
  "behavioral_changes": ["changements de comportement détectés"],
  "predictive_assessment": "évaluation prédictive des risques futurs",
  "recommendations": ["recommandations"]
}}""",

    "CROSS_ENTITY": """Tu es un analyste de renseignement.

Analyse les connexions entre les deux entités suivantes et détermine la nature et la force de leur relation dans le contexte criminel.

ENTITÉ A:
- Type: {entity_type1}
- Id: {id1}
- Nom: {label1}

ENTITÉ B:
- Type: {entity_type2}
- Id: {id2}
- Nom: {label2}

CONTEXTE DU GRAPHE:
{context}

Réponds UNIQUEMENT en JSON valide avec:
{{
  "executive_summary": "résumé de la relation",
  "connection_type": "DIRECT|INDIRECT|NONE",
  "path_length": N,
  "path_description": "description du chemin",
  "relationship_strength": 0.0,
  "shared_entities": ["entités communes"],
  "risk_if_connected": "évaluation du risque si la connexion est confirmée",
  "recommendations": ["recommandations"]
}}""",
}


class Neo4jRetriever:
    def __init__(self):
        self.driver = None

    def start(self):
        self.driver = AsyncGraphDatabase.driver(NEO4J_URI, auth=(NEO4J_USER, NEO4J_PASSWORD))
        logger.info("GraphRAG Neo4j driver connected")

    def stop(self):
        if self.driver:
            self.driver.close()

    async def retrieve(self, cypher: str, params: dict) -> dict:
        async with self.driver.session() as session:
            result = await session.run(cypher, params)
            records = await result.data()
            return self._build_context(records)

    def _build_context(self, records: list) -> dict:
        nodes = {}
        relationships = []
        for rec in records:
            for key in ("n", "source", "target", "connected", "c"):
                if key in rec and rec[key]:
                    node = rec[key]
                    nid = str(node.element_id) if hasattr(node, 'element_id') else str(node.get("id") or node.get("niu", ""))
                    nodes[nid] = {
                        "id": nid,
                        "labels": list(rec.get(f"{key}_labels", node.get("labels", node.get("_labels", [])))) or ["Unknown"],
                        "properties": dict(node.get("properties", node)) if hasattr(node, 'get') else dict(node.items()),
                    }
            if "rel_type" in rec and rec["rel_type"]:
                relationships.append({
                    "type": rec["rel_type"],
                    "source_id": str(rec.get("source", {}).get("element_id", "")),
                    "target_id": str(rec.get("target", {}).get("element_id", "")),
                    "properties": rec.get("rel_props", {}),
                })
        node_list = list(nodes.values())
        context_nodes.observe(len(node_list))
        return {"nodes": node_list[:MAX_CONTEXT_NODES], "relationships": relationships}


class LLMGenerator:
    def __init__(self):
        self.client = httpx.AsyncClient(timeout=LLM_TIMEOUT)

    async def generate(self, prompt: str, system_prompt: str = None) -> str:
        payload = {
            "model": LLM_MODEL,
            "prompt": prompt,
            "stream": False,
            "options": {"temperature": 0.1, "top_p": 0.9, "max_tokens": 4096},
        }
        if system_prompt:
            payload["system"] = system_prompt

        try:
            resp = await self.client.post(f"{LLM_ENDPOINT}/api/generate", json=payload)
            resp.raise_for_status()
            data = resp.json()
            return data.get("response", "")
        except httpx.TimeoutException:
            llm_errors.labels(stage="timeout").inc()
            logger.error("LLM request timed out")
            raise
        except Exception as e:
            llm_errors.labels(stage="request").inc()
            logger.error(f"LLM request failed: {e}")
            raise

    async def generate_chat(self, messages: list, model: str = None) -> str:
        payload = {
            "model": model or LLM_MODEL,
            "messages": messages,
            "stream": False,
            "options": {"temperature": 0.1, "max_tokens": 4096},
        }
        try:
            resp = await self.client.post(f"{LLM_ENDPOINT}/v1/chat/completions", json=payload)
            resp.raise_for_status()
            data = resp.json()
            return data.get("choices", [{}])[0].get("message", {}).get("content", "")
        except Exception as e:
            llm_errors.labels(stage="chat").inc()
            logger.error(f"LLM chat failed: {e}")
            raise


class GraphRAGEngine:
    def __init__(self):
        self.retriever = Neo4jRetriever()
        self.llm = LLMGenerator()

    def start(self):
        self.retriever.start()
        logger.info("GraphRAG Engine started")

    def stop(self):
        self.retriever.stop()

    async def generate_report(self, entity_id: str, entity_type: str = "Citizen",
                               entity_label: str = None, report_type: ReportType = ReportType.ENTITY_PROFILE,
                               depth: int = 2, since: int = None, entity_id2: str = None,
                               history: list = None) -> dict:
        start_time = time.monotonic()
        report_id = str(uuid.uuid4())
        cypher_key = report_type.value

        cypher = CYPHER_TEMPLATES.get(cypher_key)
        if not cypher:
            raise ValueError(f"Unknown report type: {report_type}")

        params = {"id": entity_id, "limit": MAX_CONTEXT_NODES}
        if report_type == ReportType.TEMPORAL_ANALYSIS:
            params["since"] = since or int((time.time() - 7776000) * 1000)
        elif report_type == ReportType.CROSS_ENTITY:
            if not entity_id2:
                raise ValueError("entity_id2 required for CROSS_ENTITY reports")
            params = {"id1": entity_id, "id2": entity_id2, "limit": MAX_CONTEXT_NODES}

        context = await self.retriever.retrieve(cypher, params)

        prompt_key = report_type.value
        if entity_label is None:
            entity_label = entity_id
        prompt = PROMPT_TEMPLATES[prompt_key].format(
            entity_type=entity_type, entity_id=entity_id, entity_label=entity_label,
            context=json.dumps(context, indent=2, default=str),
            history=json.dumps(history or [], indent=2, default=str),
            id1=entity_id, id2=entity_id2 or "", label1=entity_label, label2=entity_label,
            type1=entity_type, type2=entity_type, since=since or "toujours",
        )

        system = "Tu es un analyste de renseignement criminel senior. Réponds UNIQUEMENT en JSON valide."
        llm_response = await self.llm.generate(prompt, system)

        parsed = self._parse_response(llm_response, report_type, entity_id, entity_type, entity_label, context, report_id)

        elapsed = time.monotonic() - start_time
        report_latency.labels(type=report_type.value).observe(elapsed)
        reports_generated.labels(type=report_type.value).inc()
        logger.info(f"GraphRAG report {report_id} generated in {elapsed:.2f}s ({report_type.value})")

        return parsed

    def _parse_response(self, llm_text: str, report_type: ReportType, entity_id: str,
                        entity_type: str, entity_label: str, context: dict,
                        report_id: str = None) -> dict:
        json_match = re.search(r'\{[\s\S]*\}', llm_text)
        if json_match:
            try:
                analysis = json.loads(json_match.group())
            except json.JSONDecodeError:
                analysis = {"raw": llm_text, "parse_error": "JSON decode failed"}
        else:
            analysis = {"raw": llm_text, "parse_error": "No JSON found"}

        risk_level = analysis.get("risk_assessment", {}).get("overall_risk") or analysis.get("risk_level", "MEDIUM")

        return {
            "report_id": report_id or str(uuid.uuid4()),
            "report_type": report_type.value,
            "generated_at": datetime.utcnow().isoformat(),
            "target_entity": {
                "type": entity_type,
                "id": entity_id,
                "label": entity_label,
            },
            "executive_summary": analysis.get("executive_summary", ""),
            "key_findings": analysis.get("key_findings", []),
            "risk_assessment": {
                "overall_risk": risk_level,
                "risk_factors": analysis.get("risk_assessment", {}).get("risk_factors", []),
                "confidence": analysis.get("risk_assessment", {}).get("confidence", analysis.get("confidence", 0.7)),
            },
            "graph_context": {
                "node_count": len(context.get("nodes", [])),
                "relationship_count": len(context.get("relationships", [])),
            },
            "connections_analysis": analysis.get("connections_analysis", ""),
            "recommendations": analysis.get("recommendations", []),
            "network_metrics": analysis.get("network_metrics", {}),
            "suspicious_patterns": analysis.get("suspicious_patterns", []),
            "timeline": analysis.get("timeline", []),
            "indicators": analysis.get("indicators", {}),
            "flow_analysis": analysis.get("flow_analysis", ""),
            "aml_indicators": analysis.get("aml_indicators", []),
            "confidence_score": analysis.get("risk_assessment", {}).get("confidence", analysis.get("confidence", 0.7)),
            "model_used": LLM_MODEL,
            "model_version": "2.1",
        }


async def main():
    start_http_server(METRICS_PORT)
    engine = GraphRAGEngine()
    engine.start()
    logger.info("GraphRAG Engine ready on :{METRICS_PORT}")

    try:
        while True:
            await asyncio.sleep(3600)
    except KeyboardInterrupt:
        pass
    finally:
        engine.stop()


if __name__ == "__main__":
    asyncio.run(main())
