import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Registration Club ${run}`;
const admin = `reg-${run}`;
const ids = {};

function programme(name, series, categories) {
  return {
    name,
    series: Array.from({ length: series }, () => ({ shots: 10, maxPoints: 100 })),
    prices: categories.map((ageCategory) => ({ ageCategory, price: 10 })),
    awards: [],
  };
}

test.describe.configure({ mode: "serial" });

test.describe("registration and dashboard", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 6, disabledLanes: [6] });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Indoor 2026", startsOn: "2026-11-20" })).id;
    const ep = `${base}/events/${ids.event}/programmes`;
    ids.main = (await api(page, "POST", ep, programme("Main", 4, ["elite", "u21"]))).id;
    const section = await api(page, "POST", ep, programme("Section", 2, ["elite"]));
    await api(page, "PUT", `${base}/programmes/${section.id}/combination`, { combinableWithId: ids.main });
    await api(page, "POST", ep, { ...programme("Field match", 2, ["elite"]), variant: "Stgw 90" });
    await api(page, "POST", ep, { ...programme("Field match", 2, ["elite"]), variant: "Stgw 57" });
    ids.anna = (await api(page, "POST", `${base}/shooters`, {
      licenseNumber: "200001", firstname: "Anna", lastname: "Adult", birthday: "1990-01-01", clubNames: ["Home"],
    })).id;
    ids.yves = (await api(page, "POST", `${base}/shooters`, { firstname: "Yves", lastname: "Young", birthday: "2007-03-03" })).id;
    await page.close();
  });

  test("registers a shooter on a lane with a combination", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await expect(c(page, "[aria-label='Live updates']", "live")).toBeVisible();
    await expect(c(page, "h2", "Shooting (0)")).toBeVisible();
    await page.keyboard.press("F4");
    await expect(page.locator(":focus")).toHaveAttribute("aria-label", "Search shooter");
    await page.keyboard.type("adult");
    await c(page, "button", "Register for Home").click();

    await expect(page.locator("[data-testid=category]")).toHaveText("Elite");
    await expect(c(page, "button", /^6$/)).toBeDisabled();
    await c(page, "button", /^2$/).click();
    await c(page, "label", "Main").locator("[role=checkbox]").click();
    await c(page, "label", "Section").locator("[role=checkbox]").click();
    await c(page, "label", "take 2 series from Main").locator("[role=checkbox]").click();
    const series = page.locator("[role=group][aria-label='Series of Main']");
    await c(series, "button", /^2$/).click();
    await c(series, "button", /^4$/).click();
    await expect(c(series, "button", /^3$/)).toBeDisabled();
    await c(page, "label", "Stgw 57").locator("input").check();
    await expect(page.getByText("Total: CHF 30.00").first()).toBeVisible();
    await c(page, "button", /^Register$/).click();

    await expect(page.locator("[role=dialog]")).toHaveCount(0);
    await expect(c(page, "[data-sonner-toast]", "Anna Adult registered on lane 2")).toBeVisible();
    await expect(c(page, "h2", "Shooting (1)")).toBeVisible();
    const row = c(page, "tr", "Adult Anna");
    for (const p of ["Main", "Section", "Field match (Stgw 57)"]) await expect(row).toContainText(p);
    await expect(page.locator("[aria-label='Lane 2: Adult Anna']")).toBeVisible();
  });

  test("shows changes from elsewhere live and opens registrations by code", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await expect(c(page, "[aria-label='Live updates']", "live")).toBeVisible();
    await expect(c(page, "h2", "Waiting (0)")).toBeVisible();
    const reg = await api(page, "POST", `/api/tenants/${ids.tenant}/events/${ids.event}/registrations`, {
      shooterId: ids.yves,
      programmes: [{ programmeId: ids.main }],
    });
    await expect(c(page, "h2", "Waiting (1)")).toBeVisible();
    await expect(c(page, "tr", "Young Yves")).toContainText("U21");

    await page.locator("input[aria-label='Search shooter']").fill(`REGCODEAAAA${reg.id}AAAAL0`);
    await page.keyboard.press("Enter");
    await expect(page).toHaveURL((u) => u.search === `?results=${reg.id}`);
    await expect(c(page, "[role=dialog]", "Results: Young Yves")).toBeVisible();
    await expect(page.locator(":focus")).toHaveAttribute("aria-label", "Series 1");
    await c(page, "button", "Edit registration").click();
    const dialog = c(page, "[role=dialog]", "Total: CHF");
    await c(dialog, "button", /^3$/).click();
    await c(dialog, "button", "Save").click();
    await expect(c(page, "h2", "Shooting (2)")).toBeVisible();
  });

  test("finishes shooting and frees the lane", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await c(c(page, "tr", "Adult Anna"), "button", "Finish").click();
    await c(c(page, "[role=alertdialog]", "Finish shooting?"), "button", "Finish").click();
    await expect(c(page, "h2", "Shooting (1)")).toBeVisible();
    await expect(page.locator("[aria-label=Statistics]")).toContainText("1completed");
    await expect(page.locator("[aria-label='Lane 2: free']")).toBeVisible();

    await page.locator("input[aria-label='Search shooter']").fill("adult");
    await c(page, "button", "Register for Home").click();
    const previous = page.locator("[data-testid=previous-registrations]");
    await expect(previous).toContainText("completed");
    await expect(previous).toContainText("Main");
    const main = c(page, "label", "Main");
    await expect(main).toContainText("already registered");
    await expect(main.locator("[role=checkbox]")).toBeDisabled();
    await c(previous, "button", "Open registration").click();
    await expect(page).toHaveURL((u) => /^\?registration=\d+$/.test(u.search));
  });
});
