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
    expect(ev.modules).not.toContain("lane_sheet");
    expect(ev.modules).not.toContain("control_strip");

    await page.goto(`/t/${ids.tenant}/events/${ids.event}?register=${ids.shooter}`);
    const reg = page.getByRole("dialog");
    await expect(c(reg, "label", "Obligatory")).toBeVisible();
    await expect(c(reg, "label", "Print lane sheet")).toHaveCount(0);
    await expect(c(reg, "label", "Print control strip")).toHaveCount(0);
  });

  test("records single shots and sums them into the series", async ({ page }) => {
    await login(page, admin);
    const base = `/api/tenants/${ids.tenant}`;
    ids.shots = (await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Annual stich", showOnRankings: true, series: [{ shots: 3, maxPoints: 30, trackShots: true }, { shots: 2, maxPoints: 20 }],
      prices: [{ ageCategory: "elite", price: 0 }], awards: [],
    })).id;
    const reg = await api(page, "POST", `${base}/events/${ids.event}/registrations`, { shooterId: ids.shooter, programmes: [{ programmeId: ids.shots }] });
    await page.goto(`/t/${ids.tenant}/events/${ids.event}?results=${reg.id}`);
    const card = page.locator("section[aria-label='Annual stich']");
    const shots = card.getByRole("group", { name: "Shots of series 1" });
    await shots.getByLabel("Series 1, shot 1").fill("10");
    await shots.getByLabel("Series 1, shot 2").fill("9");
    await c(card, "button", "Save").click();
    await expect(c(card, "[role=alert]", "Series 1: enter every shot or none")).toBeVisible();
    await shots.getByLabel("Series 1, shot 3").fill("8");
    await expect(card.getByLabel("Series 1", { exact: true })).toHaveValue("27");
    await card.getByLabel("Series 2", { exact: true }).fill("18");
    await c(card, "button", "Save").click();
    await expect(card).toContainText("completed");
    const sheet = await api(page, "GET", `${base}/registrations/${reg.id}/results`);
    expect(sheet.entries[0].results).toEqual([{ points: 27, shots: [10, 9, 8] }, { points: 18 }]);
  });

  test("enters results shooter by shooter without the desk", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/entry`);
    await page.getByLabel("Programme").selectOption({ label: "Obligatory" });
    const search = page.getByPlaceholder("License number or name");
    await search.fill("500001");
    await expect(c(page, "[aria-label='Shooters found'] button", "Obligatory Olga")).toBeVisible();
    await search.press("Enter");
    const card = page.locator("section[aria-label='Obligatory']");
    await card.getByLabel("Series 1", { exact: true }).fill("45");
    await card.getByLabel("Series 1", { exact: true }).press("Enter");
    await expect(c(page, "[data-sonner-toast]", "Obligatory Olga: 45 points saved")).toBeVisible();
    await expect(c(page, "section[aria-label='Entered now'] tr", "Obligatory Olga")).toContainText("45");
    await expect(search).toBeFocused();

    await search.fill("500001");
    await search.press("Enter");
    await card.getByLabel("Series 1", { exact: true }).fill("40");
    await c(card, "button", "Save").click();
    await expect(c(card, "[role=alert]", "already has results")).toBeVisible();

    await page.goto(`/t/${ids.tenant}/events/${ids.event}/rankings`);
    await expect(c(page, "section[aria-label='Ranking Obligatory'] tr", "Obligatory Olga")).toContainText("45");
  });

  test("imports results from a CSV file", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/entry`);
    await c(page, "button", "Import CSV").click();
    const dialog = page.getByRole("dialog");
    const csv = (text) => ({ name: "results.csv", mimeType: "text/csv", buffer: Buffer.from(text) });
    const header = "firstname;lastname;birthday;club;programme;s1;s2;s1_1;s1_2;s1_3\n";
    await dialog.locator("input[type=file]").setInputFiles(csv(header + "Neo;Newman;1990-02-02;;Unknown;44;;;;\n"));
    await c(dialog, "button", /^Import$/).click();
    await expect(c(dialog, "[role=alert]", "Programme \"Unknown\" not found")).toBeVisible();

    await dialog.locator("input[type=file]").setInputFiles(csv(header + "Neo;Newman;1990-02-02;;Obligatory;44;;;;\nNeo;Newman;1990-02-02;;Annual stich;;18;9;9;9\n"));
    await c(dialog, "button", /^Import$/).click();
    await expect(c(dialog, "[role=status]", "Imported: 2 results, 1 new shooters.")).toBeVisible();
    await page.keyboard.press("Escape");

    await page.goto(`/t/${ids.tenant}/events/${ids.event}/rankings`);
    await expect(c(page, "section[aria-label='Ranking Obligatory'] tr", "Newman Neo")).toContainText("44");
    await expect(c(page, "section[aria-label='Ranking Annual stich'] tr", "Newman Neo")).toContainText("45");
  });
});
