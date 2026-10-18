import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Attempts Club ${run}`;
const admin = `attempts-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("attempts", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "50m", laneCount: 6 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Festival 2027", startsOn: "2027-06-12" })).id;
    ids.programme = (await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Lucky", showOnRankings: true, series: [{ shots: 5, maxPoints: 50 }], prices: [{ ageCategory: "elite", price: 10 }],
      awards: [{ name: "Card", kind: "card", value: 5, limits: [{ ageCategory: "elite", points: 40 }] }],
    })).id;
    ids.shooter = (await api(page, "POST", `${base}/shooters`, { licenseNumber: "700001", firstname: "Rita", lastname: "Repeat", birthday: "1985-05-05" })).id;
    await page.close();
  });

  test("a programme allows three attempts with a price for further ones", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/config/programmes`);
    await page.locator("[aria-label='Actions for Lucky']").click();
    await c(page, "[role=menuitem]", "Edit").click();
    const sheet = page.getByRole("dialog");
    await sheet.getByLabel("Attempts per shooter").fill("3");
    await sheet.getByRole("combobox", { name: "Counting attempt" }).selectOption({ label: "The best" });
    await sheet.getByLabel("Only the counting attempt receives an award").click();
    await sheet.getByLabel("Further attempts Elite").fill("4");
    await c(sheet, "button", "Save").click();
    await expect(sheet).toBeHidden();
    const p = await api(page, "GET", `/api/tenants/${ids.tenant}/programmes/${ids.programme}`);
    expect(p).toMatchObject({ attempts: 3, attemptRule: "best", awardCountingOnly: true, prices: [{ ageCategory: "elite", price: 10, repeatPrice: 4 }] });
  });

  test("registers two attempts at once and ranks the best", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}?register=${ids.shooter}`);
    const dialog = page.getByRole("dialog");
    await c(dialog, "label", "Lucky").getByRole("checkbox").click();
    await dialog.getByRole("button", { name: "More attempts of Lucky" }).click();
    await expect(dialog.getByTestId("attempts")).toHaveText("2");
    await expect(dialog).toContainText("Total: CHF 14.00");
    await c(dialog, "button", /^Register$/).click();
    await expect(dialog).toBeHidden();

    const regs = await api(page, "GET", `/api/tenants/${ids.tenant}/events/${ids.event}/registrations?all=true`);
    expect(regs[0].entries.map((e) => e.attempt)).toEqual([1, 2]);
    ids.registration = regs[0].id;
    await page.goto(`/t/${ids.tenant}/events/${ids.event}?results=${ids.registration}`);
    for (const [n, points] of [[1, "42"], [2, "47"]]) {
      const card = page.locator(`section[aria-label='Lucky, attempt ${n}']`);
      await card.getByLabel("Series 1", { exact: true }).fill(points);
      await c(card, "button", "Save").click();
      await expect(card).toContainText("completed");
    }
    const sheet = await api(page, "GET", `/api/tenants/${ids.tenant}/registrations/${ids.registration}/results`);
    expect(sheet.entries.map((e) => !!e.award)).toEqual([false, true]);

    await page.goto(`/t/${ids.tenant}/events/${ids.event}/rankings`);
    const line = c(page, "section[aria-label='Ranking Lucky'] tr", "Repeat Rita");
    await expect(line).toContainText("47");
    await expect(line.getByTestId("attempt")).toHaveText("2/2");
  });

  test("enters the third attempt without the desk and refuses a fourth", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/entry`);
    const search = page.getByPlaceholder("License number or name");
    await search.fill("700001");
    await expect(c(page, "[aria-label='Shooters found'] button", "Repeat Rita")).toBeVisible();
    await search.press("Enter");
    const card = page.locator("section[aria-label='Lucky']");
    await card.getByLabel("Series 1", { exact: true }).fill("49");
    await c(card, "button", "Save").click();
    await expect(c(page, "[data-sonner-toast]", "Repeat Rita: 49 points saved")).toBeVisible();

    await search.fill("700001");
    await search.press("Enter");
    await card.getByLabel("Series 1", { exact: true }).fill("30");
    await c(card, "button", "Save").click();
    await expect(c(card, "[role=alert]", "no attempts left")).toBeVisible();

    await page.goto(`/t/${ids.tenant}/events/${ids.event}/rankings`);
    const line = c(page, "section[aria-label='Ranking Lucky'] tr", "Repeat Rita");
    await expect(line).toContainText("49");
    await expect(line.getByTestId("attempt")).toHaveText("3/3");
    await page.goto(`/t/${ids.tenant}/events/${ids.event}?register=${ids.shooter}`);
    await expect(c(page.getByRole("dialog"), "label", "Lucky")).toContainText("no attempts left");
  });
});
