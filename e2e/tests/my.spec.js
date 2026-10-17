import { test, expect } from "@playwright/test";
import { api, c, createClub, login, mailToken, mailpit, newPage, password } from "../helpers.js";

const run = Date.now().toString(36);
const email = `felix-${run}@example.org`;
const club = `My Club ${run}`;
const clubAdmin = `myadmin-${run}`;
const license = String(700000 + Math.floor(Math.random() * 99999));
const ids = {};

async function type(page, keys) {
  for (const k of keys) await page.locator("[data-testid=keypad]").locator(`button[aria-label="${k}"]`).click();
}

const removeSeries = async (page, n) => {
  for (let i = 0; i < n; i++) await c(page, "button", "Remove series").click();
};

test.describe.configure({ mode: "serial" });

test.describe("personal training log", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, clubAdmin);

    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 4 });
    const ev = await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Example Cup", startsOn: "2026-03-07" });
    const p = await api(page, "POST", `${base}/events/${ev.id}/programmes`, {
      name: "Cup", series: [{ shots: 10, maxPoints: 100 }, { shots: 10, maxPoints: 100 }], prices: [{ ageCategory: "elite", price: 10 }], awards: [],
    });
    const sh = await api(page, "POST", `${base}/shooters`, { firstname: "Felix", lastname: "Example", birthday: "1990-05-01", licenseNumber: license });
    const reg = await api(page, "POST", `${base}/events/${ev.id}/registrations`, { shooterId: sh.id, programmes: [{ programmeId: p.id }] });
    await api(page, "PUT", `${base}/entries/${reg.entries[0].id}/results`, { series: [{ points: 96 }, { points: 94 }] });
    await api(page, "PUT", `${base}/events/${ev.id}/archived`, { archived: true });
    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem("standblatt-locale", "en"));
  });

  test("signs up and confirms the email address", async ({ page }) => {
    await page.goto("/login");
    await c(page, "a", "Create an account").click();
    await page.locator("input[name=firstname]").fill("Felix");
    await page.locator("input[name=lastname]").fill("Example");
    await page.locator("input[name=email]").fill(email);
    await page.locator("input[name=password]").fill(password);
    await c(page, "button", "Create account").click();
    await expect(page).toHaveURL((u) => u.pathname === "/my");
    await expect(c(page, "h1", "Felix")).toBeVisible();
    await expect(page.getByText("Please confirm your email address").first()).toBeVisible();
    await expect(page.getByText("Your log is empty").first()).toBeVisible();

    await page.goto(`/verify?token=${await mailToken(page, email, "Confirm your email address")}`);
    await expect(page.getByText("Email address confirmed").first()).toBeVisible();
    await page.goto("/my");
    await expect(c(page, "h1", "Felix")).toBeVisible();
    await expect(page.getByText("Please confirm your email address").first()).toHaveCount(0);
  });

  test("logs series totals and shows the result", async ({ page }) => {
    await login(page, email);
    await page.goto("/my");
    await c(page, "a", "Air pistol 10 m").click();
    await expect(page).toHaveURL((u) => u.pathname === "/my/train");
    await type(page, ["9", "3", "9", "5", "9", "1", "9", "6"]);
    await expect(page.locator("[data-testid=session-total]")).toHaveText("375");
    await type(page, ["9", "Delete", "9", "4"]);
    const dialog = page.locator("[role=dialog]");
    await expect(c(dialog, "button", "Save session")).toBeVisible();
    await c(dialog, "button", "Competition").click();
    await c(dialog, "button", "Training").click();
    await dialog.locator("button[aria-label='4']").click();
    await dialog.locator("input[name=tag]").fill("morning");
    await dialog.locator("input[name=tag]").press("Enter");
    await dialog.locator("textarea[name=notes]").fill("Calm and steady");
    await dialog.locator("input[name=place]").fill("Hall A");
    await c(dialog, "button", "Save session").click();
    await expect(page.locator("[data-testid=result-total]")).toHaveText("469");
    await expect(page.locator("[data-testid=feedback]")).toContainText("Your first session in this format");
    await c(page, "a", "Done").click();
    const rows = page.locator("[data-testid=session-row]");
    await expect(rows).toHaveCount(1);
    await expect(rows).toContainText("469");
    await expect(rows).toContainText("morning");
  });

  test("logs shots on the keypad and on the target and analyses the group", async ({ page }) => {
    await login(page, email);
    await page.goto("/my/train");
    await removeSeries(page, 3);
    await expect(page.locator("[role=tablist][aria-label=Series] [role=tab]")).toHaveCount(2);
    await c(page, "[role=tab]", "Shots").click();
    for (const v of [10, 9, 10, 10, 9, 10, 10, 10, 9, 10]) await c(page.locator("[data-testid=keypad]"), "button", new RegExp(`^${v}$`)).click();
    await expect(page.locator("[data-testid=session-total]")).toHaveText("97");
    await c(page, "[role=tab]", "Target").click();
    const target = page.locator("[aria-label='Tap where the shot hit']");
    const box = await target.boundingBox();
    for (let i = 0; i < 10; i++) await target.click({ position: { x: box.width / 2 + 8 + (i % 3), y: box.height / 2 + 8 + (i % 2) } });
    const dialog = page.locator("[role=dialog]");
    await expect(dialog.locator("input[name=place]")).toHaveValue("Hall A");
    await c(dialog, "button", "Save session").click();
    await expect(page.locator("[data-testid=result-total]")).toHaveText(/\d+/);
    const total = Number(await page.locator("[data-testid=result-total]").innerText());
    expect(total).toBeGreaterThanOrEqual(97 + 80);
    expect(total).toBeLessThanOrEqual(197);
    await c(page, "a", "Details").click();
    await expect(c(page, "h2", "Your group")).toBeVisible();
    await expect(page.getByText("Group size").first()).toBeVisible();
  });

  test("edits a session", async ({ page }) => {
    await login(page, email);
    await page.goto("/my/log");
    await c(page, "[data-testid=session-row]", "469").click();
    await page.locator("a[aria-label=Edit]").click();
    await page.locator("[role=tablist][aria-label=Series] [role=tab]").nth(0).click();
    await type(page, ["Delete", "Delete", "Delete", "9", "9"]);
    await expect(page.locator("[data-testid=session-total]")).toHaveText("475");
    await c(page, "button", "Finish").click();
    await c(page.locator("[role=dialog]"), "button", "Save session").click();
    await expect(page.locator("[data-testid=result-total]")).toHaveText("475");
  });

  test("shows the progress of a discipline", async ({ page }) => {
    await login(page, email);
    await page.goto("/my/progress");
    await expect(page.getByText("Personal best").first().locator("..")).toContainText("475");
    await expect(c(page, "h2", "Trend")).toBeVisible();
    await page.getByText("Set a goal").first().click();
    await page.locator("input[name=goal]").fill("480");
    await c(page, "button", "Save").click();
    await expect(page.getByText("480 over 50 shots").first()).toBeVisible();
    await expect(c(page, "h2", "Shot values")).toBeVisible();
    await expect(c(page, "h2", "Your group")).toBeVisible();
  });

  test("adds a competition result found by license number", async ({ page }) => {
    await login(page, email);
    await page.goto("/my/settings");
    await page.locator("input[name=birthday]").fill("1990-05-01");
    await page.locator("input[name=licenseNumber]").fill(license);
    await c(page, "button", "Save").click();
    await expect(page.getByText("Profile saved").first()).toBeVisible();
    await page.goto("/my");
    await page.getByText("We found 1 competition result").first().click();
    const found = c(page, "section", "Your competition results");
    await expect(found).toContainText("Example Cup");
    await expect(found).toContainText("Rank 1 of 1");
    await c(found, "button", "Add").click();
    await expect(page.getByText("Competition added to your log").first()).toBeVisible();
    await expect(c(page, "[data-testid=session-row]", "Example Cup")).toContainText("190");
  });

  test("keeps a session typed offline and uploads it later", async ({ page, context }) => {
    await login(page, email);
    await page.goto("/my/train");
    await expect(page.locator("[data-testid=keypad]")).toBeVisible();
    await context.setOffline(true);
    await type(page, ["8", "8"]);
    await type(page, ["9", "0", "9", "1", "9", "2", "8", "9"]);
    await c(page.locator("[role=dialog]"), "button", "Save session").click();
    await expect(page.getByText("Saved on this device").first()).toBeVisible();
    await context.setOffline(false);
    await page.goto("/my");
    await expect(page.locator("[data-testid=session-row]").first()).toContainText("450");
    expect(await page.evaluate(async () => (await navigator.serviceWorker.ready).scope)).toMatch(/\/my$/);
  });

  test("compares two weapons", async ({ page }) => {
    await login(page, email);
    await page.goto("/my/settings");
    await expect(page.locator("[data-testid=version]")).not.toBeEmpty();
    await expect(page.locator("aside").getByText(/^Standblatt /)).toBeVisible();
    for (const name of ["Steyr EVO 10", "Walther LP500"]) {
      await c(page, "button", "Add weapon").click();
      await page.locator("input[name=weaponName]").fill(name);
      await c(page.locator("[role=dialog]"), "button", "Save").click();
      await expect(c(page, "li", name)).toBeVisible();
    }
    await page.goto("/my/train");
    await removeSeries(page, 4);
    await type(page, ["9", "5"]);
    const dialog = page.locator("[role=dialog]");
    const walther = c(dialog, "[role=radio]", "Walther LP500");
    await walther.click();
    await expect(walther).toHaveAttribute("aria-checked", "true");
    await c(dialog, "button", "Save session").click();
    await expect(page.locator("[data-testid=result-total]")).toHaveText("95");
    const my = await api(page, "GET", "/api/my");
    await api(page, "POST", "/api/my/trainings", {
      disciplineId: my.disciplines.find((d) => d.key === "air_pistol").id,
      kind: "training",
      startedAt: new Date().toISOString(),
      weaponId: my.weapons.find((w) => w.name === "Steyr EVO 10").id,
      series: [{ shots: 10, total: 90 }],
    });
    await page.goto("/my/log");
    await expect(page.locator("[data-testid=session-row]").first()).toContainText("Steyr EVO 10");
    await page.goto("/my/progress");
    const weapons = page.locator("[data-testid=weapons] button");
    await expect(weapons).toHaveCount(2);
    await expect(weapons.first()).toContainText("Walther LP500");
    await expect(weapons.first()).toContainText("9.50");
    await weapons.first().click();
    await expect(page).toHaveURL(/[?&]w=\d+/);
    await expect(page.locator("select[aria-label=Weapon] option:checked")).toHaveText("Walther LP500");
    await expect(c(page, "div", /^Sessions$/).locator("..")).toContainText("1");
  });

  test("creates an API token, uses it and revokes it", async ({ page, playwright }) => {
    await login(page, email);
    await page.goto("/my/settings");
    await c(page, "button", "New token").click();
    const dialog = page.locator("[role=dialog]");
    await dialog.locator("input[name=tokenName]").fill("import script");
    await c(dialog, "label", "Log and change sessions").locator("input").check();
    await c(dialog, "button", "Create token").click();
    const token = await dialog.locator("[data-testid=api-token-secret]").innerText();
    expect(token).toMatch(/^sb_[0-9a-f]{40}$/);
    await page.keyboard.press("Escape");
    await expect(c(page, "[data-testid=api-token]", "import script")).toContainText("never used");

    const client = await playwright.request.newContext({ baseURL: new URL(page.url()).origin, extraHTTPHeaders: { Authorization: `Bearer ${token}` } });
    expect((await client.get("/api/my/trainings")).status()).toBe(200);
    expect((await client.get("/api/my/tokens")).status()).toBe(403);
    await page.reload();
    await expect(c(page, "[data-testid=api-token]", "import script")).toContainText("last used");

    await page.locator("[aria-label='Revoke import script']").click();
    await c(page.locator("[role=alertdialog]"), "button", "Revoke").click();
    await expect(page.locator("[data-testid=api-token]")).toHaveCount(0);
    expect((await client.get("/api/my/trainings")).status()).toBe(401);
    await client.dispose();
  });

  test("deletes a session without waiting", async ({ page }) => {
    await login(page, email);
    const my = await api(page, "GET", "/api/my");
    const saved = await api(page, "POST", "/api/my/trainings", {
      disciplineId: my.disciplines.find((d) => d.key === "air_pistol").id,
      kind: "training",
      startedAt: "2020-01-01T10:00:00Z",
      title: "To be deleted",
      series: [{ shots: 10, total: 77 }],
    });
    await page.goto(`/my/log/${saved.session.id}`);
    await page.locator("button[aria-label=Delete]").click();
    await c(page.locator("[role=alertdialog]"), "button", "Delete").click();
    await expect(page).toHaveURL((u) => u.pathname === "/my/log", { timeout: 2000 });
    await expect(c(page, "[data-testid=session-row]", "To be deleted")).toHaveCount(0);
  });

  test("logs out from the sidebar", async ({ page }) => {
    await login(page, email);
    await page.goto("/my");
    await page.locator("aside button[aria-label='Log out']").click();
    await expect(page).toHaveURL((u) => u.pathname === "/login");
    expect((await page.request.get("/api/my")).status()).toBe(401);
  });

  test("deletes the account on the account page", async ({ page }) => {
    await login(page, email);
    await page.goto("/account");
    await c(page, "button", "Delete account").click();
    await page.locator("[role=dialog] input[name=password]").fill("wrong-password-1");
    await c(page, "[role=dialog] button", "Delete account").click();
    await expect(page.getByText("Wrong password.").first()).toBeVisible();
    await page.locator("[role=dialog] input[name=password]").fill(password);
    await c(page, "[role=dialog] button", "Delete account").click();
    await expect(page).toHaveURL((u) => u.pathname === "/login");
    expect((await page.request.post("/api/auth/login", { data: { name: email, password } })).status()).toBe(401);
  });
});
