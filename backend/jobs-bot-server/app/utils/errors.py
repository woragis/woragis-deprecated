ERRORS = {
    "API001": "Invalid run request",
    "API002": "Run not found",
    "API003": "Run already in progress",
    "JOBS001": "Job queue empty",
    "JOBS002": "Job application failed",
    "LLM001": "LLM fit check failed",
    "LLM002": "LLM response parsing error",
    "BROWSER001": "Browser session crashed",
    "BROWSER002": "Login verification failed",
    "DB001": "Database connection error",
    "DB002": "Database write failed",
    # Add more codes as needed
}


class JobsBotError(Exception):
    def __init__(self, code, context=None):
        self.code = code
        self.message = ERRORS.get(code, "Unknown error")
        self.context = context
        super().__init__(f"{code}: {self.message} | {context}")
