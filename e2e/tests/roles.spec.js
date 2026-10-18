import { test, expect } from "@playwright/test";
import { addUser, api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Roles Club ${run}`;
const admin = `roles-${run}`;
const member = `member-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("roles", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    ids.event = (await api(page, "POST", `/api/tenants/${ids.tenant}/events`, { name: "Club championship", startsOn: "2027-03-01" })).id;
    await addUser(page, ids.tenant, member, "member");
    await page.close();
  });

  test("shows a member the rankings and nothing else", async ({ page }) => {
    await login(page, member);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    await expect(page).toHaveURL((u) => u.pathname.endsWith("/rankings"));
    await expect(c(page, "h1", "Rankings")).toBeVisible();
    const sidebar = page.locator("[data-sidebar=sidebar]");
    await expect(sidebar).toContainText("Member");
    for (const hidden of ["Dashboard", "Enter results", "Statistics", "Configuration", "Shooters", "Users", "Ranges"]) {
      await expect(sidebar).not.toContainText(hidden);
    }
    await expect(page.locator("[data-sonner-toast]")).toHaveCount(0);
    await page.goto(`/t/${ids.tenant}/shooters`);
    await expect(page).toHaveURL((u) => u.pathname === `/t/${ids.tenant}`);
  });

  test("makes a role in the matrix and gives it to a member", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/users`);
    const matrix = page.locator("section[aria-label=Roles]");
    await expect(matrix.getByRole("checkbox", { name: "Run events for Admin" })).toBeDisabled();
    await expect(matrix.getByRole("checkbox", { name: "Run events for Helper" })).toBeChecked();
    await c(matrix, "button", "New role").click();
    const dialog = page.locator("[role=dialog]");
    await c(dialog, "label", "Name").locator("input").fill("Treasurer");
    await c(dialog, "button", "Save").click();
    await expect(c(page, "[data-sonner-toast]", "Role Treasurer created")).toBeVisible();
    await matrix.getByRole("checkbox", { name: "Statistics and reports for Treasurer" }).click();
    await expect(matrix.getByRole("checkbox", { name: "Statistics and reports for Treasurer" })).toBeChecked();
    await c(page, "tr", member).locator("select").selectOption({ label: "Treasurer" });
    await expect(c(page, "[data-sonner-toast]", `${member} is now Treasurer`)).toBeVisible();
    await expect(c(matrix, "th", "Treasurer")).toContainText("1 member");

    await matrix.getByRole("button", { name: "Delete Treasurer" }).click();
    await c(page, "[role=alertdialog] button", "Delete").click();
    await expect(c(page, "[data-sonner-toast]", "another role first")).toBeVisible();

    await login(page, member);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/rankings`);
    const sidebar = page.locator("[data-sidebar=sidebar]");
    await expect(sidebar).toContainText("Treasurer");
    await expect(sidebar).toContainText("Statistics");
    await expect(sidebar).not.toContainText("Dashboard");
  });

  test("shows a manager without the admin role only what it may give", async ({ page }) => {
    await login(page, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const manager = await api(page, "POST", `${base}/roles`, { name: "Manager", permissions: ["club.manage", "events.view"] });
    const name = `manager-${run}`;
    await api(page, "POST", `${base}/users`, { name, roleId: manager.id, password: "e2e-password-123" });

    await login(page, name);
    await page.goto(`/t/${ids.tenant}/users`);
    const boss = c(page, "tr", admin);
    await expect(boss).toContainText("Admin");
    await expect(boss.locator("select")).toHaveCount(0);
    await expect(boss.getByRole("button", { name: "Remove" })).toHaveCount(0);
    const options = c(page, "tr", name).locator("select option");
    await expect(options.filter({ hasText: "Member" })).toHaveCount(1);
    await expect(options.filter({ hasText: /^Admin$/ })).toHaveCount(0);
    await expect(options.filter({ hasText: "Board" })).toHaveCount(0);
    const matrix = page.locator("section[aria-label=Roles]");
    await expect(matrix.getByRole("button", { name: "Rename Admin" })).toHaveCount(0);
    await expect(matrix.getByRole("button", { name: "Delete Board" })).toHaveCount(0);
    await expect(matrix.getByRole("checkbox", { name: "Run events for Member" })).toBeDisabled();
    await expect(matrix.getByRole("checkbox", { name: "See events and rankings for Member" })).toBeEnabled();
  });
});
