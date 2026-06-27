from app.utils.errors import JobsBotError
from .base import JobPlatform


class GupyPlatform(JobPlatform):
    def login(self, browser, profile):
        try:
            if not browser:
                raise JobsBotError("VALID701", {"browser": browser})
            if not profile:
                raise JobsBotError("VALID702", {"profile": profile})
            # TODO: Implement login flow for Gupy
            pass
        except Exception as e:
            raise JobsBotError(
                "PLATFORM401", {"platform": "Gupy", "error": str(e)})

    def collect_jobs(self, browser):
        try:
            if not browser:
                raise JobsBotError("VALID703", {"browser": browser})
            # TODO: Scrape job URLs from Gupy
            return []
        except Exception as e:
            raise JobsBotError(
                "PLATFORM402", {"platform": "Gupy", "error": str(e)})

    def apply_to_job(self, browser, job_url, application_data):
        try:
            if not browser:
                raise JobsBotError("VALID704", {"browser": browser})
            if not job_url or not isinstance(job_url, str):
                raise JobsBotError("VALID705", {"job_url": job_url})
            if not application_data or not isinstance(application_data, dict):
                raise JobsBotError(
                    "VALID706", {"application_data": application_data})
            # TODO: Fill and submit application form
            pass
        except Exception as e:
            raise JobsBotError(
                "PLATFORM403", {"platform": "Gupy", "job_url": job_url, "error": str(e)})
