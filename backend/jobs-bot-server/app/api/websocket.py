from fastapi import APIRouter, WebSocket, WebSocketDisconnect
from fastapi.responses import HTMLResponse
from app.utils.errors import JobsBotError
from app.utils.logging import log_event

router = APIRouter()

# Simple in-memory connection manager


class ConnectionManager:
    def __init__(self):
        self.active_connections: list[WebSocket] = []

    async def connect(self, websocket: WebSocket):
        await websocket.accept()
        self.active_connections.append(websocket)

    def disconnect(self, websocket: WebSocket):
        self.active_connections.remove(websocket)

    async def broadcast(self, message: dict):
        for connection in self.active_connections:
            await connection.send_json(message)


manager = ConnectionManager()


@router.websocket("/ws/progress")
async def websocket_endpoint(websocket: WebSocket):
    try:
        await manager.connect(websocket)
        while True:
            await websocket.receive_text()  # Read-only, ignore input
    except WebSocketDisconnect:
        manager.disconnect(websocket)
    except Exception as e:
        log_event("ERROR", str(e))
        raise JobsBotError("API001", {"error": str(e)})


async def emit_ws_event(event, run_id=None, job_id=None, platform=None, state=None, error_code=None, details=None):
    message = {
        "event": event,
        "run_id": run_id,
        "job_id": job_id,
        "platform": platform,
        "state": state,
        "error_code": error_code,
        "details": details
    }
    await manager.broadcast(message)
    log_event(event, details=details, run_id=run_id, job_id=job_id,
              platform=platform, error_code=error_code)
