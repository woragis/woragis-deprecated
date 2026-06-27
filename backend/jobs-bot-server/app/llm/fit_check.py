import openai
import os
import asyncio
from app.utils.errors import JobsBotError

OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")


def validate_llm_input(job_description, cv_data):
    if not job_description or not isinstance(job_description, str):
        raise JobsBotError("VALID101", {"job_description": job_description})
    if not cv_data or not isinstance(cv_data, str):
        raise JobsBotError("VALID102", {"cv_data": cv_data})


async def should_apply(job_description, cv_data):
    try:
        validate_llm_input(job_description, cv_data)
        prompt = f"Job description: {job_description}\nCV: {cv_data}\nShould I apply? Respond with JSON: {{'apply': true/false, 'score': float, 'reason': string}}"
        response = await openai.ChatCompletion.acreate(
            model="gpt-3.5-turbo",
            messages=[{"role": "user", "content": prompt}],
            api_key=OPENAI_API_KEY
        )
        content = response.choices[0].message.content
        import json
        result = json.loads(content)
        if "apply" not in result or "score" not in result or "reason" not in result:
            raise JobsBotError("VALID103", {"llm_response": result})
        if not isinstance(result["apply"], bool):
            raise JobsBotError("VALID104", {"apply": result["apply"]})
        if not (0 <= result["score"] <= 1):
            raise JobsBotError("VALID105", {"score": result["score"]})
        if not isinstance(result["reason"], str):
            raise JobsBotError("VALID106", {"reason": result["reason"]})
        return result
    except Exception as e:
        raise JobsBotError(
            "LLM001", {"job_description": job_description, "error": str(e)})
