import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Stats Club ${run}`;
const admin = `stats-${run}`;
const ids = {};

async function shoot(page, name, points) {
  const base = `/api/tenants/${ids.tenant}`;
  const sh = await api(page, "POST", `${base}/shooters`, { firstname: name, lastname: "Shooter", birthday: "1990-01-01", clubNames: ["Home"] });
  const reg = await api(page, "POST", `${base}/events/${ids.event}/registrations`, {
    shooterId: sh.id, clubId: sh.clubs[0].id, programmes: [{ programmeId: ids.main }],
  });
  if (points) await api(page, "PUT", `${base}/entries/${reg.entries[0].id}/results`, { series: points.map((p) => ({ points: p })) });
}

async function download(page, item, kind) {
  const href = await c(page.locator("[aria-label=Exports]"), "li", item).locator("a", { hasText: kind }).getAttribute("href");
  return page.request.get(href);
}

test.describe.configure({ mode: "serial" });

test.describe("statistics and exports", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 4 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Indoor 2026", startsOn: "2026-11-20" })).id;
    const main = await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Main", countsForClubCompetition: true, showOnRankings: true,
      series: [{ shots: 10, maxPoints: 100 }, { shots: 10, maxPoints: 100 }],
      prices: [{ ageCategory: "elite", price: 20 }],
      awards: [],
    });
    ids.main = main.id;
    for (const [i, p] of [95, 94, 93, 92, 91].entries()) await shoot(page, `S${i}`, [p, p]);
    await page.close();
  });

  test("shows statistics, the club competition and exports", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/statistics`);
    const figures = page.locator("[aria-label='Key figures']");
    await expect(figures).toContainText("100%");
    await expect(figures).toContainText("CHF 100.00");
    const first = page.locator("[aria-label='Club competition'] tbody tr").first();
    await expect(first).toContainText("1");
    await expect(first).toContainText("Home");
    await expect(first).toContainText("186.000");

    const finance = await download(page, "Financial report", "CSV");
    expect(finance.headers()["content-type"]).toContain("text/csv");
    expect(await finance.text()).toContain("Revenue;CHF 100.00");
    expect((await download(page, "Club competition", "PDF")).headers()["content-type"]).toBe("application/pdf");
    expect(await (await download(page, "Ranking Main", "CSV")).text()).toContain("1;Shooter;S0;1990;Elite;Home;95;95;190;0");

    await shoot(page, "Late", null);
    await expect(figures).toContainText("83.3%");
    await expect(figures).toContainText("CHF 120.00");
  });
});
