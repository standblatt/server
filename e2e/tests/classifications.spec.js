import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Classes Club ${run}`;
const admin = `classes-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("classifications", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "300m", laneCount: 6 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Section 2027", startsOn: "2027-05-20" })).id;
    const programme = (name) => ({ name, showOnRankings: true, series: [{ shots: 10, maxPoints: 100 }], awards: [],
      prices: [{ ageCategory: "elite", price: 0 }, { ageCategory: "u21", price: 0 }] });
    ids.field = (await api(page, "POST", `${base}/events/${ids.event}/programmes`, programme("Field match"))).id;
    ids.section = (await api(page, "POST", `${base}/events/${ids.event}/programmes`, programme("Section stich"))).id;
    const anna = (await api(page, "POST", `${base}/shooters`, { firstname: "Anna", lastname: "Adult", birthday: "1985-01-01" })).id;
    const yves = (await api(page, "POST", `${base}/shooters`, { firstname: "Yves", lastname: "Young", birthday: "2008-01-01" })).id;
    const enter = (shooterId, programmeId, points) =>
      api(page, "POST", `${base}/events/${ids.event}/results`, { shooterId, programmeId, series: [{ points }] });
    await enter(anna, ids.field, 90);
    await enter(anna, ids.section, 85);
    await enter(yves, ids.field, 95);
    await enter(yves, ids.section, 91);
    await page.close();
  });

  test("creates a classification with classes in the setup", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/config/classifications`);
    await c(page, "button", "New classification").click();
    const sheet = page.getByRole("dialog");
    await sheet.getByLabel("Name", { exact: true }).fill("Section championship");
    await c(sheet, "button", "Add programme").click();
    await c(sheet, "button", "Add programme").click();
    await c(sheet, "button", "Add class").click();
    await sheet.getByLabel("Class 1", { exact: true }).fill("Adults");
    await sheet.getByRole("checkbox", { name: "Elite in class 1" }).click();
    await c(sheet, "button", "Add class").click();
    await sheet.getByLabel("Class 2", { exact: true }).fill("Juniors");
    await sheet.getByRole("checkbox", { name: "U21 in class 2" }).click();
    await c(sheet, "button", /^Create$/).click();
    await expect(c(page, "[data-sonner-toast]", "Classification \"Section championship\" saved")).toBeVisible();
    await expect(c(page, "li:not([data-sonner-toast])", "Section championship")).toContainText("Field match + Section stich");

    const [saved] = await api(page, "GET", `/api/tenants/${ids.tenant}/events/${ids.event}/classifications`);
    expect(saved).toMatchObject({ name: "Section championship", tieBreak: ["series", "mouches", "age"], showOnKiosk: true,
      parts: [{ programmeId: ids.field, count: 1 }, { programmeId: ids.section, count: 1 }],
      classes: [{ name: "Adults", categories: ["elite"] }, { name: "Juniors", categories: ["u21"] }] });
  });

  test("shows the classification in the rankings, the report and on kiosk displays", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/rankings`);
    await expect(c(page, "section[aria-label='Classification Section championship – Juniors'] tr", "Young Yves")).toContainText("186");
    await expect(c(page, "section[aria-label='Classification Section championship – Adults'] tr", "Adult Anna")).toContainText("175");

    const report = await page.request.get(`/api/tenants/${ids.tenant}/events/${ids.event}/reports/classifications?format=csv`);
    expect(report.ok()).toBeTruthy();
    expect(await report.text()).toContain("Section championship – Juniors");

    const { secret } = await api(page, "POST", `/api/tenants/${ids.tenant}/events/${ids.event}/kiosk-tokens`, { name: "Hall" });
    await page.goto(`/kiosk/${secret}`);
    await expect(c(page, "section[aria-label='Classification Section championship – Adults'] tr", "Adult Anna")).toContainText("175");
  });
});
