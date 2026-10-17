import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Multi Club ${run}`;
const admin = `multi-${run}`;
const ids = {};

function programme(name, series, categories) {
  return {
    name,
    showOnRankings: true,
    series: Array.from({ length: series }, () => ({ shots: 10, maxPoints: 100 })),
    prices: categories.map((ageCategory) => ({ ageCategory, price: 10 })),
    awards: [],
  };
}

test.describe.configure({ mode: "serial" });

test.describe("events on several ranges and days", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const r300 = await api(page, "POST", `${base}/ranges`, { name: "300m", laneCount: 4 });
    const r50 = await api(page, "POST", `${base}/ranges`, { name: "50m", laneCount: 2 });
    ids.r300 = r300.id;
    ids.r50 = r50.id;
    ids.event = (await api(page, "POST", `${base}/events`, {
      rangeIds: [r300.id, r50.id], name: "Fest 2026", startsOn: "2026-11-21", schedule: { days: [{ date: "2026-11-18" }, { date: "2026-11-22" }] },
    })).id;
    const ep = `${base}/events/${ids.event}/programmes`;
    await api(page, "POST", ep, programme("Rifle", 2, ["elite"]));
    const pistol = await api(page, "POST", ep, programme("Pistol", 2, ["elite"]));
    await api(page, "PUT", `${base}/programmes/${pistol.id}/range`, { rangeId: r50.id });
    await api(page, "POST", `${base}/shooters`, { licenseNumber: "300001", firstname: "Anna", lastname: "Adult", birthday: "1990-01-01" });
    await page.close();
  });

  test("shows the days, groups lanes by range and follows the programme's range", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await expect(c(page, "span", "18.11.2026–22.11.2026 · 300m, 50m")).toBeVisible();
    await expect(c(page, "h2", "300m")).toBeVisible();
    await expect(c(page, "h2", "50m")).toBeVisible();

    await page.keyboard.press("F4");
    await page.keyboard.type("adult");
    await c(page, "button", /^Register$/).click();
    await expect(page.locator("[role=group][aria-label='300m']")).toBeVisible();
    await expect(page.locator("[role=group][aria-label='50m']")).toBeVisible();
    await c(page, "label", "Pistol").locator("[role=checkbox]").click();
    await expect(page.locator("[role=group][aria-label='300m']")).toHaveCount(0);
    await expect(page.getByText("The chosen programmes are shot on 50m.")).toBeVisible();
    await c(page.locator("[role=group][aria-label='50m']"), "button", /^2$/).click();
    await c(page, "button", /^Register$/).click();
    await expect(c(page, "[data-sonner-toast]", "Anna Adult registered on 50m, lane 2")).toBeVisible();
    await expect(page.locator("[aria-label='Lane 2: Adult Anna']")).toBeVisible();
  });

  test("refuses to drop a range in use", async ({ page }) => {
    await login(page, admin);
    const ev = await api(page, "GET", `/api/tenants/${ids.tenant}/events/${ids.event}`);
    const res = await page.request.put(`/api/tenants/${ids.tenant}/events/${ids.event}`, {
      data: { rangeIds: [ids.r300], name: ev.name, startsOn: ev.startsOn },
    });
    expect(res.status()).toBe(409);
  });

  test("a kiosk link limited to one range shows only its lanes and programmes", async ({ page }) => {
    await login(page, admin);
    const base = `/api/tenants/${ids.tenant}`;
    await page.request.put(`${base}/events/${ids.event}`, {
      data: { rangeIds: [ids.r300, ids.r50], name: "Fest 2026", startsOn: "2026-11-21", schedule: { days: [{ date: "2026-11-18" }, { date: "2026-11-22" }] }, settings: { showLaneStatus: true } },
    });
    const token = await api(page, "POST", `${base}/events/${ids.event}/kiosk-tokens`, { name: "50m hall" });
    await page.goto(`/kiosk/${token.secret}?range=${ids.r50}`);
    await expect(c(page, "h1", "Fest 2026")).toBeVisible();
    await expect(c(page, "p", "18.11.2026–22.11.2026 · 50m")).toBeVisible();
    await expect(page.locator("[aria-label='Lane 2: Adult Anna']")).toBeVisible();
    await expect(page.locator("[aria-label='Lane 4: free']")).toHaveCount(0);
    await expect(c(page, "h2", "Pistol")).toBeVisible();
    await expect(c(page, "h2", "Rifle")).toBeVisible();
  });
});
