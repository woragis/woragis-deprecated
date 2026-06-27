# Configuration: env vars, limits, feature flags

from app.utils.errors import JobsBotError
import os


class Config:
    def __init__(self):
        self.MAX_APPLICATIONS_PER_DAY = self.get_int(
            "MAX_APPLICATIONS_PER_DAY", 10, min_value=1, code="CONFIG003")
        self.OPENAI_API_KEY = self.get_str(
            "OPENAI_API_KEY", required=True, code="CONFIG004")
        self.DATABASE_URL = self.get_str(
            "DATABASE_URL", default="sqlite:///jobsbot.db", code="CONFIG005")
        self.FEATURE_FLAGS = {}

    def get_str(self, key, default=None, required=False, code=None):
        value = os.getenv(key, default)
        if required and (value is None or value == ""):
            raise JobsBotError(code or "CONFIG001", {
                               "env_var": key, "error": "Missing required environment variable"})
        return value

    def get_int(self, key, default=None, min_value=None, max_value=None, code=None):
        value = os.getenv(key, default)
        try:
            value = int(value)
        except Exception:
            raise JobsBotError(code or "CONFIG002", {
                               "env_var": key, "value": value, "error": "Not an integer"})
        if min_value is not None and value < min_value:
            raise JobsBotError(code or "CONFIG006", {
                               "env_var": key, "value": value, "error": f"Below min {min_value}"})
        if max_value is not None and value > max_value:
            raise JobsBotError(code or "CONFIG007", {
                               "env_var": key, "value": value, "error": f"Above max {max_value}"})
        return value


config = Config()
