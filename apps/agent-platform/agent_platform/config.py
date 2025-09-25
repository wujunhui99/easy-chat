from functools import lru_cache
from pydantic import BaseSettings, Field


class Settings(BaseSettings):
    app_name: str = Field(default="agent-platform")
    database_url: str = Field(default="sqlite:///./agent_platform.db")
    user_rpc_target: str = Field(default="127.0.0.1:9000")
    user_rpc_timeout_seconds: float = Field(default=3.0)

    class Config:
        env_file = ".env"


@lru_cache()
def get_settings() -> Settings:
    return Settings()
