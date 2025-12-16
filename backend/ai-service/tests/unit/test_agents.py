"""
Unit tests for agent registry and agent building.
"""
import pytest
import os
from unittest.mock import Mock, patch, MagicMock
from langchain_core.language_models import BaseChatModel
from langchain_core.runnables import Runnable

from app.agents import (
    get_agent_names,
    get_agent,
    build_agent_with_model,
    build_system_message
)


class TestAgentRegistry:
    """Tests for agent registry functions."""

    def test_get_agent_names(self):
        """Test getting list of available agent names."""
        names = get_agent_names()
        assert isinstance(names, list)
        assert len(names) > 0
        assert "economist" in names
        assert "strategist" in names
        assert "entrepreneur" in names
        assert "startup" in names
        # Should be sorted
        assert names == sorted(names)

    @patch.dict(os.environ, {"OPENAI_API_KEY": "test-key"})
    @patch('app.agents.registry.ChatOpenAI')
    def test_get_agent_valid_name(self, mock_chat_openai, mock_chat_model):
        """Test getting an agent with a valid name."""
        mock_chat_openai.return_value = mock_chat_model
        agent = get_agent("startup")
        assert agent is not None
        assert isinstance(agent, Runnable)

    @patch.dict(os.environ, {"OPENAI_API_KEY": "test-key"})
    @patch('app.agents.registry.ChatOpenAI')
    def test_get_agent_case_insensitive(self, mock_chat_openai, mock_chat_model):
        """Test that agent names are case-insensitive."""
        mock_chat_openai.return_value = mock_chat_model
        agent1 = get_agent("STARTUP")
        agent2 = get_agent("startup")
        agent3 = get_agent("Startup")
        assert agent1 is not None
        assert agent2 is not None
        assert agent3 is not None

    @patch.dict(os.environ, {"OPENAI_API_KEY": "test-key"})
    @patch('app.agents.registry.ChatOpenAI')
    def test_get_agent_with_whitespace(self, mock_chat_openai, mock_chat_model):
        """Test that agent names are trimmed."""
        mock_chat_openai.return_value = mock_chat_model
        agent = get_agent("  startup  ")
        assert agent is not None

    def test_get_agent_invalid_name(self):
        """Test getting an agent with an invalid name."""
        agent = get_agent("nonexistent")
        assert agent is None

    @pytest.mark.parametrize("agent_name", ["economist", "strategist", "entrepreneur", "startup"])
    @patch.dict(os.environ, {"OPENAI_API_KEY": "test-key"})
    @patch('app.agents.registry.ChatOpenAI')
    def test_get_agent_all_agents(self, mock_chat_openai, agent_name, mock_chat_model):
        """Test that all defined agents can be retrieved."""
        mock_chat_openai.return_value = mock_chat_model
        agent = get_agent(agent_name)
        assert agent is not None
        assert isinstance(agent, Runnable)


class TestBuildAgentWithModel:
    """Tests for building agents with custom models."""

    def test_build_agent_with_model_valid(self, mock_chat_model):
        """Test building an agent with a valid model."""
        agent = build_agent_with_model("startup", mock_chat_model)
        assert agent is not None
        assert isinstance(agent, Runnable)

    def test_build_agent_with_model_invalid_name(self, mock_chat_model):
        """Test building an agent with invalid name."""
        agent = build_agent_with_model("nonexistent", mock_chat_model)
        assert agent is None

    def test_build_agent_with_model_case_insensitive(self, mock_chat_model):
        """Test that agent names are case-insensitive."""
        agent1 = build_agent_with_model("STARTUP", mock_chat_model)
        agent2 = build_agent_with_model("startup", mock_chat_model)
        assert agent1 is not None
        assert agent2 is not None

    @pytest.mark.parametrize("agent_name", ["economist", "strategist", "entrepreneur", "startup"])
    def test_build_agent_with_model_all_agents(self, agent_name, mock_chat_model):
        """Test building all agents with custom model."""
        agent = build_agent_with_model(agent_name, mock_chat_model)
        assert agent is not None
        assert isinstance(agent, Runnable)


class TestBuildSystemMessage:
    """Tests for building system messages."""

    def test_build_system_message_valid(self):
        """Test building system message for valid agent."""
        message = build_system_message("startup")
        assert message is not None
        assert isinstance(message, str)
        assert "Startup Agent" in message
        assert len(message) > 0

    def test_build_system_message_invalid(self):
        """Test building system message for invalid agent."""
        message = build_system_message("nonexistent")
        assert message is None

    def test_build_system_message_case_insensitive(self):
        """Test that agent names are case-insensitive."""
        msg1 = build_system_message("STARTUP")
        msg2 = build_system_message("startup")
        assert msg1 is not None
        assert msg2 is not None
        assert msg1 == msg2

    @pytest.mark.parametrize("agent_name", ["economist", "strategist", "entrepreneur", "startup"])
    def test_build_system_message_all_agents(self, agent_name):
        """Test building system messages for all agents."""
        message = build_system_message(agent_name)
        assert message is not None
        assert isinstance(message, str)
        assert agent_name.title() in message or agent_name.lower() in message.lower()

    def test_build_system_message_contains_persona(self):
        """Test that system message contains persona description."""
        message = build_system_message("startup")
        assert "startup advisor" in message.lower() or "mentor" in message.lower()
