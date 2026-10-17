import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Reservation Club ${run}`;
const admin = `resv-${run}`;
const ids = {};
const today = new Date().toISOString().slice(0, 10);

test.describe.configure({ mode: "serial" });

test.describe("lane reservations", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 6 });
    ids.range = range.id;
    const ev = await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Today", startsOn: today });
    ids.event = ev.id;
    await api(page, "POST", `${base}/events/${ev.id}/programmes`, {
      name: "Main", series: [{ shots: 10, maxPoints: 100 }], prices: [{ ageCategory: "elite", price: 10 }], awards: [],
    });
    ids.anna = (await api(page, "POST", `${base}/shooters`, { firstname: "Anna", lastname: "Adult", birthday: "1990-01-01" })).id;
    await page.close();
  });

  test("creates reservations and shows reserved lanes", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/reservations`);
    await c(page, "button", "New reservation").click();
    await c(page, "label", "Description").locator("input").fill("Junior training");
    await c(page, "label", "Lanes").locator("input").fill("2-3");
    await c(page, "button", "Save").click();
    const junior = page.locator("[aria-label='Reservation Junior training']");
    await expect(junior).toContainText("active");
    await expect(junior).toContainText("lanes 2, 3");

    await c(page, "button", "New reservation").click();
    await c(page, "label", "Description").locator("input").fill("Club evening");
    await c(page, "label", "Lanes").locator("input").fill("3, 4");
    await c(page, "button", "Save").click();
    await expect(c(page, "[role=alert]", 'Overlaps with "Junior training" (lanes 2, 3)')).toBeVisible();
    await c(page, "[role=dialog] button", "Close").click();

    await c(page, "a", "Dashboard").click();
    await expect(page.locator("[aria-label='Lane 2: reserved']")).toContainText("Junior training");
    await expect(page.locator("[aria-label='Lane 4: free']")).toBeVisible();

    await api(page, "POST", `/api/tenants/${ids.tenant}/ranges/${ids.range}/reservations`, {
      lanes: "5", startsAt: new Date().toISOString(), endsAt: new Date(Date.now() + 3600000).toISOString(), leadMinutes: 0, description: "Visitors",
    });
    await expect(page.locator("[aria-label='Lane 5: reserved']")).toContainText("Visitors");
  });

  test("asks before registering on a reserved lane", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}?register=${ids.anna}`);
    await c(page, "button", /^2$/).click();
    const ask = page.locator("[role=alertdialog]");
    await expect(ask).toContainText('Reserved for "Junior training"');
    await c(ask, "button", "Cancel").click();
    await expect(page.getByText("No lane: the shooter waits.").first()).toBeVisible();
  });
});
