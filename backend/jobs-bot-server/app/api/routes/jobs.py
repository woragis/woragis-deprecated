from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List, Optional

router = APIRouter()


class Job(BaseModel):
    id: int
    title: str
    description: Optional[str] = None
    status: str


jobs_db = []


@router.get("/jobs", response_model=List[Job])
def list_jobs():
    return jobs_db


@router.post("/jobs", response_model=Job)
def create_job(job: Job):
    jobs_db.append(job)
    return job


@router.get("/jobs/{job_id}", response_model=Job)
def get_job(job_id: int):
    for job in jobs_db:
        if job.id == job_id:
            return job
    raise HTTPException(status_code=404, detail="Job not found")


@router.put("/jobs/{job_id}", response_model=Job)
def update_job(job_id: int, job: Job):
    for idx, j in enumerate(jobs_db):
        if j.id == job_id:
            jobs_db[idx] = job
            return job
    raise HTTPException(status_code=404, detail="Job not found")


@router.delete("/jobs/{job_id}")
def delete_job(job_id: int):
    for idx, j in enumerate(jobs_db):
        if j.id == job_id:
            del jobs_db[idx]
            return {"detail": "Job deleted"}
    raise HTTPException(status_code=404, detail="Job not found")
