from fastapi import APIRouter
from app.utils.errors import JobsBotError

router = APIRouter()


@router.get("/health")
def health():
    try:
        return {"status": "ok"}
    except Exception as e:
        log_event("ERROR", str(e))
        raise JobsBotError("API001", {"error": str(e)})
