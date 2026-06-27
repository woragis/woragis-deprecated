from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
import os
from app.utils.errors import JobsBotError

DATABASE_URL = os.getenv("DATABASE_URL", "sqlite:///jobsbot.db")

try:
    engine = create_engine(DATABASE_URL, connect_args={
                           "check_same_thread": False})
    SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
except Exception as e:
    raise JobsBotError("DB001", {"url": DATABASE_URL, "error": str(e)})
