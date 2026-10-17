import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Update Club ${run}`;
const admin = `update-${run}`;
const ids = {};

const newVersion = (page) =>
  page.route("**/_app/version.json", (route) => route.fulfill({ json: { version: `next-${run}` } }));

test.describe.configure({ mode: "serial" });

test.describe("new release", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 4 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Indoor 2026", startsOn: "2026-11-20" })).id;
    ids.kiosk = (await api(page, "POST", `${base}/events/${ids.event}/kiosk-tokens`, { name: "Hall" })).secret;
    await page.close();
  });

  test("offers a reload and reloads on the next navigation", async ({ page }) => {
    await page.clock.install();
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await expect(c(page, "header", "Indoor 2026")).toBeVisible();
    await page.evaluate(() => (window.marker = true));
    await newVersion(page);
    await page.clock.fastForward(61_000);
    const toast = c(page, "[data-sonner-toast]", "A new version of Standblatt is available.");
    await expect(toast).toBeVisible();
    await expect(c(toast, "button", "Reload")).toBeVisible();
    await c(page, "[data-sidebar=sidebar] a", "Shooters").click();
    await expect(page).toHaveURL((u) => u.pathname.endsWith("/shooters"));
    expect(await page.evaluate(() => window.marker)).toBeUndefined();
  });

  test("reloads a kiosk display by itself", async ({ page }) => {
    await page.clock.install();
    await page.goto(`/kiosk/${ids.kiosk}`);
    await expect(c(page, "h1", "Indoor 2026")).toBeVisible();
    await page.evaluate(() => (window.marker = true));
    await newVersion(page);
    await page.clock.fastForward(61_000);
    await expect.poll(() => page.evaluate(() => window.marker).catch(() => "reloading")).toBeUndefined();
  });
});
