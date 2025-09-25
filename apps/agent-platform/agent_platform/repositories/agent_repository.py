from typing import List, Optional

from sqlalchemy import select
from sqlalchemy.orm import Session

from agent_platform.core import models


class AgentRepository:
    def __init__(self, session: Session) -> None:
        self._session = session

    def list_agents(self) -> List[models.Agent]:
        statement = select(models.Agent).order_by(models.Agent.created_at.desc())
        return list(self._session.execute(statement).scalars())

    def get_by_id(self, agent_id: str) -> Optional[models.Agent]:
        return self._session.get(models.Agent, agent_id)

    def get_by_code(self, code: str) -> Optional[models.Agent]:
        statement = select(models.Agent).where(models.Agent.code == code)
        return self._session.execute(statement).scalar_one_or_none()

    def add(self, agent: models.Agent) -> models.Agent:
        self._session.add(agent)
        return agent

    def add_version(self, version: models.AgentVersion) -> models.AgentVersion:
        self._session.add(version)
        return version

    def next_version(self, agent_id: str) -> int:
        statement = select(models.AgentVersion.version).where(models.AgentVersion.agent_id == agent_id)
        result = self._session.execute(statement).scalars().all()
        return max(result, default=0) + 1
