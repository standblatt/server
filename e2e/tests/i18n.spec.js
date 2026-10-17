import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `I18n Club ${run}`;
const admin = `i18n-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("languages", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 2 });
    const ev = await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Tir 2026", startsOn: "2026-11-20", locale: "fr" });
    ids.kiosk = (await api(page, "POST", `${base}/events/${ev.id}/kiosk-tokens`, { name: "Hall" })).secret;
    await page.close();
  });

  test("switches the UI language per user and keeps it", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}`);
    await expect(c(page, "h1", "Events")).toBeVisible();
    await expect(c(page, "[data-sidebar=footer] p", /^Standblatt \S+$/)).toBeVisible();
    await page.locator("button[aria-label='User menu']").click();
    await c(page, "[role=menuitem]", "Language").click();
    await c(page, "[role=menuitemradio]", "Deutsch").click();
    await expect(c(page, "h1", "Anlässe")).toBeVisible();
    await expect(c(page, "[data-sidebar=menu-button]", "Schützen")).toBeVisible();

    expect((await api(page, "GET", "/api/me")).user.locale).toBe("de");
    await page.evaluate(() => localStorage.clear());
    await page.reload();
    await expect(c(page, "h1", "Anlässe")).toBeVisible();
  });

  test("shows the kiosk in the language of the event", async ({ page }) => {
    await page.goto(`/kiosk/${ids.kiosk}`);
    await expect(c(page, "h1", "Tir 2026")).toBeVisible();
    await expect(page).toHaveTitle("Tir 2026 – Classements");

    await page.goto(`/kiosk/${ids.kiosk}?language=de`);
    await expect(page).toHaveTitle("Tir 2026 – Ranglisten");
  });
});
