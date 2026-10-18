import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Helper Club ${run}`;
const admin = `helper-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("helper links", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 6 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Helper match", startsOn: "2026-11-20", modules: ["reservations", "lane_sheet", "control_strip", "ammo", "awards", "kiosk", "helpers"] })).id;
    ids.main = (await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Main", series: [{ shots: 10, maxPoints: 100 }, { shots: 10, maxPoints: 100 }], prices: [{ ageCategory: "elite", price: 10 }], awards: [],
    })).id;
    const shooter = async (name, lane) => {
      const sh = await api(page, "POST", `${base}/shooters`, { firstname: name, lastname: `${name}son`, birthday: "1990-01-01" });
      return api(page, "POST", `${base}/events/${ids.event}/registrations`, { shooterId: sh.id, lane, programmes: [{ programmeId: ids.main }] });
    };
    await shooter("Anna", 2);
    await shooter("Bert", 5);
    await shooter("Carl");
    await page.close();
  });

  test("creates a scorer link in the setup", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/config/helpers`);
    await page.locator("input[aria-label='Helper name']").fill("Scorer A");
    await page.locator("select[aria-label='Role']").selectOption("scorer");
    await page.locator("input[aria-label='Lanes']").fill("1-3");
    await c(page, "button", "Create helper link").click();
    ids.scorer = new URL(await page.locator("[data-testid=helper-url]").innerText()).pathname;
    await expect(c(page, "li", "Scorer A")).toContainText("Scorer · 1, 2, 3");

    await page.locator("input[aria-label='Helper name']").fill("Desk");
    await page.locator("select[aria-label='Role']").selectOption("ammunition");
    await c(page, "button", "Create helper link").click();
    await expect(c(page, "li", "Desk")).toContainText("Ammunition desk");
    ids.ammo = new URL(await page.locator("[data-testid=helper-url]").innerText()).pathname;
    expect(ids.ammo).not.toBe(ids.scorer);
  });

  test("lets the scorer enter results on their lanes only", async ({ page }) => {
    await page.goto(ids.scorer);
    await expect(c(page, "h1", "Helper match")).toBeVisible();
    await expect(page.locator("[aria-label='Lane 2: Annason Anna']")).toBeVisible();
    await expect(page.locator("[aria-label^='Lane 5']")).toHaveCount(0);
    await expect(c(page, "h2", "Shooting (1)")).toBeVisible();
    await expect(c(page, "tr", "Annason")).toBeVisible();
    await expect(c(page, "tr", "Bertson")).toHaveCount(0);
    await expect(page.locator("input[aria-label='Search shooter']")).toHaveCount(0);
    await page.locator("[aria-label='Lane 2: Annason Anna']").click();
    const inputs = page.locator("[role=dialog] input[data-score]");
    await inputs.nth(0).fill("95");
    await inputs.nth(1).fill("96");
    await c(page, "[role=dialog] button", /^Save$/).click();
    await expect(c(page, "[role=dialog]", "191")).toBeVisible();
    await expect(c(page, "[role=dialog] button", "Edit registration")).toHaveCount(0);
  });

  test("lets the ammunition desk issue ammunition and assign a lane", async ({ page }) => {
    await page.goto(ids.ammo);
    const carl = c(page, "tr", "Carlson");
    await c(carl, "button", "Ammo").click();
    await expect(c(page, "[data-sonner-toast]", "Ammunition issued to Carl Carlson")).toBeVisible();
    await expect(c(carl, "span", "ammo")).toBeVisible();
    await c(carl, "button", "Assign lane").click();
    await expect(page.locator("[role=dialog] [role=checkbox]").first()).toBeDisabled();
    await c(page, "[role=dialog] button", /^4$/).click();
    await c(page, "[role=dialog] button", /^Save$/).click();
    await expect(c(page, "[data-sonner-toast]", "Registration saved")).toBeVisible();
    await expect(page.locator("[aria-label='Lane 4: Carlson Carl']")).toBeVisible();
    await expect(c(page, "button", "Results")).toHaveCount(0);
  });

  test("stops working when the link is revoked", async ({ page, browser }) => {
    const helper = await newPage(browser);
    await helper.goto(ids.scorer);
    await expect(c(helper, "h1", "Helper match")).toBeVisible();
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/config/helpers`);
    const row = c(page, "li", "Scorer A");
    await expect(row).toContainText("last seen");
    await c(row, "button", "revoke").click();
    await c(c(page, "[role=alertdialog]", "Revoke helper link?"), "button", "Revoke").click();
    await expect(c(page, "li", "Scorer A")).toHaveCount(0);
    await expect(c(helper, "[role=alert]", "not valid")).toBeVisible();
    await helper.close();
    await page.context().clearCookies();
    await page.goto(ids.scorer);
    await expect(c(page, "[role=alert]", "not valid")).toBeVisible();
  });
});
