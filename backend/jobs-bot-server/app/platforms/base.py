from abc import ABC, abstractmethod
from app.utils.errors import JobsBotError


class JobPlatform(ABC):
    @abstractmethod
    def login(self, browser, profile):
        try:
            if not browser:
                raise JobsBotError("VALID301", {"browser": browser})
            if not profile:
                raise JobsBotError("VALID302", {"profile": profile})
            pass
        except Exception as e:
            raise JobsBotError(
                "PLATFORM001", {"platform": self.__class__.__name__, "error": str(e)})

    @abstractmethod
    def collect_jobs(self, browser):
        try:
            if not browser:
                raise JobsBotError("VALID303", {"browser": browser})
            pass
        except Exception as e:
            raise JobsBotError(
                "PLATFORM002", {"platform": self.__class__.__name__, "error": str(e)})

    @abstractmethod
    def apply_to_job(self, browser, job_url, application_data):
        try:
            if not browser:
                raise JobsBotError("VALID304", {"browser": browser})
            if not job_url or not isinstance(job_url, str):
                raise JobsBotError("VALID305", {"job_url": job_url})
            if not application_data or not isinstance(application_data, dict):
                raise JobsBotError(
                    "VALID306", {"application_data": application_data})
            pass
        except Exception as e:
            raise JobsBotError("PLATFORM003", {
                               "platform": self.__class__.__name__, "job_url": job_url, "error": str(e)})
