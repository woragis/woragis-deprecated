import os
from typing import Dict, Callable
from langchain_core.language_models import BaseChatModel
from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import Runnable


def _persona_prompt(persona_description: str) -> ChatPromptTemplate:
    system = (
        "You are {agent_name}. "
        + persona_description
        + " Be concise, structured, and provide actionable insights. "
          "When assumptions are needed, state them explicitly."
    )
    return ChatPromptTemplate.from_messages(
        [
            ("system", system),
            ("human", "{input}"),
        ]
    )


def _build_chain(agent_name: str, persona_description: str, model: BaseChatModel) -> Runnable:
    prompt = _persona_prompt(persona_description)
    return prompt | model


def _model() -> BaseChatModel:
    # Default OpenAI model; used when no provider specified
    return ChatOpenAI(
        model=os.getenv("OPENAI_MODEL", "gpt-4o-mini"),
        temperature=float(os.getenv("OPENAI_TEMPERATURE", "0.3")),
        timeout=60,
    )


AGENTS: Dict[str, Callable[[], Runnable]] = {
    "economist": lambda: _build_chain(
        "Economist Agent",
        "A seasoned macro- and micro-economist. Analyze markets, policy impacts, "
        "unit economics, and risks. Support with data references when possible.",
        _model(),
    ),
    "strategist": lambda: _build_chain(
        "Strategist Agent",
        "An experienced business strategist. Craft positioning, go-to-market plans, "
        "competitive analysis, and prioritization frameworks.",
        _model(),
    ),
    "entrepreneur": lambda: _build_chain(
        "Entrepreneur Agent",
        "A pragmatic startup founder. Bias to action, validate quickly, propose MVP scope, "
        "scrappy tactics, and milestone-based plans.",
        _model(),
    ),
    "startup": lambda: _build_chain(
        "Startup Advisor Agent",
        "A startup advisor and mentor. Blend product thinking, growth, fundraising, "
        "and execution advice tailored to stage.",
        _model(),
    ),
}


def get_agent_names() -> list[str]:
    return sorted(AGENTS.keys())


def get_agent(name: str) -> Runnable | None:
    key = name.lower().strip()
    factory = AGENTS.get(key)
    if not factory:
        return None
    return factory()


# Build an agent chain using a provided model (supports multi-provider callers)
def build_agent_with_model(agent_name: str, model: BaseChatModel) -> Runnable | None:
    key = agent_name.lower().strip()
    personas = {
        "economist": "A seasoned macro- and micro-economist. Analyze markets, policy impacts, "
        "unit economics, and risks. Support with data references when possible.",
        "strategist": "An experienced business strategist. Craft positioning, go-to-market plans, "
        "competitive analysis, and prioritization frameworks.",
        "entrepreneur": "A pragmatic startup founder. Bias to action, validate quickly, propose MVP scope, "
        "scrappy tactics, and milestone-based plans.",
        "startup": "A startup advisor and mentor. Blend product thinking, growth, fundraising, "
        "and execution advice tailored to stage.",
    }
    desc = personas.get(key)
    if not desc:
        return None
    return _build_chain(agent_name.title() + " Agent", desc, model)


def build_system_message(agent_name: str) -> str | None:
    key = agent_name.lower().strip()
    personas = {
        "economist": "A seasoned macro- and micro-economist. Analyze markets, policy impacts, "
        "unit economics, and risks. Support with data references when possible.",
        "strategist": "An experienced business strategist. Craft positioning, go-to-market plans, "
        "competitive analysis, and prioritization frameworks.",
        "entrepreneur": "A pragmatic startup founder. Bias to action, validate quickly, propose MVP scope, "
        "scrappy tactics, and milestone-based plans.",
        "startup": "A startup advisor and mentor. Blend product thinking, growth, fundraising, "
        "and execution advice tailored to stage.",
    }
    desc = personas.get(key)
    if not desc:
        return None
    return (
        f"You are {agent_name.title()} Agent. "
        + desc
        + " Be concise, structured, and provide actionable insights. "
          "When assumptions are needed, state them explicitly."
    )


