import pytest
from app.llm.fit_check import should_apply
from app.llm.text_gen import generate_application_text
import asyncio


@pytest.mark.asyncio
async def test_should_apply_validation():
    with pytest.raises(Exception) as e:
        await should_apply(None, "cv")
    assert "VALID101" in str(e.value)

    with pytest.raises(Exception) as e:
        await should_apply("desc", None)
    assert "VALID102" in str(e.value)


@pytest.mark.asyncio
async def test_generate_application_text_validation():
    with pytest.raises(Exception) as e:
        await generate_application_text(None, "cv", [])
    assert "VALID201" in str(e.value)

    with pytest.raises(Exception) as e:
        await generate_application_text("desc", None, [])
    assert "VALID202" in str(e.value)

    with pytest.raises(Exception) as e:
        await generate_application_text("desc", "cv", "notalist")
    assert "VALID203" in str(e.value)
