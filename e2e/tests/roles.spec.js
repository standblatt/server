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
});
