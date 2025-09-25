from contextlib import contextmanager
from typing import Generator, Iterator

from sqlalchemy import create_engine
from sqlalchemy.orm import Session, sessionmaker

from agent_platform.config import get_settings

_engine = create_engine(get_settings().database_url, pool_pre_ping=True, future=True)
_Session = sessionmaker(bind=_engine, class_=Session, expire_on_commit=False, autoflush=False)


@contextmanager
def session_scope() -> Iterator[Session]:
    session = _Session()
    try:
        yield session
        session.commit()
    except Exception:
        session.rollback()
        raise
    finally:
        session.close()


def get_session() -> Generator[Session, None, None]:
    session = _Session()
    try:
        yield session
    finally:
        session.close()


def get_engine():
    return _engine
