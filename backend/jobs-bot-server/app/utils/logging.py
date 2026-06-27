import logging
from app.persistence.models import Log
from app.persistence.db import SessionLocal
from datetime import datetime

logging.basicConfig(
    filename='logs/runs.log',
    level=logging.INFO,
    format='%(asctime)s %(levelname)s %(message)s'
)


def log_event(event, details=None, run_id=None, job_id=None, platform=None, error_code=None):
    context = {
        "run_id": run_id,
        "job_id": job_id,
        "platform": platform,
        "error_code": error_code,
        "details": details,
        "timestamp": datetime.utcnow().isoformat()
    }
    logging.info(f"{event}: {context}")
    # Persist to DB
    try:
        session = SessionLocal()
        log = Log(run_id=run_id, event=event,
                  timestamp=datetime.utcnow(), details=str(context))
        log.save(session)
        session.close()
    except Exception as e:
        logging.error(f"Failed to persist log: {e}")
