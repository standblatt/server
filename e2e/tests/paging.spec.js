import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const admin = `paging-${run}`;
const ids = {};
const N = 55;

test.describe.configure({ mode: "serial" });

test.describe("long lists are paged", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, `Paging Club ${run}`, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 10 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Paging", startsOn: "2026-11-21" })).id;
    const p = await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Main", showOnRankings: true, series: [{ shots: 10, maxPoints: 100 }], prices: [{ ageCategory: "elite", price: 10 }], awards: [],
    });
    await Promise.all(Array.from({ length: N }, async (_, i) => {
      const n = String(i).padStart(2, "0");
      const club = await api(page, "POST", `${base}/clubs`, { name: `Club ${n}` });
      const sh = await api(page, "POST", `${base}/shooters`, { firstname: "Pat", lastname: `Page${n}`, birthday: "1990-01-01", clubIds: [club.id] });
      const r = await api(page, "POST", `${base}/events/${ids.event}/registrations`, { shooterId: sh.id, programmes: [{ programmeId: p.id }] });
      await page.request.post(`${base}/registrations/${r.id}/finish`);
    }));
    await page.close();
  });

  test("shooters page through the server", async ({ page }) => {
    await login(page, admin);
    const res = await page.request.get(`/api/tenants/${ids.tenant}/shooters?q=page&limit=20&offset=40`);
    expect(res.headers()["x-total-count"]).toBe(String(N));
    expect((await res.json()).length).toBe(15);

    await page.goto(`/t/${ids.tenant}/shooters`);
    await expect(page.getByText(`1–50 of ${N}`)).toBeVisible();
    await expect(c(page, "a", "Page00 Pat")).toBeVisible();
    await page.getByRole("button", { name: "Next page" }).click();
    await expect(page.getByText(`51–${N} of ${N}`)).toBeVisible();
    await expect(c(page, "a", "Page54 Pat")).toBeVisible();
    await page.getByPlaceholder("License number or name").fill("page5");
    await expect(page.getByText("5 shooters")).toBeVisible();
    await expect(page.getByRole("button", { name: "Next page" })).toHaveCount(0);
  });

  test("clubs and the completed list of the dashboard page in the browser", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/clubs`);
    await expect(page.getByText(`1–50 of ${N}`)).toBeVisible();
    await page.getByRole("button", { name: "Next page" }).click();
    await expect(c(page, "td", "Club 54")).toBeVisible();

    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await expect(page.getByText(`1–50 of ${N}`)).toBeVisible();
    await expect(page.locator("tr", { hasText: "Pat" })).toHaveCount(50);
    await page.getByRole("button", { name: "Next page" }).click();
    await expect(page.locator("tr", { hasText: "Pat" })).toHaveCount(5);
  });
});
