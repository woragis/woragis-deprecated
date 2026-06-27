import openai
import os
import asyncio
from app.utils.errors import JobsBotError

OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")


def validate_llm_text_input(job_description, cv_data, questions):
    if not job_description or not isinstance(job_description, str):
        raise JobsBotError("VALID201", {"job_description": job_description})
    if not cv_data or not isinstance(cv_data, str):
        raise JobsBotError("VALID202", {"cv_data": cv_data})
    if not isinstance(questions, list):
        raise JobsBotError("VALID203", {"questions": questions})


async def generate_application_text(job_description, cv_data, questions):
    try:
        validate_llm_text_input(job_description, cv_data, questions)
        prompt = f"Job description: {job_description}\nCV: {cv_data}\nQuestions: {questions}\nGenerate cover letter and answers as JSON: {{'cover_letter': string, 'answers': list}}"
        response = await openai.ChatCompletion.acreate(
            model="gpt-3.5-turbo",
            messages=[{"role": "user", "content": prompt}],
            api_key=OPENAI_API_KEY
        )
        content = response.choices[0].message.content
        import json
        result = json.loads(content)
        if "cover_letter" not in result or "answers" not in result:
            raise JobsBotError("VALID204", {"llm_response": result})
        if not isinstance(result["cover_letter"], str):
            raise JobsBotError(
                "VALID205", {"cover_letter": result["cover_letter"]})
        if not isinstance(result["answers"], list):
            raise JobsBotError("VALID206", {"answers": result["answers"]})
        return result
    except Exception as e:
        raise JobsBotError("LLM002", {
                           "job_description": job_description, "questions": questions, "error": str(e)})
