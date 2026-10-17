import { test, expect } from "@playwright/test";
import { c, createClub, login, mailToken, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Invite Club ${run}`;
const admin = `inv-${run}`;
const guest = `guest-${run}@example.org`;

test.describe.configure({ mode: "serial" });

test.describe("email invitations", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    await createClub(page, club, admin);
    await page.close();
  });

  test("invites a new person and an invalid entry", async ({ page }) => {
    await login(page, admin);
    await page.goto("/");
    await c(page, "a", club).click();
    await c(page, "a", "Users").click();
    await c(page, "button", "Invite by email").click();
    await page.locator("[role=dialog] textarea[name=emails]").fill(`${guest}, not-an-address`);
    await c(page, "[role=dialog] button", "Send invitations").click();
    const results = page.locator("[data-testid=invite-results]");
    await expect(c(results, "li", guest)).toContainText("invited by mail");
    await expect(c(results, "li", "not-an-address")).toContainText("not an email address");
    await page.keyboard.press("Escape");
    await expect(c(page, "tr", guest)).toBeVisible();
  });

  test("sets the password from the mail and logs in", async ({ page }) => {
    const token = await mailToken(page, guest, "invites you");
    await page.goto(`/reset-password?token=${token}`);
    await page.locator("input[name=password]").fill("my-chosen-password-1");
    await c(page, "button", "Save password").click();
    await expect(page.getByText("Your password is changed").first()).toBeVisible();
    await page.goto("/login");
    await page.locator("input[name=name]").fill(guest);
    await page.locator("input[name=password]").fill("my-chosen-password-1");
    await c(page, "button", "Log in").click();
    await expect(c(page, "a", club)).toBeVisible();
  });
});
