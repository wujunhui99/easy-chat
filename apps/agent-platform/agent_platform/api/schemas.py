from datetime import datetime
from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field


class AgentCreateRequest(BaseModel):
    code: str = Field(..., regex=r"^[a-zA-Z0-9_\-]+$")
    name: str
    description: Optional[str] = None
    prompt: str
    model: str = Field(default="gpt-3.5-turbo")
    tools: List[Dict[str, Any]] = Field(default_factory=list)
    memory_strategy: str = Field(default="recent")
    config: Dict[str, Any] = Field(default_factory=dict)
    creator_user_id: str
    nickname: Optional[str] = None
    avatar: str = Field(default="")
    phone: str = Field(default="")
    sex: int = Field(default=0, ge=0, le=2)


class AgentResponse(BaseModel):
    id: str
    user_id: str
    code: str
    name: str
    description: Optional[str]
    status: str
    model: str
    prompt: str
    tools: List[Dict[str, Any]]
    memory_strategy: str
    config: Dict[str, Any]
    created_by: str
    updated_by: Optional[str]
    created_at: datetime
    updated_at: datetime

    class Config:
        orm_mode = True


class AgentListResponse(BaseModel):
    items: List[AgentResponse]
