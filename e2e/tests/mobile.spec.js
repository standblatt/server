import { test, expect } from "@playwright/test";
import { c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const user = `mobile-${run}`;
const ids = {};
const statusBar = 47;

test.describe.configure({ mode: "serial" });
test.use({ viewport: { width: 390, height: 844 }, isMobile: true, hasTouch: true });

/** Top edge of a locator in CSS pixels. */
const top = async (locator) => (await locator.boundingBox()).y;

test.describe("phone", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, `Mobile Club ${run}`, user);
    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem("standblatt-locale", "en"));
  });

  test("closes the sidebar after choosing a page", async ({ page }) => {
    await login(page, user);
    await page.goto(`/t/${ids.tenant}`);
    const sheet = page.locator("[data-mobile=true]");
    await page.locator("[data-slot=sidebar-trigger]").click();
    await expect(sheet).toBeVisible();
    await c(sheet, "a", "Shooters").click();
    await expect(page).toHaveURL((u) => u.pathname.endsWith("/shooters"));
    await expect(sheet).toHaveCount(0);
  });

  test("keeps content below the status bar of an installed iPhone app", async ({ page, context }) => {
    const cdp = await context.newCDPSession(page);
    await cdp.send("Emulation.setSafeAreaInsetsOverride", { insets: { top: statusBar } });
    await page.goto("/login");
    expect(await top(page.locator("input[name=name]"))).toBeGreaterThanOrEqual(statusBar);
    const strip = await page.evaluate(() => {
      const s = getComputedStyle(document.body, "::before");
      return { height: s.height, position: s.position, zIndex: s.zIndex };
    });
    expect(strip).toEqual({ height: `${statusBar}px`, position: "fixed", zIndex: "100" });

    await login(page, user);
    await page.goto(`/t/${ids.tenant}/shooters`);
    const trigger = page.locator("[data-slot=sidebar-trigger]");
    expect(await top(trigger)).toBeGreaterThanOrEqual(statusBar);
    await trigger.click();
    await expect.poll(() => top(page.locator("[data-mobile=true]"))).toBe(statusBar);
    await page.keyboard.press("Escape");

    await page.goto("/my");
    expect(await top(page.locator("header").first())).toBe(0);
    expect(await top(page.locator("header").first().locator("a, button").first())).toBeGreaterThanOrEqual(statusBar);
    await page.goto("/my/train");
    expect(await top(page.getByRole("button", { name: "Close" }))).toBeGreaterThanOrEqual(statusBar);
    await page.goto("/account");
    expect(await top(page.locator("main").locator("h1, h2").first())).toBeGreaterThanOrEqual(statusBar);
  });
});
