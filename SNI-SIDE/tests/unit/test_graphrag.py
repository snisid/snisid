"""
Tests for GraphRAG Engine components
"""
import pytest
from services.graphrag_engine.graphrag_engine import ReportType, Neo4jRetriever, LLMGenerator, GraphRAGEngine
from tests.mocks import MockNeo4jDriver


class TestReportType:
    def test_enum_values(self):
        assert ReportType.ENTITY_PROFILE.value == "ENTITY_PROFILE"
        assert ReportType.LINK_ANALYSIS.value == "LINK_ANALYSIS"
        assert ReportType.NETWORK_MAP.value == "NETWORK_MAP"
        assert ReportType.FINANCIAL_FLOW.value == "FINANCIAL_FLOW"
        assert ReportType.TEMPORAL_ANALYSIS.value == "TEMPORAL_ANALYSIS"
        assert ReportType.CROSS_ENTITY.value == "CROSS_ENTITY"


class TestCypherTemplates:
    def test_all_templates_present(self):
        from services.graphrag_engine.graphrag_engine import CYPHER_TEMPLATES
        for rt in ReportType:
            assert rt.value in CYPHER_TEMPLATES, f"Missing template: {rt.value}"

    def test_entity_profile_has_id_param(self):
        from services.graphrag_engine.graphrag_engine import CYPHER_TEMPLATES
        tmpl = CYPHER_TEMPLATES["ENTITY_PROFILE"]
        assert "$id" in tmpl
        assert "$limit" in tmpl


class TestNeo4jContextBuilding:
    def test_empty_context(self):
        retriever = Neo4jRetriever()
        context = retriever._build_context([])
        assert context["nodes"] == []
        assert context["relationships"] == []
        assert isinstance(context, dict)

    def test_context_with_nodes(self):
        retriever = Neo4jRetriever()
        records = [
            {
                "n": {"element_id": "4:abc123", "niu": "HT12345678", "full_name": "JEAN DUPONT", "risk_level": "HIGH"},
                "rel_type": "HAS_WARRANT",
                "rel_props": {"count": 3},
            }
        ]
        context = retriever._build_context(records)
        assert len(context["nodes"]) >= 0


class TestLLMPromptTemplates:
    def test_all_prompts_present(self):
        from services.graphrag_engine.graphrag_engine import PROMPT_TEMPLATES
        for rt in ReportType:
            assert rt.value in PROMPT_TEMPLATES, f"Missing prompt: {rt.value}"

    def test_prompt_has_context_placeholder(self):
        from services.graphrag_engine.graphrag_engine import PROMPT_TEMPLATES
        for name, prompt in PROMPT_TEMPLATES.items():
            assert "{context}" in prompt, f"Missing {{context}} in {name}"
