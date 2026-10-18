import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Modules Club ${run}`;
const admin = `modules-${run}`;
const staff = `modules-staff-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("modules", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    ids.event = (await api(page, "POST", `${base}/events`, { name: "Annual programme", startsOn: "2027-03-01", modules: ["awards"] })).id;
    await api(page, "POST", `${base}/users`, { name: staff, role: "staff", password: "e2e-password-123" });
    await page.close();
  });

  const event = () => `/t/${ids.tenant}/events/${ids.event}`;

  test("switches a module on and keeps it on while it holds data", async ({ page }) => {
    await login(page, admin);
    await page.goto(`${event()}/config/modules`);
    const sidebar = page.locator("[data-sidebar=sidebar]");
    await expect(sidebar.getByRole("link", { name: "Modules" })).toBeVisible();
    await expect(sidebar.getByRole("link", { name: "Kiosk displays" })).toHaveCount(0);
    await expect(sidebar.getByRole("link", { name: "Helper links" })).toHaveCount(0);

    const kiosk = page.getByRole("switch", { name: "Kiosk displays", exact: true });
    await kiosk.click();
    await expect(c(page, "[data-sonner-toast]", "Kiosk displays switched on")).toBeVisible();
    await sidebar.getByRole("link", { name: "Kiosk displays" }).click();
    await expect(page).toHaveURL((u) => u.pathname.endsWith("/config/kiosk"));
    await page.locator("input[aria-label='Display name']").fill("Hall");
    await c(page, "button", "Create kiosk link").click();
    await expect(page.locator("[data-testid=kiosk-url]")).toBeVisible();

    await page.goto(`${event()}/config/modules`);
    await expect(kiosk).toBeChecked();
    await expect(kiosk).toBeDisabled();
    await expect(page.locator("[data-testid=module-kiosk]")).toContainText("In use: remove its data first.");
  });

  test("lists only the reports of the modules that are on", async ({ page }) => {
    await login(page, admin);
    await page.goto(`${event()}/statistics`);
    const exports = page.locator("section[aria-label=Exports]");
    await expect(exports).toContainText("Financial report");
    await expect(exports).toContainText("Awards");
    await expect(exports).not.toContainText("Club competition");
  });

  test("sends old setup links to the configuration", async ({ page }) => {
    await login(page, admin);
    await page.goto(`${event()}/setup?tab=kiosk`);
    await expect(page).toHaveURL((u) => u.pathname.endsWith("/config/kiosk"));
    await expect(c(page, "li", "Hall")).toBeVisible();
  });

  test("keeps staff out of the configuration", async ({ page }) => {
    await login(page, staff);
    await page.goto(`${event()}/config`);
    await expect(page).toHaveURL((u) => u.pathname === event());
    await expect(page.locator("[data-sidebar=sidebar]").getByRole("link", { name: "Configuration" })).toHaveCount(0);
  });
});
