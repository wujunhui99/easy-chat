from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from agent_platform.api.schemas import AgentCreateRequest, AgentListResponse, AgentResponse
from agent_platform.core.db import get_session
from agent_platform.core.models import agent_to_dict
from agent_platform.services.agent_service import AgentService

router = APIRouter(prefix="/agents", tags=["agents"])


def get_agent_service(session: Session = Depends(get_session)) -> AgentService:
    return AgentService(session=session)


@router.get("", response_model=AgentListResponse)
def list_agents(service: AgentService = Depends(get_agent_service)) -> AgentListResponse:
    agents = service.list_agents()
    return AgentListResponse(items=[AgentResponse.from_orm(agent) for agent in agents])


@router.post("", response_model=AgentResponse, status_code=201)
def create_agent(
    request: AgentCreateRequest,
    service: AgentService = Depends(get_agent_service),
) -> AgentResponse:
    try:
        agent = service.create_agent(
            code=request.code,
            name=request.name,
            description=request.description or "",
            prompt=request.prompt,
            model=request.model,
            tools=request.tools,
            memory_strategy=request.memory_strategy,
            config=request.config,
            creator_user_id=request.creator_user_id,
            nickname=request.nickname or request.name,
            avatar=request.avatar,
            phone=request.phone,
            sex=request.sex,
        )
    except ValueError as exc:
        raise HTTPException(status_code=400, detail=str(exc)) from exc

    return AgentResponse.from_orm(agent)


@router.get("/{agent_id}", response_model=AgentResponse)
def get_agent(
    agent_id: str,
    service: AgentService = Depends(get_agent_service),
) -> AgentResponse:
    agent = service.get_agent(agent_id)
    if not agent:
        raise HTTPException(status_code=404, detail="Agent not found")
    return AgentResponse.from_orm(agent)
