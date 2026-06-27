from app.utils.errors import JobsBotError
from app.workers.states import RunState
from typing import Optional


class Run:
    def __init__(self, run_id: str, max_applications: int, platforms: list[str]):
        try:
            self.run_id = run_id
            self.max_applications = max_applications
            self.platforms = platforms
            self.state = RunState.PENDING
            self.progress = 0
            self.total = max_applications
        except Exception as e:
            raise JobsBotError("API001", {"error": str(e)})

    def start(self):
        try:
            self.state = RunState.RUNNING
        except Exception as e:
            raise JobsBotError("API001", {"error": str(e)})

    def complete(self):
        try:
            self.state = RunState.COMPLETED
        except Exception as e:
            raise JobsBotError("API001", {"error": str(e)})

    def fail(self):
        try:
            self.state = RunState.FAILED
        except Exception as e:
            raise JobsBotError("API001", {"error": str(e)})

    def cancel(self):
        try:
            self.state = RunState.CANCELLED
        except Exception as e:
            raise JobsBotError("API001", {"error": str(e)})
