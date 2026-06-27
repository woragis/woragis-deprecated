from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List, Optional
from app.utils.errors import JobsBotError
from app.utils.logging import log_event

ALLOWED_PLATFORMS = {"greenhouse", "lever", "linkedin", "gupy"}

router = APIRouter()


class RunRequest(BaseModel):
    max_applications: int
    platforms: List[str]

    @classmethod
    def validate_platforms(cls, platforms):
        invalid = [p for p in platforms if p not in ALLOWED_PLATFORMS]
        if invalid:
            raise JobsBotError("VALID001", {"invalid_platforms": invalid})

    @classmethod
    def validate_max_applications(cls, max_applications):
        if not isinstance(max_applications, int) or max_applications <= 0:
            raise JobsBotError(
                "VALID002", {"max_applications": max_applications})

    @classmethod
    def validate(cls, values):
        cls.validate_platforms(values.get("platforms", []))
        cls.validate_max_applications(values.get("max_applications", 0))
        return values

    class Config:
        extra = "forbid"

        @staticmethod
        def schema_extra(schema, model):
            schema["example"] = {
                "max_applications": 10,
                "platforms": ["greenhouse", "lever"]
            }

    @classmethod
    def __get_validators__(cls):
        yield cls.validate


class RunStatus(BaseModel):
    run_id: str
    state: str
    progress: Optional[int] = 0
    total: Optional[int] = 0


# In-memory placeholder for runs
runs = {}

# List all runs


@router.get("/runs", response_model=List[RunStatus])
def list_runs():
    return [RunStatus(run_id=run_id, state=run["state"], progress=run["progress"], total=run["total"]) for run_id, run in runs.items()]

# Update a run (e.g., change state, progress, total)


class RunUpdate(BaseModel):
    state: Optional[str] = None
    progress: Optional[int] = None
    total: Optional[int] = None


@router.put("/runs/{run_id}", response_model=RunStatus)
def update_run(run_id: str, update: RunUpdate):
    run = runs.get(run_id)
    if not run:
        raise HTTPException(status_code=404, detail="Run not found")
    if update.state is not None:
        run["state"] = update.state
    if update.progress is not None:
        run["progress"] = update.progress
    if update.total is not None:
        run["total"] = update.total
    return RunStatus(run_id=run_id, state=run["state"], progress=run["progress"], total=run["total"])

# Delete a run


@router.delete("/runs/{run_id}")
def delete_run(run_id: str):
    if run_id in runs:
        del runs[run_id]
        return {"detail": "Run deleted"}
    raise HTTPException(status_code=404, detail="Run not found")


@router.post("/runs", response_model=RunStatus)
def start_run(request: RunRequest):
    try:
        # Validate and create a new run (stub)
        run_id = "run_1"  # TODO: generate unique
        if run_id in runs:
            raise JobsBotError("API003", {"run_id": run_id})
        runs[run_id] = {"state": "PENDING", "progress": 0,
                        "total": request.max_applications}
        return RunStatus(run_id=run_id, state="PENDING", progress=0, total=request.max_applications)
    except JobsBotError as e:
        log_event("ERROR", str(e))
        raise HTTPException(status_code=400, detail={
                            "code": e.code, "message": e.message, "context": e.context})


@router.get("/runs/{run_id}", response_model=RunStatus)
def get_run_status(run_id: str):
    try:
        run = runs.get(run_id)
        if not run:
            raise JobsBotError("API002", {"run_id": run_id})
        return RunStatus(run_id=run_id, state=run["state"], progress=run["progress"], total=run["total"])
    except JobsBotError as e:
        log_event("ERROR", str(e))
        raise HTTPException(status_code=404, detail={
                            "code": e.code, "message": e.message, "context": e.context})


@router.post("/runs/{run_id}/stop")
def stop_run(run_id: str):
    try:
        run = runs.get(run_id)
        if not run:
            raise JobsBotError("API002", {"run_id": run_id})
        run["state"] = "CANCELLED"
        return {"status": "stopped"}
    except JobsBotError as e:
        log_event("ERROR", str(e))
        raise HTTPException(status_code=404, detail={
                            "code": e.code, "message": e.message, "context": e.context})
