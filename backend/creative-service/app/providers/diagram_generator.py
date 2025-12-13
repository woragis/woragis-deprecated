import base64
import subprocess
import tempfile
import os
from typing import Literal, Optional
from openai import AsyncOpenAI
import anthropic
import httpx
from app.config import settings


class DiagramGenerator:
    """
    Generate technical diagrams using AI to create diagram code (Mermaid/PlantUML/Graphviz)
    and then render them to images.
    """
    
    def __init__(self, provider: Literal["openai", "anthropic"] = "openai"):
        self.provider = provider
        if provider == "openai":
            self.client = AsyncOpenAI(api_key=settings.OPENAI_API_KEY)
        else:
            self.client = anthropic.AsyncAnthropic(api_key=settings.ANTHROPIC_API_KEY)
    
    async def _generate_diagram_code(
        self,
        description: str,
        diagram_type: Literal["mermaid", "plantuml", "graphviz"] = "mermaid",
        diagram_kind: str = "flowchart",  # flowchart, sequence, er, etc.
    ) -> str:
        """Use AI to generate diagram code from description."""
        
        system_prompt = f"""You are an expert at creating {diagram_type} diagrams for technical documentation.
Generate ONLY the {diagram_type} code, nothing else - no markdown, no explanations, just the raw code.
Make it clear, professional, and suitable for technical blog posts about software architecture, microservices, and workflows.

Diagram type: {diagram_kind}
Format: {diagram_type}"""
        
        user_prompt = f"""Create a {diagram_type} {diagram_kind} diagram for the following:
{description}

Return ONLY the {diagram_type} code, no markdown formatting, no explanations."""
        
        if self.provider == "openai":
            response = await self.client.chat.completions.create(
                model=settings.OPENAI_MODEL,
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_prompt},
                ],
                temperature=0.3,
            )
            code = response.choices[0].message.content.strip()
        else:
            response = await self.client.messages.create(
                model=settings.ANTHROPIC_MODEL,
                max_tokens=2000,
                system=system_prompt,
                messages=[
                    {"role": "user", "content": user_prompt},
                ],
            )
            code = response.content[0].text.strip()
        
        # Clean up if wrapped in code blocks
        if code.startswith("```"):
            lines = code.split("\n")
            # Remove first line (```mermaid or ```)
            if lines[0].startswith("```"):
                lines = lines[1:]
            # Remove last line (```)
            if lines and lines[-1].strip() == "```":
                lines = lines[:-1]
            code = "\n".join(lines)
        
        return code
    
    async def generate_mermaid(
        self,
        description: str,
        diagram_kind: str = "flowchart",
        output_format: str = "png",
        width: int = 1200,
        height: int = 800,
    ) -> dict:
        """
        Generate a Mermaid diagram from description.
        
        Returns dict with 'b64_json' (base64 image) and 'code' (mermaid code).
        """
        # Generate Mermaid code
        code = await self._generate_diagram_code(description, "mermaid", diagram_kind)
        
        # Render using Mermaid CLI or API
        # Option 1: Use mermaid.ink API (simpler, no local install needed)
        try:
            import urllib.parse
            encoded = urllib.parse.quote(code)
            api_url = f"https://mermaid.ink/img/{encoded}"
            
            async with httpx.AsyncClient() as client:
                response = await client.get(api_url, params={"type": output_format})
                if response.status_code == 200:
                    b64_data = base64.b64encode(response.content).decode('utf-8')
                    return {
                        "b64_json": b64_data,
                        "code": code,
                        "format": output_format,
                    }
        except Exception as e:
            # Fallback: Try local mermaid-cli if available
            pass
        
        # Option 2: Try local mermaid-cli (requires npm install -g @mermaid-js/mermaid-cli)
        try:
            with tempfile.NamedTemporaryFile(mode='w', suffix='.mmd', delete=False) as f:
                f.write(code)
                temp_mmd = f.name
            
            temp_output = temp_mmd.replace('.mmd', f'.{output_format}')
            
            subprocess.run(
                ["mmdc", "-i", temp_mmd, "-o", temp_output, "-w", str(width), "-H", str(height)],
                check=True,
                capture_output=True,
            )
            
            with open(temp_output, 'rb') as f:
                b64_data = base64.b64encode(f.read()).decode('utf-8')
            
            os.unlink(temp_mmd)
            os.unlink(temp_output)
            
            return {
                "b64_json": b64_data,
                "code": code,
                "format": output_format,
            }
        except Exception:
            raise Exception("Could not render Mermaid diagram. Install mermaid-cli: npm install -g @mermaid-js/mermaid-cli")
    
    async def generate_graphviz(
        self,
        description: str,
        engine: str = "dot",  # dot, neato, fdp, etc.
        output_format: str = "png",
    ) -> dict:
        """
        Generate a Graphviz diagram from description.
        
        Returns dict with 'b64_json' and 'code' (DOT code).
        """
        code = await self._generate_diagram_code(description, "graphviz", engine)
        
        try:
            import graphviz
            
            # Create graph from DOT code
            graph = graphviz.Source(code, engine=engine, format=output_format)
            
            # Render to bytes
            output_bytes = graph.pipe(format=output_format)
            b64_data = base64.b64encode(output_bytes).decode('utf-8')
            
            return {
                "b64_json": b64_data,
                "code": code,
                "format": output_format,
            }
        except ImportError:
            raise Exception("graphviz package required. Install: pip install graphviz (and system graphviz)")
        except Exception as e:
            raise Exception(f"Graphviz rendering error: {str(e)}")
    
    async def generate(
        self,
        description: str,
        diagram_type: Literal["mermaid", "graphviz"] = "mermaid",
        diagram_kind: str = "flowchart",
        output_format: str = "png",
        **kwargs
    ) -> dict:
        """Main entry point for diagram generation."""
        if diagram_type == "mermaid":
            return await self.generate_mermaid(description, diagram_kind, output_format, **kwargs)
        elif diagram_type == "graphviz":
            return await self.generate_graphviz(description, diagram_kind, output_format, **kwargs)
        else:
            raise ValueError(f"Unsupported diagram type: {diagram_type}")

