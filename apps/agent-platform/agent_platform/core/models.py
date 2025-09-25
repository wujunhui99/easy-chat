from datetime import datetime
from typing import Any, Dict, Optional

from sqlalchemy import JSON, Column, DateTime, Integer, String, Text
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class Agent(Base):
    __tablename__ = "agents"

    id = Column(String(32), primary_key=True)
    user_id = Column(String(24), nullable=False, unique=True)
    code = Column(String(64), nullable=False, unique=True)
    name = Column(String(128), nullable=False)
    description = Column(Text, nullable=True)
    status = Column(String(32), nullable=False, default="draft")
    model = Column(String(128), nullable=False, default="gpt-3.5-turbo")
    prompt = Column(Text, nullable=False)
    tools = Column(JSON, nullable=False, default=list)
    memory_strategy = Column(String(64), nullable=False, default="recent")
    config = Column(JSON, nullable=False, default=dict)
    created_by = Column(String(24), nullable=False)
    updated_by = Column(String(24), nullable=True)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=False, default=datetime.utcnow, onupdate=datetime.utcnow)


class AgentVersion(Base):
    __tablename__ = "agent_versions"

    id = Column(Integer, primary_key=True, autoincrement=True)
    agent_id = Column(String(32), nullable=False, index=True)
    version = Column(Integer, nullable=False)
    prompt = Column(Text, nullable=False)
    tools = Column(JSON, nullable=False, default=list)
    config = Column(JSON, nullable=False, default=dict)
    created_by = Column(String(24), nullable=False)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)


def agent_to_dict(agent: Agent) -> Dict[str, Any]:
    return {
        "id": agent.id,
        "user_id": agent.user_id,
        "code": agent.code,
        "name": agent.name,
        "description": agent.description,
        "status": agent.status,
        "model": agent.model,
        "prompt": agent.prompt,
        "tools": agent.tools,
        "memory_strategy": agent.memory_strategy,
        "config": agent.config,
        "created_by": agent.created_by,
        "updated_by": agent.updated_by,
        "created_at": agent.created_at,
        "updated_at": agent.updated_at,
    }


def version_to_dict(version: AgentVersion) -> Dict[str, Any]:
    return {
        "id": version.id,
        "agent_id": version.agent_id,
        "version": version.version,
        "prompt": version.prompt,
        "tools": version.tools,
        "config": version.config,
        "created_by": version.created_by,
        "created_at": version.created_at,
    }
