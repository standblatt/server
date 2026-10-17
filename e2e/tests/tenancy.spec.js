import { test, expect } from "@playwright/test";
import { adminLogin, c, password } from "../helpers.js";

const run = Date.now().toString(36);
const club = `E2E Club ${run}`;
const admin = `boss-${run}`;
const staff = `helper-${run}`;

async function login(page, name) {
  await page.goto("/login");
  await page.locator("input[name=name]").fill(name);
  await page.locator("input[name=password]").fill(password);
  await c(page, "button", "Log in").click();
}

test.describe.configure({ mode: "serial" });

test.describe("admin panel and tenancy", () => {
  test("rejects a wrong admin token", async ({ page }) => {
    await page.goto("/admin");
    await page.locator("input[name=token]").fill("wrong");
    await c(page, "button", "Log in").click();
    await expect(page.getByText("Wrong token").first()).toBeVisible();
  });

  test("creates a club with an admin user", async ({ page }) => {
    await adminLogin(page);
    await expect(c(page, "h1", "Administration")).toBeVisible();
    await c(page, "nav a", "Clubs").click();
    await page.getByPlaceholder("Club name").fill(club);
    await c(page, "button", "Create club").click();
    await expect(c(page, "h1", club)).toBeVisible();
    const clubUrl = page.url();
    await c(page, "nav a", "Users").click();
    await page.locator("form[action='/admin/users'] input[name=name]").fill(admin);
    await page.locator("form[action='/admin/users'] input[name=password]").fill(password);
    await c(page, "button", "Create user").click();
    await expect(c(page, "h1", admin)).toBeVisible();
    await page.locator("form[action$='/email'] input[name=email]").fill(`${admin}@example.org`);
    await c(page, "button", "Save").click();
    await expect(page.locator("form[action$='/email'] input[name=email]")).toHaveValue(`${admin}@example.org`);
    await page.goto(clubUrl);
    await page.locator("input[name=user]").fill(`${admin}@example.org`);
    await page.locator("select[name=role]").selectOption("admin");
    await c(page, "button", "assign").click();
    await expect(c(page, "tr", admin).locator(".badge")).toHaveText("admin");
  });

  test("lets the club admin set up a range, an event and a staff user", async ({ page }) => {
    await login(page, admin);
    await expect(c(page, "a", "My training")).toBeVisible();
    await c(page, "a", club).click();
    await expect(c(page, "header", club)).toBeVisible();
    const start = page.locator("section[aria-label='Getting started']");
    await c(start, "a", "Ranges").click();
    await c(page, "button", "New range").click();
    await c(page, "label", "Name").locator("input").fill("10m hall");
    await c(page, "label", "Number of lanes").locator("input").fill("12");
    await c(page, "label", "Out of service").locator("input").fill("3, 7");
    await c(page, "button", "Save").click();
    await expect(c(page, "tr", "10m hall")).toContainText("3, 7");

    await c(page, "a", "Events").click();
    await c(page, "button", "New event").click();
    await c(page, "label", "Name").locator("input").fill("Indoor 2026");
    await c(page, "label", "Date").locator("input").fill("2026-11-20");
    await c(page, "label", "Organizer").locator("input").fill(club);
    await c(page, "button", "Create").click();
    await expect(c(page, "header", "Indoor 2026")).toBeVisible();
    await expect(page.getByText("20.11.2026 · 10m hall").first()).toBeVisible();
    await expect(start.locator("li[data-done=true]")).toHaveCount(2);
    await expect(c(start, "a", "Setup")).toBeVisible();

    await c(page, "a", "Users").click();
    await c(page, "button", "Add user").click();
    const dialog = page.locator("[role=dialog]");
    await c(dialog, "label", "User name").locator("input").fill(staff);
    await c(dialog, "label", "Password").locator("input").fill(password);
    await c(dialog, "button", /^Add$/).click();
    await expect(c(page, "tr", staff).locator("select")).toHaveValue("staff");
  });

  test("gives staff read access without admin pages", async ({ page }) => {
    await login(page, staff);
    await c(page, "a", club).click();
    await c(page, "a", "Indoor 2026").click();
    await expect(c(page, "header", "Indoor 2026")).toBeVisible();
    await expect(page.locator("button", { hasText: "Edit" })).toHaveCount(0);
    const sidebar = page.locator("[data-sidebar=sidebar]");
    await expect(sidebar).toContainText("Shooters");
    await expect(sidebar).not.toContainText("Ranges");
    await expect(sidebar).not.toContainText("Users");
    await page.locator("[aria-label='User menu']").click();
    await c(page, "[role=menuitem]", "My training").click();
    await expect(c(page, "h1", staff)).toBeVisible();
    await page.goBack();
    await page.locator("[aria-label='User menu']").click();
    await c(page, "[role=menuitem]", "Log out").click();
    await expect(page).toHaveURL((u) => u.pathname === "/login");
  });

  test("shows a login error for a wrong password", async ({ page }) => {
    await page.goto("/login");
    await page.locator("input[name=name]").fill(admin);
    await page.locator("input[name=password]").fill("wrong-wrong");
    await c(page, "button", "Log in").click();
    await expect(c(page, "[role=alert]", "Wrong user name or password")).toBeVisible();
  });
});
