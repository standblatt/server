import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Season Club ${run}`;
const admin = `season-${run}`;
const ids = {};

/** YYYY-MM-DD n days from today. */
const day = (n) => new Date(Date.now() + n * 86400000).toISOString().slice(0, 10);

test.describe.configure({ mode: "serial" });

test.describe("seasons and lane groups", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    await api(page, "POST", `/api/tenants/${ids.tenant}/ranges`, { name: "300m", laneCount: 6 });
    ids.shooter = (await api(page, "POST", `/api/tenants/${ids.tenant}/shooters`, { licenseNumber: "400001", firstname: "Paula", lastname: "Pistol", birthday: "1990-01-01" })).id;
    await page.close();
  });

  test("an event over a season with weekly shooting times stays under Today", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}`);
    await c(page, "button", "New event").click();
    const dialog = page.getByRole("dialog");
    await dialog.getByLabel("Name").fill("Feldschiessen");
    await dialog.locator("label", { hasText: "Period (season)" }).getByRole("switch").click();
    await dialog.getByLabel("First day").fill(day(-30));
    await dialog.getByLabel("Last day").fill(day(30));
    const times = dialog.getByRole("group", { name: "Shooting times" });
    await times.locator("label", { hasText: "Wednesday" }).getByRole("checkbox").click();
    await times.getByRole("button", { name: "Time" }).click();
    await dialog.getByLabel("Wednesday from").fill("17:30");
    await dialog.getByLabel("Wednesday until").fill("19:30");
    await expect(dialog.getByText(/^\d+ shooting days$/)).toBeVisible();
    await c(dialog, "button", "Create").click();
    await expect(c(page, "[data-sonner-toast]", 'Event "Feldschiessen" created')).toBeVisible();

    await page.waitForURL(/\/events\/\d+$/);
    const id = Number(page.url().split("/").pop());
    const ev = await api(page, "GET", `/api/tenants/${ids.tenant}/events/${id}`);
    expect(ev.schedule.weekly).toEqual([{ weekday: 3, times: [{ from: "17:30", to: "19:30" }] }]);
    expect(ev.days.length).toBeGreaterThan(7);
    await page.goto(`/t/${ids.tenant}`);
    await expect(page.locator("section", { has: c(page, "h2", "Today") }).getByText("Feldschiessen")).toBeVisible();

  });

  test("the setup form opens events whose days and weekdays have no shooting times", async ({ page }) => {
    await login(page, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = (await api(page, "GET", `${base}/ranges`))[0];
    const ev = await api(page, "POST", `${base}/events`, {
      rangeIds: [range.id], name: "Without times", startsOn: day(3),
      schedule: { days: [{ date: day(4) }], from: day(-10), to: day(10), weekly: [{ weekday: 6 }] },
    });
    const errors = [];
    page.on("pageerror", (e) => errors.push(e.message));
    await page.goto(`/t/${ids.tenant}/events/${ev.id}/setup?tab=event`);
    await expect(page.getByLabel("Shooting day")).toHaveValue(day(4));
    await page.getByRole("group", { name: "Shooting times" }).getByRole("button", { name: "Time" }).click();
    await expect(page.getByLabel("Saturday from")).toBeVisible();
    expect(errors).toEqual([]);
  });

  test("a pistol range with 25 m and 50 m lane groups", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/ranges`);
    await c(page, "button", "New range").click();
    const dialog = page.getByRole("dialog");
    await dialog.getByLabel("Name", { exact: true }).fill("Pistol");
    await dialog.getByLabel("Number of lanes").fill("15");
    await c(dialog, "button", "Add group").click();
    await dialog.getByLabel("Group name").fill("25 m");
    await dialog.getByLabel("Lanes of 25 m").fill("1-10");
    await c(dialog, "button", "Add group").click();
    await dialog.getByLabel("Group name").nth(1).fill("50 m");
    await dialog.getByLabel("Lanes of 50 m").fill("11-15");
    await c(dialog, "button", "Save").click();
    await expect(c(page, "td", "25 m: 1-10 · 50 m: 11-15")).toBeVisible();

    const base = `/api/tenants/${ids.tenant}`;
    const range = (await api(page, "GET", `${base}/ranges`)).find((r) => r.name === "Pistol");
    const event = await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Pistol day", startsOn: day(1) });
    const p50 = await api(page, "POST", `${base}/events/${event.id}/programmes`, {
      name: "Precision", showOnRankings: true, series: [{ shots: 10, maxPoints: 100 }], prices: [{ ageCategory: "elite", price: 20 }], awards: [],
    });
    await api(page, "PUT", `${base}/programmes/${p50.id}/range`, { groupId: range.groups[1].id });
    const wrong = await page.request.post(`${base}/events/${event.id}/registrations`, {
      data: { shooterId: ids.shooter, rangeId: range.id, lane: 3, programmes: [{ programmeId: p50.id }] },
    });
    expect(wrong.status()).toBe(422);
    expect(await wrong.text()).toContain("Precision is shot on 50 m, not on lane 3");

    await page.goto(`/t/${ids.tenant}/events/${event.id}`);
    await expect(c(page, "h3", "25 m")).toBeVisible();
    await expect(c(page, "h3", "50 m")).toBeVisible();
    await page.keyboard.press("F4");
    await page.keyboard.type("pistol");
    await c(page, "button", /^Register$/).click();
    await c(page, "label", "Precision").locator("[role=checkbox]").click();
    await expect(page.getByText("The chosen programmes are shot on Pistol · 50 m.")).toBeVisible();
    const grid = page.getByRole("group", { name: "Pistol 50 m" });
    await expect(grid.getByRole("button")).toHaveCount(5);
    await c(grid, "button", /^12$/).click();
    await c(page, "button", /^Register$/).click();
    await expect(c(page, "[data-sonner-toast]", "Paula Pistol registered")).toBeVisible();
  });
});
