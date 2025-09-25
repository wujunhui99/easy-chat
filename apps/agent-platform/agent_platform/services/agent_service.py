import uuid
from typing import List, Optional

from sqlalchemy.orm import Session

from agent_platform.core import models
from agent_platform.repositories.agent_repository import AgentRepository
from agent_platform.services.user_rpc import UserRpcClient


class AgentService:
    def __init__(self, session: Session, user_client: Optional[UserRpcClient] = None) -> None:
        self._repository = AgentRepository(session)
        self._session = session
        self._user_client = user_client or UserRpcClient()

    def list_agents(self) -> List[models.Agent]:
        return self._repository.list_agents()

    def get_agent(self, agent_id: str) -> Optional[models.Agent]:
        return self._repository.get_by_id(agent_id)

    def create_agent(
        self,
        *,
        code: str,
        name: str,
        description: str,
        prompt: str,
        model: str,
        tools: list,
        memory_strategy: str,
        config: dict,
        creator_user_id: str,
        nickname: str,
        avatar: str,
        phone: str,
        sex: int = 0,
    ) -> models.Agent:
        if self._repository.get_by_code(code):
            raise ValueError(f"Agent with code {code} already exists")

        agent_id = uuid.uuid4().hex
        user_id = self._user_client.create_agent(
            nickname=nickname or name,
            avatar=avatar,
            phone=phone,
            sex=sex,
        )

        agent = models.Agent(
            id=agent_id,
            user_id=user_id,
            code=code,
            name=name,
            description=description,
            status="draft",
            model=model,
            prompt=prompt,
            tools=tools,
            memory_strategy=memory_strategy,
            config=config,
            created_by=creator_user_id,
        )
        self._repository.add(agent)

        version = models.AgentVersion(
            agent_id=agent.id,
            version=self._repository.next_version(agent.id),
            prompt=prompt,
            tools=tools,
            config=config,
            created_by=creator_user_id,
        )
        self._repository.add_version(version)
        return agent
