from app.utils.errors import JobsBotError
from .base import JobPlatform


class GreenhousePlatform(JobPlatform):
    def login(self, browser, profile):
        try:
            if not browser:
                raise JobsBotError("VALID401", {"browser": browser})
            if not profile:
                raise JobsBotError("VALID402", {"profile": profile})
            # TODO: Implement login flow for Greenhouse
            pass
        except Exception as e:
            raise JobsBotError(
                "PLATFORM101", {"platform": "Greenhouse", "error": str(e)})

    def collect_jobs(self, browser):
        try:
            if not browser:
                raise JobsBotError("VALID403", {"browser": browser})
            # TODO: Scrape job URLs from Greenhouse
            return []
        except Exception as e:
            raise JobsBotError(
                "PLATFORM102", {"platform": "Greenhouse", "error": str(e)})

    def apply_to_job(self, browser, job_url, application_data):
        try:
            if not browser:
                raise JobsBotError("VALID404", {"browser": browser})
            if not job_url or not isinstance(job_url, str):
                raise JobsBotError("VALID405", {"job_url": job_url})
            if not application_data or not isinstance(application_data, dict):
                raise JobsBotError(
                    "VALID406", {"application_data": application_data})
            # TODO: Fill and submit application form
            pass
        except Exception as e:
            raise JobsBotError(
                "PLATFORM103", {"platform": "Greenhouse", "job_url": job_url, "error": str(e)})
