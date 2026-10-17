import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Help Club ${run}`;
const admin = `help-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("help", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const range = await api(page, "POST", `/api/tenants/${ids.tenant}/ranges`, { name: "10m", laneCount: 4 });
    ids.event = (await api(page, "POST", `/api/tenants/${ids.tenant}/events`, { rangeIds: [range.id], name: "Help Cup", startsOn: "2027-03-13" })).id;
    await page.close();
  });

  test("every page, internal link and image works without login", async ({ page }) => {
    // About 80 page loads: 5 s locally, more than 30 s on a busy CI runner.
    test.setTimeout(180_000);
    await page.goto("/help?lang=en");
    const nav = page.getByRole("navigation", { name: "Help pages" });
    await expect(nav.getByRole("link").first()).toBeVisible();
    const pages = await nav.getByRole("link").evaluateAll((links) => links.map((a) => a.getAttribute("href")));
    expect(pages.length).toBeGreaterThan(20);
    const targets = new Set();
    for (const href of pages) {
      await page.goto(`${href}?lang=en`);
      await expect(page.locator("article.help-prose h1"), href).toBeVisible();
      for (const t of await page.locator("article a[href^='/help']").evaluateAll((links) => links.map((a) => a.getAttribute("href")))) targets.add(t);
      for (const src of await page.locator("article img").evaluateAll((imgs) => imgs.map((i) => i.getAttribute("src")))) {
        expect((await page.request.get(src)).status(), `${href}: ${src}`).toBe(200);
      }
    }
    for (const t of targets) {
      const [path, hash] = t.split("#");
      await page.goto(`${path}?lang=en`);
      await expect(page.locator("article.help-prose h1"), t).toBeVisible();
      if (hash) await expect(page.locator(`article [id='${hash}']`), t).toHaveCount(1);
    }
  });

  test("shows the language of the page and falls back to English", async ({ page }) => {
    await page.goto("/help/user-guide/reservations?lang=de");
    await expect(c(page, "article h1", "Reservationen")).toBeVisible();
    await expect(page.getByText("Diese Seite ist noch nicht übersetzt.")).toHaveCount(0);
    await page.getByRole("link", { name: "Changelog" }).click();
    await expect(c(page, "article h1", "Changelog")).toBeVisible();
    await expect(page.getByText("Diese Seite ist noch nicht übersetzt.")).toBeVisible();
    await page.getByLabel("Hilfe durchsuchen").fill("Anmeldeschluss");
    await c(page, "[aria-label=Suchergebnisse] a", "Online-Reservation").first().click();
    await expect(c(page, "article h1", "Online-Reservation")).toBeVisible();
  });

  test("the help button opens the page next to the work", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/reservations`);
    await page.getByRole("button", { name: "Help", exact: true }).click();
    const sheet = page.getByRole("dialog");
    await expect(c(sheet, "article h1", "Reservations")).toBeVisible();
    await expect(sheet.getByRole("link", { name: "Open the full help" })).toHaveAttribute("href", "/help/user-guide/reservations");
  });
});
