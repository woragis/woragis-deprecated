import os
from dataclasses import dataclass, field
from dotenv import load_dotenv


# Load variables from .env if present
load_dotenv()


@dataclass(frozen=True)
class Settings:
    # General
    CORS_ENABLED: bool = (os.getenv("CORS_ENABLED", "true").lower() == "true")
    CORS_ALLOWED_ORIGINS: str = os.getenv("CORS_ALLOWED_ORIGINS", "*")
    
    # Docs configuration
    DOCS_ROOT: str = os.getenv("DOCS_ROOT", "/app/docs")
    DOCS_EXTENSIONS: list[str] = field(default_factory=lambda: [".md", ".markdown"])
    
    # Markdown extensions
    MARKDOWN_EXTENSIONS: str = os.getenv(
        "MARKDOWN_EXTENSIONS",
        "fenced_code,codehilite,tables,toc,extra"
    )


settings = Settings()
