# Agent Platform

Python-based microservice that manages conversational agents, their personas, and integration with the Go user service.

## Features

- FastAPI management API for CRUD operations on agents.
- SQLAlchemy models for `agents` and `agent_versions` tables.
- gRPC client that calls the Go user RPC `CreateAgent` to provision IM accounts.

## Getting Started

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn agent_platform.api.app:app --reload
```

Configure connection details with environment variables:

- `DATABASE_URL` – SQLAlchemy connection string (default `sqlite:///./agent_platform.db`).
- `USER_RPC_TARGET` – gRPC endpoint of the Go user service.
- `USER_RPC_TIMEOUT_SECONDS` – timeout for RPC calls.

Before running in production, apply SQL schemas located in `deploy/sql/agent_platform.sql` to your MySQL instance.
