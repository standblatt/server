import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `No Desk Club ${run}`;
const admin = `nodesk-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("events without a desk", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    await api(page, "POST", `${base}/ranges`, { name: "300m", laneCount: 6 });
    ids.shooter = (await api(page, "POST", `${base}/shooters`, { licenseNumber: "500001", firstname: "Olga", lastname: "Obligatory", birthday: "1985-05-05" })).id;
    await page.close();
  });

  test("an event without a range has no lanes and no reservations", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}`);
    await c(page, "button", "New event").click();
    const dialog = page.getByRole("dialog");
    await dialog.getByLabel("Name").fill("Annual programme");
    await dialog.locator("label", { hasText: "300m (6 lanes)" }).getByRole("checkbox").click();
    await c(dialog, "button", "Create").click();
    await page.waitForURL(/\/events\/\d+$/);
    ids.event = Number(page.url().split("/").pop());
    const base = `/api/tenants/${ids.tenant}`;
    expect((await api(page, "GET", `${base}/events/${ids.event}`)).ranges).toEqual([]);
    await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Obligatory", showOnRankings: true, series: [{ shots: 5, maxPoints: 50 }], prices: [{ ageCategory: "elite", price: 0 }], awards: [],
    });

    const sidebar = page.locator("[data-sidebar=sidebar]");
    await expect(sidebar.getByRole("link", { name: "Dashboard" })).toBeVisible();
    await expect(sidebar.getByRole("link", { name: "Reservations" })).toHaveCount(0);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}?register=${ids.shooter}`);
    const reg = page.getByRole("dialog");
    await expect(c(reg, "label", "Obligatory")).toBeVisible();
    await expect(c(reg, "h3", "Lane")).toHaveCount(0);
    await expect(c(reg, "label", "Print lane sheet")).toBeVisible();
  });

  test("lane sheet and control strip can be switched off per event", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/setup?tab=event`);
    await page.getByRole("checkbox", { name: "Lane sheet", exact: true }).click();
    await page.getByRole("checkbox", { name: "Control strip", exact: true }).click();
    await c(page, "button", "Save").click();
    await expect(page.locator("[data-sonner-toast]").first()).toBeVisible();
    const ev = await api(page, "GET", `/api/tenants/${ids.tenant}/events/${ids.event}`);
    expect(ev.settings).toMatchObject({ noLaneSheet: true, noControlStrip: true });

    await page.goto(`/t/${ids.tenant}/events/${ids.event}?register=${ids.shooter}`);
    const reg = page.getByRole("dialog");
    await expect(c(reg, "label", "Obligatory")).toBeVisible();
    await expect(c(reg, "label", "Print lane sheet")).toHaveCount(0);
    await expect(c(reg, "label", "Print control strip")).toHaveCount(0);
  });
});
