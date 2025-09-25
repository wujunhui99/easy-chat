from fastapi import FastAPI

from agent_platform.api import routes
from agent_platform.core.db import get_engine
from agent_platform.core.models import Base


def create_app() -> FastAPI:
    app = FastAPI(title="Agent Platform")

    # Ensure tables exist for quick start scenarios.
    engine = get_engine()
    Base.metadata.create_all(bind=engine)

    app.include_router(routes.router)
    return app


app = create_app()
