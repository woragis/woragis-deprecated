from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import Optional

router = APIRouter()


class Settings(BaseModel):
    notifications_enabled: bool = True
    max_concurrent_runs: int = 1
    default_platform: Optional[str] = None


# In-memory settings store
settings_db = Settings()


@router.get("/settings", response_model=Settings)
def get_settings():
    return settings_db


@router.put("/settings", response_model=Settings)
def update_settings(new_settings: Settings):
    global settings_db
    settings_db = new_settings
    return settings_db
