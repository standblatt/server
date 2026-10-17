import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Archive Club ${run}`;
const admin = `archive-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("closed events", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 2 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "41. Old match 2026", startsOn: "2026-01-02" })).id;
    const p = await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Main", series: [{ shots: 10, maxPoints: 100 }], prices: [{ ageCategory: "elite", price: 10 }], awards: [],
    });
    const sh = await api(page, "POST", `${base}/shooters`, { firstname: "Anna", lastname: "Adult", birthday: "1990-01-01" });
    const reg = await api(page, "POST", `${base}/events/${ids.event}/registrations`, { shooterId: sh.id, programmes: [{ programmeId: p.id }] });
    ids.registration = reg.id;
    ids.entry = reg.entries[0].id;
    await api(page, "PUT", `${base}/entries/${ids.entry}/results`, { series: [{ points: 88 }] });
    await page.close();
  });

  test("locks a closed event and unlocks it again", async ({ page }) => {
    await login(page, admin);
    const event = `/t/${ids.tenant}/events/${ids.event}`;
    await page.goto(`${event}/setup`);
    await expect(c(page, "td", "Main")).toBeVisible();
    await c(page, "[role=tab]", "Copy & close").click();
    await c(c(page, "section", "Close this event"), "button", "Close event").click();
    await c(c(page, "[role=alertdialog]", "Close event?"), "button", "Close event").click();
    await expect(page.locator("[data-testid=archived]")).toContainText("nothing can be changed");
    await expect(c(page, "button", "Reopen")).toBeVisible();

    await page.goto(event);
    await expect(page.locator("[data-testid=archived]")).toBeVisible();
    await expect(page.locator("button", { hasText: "New shooter" })).toHaveCount(0);
    await page.locator("input[aria-label='Search shooter']").fill("adult");
    await expect(c(page, "li", "Adult Anna")).toBeVisible();
    await expect(c(page, "li", "Adult Anna").locator("button", { hasText: "Register" })).toHaveCount(0);

    await page.goto(`${event}?results=${ids.registration}`);
    const series = page.locator("input[aria-label='Series 1']");
    await expect(series).toBeDisabled();
    await expect(series).toHaveValue("88");
    const res = await page.request.put(`/api/tenants/${ids.tenant}/entries/${ids.entry}/results`, { data: { series: [{ points: 90 }] } });
    expect(res.status()).toBe(409);

    await page.goto(`${event}/setup?tab=more`);
    await c(c(page, "section", "Reopen this event"), "button", "Reopen").click();
    await expect(c(page, "[data-sonner-toast]", "Event reopened")).toBeVisible();
    await expect(page.locator("[data-testid=archived]")).toHaveCount(0);
    await page.goto(`${event}?results=${ids.registration}`);
    await expect(page.locator("input[aria-label='Series 1']")).toBeEnabled();
  });

  test("says why an event with registrations cannot be deleted", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/setup?tab=more`);
    await c(c(page, "section", "Delete this event"), "button", "Delete").click();
    await c(c(page, "[role=alertdialog]", "Delete event?"), "button", "Delete").click();
    await expect(c(page, "[data-sonner-toast]", "Events with registrations cannot be deleted")).toBeVisible();
    await expect(page).toHaveURL((u) => u.pathname.endsWith("/setup"));
  });

  test("copies the event for next year", async ({ page }) => {
    await login(page, admin);
    const setup = `/t/${ids.tenant}/events/${ids.event}/setup`;
    await page.goto(setup);
    await expect(c(page, "td", "Main")).toBeVisible();
    await c(page, "[role=tab]", "Copy & close").click();
    await c(c(page, "section", "Copy for next year"), "button", "Copy").click();
    await expect(page.locator("[role=dialog] input[name=name]")).toHaveValue("42. Old match 2027");
    await expect(page.locator("[role=dialog] input[name=startsOn]")).toHaveValue("2027-01-02");
    await c(page, "[role=dialog] button", "Create copy").click();
    await expect(c(page, "[data-sonner-toast]", 'Event "42. Old match 2027" created')).toBeVisible();
    await expect(page).toHaveURL((u) => /\/events\/\d+\/setup$/.test(u.pathname) && u.pathname !== setup);
    await expect(c(page, "td", "Main")).toBeVisible();
    await expect(c(page, "[data-sidebar=event]", "42. Old match 2027")).toBeVisible();
  });
});
