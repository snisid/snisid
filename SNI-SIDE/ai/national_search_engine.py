"""
SNI-SIDE: National Sovereign Search Engine
===========================================
Moteur de recherche fédéré unifié couvrant les 15 bases nationales.
Capacités: recherche textuelle, biométrique, ADN, véhicule, document,
criminelle, financière, cyber, géospatiale, multimodale.
"""

import asyncio
import time
import hashlib
from typing import List, Dict, Optional, Any, Set
from dataclasses import dataclass, field
from enum import Enum
import logging

logger = logging.getLogger("sniside.search")


# ============ TYPES DE RECHERCHE ============
class SearchType(Enum):
    ALL = "ALL"
    NAME = "NAME"
    NIU = "NIU"
    PHOTO = "PHOTO"
    FINGERPRINT = "FINGERPRINT"
    DNA = "DNA"
    PHONE = "PHONE"
    EMAIL = "EMAIL"
    ADDRESS = "ADDRESS"
    PLATE = "PLATE"
    VIN = "VIN"
    PASSPORT = "PASSPORT"
    CASE = "CASE"
    IP = "IP"
    DOMAIN = "DOMAIN"
    WALLET = "WALLET"
    DOCUMENT = "DOCUMENT"
    CRYPTO = "CRYPTO"
    LOCATION = "LOCATION"
    UNKNOWN = "UNKNOWN"


class Database(Enum):
    NCID = "NCID"
    BIOMETRICS = "BIOMETRICS"
    CODIS = "CODIS"
    MISSING = "MISSING"
    VEHICLE = "VEHICLE"
    ALPR = "ALPR"
    FIREARMS = "FIREARMS"
    BORDER = "BORDER"
    NARCOTICS = "NARCOTICS"
    FINANCIAL = "FINANCIAL"
    CYBER = "CYBER"
    WATCHLIST = "WATCHLIST"
    DOCUMENTS = "DOCUMENTS"
    GEOINT = "GEOINT"
    EVIDENCE = "EVIDENCE"


# ============ MODÈLES DE DONNÉES ============
@dataclass
class SearchQuery:
    """Requête de recherche unifiée"""
    raw_query: str
    query_type: SearchType = SearchType.ALL
    query_fingerprint: str = ""
    databases: List[Database] = field(default_factory=lambda: list(Database))
    page: int = 1
    limit: int = 20
    fuzzy: bool = True
    agency_context: str = ""
    user_clearance: str = "RESTRICTED"
    timeout_ms: int = 5000

    def __post_init__(self):
        if not self.query_fingerprint:
            self.query_fingerprint = hashlib.sha256(
                self.raw_query.encode()
            ).hexdigest()[:16]


@dataclass
class SearchResultItem:
    """Élément de résultat de recherche"""
    database: Database
    result_type: str
    id: str
    title: str
    description: str = ""
    score: float = 0.0
    match_field: str = ""
    match_confidence: float = 0.0
    risk_score: float = 0.0
    metadata: Dict[str, Any] = field(default_factory=dict)
    source_url: str = ""
    created_at: Optional[float] = None


@dataclass
class SearchResponse:
    """Réponse de recherche unifiée"""
    query: str
    results: Dict[Database, List[SearchResultItem]] = field(default_factory=dict)
    total_results: int = 0
    page: int = 1
    limit: int = 20
    search_duration_ms: float = 0.0
    databases_searched: int = 0
    graph_context: Optional[Dict] = None
    suggested_queries: List[str] = field(default_factory=list)


# ============ DÉTECTION AUTOMATIQUE DU TYPE DE RECHERCHE ============
class QueryTypeDetector:
    """Détecte automatiquement le type de recherche basé sur le format de la requête"""

    PATTERNS = {
        SearchType.NIU: [
            r'^HT[A-Z0-9]{8}$', r'^NIU[-\s]?[A-Z0-9]{10}$', r'^HT-\d{8}$',
            r'^\d{10}$',
        ],
        SearchType.PLATE: [
            r'^[A-Za-z]{2}[- ]?\d{3,4}[- ]?[A-Za-z]{2}$',
            r'^\d{3,4}[- ]?[A-Za-z]{2}$', r'^[A-Za-z]{1,3}[- ]?\d{1,4}$'
        ],
        SearchType.VIN: [
            r'^[A-HJ-NPR-Z0-9]{17}$', r'^[a-hj-npr-z0-9]{17}$'
        ],
        SearchType.PASSPORT: [
            r'^[A-Za-z]{2}\d{6,9}$', r'^HT\d{6,9}$', r'^P[A-Za-z]{1,2}\d{6,8}$'
        ],
        SearchType.PHONE: [
            r'^\+?509\d{8}$', r'^\+?\d{10,15}$', r'^0\d{8,9}$'
        ],
        SearchType.EMAIL: [
            r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        ],
        SearchType.IP: [
            r'^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$',
            r'^[0-9a-fA-F:]+$'
        ],
        SearchType.WALLET: [
            r'^1[a-km-zA-HJ-NP-Z1-9]{25,34}$',
            r'^0x[a-fA-F0-9]{40}$',
            r'^[13][a-km-zA-HJ-NP-Z1-9]{33}$',
            r'^bc1[a-zA-HJ-NP-Z0-9]{39,59}$',
        ],
        SearchType.DOMAIN: [
            r'^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$'
        ],
        SearchType.CASE: [
            r'^[Cc][Aa][Ss][Ee][- ]?\d{4}[- ]?\d{0,6}$', r'^NC[- ]?\d{4,10}$'
        ],
        SearchType.DOCUMENT: [
            r'^\d{10,15}$', r'^[A-Za-z]{2}\d{8,12}$'
        ]
    }

    @classmethod
    def detect(cls, query: str) -> SearchType:
        """Détecte le type de recherche à partir de la requête brute"""
        import re
        q = query.strip()
        if not q or (q.isdigit() and len(q) < 10):
            return SearchType.UNKNOWN
        for search_type, patterns in cls.PATTERNS.items():
            for pattern in patterns:
                if re.match(pattern, q):
                    return search_type
        return SearchType.NAME


def detect_query_type(query: str) -> str:
    """Module-level convenience wrapper — used by tests"""
    return QueryTypeDetector.detect(query).value


# ============ CONNECTEURS DE BASE DE DONNÉES ============
class DatabaseConnector:
    """Classe de base pour les connecteurs de base de données"""

    def __init__(self, database: Database):
        self.database = database
        self.search_timeout_ms = 3000

    async def search(self, query: SearchQuery) -> List[SearchResultItem]:
        """Recherche dans la base - à surcharger par chaque connecteur"""
        raise NotImplementedError

    def can_handle(self, query: SearchQuery) -> bool:
        """Détermine si ce connecteur peut gérer ce type de requête"""
        return True


class NCIDConnector(DatabaseConnector):
    """Connecteur pour la base criminelle NCID"""

    def __init__(self):
        super().__init__(Database.NCID)

    def can_handle(self, query: SearchQuery) -> bool:
        ncid_types = {SearchType.ALL, SearchType.NAME, SearchType.NIU,
                      SearchType.PHONE, SearchType.EMAIL, SearchType.CASE,
                      SearchType.ADDRESS}
        return query.query_type in ncid_types or query.query_type == SearchType.ALL

    async def search(self, query: SearchQuery) -> List[SearchResultItem]:
        results = []
        # NCID Query Simulation:
        # SELECT * FROM snisid_ncid.wanted_persons
        # WHERE to_tsvector('french', full_name) @@ plainto_tsquery('french', $1)
        #    OR alias ILIKE $2
        #    OR niu = $3
        # UNION
        # SELECT ... FROM criminal_cases ...
        # UNION
        # SELECT ... FROM criminal_aliases ...

        # Simulated results for demonstration
        mock = SearchResultItem(
            database=self.database,
            result_type="wanted_person",
            id="mock-ncid-001",
            title="Simulated NCID Result",
            score=0.95,
            match_field=query.raw_query,
            match_confidence=0.92,
            risk_score=0.85,
            metadata={"niu": "0000000001", "case_count": 3, "warrant_count": 2}
        )
        results.append(mock)
        return results


class BiometricConnector(DatabaseConnector):
    """Connecteur pour la base biométrique HN-NGI via Milvus"""

    def __init__(self):
        super().__init__(Database.BIOMETRICS)

    def can_handle(self, query: SearchQuery) -> bool:
        return query.query_type in {SearchType.PHOTO, SearchType.FINGERPRINT,
                                    SearchType.NIU, SearchType.ALL}

    async def search(self, query: SearchQuery) -> List[SearchResultItem]:
        # Milvus vector search:
        # collection_name = "sniside_face_embeddings"
        # search_params = {"metric_type": "IP", "params": {"nprobe": 10}}
        # results = milvus.search(
        #     collection_name=collection_name,
        #     data=[query_embedding],
        #     anns_field="embedding",
        #     param=search_params,
        #     limit=query.limit
        # )
        return []


class BorderConnector(DatabaseConnector):
    """Connecteur pour la base frontalière"""

    def __init__(self):
        super().__init__(Database.BORDER)

    def can_handle(self, query: SearchQuery) -> bool:
        return query.query_type in {SearchType.PASSPORT, SearchType.NAME,
                                    SearchType.NIU, SearchType.ALL}

    async def search(self, query: SearchQuery) -> List[SearchResultItem]:
        return []


class WatchlistConnector(DatabaseConnector):
    """Connecteur universel pour la watchlist"""

    def __init__(self):
        super().__init__(Database.WATCHLIST)

    def can_handle(self, query: SearchQuery) -> bool:
        return True  # Watchlist peut matcher n'importe quel type

    async def search(self, query: SearchQuery) -> List[SearchResultItem]:
        return []


# ============ ORCHESTRATEUR DE RECHERCHE FÉDÉRÉE ============
class FederatedSearchOrchestrator:
    """
    Orchestrateur de recherche fédérée.
    Distribue la requête à toutes les bases pertinentes en parallèle,
    agrège les résultats, et enrichit avec le contexte graphique Neo4j.
    """

    def __init__(self):
        self.connectors: Dict[Database, DatabaseConnector] = {
            Database.NCID: NCIDConnector(),
            Database.BIOMETRICS: BiometricConnector(),
            Database.BORDER: BorderConnector(),
            Database.WATCHLIST: WatchlistConnector(),
        }
        self.cache = {}  # Redis-like cache
        self.max_parallel = 15
        self.semaphore = asyncio.Semaphore(self.max_parallel)

    async def search(self, query: SearchQuery) -> SearchResponse:
        """Point d'entrée principal pour la recherche fédérée"""
        start_time = time.monotonic()
        response = SearchResponse(
            query=query.raw_query,
            page=query.page,
            limit=query.limit
        )

        # 1. Détection automatique du type de requête
        detected_type = QueryTypeDetector.detect(query.raw_query)
        if query.query_type == SearchType.ALL:
            query.query_type = detected_type

        # 2. Sélection des bases pertinentes
        target_databases = self._select_databases(query)

        # 3. Recherche parallèle dans toutes les bases
        tasks = []
        for db in target_databases:
            connector = self.connectors.get(db)
            if connector and connector.can_handle(query):
                tasks.append(self._search_with_timeout(connector, query))

        results_lists = await asyncio.gather(*tasks, return_exceptions=True)

        # 4. Agrégation et scoring des résultats
        all_results = []
        for db_results in results_lists:
            if isinstance(db_results, list):
                for item in db_results:
                    if item.database not in response.results:
                        response.results[item.database] = []
                    response.results[item.database].append(item)
                    all_results.append(item)

        # 5. Tri global par score
        all_results.sort(key=lambda r: r.score, reverse=True)

        # 6. Enrichissement graphique (Neo4j)
        if all_results:
            response.graph_context = await self._enrich_with_graph(query)

        # 7. Suggestions de recherche
        response.suggested_queries = self._generate_suggestions(query, all_results)

        response.total_results = len(all_results)
        response.databases_searched = len(target_databases)
        response.search_duration_ms = (time.monotonic() - start_time) * 1000

        logger.info(
            f"Search | query={query.raw_query} type={query.query_type.value} "
            f"databases={len(target_databases)} results={response.total_results} "
            f"duration={response.search_duration_ms:.1f}ms"
        )

        return response

    def _select_databases(self, query: SearchQuery) -> List[Database]:
        """Sélectionne les bases de données pertinentes selon le type de requête"""
        mapping = {
            SearchType.ALL: list(Database),
            SearchType.NIU: [Database.NCID, Database.BIOMETRICS, Database.CODIS,
                             Database.MISSING, Database.BORDER, Database.WATCHLIST,
                             Database.FINANCIAL, Database.DOCUMENTS],
            SearchType.NAME: [Database.NCID, Database.MISSING, Database.BORDER,
                              Database.WATCHLIST, Database.DOCUMENTS],
            SearchType.PHOTO: [Database.BIOMETRICS, Database.EVIDENCE],
            SearchType.FINGERPRINT: [Database.BIOMETRICS],
            SearchType.DNA: [Database.CODIS],
            SearchType.PHONE: [Database.NCID, Database.CYBER, Database.WATCHLIST],
            SearchType.EMAIL: [Database.CYBER, Database.WATCHLIST],
            SearchType.PLATE: [Database.ALPR, Database.VEHICLE, Database.WATCHLIST],
            SearchType.VIN: [Database.VEHICLE, Database.FIREARMS],
            SearchType.PASSPORT: [Database.BORDER, Database.DOCUMENTS, Database.WATCHLIST],
            SearchType.IP: [Database.CYBER],
            SearchType.DOMAIN: [Database.CYBER],
            SearchType.WALLET: [Database.CYBER, Database.FINANCIAL],
            SearchType.CASE: [Database.NCID, Database.MISSING, Database.FIREARMS],
            SearchType.LOCATION: [Database.GEOINT, Database.ALPR, Database.NARCOTICS],
        }
        return mapping.get(query.query_type, list(Database))

    async def _search_with_timeout(self, connector: DatabaseConnector,
                                   query: SearchQuery) -> List[SearchResultItem]:
        """Recherche avec timeout pour une base"""
        try:
            async with self.semaphore:
                return await asyncio.wait_for(
                    connector.search(query),
                    timeout=query.timeout_ms / 1000
                )
        except asyncio.TimeoutError:
            logger.warning(f"Timeout on {connector.database.value} for {query.raw_query}")
            return []
        except Exception as e:
            logger.error(f"Error on {connector.database.value}: {e}")
            return []

    async def _enrich_with_graph(self, query: SearchQuery) -> Dict:
        """Enrichit les résultats avec le contexte du graphe Neo4j"""
        # Neo4j query:
        # MATCH (n)
        # WHERE n.niu = $query OR n.full_name CONTAINS $query
        #       OR n.plate = $query OR n.phone = $query
        # OPTIONAL MATCH (n)-[r]-(connected)
        # RETURN n, collect(r), collect(connected)
        return {
            "nodes": [],
            "edges": [],
            "depth": 2,
            "centrality": 0.0
        }

    def _generate_suggestions(self, query: SearchQuery,
                              results: List[SearchResultItem]) -> List[str]:
        """Génère des suggestions de recherche basées sur les résultats"""
        suggestions = []
        if len(results) > 0:
            # Suggérer des recherches connexes basées sur le type de résultat dominant
            suggestions.append(f"Graph search for: {query.raw_query}")
            suggestions.append(f"Biometric match for: {query.raw_query}")
        return suggestions


# ============ RECHERCHE GRAPHIQUE (NEO4J) ============
class GraphSearchEngine:
    """Moteur de recherche graphique Neo4j pour l'intelligence de liens"""

    def __init__(self, neo4j_uri: str = "bolt://sniside-neo4j:7687"):
        self.neo4j_uri = neo4j_uri

    async def search(self, entity_id: str, entity_type: str, depth: int = 2,
                     relationship_types: List[str] = None) -> Dict:
        """
        Recherche graphique étendue avec détection de patterns criminels.

        Exemple de requête Cypher générée:
        MATCH path = (start {niu: $entity_id})
        CALL apoc.path.subgraph(start, {
            maxLevel: $depth,
            relationshipFilter: 'OWNS|USES|ASSOCIATED_WITH|FINANCED_BY|LINKED_TO'
        })
        YIELD path
        RETURN path
        """
        return {
            "entity_id": entity_id,
            "entity_type": entity_type,
            "nodes": [],
            "edges": [],
            "detected_patterns": [],
            "centrality": 0.0,
            "network_size": 0
        }

    def detect_patterns(self, graph: Dict) -> List[str]:
        """Détecte les patterns criminels dans le graphe"""
        patterns = []
        # Pattern: Money Laundering Cycle
        # Pattern: Criminal Network Cluster
        # Pattern: Human Trafficking Route
        # Pattern: Narcotics Supply Chain
        return patterns


# ============ API RAPIDE POUR LE MOTEUR DE RECHERCHE ============
class NationalSearchAPI:
    """API du Moteur de Recherche National Souverain"""

    def __init__(self):
        self.federated = FederatedSearchOrchestrator()
        self.graph = GraphSearchEngine()
        self.query_detector = QueryTypeDetector()

    async def unified_search(self, query_str: str, search_type: str = "ALL",
                             databases: List[str] = None, page: int = 1,
                             limit: int = 20, fuzzy: bool = True) -> SearchResponse:
        """Point d'entrée API pour la recherche unifiée"""
        query = SearchQuery(
            raw_query=query_str,
            query_type=SearchType(search_type.upper()),
            databases=[Database(db) for db in databases] if databases else [],
            page=page,
            limit=limit,
            fuzzy=fuzzy
        )
        return await self.federated.search(query)

    async def graph_search(self, niu: str = None, phone: str = None,
                           plate: str = None, depth: int = 2,
                           relationship_types: List[str] = None) -> Dict:
        """Recherche graphique"""
        entity_id = niu or phone or plate
        entity_type = "NIU" if niu else "PHONE" if phone else "PLATE"
        return await self.graph.search(entity_id, entity_type, depth,
                                       relationship_types)

    async def cross_reference(self, entity_id: str, entity_type: str,
                              depth: int = 2) -> Dict:
        """Recherche de référence croisée à travers toutes les bases"""
        # 1. Graph search
        graph_results = await self.graph.search(entity_id, entity_type, depth)

        # 2. Federated search on all known identifiers
        federated_results = await self.federated.search(
            SearchQuery(raw_query=entity_id, databases=list(Database))
        )

        return {
            "entity_id": entity_id,
            "entity_type": entity_type,
            "graph": graph_results,
            "federated": {
                "total_results": federated_results.total_results,
                "by_database": {
                    db.value: len(items)
                    for db, items in federated_results.results.items()
                }
            }
        }

    def health(self) -> Dict:
        """Vérification de santé du moteur de recherche"""
        return {
            "status": "healthy",
            "connectors": len(self.federated.connectors),
            "databases_supported": len(Database),
            "search_types": len(SearchType),
            "cache_available": len(self.federated.cache) > 0
        }


# ============ POINT D'ENTRÉE ============
search_engine = NationalSearchAPI()

# Exemple d'utilisation:
# import asyncio
# result = asyncio.run(search_engine.unified_search("Jean Dupont", "NAME"))
# print(f"Found {result.total_results} results in {result.search_duration_ms:.1f}ms")
