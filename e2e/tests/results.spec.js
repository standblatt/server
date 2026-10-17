import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Results Club ${run}`;
const admin = `res-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("results and awards", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 4 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Indoor 2026", startsOn: "2026-11-20" })).id;
    const ep = `${base}/events/${ids.event}/programmes`;
    const main = await api(page, "POST", ep, {
      name: "Main",
      series: [
        { shots: 10, maxPoints: 100 },
        { shots: 10, maxPoints: 100, trackMouches: true },
        { shots: 10, maxPoints: 100 },
        { shots: 10, maxPoints: 100 },
      ],
      prices: [{ ageCategory: "elite", price: 20 }],
      awards: [
        { name: "Badge", kind: "natural", value: 0, limits: [{ ageCategory: "elite", points: 300 }] },
        { name: "Wreath", kind: "card", value: 25, limits: [{ ageCategory: "elite", points: 350 }] },
      ],
    });
    const section = await api(page, "POST", ep, {
      name: "Section",
      series: [{ shots: 10, maxPoints: 100 }, { shots: 10, maxPoints: 100 }],
      prices: [{ ageCategory: "elite", price: 5 }],
      awards: [],
    });
    await api(page, "PUT", `${base}/programmes/${section.id}/combination`, { combinableWithId: main.id });
    const anna = await api(page, "POST", `${base}/shooters`, { firstname: "Anna", lastname: "Adult", birthday: "1990-01-01" });
    await api(page, "POST", `${base}/events/${ids.event}/registrations`, {
      shooterId: anna.id,
      lane: 1,
      programmes: [{ programmeId: main.id }, { programmeId: section.id, combinedSeries: [2, 4] }],
    });
    await page.close();
  });

  test("enters results by keyboard and chooses an award", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await c(c(page, "tr", "Adult Anna"), "a", "Results").click();
    await expect(c(page, "section[aria-label=Section]", "Waiting for the results of Main")).toBeVisible();

    const main = page.locator("section[aria-label=Main]");
    const focused = page.locator(":focus");
    await main.locator("input[aria-label='Series 1']").pressSequentially("95");
    for (const [label, keys] of [["Series 2", "88"], ["Mouches series 2", "3"], ["Series 3", "10"], ["Series 3", "0"], ["Series 4", "97"]]) {
      await expect(focused).toHaveAttribute("aria-label", label);
      await page.keyboard.type(keys);
    }
    await expect(focused).toHaveText("Save");
    await focused.click();
    await expect(main.getByText("Total 380 / 400").first()).toBeVisible();
    await expect(c(main, "label", "Wreath (CHF 25.00)").locator("input")).toBeChecked();
    await expect(c(main, "legend", "Award: Wreath")).toBeVisible();
    await c(main, "label", "Badge").locator("input").check();
    await expect(c(main, "legend", "Award: Badge")).toBeVisible();
    await c(main, "label", "Wreath (CHF 25.00)").locator("input").check();
    await expect(c(main, "legend", "Award: Wreath")).toBeVisible();
    const section = page.locator("section[aria-label=Section]");
    await expect(section).toContainText("S1: 88 · S2: 97");
    await expect(section).toContainText("Total 185");
    await expect(c(page, "[role=status]", "All programmes have results")).toBeVisible();

    await c(page, "[role=dialog] button", "Close").click();
    await expect(page.locator("[role=dialog]")).toHaveCount(0);
    await expect(c(page, "h2", "Completed (1)")).toBeVisible();
    await expect(c(page, "tr", "Adult Anna")).toContainText("complete");
    await expect(page.locator("[aria-label='Lane 1: free']")).toBeVisible();
  });

  test("hands the awards out and takes it back", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    const row = c(page, "tr", "Adult Anna");
    await c(row, "button", "awards to hand out").click();
    await c(page, "[role=dialog] button", "Awards handed out").click();
    await expect(c(page, "[role=dialog] span", "awards handed out")).toBeVisible();
    await c(page, "[role=dialog] button", "Close").click();
    await expect(c(row, "span", "awards handed out")).toBeVisible();
    await c(row, "button", "awards handed out").click();
    await c(c(page, "[role=alertdialog]", "Take the awards handout back?"), "button", "Take back").click();
    await expect(c(row, "span", "awards to hand out")).toBeVisible();
  });

  test("rejects points above the maximum in the browser", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await c(c(page, "tr", "Adult Anna"), "a", "Results").click();
    const main = page.locator("section[aria-label=Main]");
    await main.locator("input[aria-label='Series 1']").fill("101");
    await expect(c(main, "button", "Save")).toBeDisabled();
  });
});
