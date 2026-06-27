import random
import time
from app.utils.errors import JobsBotError

# Human-like delays


def human_delay(min_ms=300, max_ms=1200):
    try:
        time.sleep(random.uniform(min_ms, max_ms) / 1000)
    except Exception as e:
        raise JobsBotError("BROWSER001", {"error": str(e)})

# Scroll to element


def scroll_to_element(page, selector):
    try:
        page.eval_on_selector(
            selector, "element => element.scrollIntoView({behavior: 'smooth'})")
        human_delay()
    except Exception as e:
        raise JobsBotError(
            "BROWSER001", {"selector": selector, "error": str(e)})

# Simulate mouse movement


def move_mouse(page, x, y):
    try:
        page.mouse.move(x, y)
        human_delay()
    except Exception as e:
        raise JobsBotError("BROWSER001", {"x": x, "y": y, "error": str(e)})

# Randomize typing speed


def human_type(page, selector, text):
    try:
        for char in text:
            page.fill(selector, char)
            human_delay(50, 200)
    except Exception as e:
        raise JobsBotError(
            "BROWSER001", {"selector": selector, "text": text, "error": str(e)})
