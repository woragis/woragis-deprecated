# FastAPI entrypoint

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.routes import runs, health, jobs, settings
from app.api import websocket

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(runs.router)
app.include_router(health.router)
app.include_router(jobs.router)
app.include_router(settings.router)
app.include_router(websocket.router)
