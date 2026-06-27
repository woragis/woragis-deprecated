from queue import Queue
from app.utils.errors import JobsBotError


class JobQueue:
    def __init__(self):
        try:
            self.queue = Queue()
        except Exception as e:
            raise JobsBotError("JOBS001", {"error": str(e)})

    def add_job(self, job):
        try:
            if not job or not isinstance(job, dict):
                raise JobsBotError("VALID801", {"job": job})
            self.queue.put(job)
        except Exception as e:
            raise JobsBotError("JOBS002", {"job": job, "error": str(e)})

    def get_job(self):
        try:
            if self.queue.empty():
                raise JobsBotError("VALID802", {"error": "Queue is empty"})
            return self.queue.get()
        except Exception as e:
            raise JobsBotError("JOBS001", {"error": str(e)})

    def is_empty(self):
        try:
            return self.queue.empty()
        except Exception as e:
            raise JobsBotError("JOBS001", {"error": str(e)})

    def size(self):
        try:
            return self.queue.qsize()
        except Exception as e:
            raise JobsBotError("JOBS001", {"error": str(e)})
