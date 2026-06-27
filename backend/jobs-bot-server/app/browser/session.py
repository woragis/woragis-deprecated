from playwright.sync_api import sync_playwright
from app.utils.errors import JobsBotError


class BrowserSession:
    def __init__(self, profile_path=None):
        self.profile_path = profile_path
        self.playwright = None
        self.browser = None
        self.context = None

    def start(self):
        try:
            self.playwright = sync_playwright().start()
            self.browser = self.playwright.chromium.launch(headless=False)
            self.context = self.browser.new_context(
                storage_state=self.profile_path)
            return self.context
        except Exception as e:
            raise JobsBotError("BROWSER001", {"error": str(e)})

    def stop(self):
        try:
            if self.context:
                self.context.close()
            if self.browser:
                self.browser.close()
            if self.playwright:
                self.playwright.stop()
        except Exception as e:
            raise JobsBotError("BROWSER001", {"error": str(e)})

    def verify_login(self, url, login_check_selector):
        try:
            page = self.context.new_page()
            page.goto(url)
            return page.is_visible(login_check_selector)
        except Exception as e:
            raise JobsBotError(
                "BROWSER002", {"url": url, "selector": login_check_selector, "error": str(e)})
