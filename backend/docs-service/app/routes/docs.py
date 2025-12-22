"""
Routes for serving documentation.
"""
import os
from pathlib import Path
from typing import Optional
from fastapi import APIRouter, HTTPException, Query
from fastapi.responses import JSONResponse, HTMLResponse
from pydantic import BaseModel
import markdown
from pygments import highlight
from pygments.lexers import get_lexer_by_name
from pygments.formatters import HtmlFormatter
from pygments.util import ClassNotFound

from app.config import settings
from app.logger import get_logger

logger = get_logger()
router = APIRouter(prefix="/api/v1/docs", tags=["docs"])

# Configure markdown extensions
MARKDOWN_EXTENSIONS = settings.MARKDOWN_EXTENSIONS.split(",")
MARKDOWN_EXTENSIONS = [ext.strip() for ext in MARKDOWN_EXTENSIONS if ext.strip()]

# Configure code highlighting
# Use 'default' style which is always available in pygments
formatter = HtmlFormatter(style="default", cssclass="codehilite")


class DocResponse(BaseModel):
    """Response model for a documentation file."""
    path: str
    title: str
    content: str
    html: str
    metadata: Optional[dict] = None


class DocListResponse(BaseModel):
    """Response model for listing documentation files."""
    files: list[dict]
    total: int


def get_docs_root() -> Path:
    """Get the docs root directory."""
    return Path(settings.DOCS_ROOT)


def find_doc_file(path: str) -> Optional[Path]:
    """
    Find a documentation file by path.
    
    Args:
        path: Relative path to the doc file (e.g., "architecture/system-overview.md")
        
    Returns:
        Path object if found, None otherwise
    """
    docs_root = get_docs_root()
    
    # Normalize path
    path = path.strip("/")
    if not path:
        return None
    
    # Try with .md extension
    doc_path = docs_root / path
    if doc_path.exists() and doc_path.is_file():
        return doc_path
    
    # Try with .md extension if not provided
    if not path.endswith((".md", ".markdown")):
        doc_path = docs_root / f"{path}.md"
        if doc_path.exists() and doc_path.is_file():
            return doc_path
    
    # Try as directory with index.md
    if doc_path.is_dir():
        index_path = doc_path / "README.md"
        if index_path.exists():
            return index_path
        index_path = doc_path / "index.md"
        if index_path.exists():
            return index_path
    
    return None


def parse_markdown(content: str) -> tuple[str, dict]:
    """
    Parse markdown content and extract metadata.
    
    Args:
        content: Markdown content
        
    Returns:
        Tuple of (html_content, metadata)
    """
    # Try to extract frontmatter if present
    metadata = {}
    
    if content.startswith("---"):
        parts = content.split("---", 2)
        if len(parts) >= 3:
            # Parse frontmatter (simple YAML-like parsing)
            frontmatter = parts[1].strip()
            content = parts[2].strip()
            
            # Simple key-value parsing
            for line in frontmatter.split("\n"):
                if ":" in line:
                    key, value = line.split(":", 1)
                    metadata[key.strip()] = value.strip().strip('"').strip("'")
    
    # Convert markdown to HTML
    md = markdown.Markdown(extensions=MARKDOWN_EXTENSIONS)
    html = md.convert(content)
    
    # Add CSS for code highlighting
    css = formatter.get_style_defs()
    html = f'<style>{css}</style>\n{html}'
    
    return html, metadata


def extract_title(content: str) -> str:
    """Extract title from markdown content (first H1 or filename)."""
    lines = content.split("\n")
    for line in lines:
        if line.startswith("# "):
            return line[2:].strip()
    return "Documentation"


@router.get("/", response_model=DocListResponse)
async def list_docs(
    category: Optional[str] = Query(None, description="Filter by category (e.g., 'architecture', 'adr', 'runbooks')")
):
    """
    List all available documentation files.
    
    Args:
        category: Optional category filter (subdirectory name)
        
    Returns:
        List of documentation files with metadata
    """
    docs_root = get_docs_root()
    
    if not docs_root.exists():
        raise HTTPException(status_code=500, detail="Docs directory not found")
    
    files = []
    
    # Find all markdown files
    pattern = "**/*.md"
    md_files = list(docs_root.glob(pattern)) + list(docs_root.glob("**/*.markdown"))
    
    for md_file in md_files:
        # Skip if in hidden directories
        if any(part.startswith(".") for part in md_file.parts):
            continue
        
        # Filter by category if specified
        if category:
            if category not in md_file.parts:
                continue
        
        # Get relative path
        rel_path = md_file.relative_to(docs_root)
        path_str = str(rel_path).replace("\\", "/")
        
        # Read first few lines to get title
        try:
            with open(md_file, "r", encoding="utf-8") as f:
                content = f.read()
                title = extract_title(content)
        except Exception as e:
            logger.warning("failed_to_read_doc", path=path_str, error=str(e))
            title = md_file.stem
        
        files.append({
            "path": path_str,
            "title": title,
            "size": md_file.stat().st_size,
            "category": rel_path.parts[0] if len(rel_path.parts) > 1 else "root",
        })
    
    # Sort by path
    files.sort(key=lambda x: x["path"])
    
    logger.info("listed_docs", count=len(files), category=category)
    
    return DocListResponse(files=files, total=len(files))


@router.get("/{path:path}", response_model=DocResponse)
async def get_doc(
    path: str,
    format: Optional[str] = Query("json", description="Response format: 'json' or 'html'")
):
    """
    Get a documentation file by path.
    
    Args:
        path: Relative path to the doc file (e.g., "architecture/system-overview.md")
        format: Response format - 'json' returns structured data, 'html' returns rendered HTML
        
    Returns:
        Documentation content in requested format
    """
    doc_file = find_doc_file(path)
    
    if not doc_file:
        raise HTTPException(status_code=404, detail=f"Documentation file not found: {path}")
    
    try:
        with open(doc_file, "r", encoding="utf-8") as f:
            content = f.read()
    except Exception as e:
        logger.exception("failed_to_read_doc_file", path=path, error=str(e))
        raise HTTPException(status_code=500, detail=f"Failed to read documentation file: {str(e)}")
    
    # Parse markdown
    html, metadata = parse_markdown(content)
    title = extract_title(content)
    
    # Get relative path
    docs_root = get_docs_root()
    rel_path = doc_file.relative_to(docs_root)
    path_str = str(rel_path).replace("\\", "/")
    
    logger.info("served_doc", path=path_str, format=format)
    
    response_data = DocResponse(
        path=path_str,
        title=title,
        content=content,
        html=html,
        metadata=metadata if metadata else None,
    )
    
    if format == "html":
        # Return full HTML page
        html_page = f"""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{title} - Woragis Docs</title>
    <style>
        body {{
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 900px;
            margin: 0 auto;
            padding: 2rem;
            line-height: 1.6;
            color: #333;
        }}
        h1, h2, h3, h4, h5, h6 {{
            margin-top: 2rem;
            margin-bottom: 1rem;
        }}
        code {{
            background: #f4f4f4;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-size: 0.9em;
        }}
        pre {{
            background: #f4f4f4;
            padding: 1rem;
            border-radius: 5px;
            overflow-x: auto;
        }}
        a {{
            color: #0066cc;
            text-decoration: none;
        }}
        a:hover {{
            text-decoration: underline;
        }}
        table {{
            border-collapse: collapse;
            width: 100%;
            margin: 1rem 0;
        }}
        th, td {{
            border: 1px solid #ddd;
            padding: 0.5rem;
            text-align: left;
        }}
        th {{
            background: #f4f4f4;
            font-weight: 600;
        }}
    </style>
</head>
<body>
    {html}
</body>
</html>"""
        return HTMLResponse(content=html_page)
    
    return response_data
