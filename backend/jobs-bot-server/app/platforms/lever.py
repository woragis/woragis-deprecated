from app.utils.errors import JobsBotError
from .base import JobPlatform


class LeverPlatform(JobPlatform):
    def login(self, browser, profile):
        try:
            if not browser:
                raise JobsBotError("VALID501", {"browser": browser})
            if not profile:
                raise JobsBotError("VALID502", {"profile": profile})
            # TODO: Implement login flow for Lever
            pass
        except Exception as e:
            raise JobsBotError(
                "PLATFORM201", {"platform": "Lever", "error": str(e)})

    def collect_jobs(self, browser):
        try:
            if not browser:
                raise JobsBotError("VALID503", {"browser": browser})
            # TODO: Scrape job URLs from Lever
            return []
        except Exception as e:
            raise JobsBotError(
                "PLATFORM202", {"platform": "Lever", "error": str(e)})

    def apply_to_job(self, browser, job_url, application_data):
        try:
            if not browser:
                raise JobsBotError("VALID504", {"browser": browser})
            if not job_url or not isinstance(job_url, str):
                raise JobsBotError("VALID505", {"job_url": job_url})
            if not application_data or not isinstance(application_data, dict):
                raise JobsBotError(
                    "VALID506", {"application_data": application_data})
            # TODO: Fill and submit application form
            pass
        except Exception as e:
            raise JobsBotError(
                "PLATFORM203", {"platform": "Lever", "job_url": job_url, "error": str(e)})
