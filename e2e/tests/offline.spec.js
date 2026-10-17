import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage, password } from "../helpers.js";

const run = Date.now().toString(36);
const userA = `offline-a-${run}`;
const userB = `offline-b-${run}`;
const ids = {};

async function type(page, keys) {
  for (const k of keys) await page.locator("[data-testid=keypad]").locator(`button[aria-label="${k}"]`).click();
}

/** Types a session on /my/train while offline and saves it on the device. */
async function trainOffline(page, context) {
  await page.goto("/my/train");
  await expect(page.locator("[data-testid=keypad]")).toBeVisible();
  await context.setOffline(true);
  await type(page, ["8", "8"]);
  await type(page, ["9", "0", "9", "1", "9", "2", "8", "9"]);
  await c(page.locator("[role=dialog]"), "button", "Save session").click();
  await expect(page.getByText("Saved on this device").first()).toBeVisible();
  await context.setOffline(false);
}

/** Logs in through the login form. */
async function loginForm(page, name) {
  await expect(page).toHaveURL((u) => u.pathname === "/login");
  await page.locator("input[name=name]").fill(name);
  await page.locator("input[name=password]").fill(password);
  await c(page, "button", "Log in").click();
  await expect(page).not.toHaveURL((u) => u.pathname === "/login");
}

const rows = (page) => page.locator("[data-testid=session-row]");
const queued = (page) => page.getByText("saved on this device and will be uploaded");

test.describe.configure({ mode: "serial" });

test.describe("offline", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, `Offline Club ${run}`, userA);
    await page.close();
    const other = await newPage(browser);
    await createClub(other, `Offline Other ${run}`, userB);
    await other.close();
  });

  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem("standblatt-locale", "en"));
  });

  test("uploads a session saved offline after the login expired", async ({ page, context }) => {
    await login(page, userA);
    await trainOffline(page, context);
    await context.clearCookies();
    await page.goto("/my");
    await loginForm(page, userA);
    await page.goto("/my");
    await expect(rows(page)).toHaveCount(1);
    await expect(queued(page)).toHaveCount(0);
  });

  test("keeps a session saved offline for the user who typed it", async ({ page, context }) => {
    await login(page, userA);
    await trainOffline(page, context);
    await context.clearCookies();
    await login(page, userB);
    await page.goto("/my");
    await expect(page.getByText("Your log is empty").first()).toBeVisible();
    await expect(queued(page)).toHaveCount(0);
    await context.clearCookies();
    await login(page, userA);
    await page.goto("/my");
    await expect(rows(page)).toHaveCount(2);
    await expect(queued(page)).toHaveCount(0);
  });

  test("shows a refused upload and discards it", async ({ page }) => {
    await login(page, userA);
    const my = await api(page, "GET", "/api/my");
    await page.goto("/my");
    await expect(rows(page)).toHaveCount(2);
    await page.evaluate((userId) => {
      const input = { clientId: crypto.randomUUID(), disciplineId: 999999999, kind: "training", startedAt: new Date().toISOString(), series: [{ shots: [9] }] };
      localStorage.setItem("standblatt-outbox", JSON.stringify([{ userId, input }]));
    }, my.userId);
    await page.reload();
    const failed = page.locator("[data-testid=outbox-failed]");
    await expect(failed).toContainText("could not be uploaded");
    await page.reload();
    await expect(failed).toBeVisible();
    await c(failed, "button", "Discard").click();
    await c(page.locator("[role=alertdialog]"), "button", "Discard").click();
    await expect(failed).toHaveCount(0);
    await expect(rows(page)).toHaveCount(2);
  });

  test("opens the log offline", async ({ page, context }) => {
    await login(page, userA);
    await page.goto("/my");
    await expect(rows(page)).toHaveCount(2);
    await page.evaluate(async () => {
      await navigator.serviceWorker.ready;
      if (!navigator.serviceWorker.controller) await new Promise((r) => navigator.serviceWorker.addEventListener("controllerchange", r));
    });
    await context.setOffline(true);
    await page.reload();
    await expect(rows(page)).toHaveCount(2);
    await context.setOffline(false);
  });

  test("warns on club pages without a connection and keeps the input", async ({ page, context }) => {
    await login(page, userA);
    await page.goto(`/t/${ids.tenant}/shooters`);
    await c(page, "button", "New shooter").click();
    await c(page, "label", "First name").locator("input").fill("Nora");
    await c(page, "label", "Last name").locator("input").fill("Example");
    await context.setOffline(true);
    await expect(page.locator("[data-testid=offline]")).toBeVisible();
    await c(page, "button", "Create").click();
    const dialog = page.locator("[role=dialog]");
    await expect(dialog).toContainText("No connection. Nothing was saved.");
    await expect(c(dialog, "label", "First name").locator("input")).toHaveValue("Nora");
    await context.setOffline(false);
    await expect(page.locator("[data-testid=offline]")).toHaveCount(0);
    await c(page, "button", "Create").click();
    await expect(page.getByText("Nora Example created").first()).toBeVisible();
  });
});
