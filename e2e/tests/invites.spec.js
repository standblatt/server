import { test, expect } from "@playwright/test";
import { api, c, createClub, login, mailToken, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Invite Club ${run}`;
const other = `Other Club ${run}`;
const admin = `inv-${run}`;
const otherAdmin = `inv-other-${run}`;
const guest = `guest-${run}@example.org`;
const ids = {};
const invitePattern = /\/invite\/(inv_[A-Za-z0-9_-]+)/;

test.describe.configure({ mode: "serial" });

test.describe("email invitations", () => {
  test.beforeAll(async ({ browser }) => {
    for (const [key, name, user] of [["club", club, admin], ["other", other, otherAdmin]]) {
      const page = await newPage(browser);
      ids[key] = await createClub(page, name, user);
      await page.close();
    }
  });

  test("invites a new person and an invalid entry", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.club}/users`);
    await c(page, "button", "Invite by email").click();
    await page.locator("[role=dialog] textarea[name=emails]").fill(`${guest}, not-an-address`);
    await c(page, "[role=dialog] button", "Send invitations").click();
    const results = page.locator("[data-testid=invite-results]");
    await expect(c(results, "li", guest)).toContainText("invitation sent");
    await expect(c(results, "li", "not-an-address")).toContainText("not an email address");
    await page.keyboard.press("Escape");
    await expect(c(page, "section[aria-label=Invitations] li", guest)).toContainText("Helper");
    await expect(c(page, "tr", guest)).toHaveCount(0);
  });

  test("accepts with a new account", async ({ page }) => {
    const token = await mailToken(page, guest, "invites you", invitePattern);
    await page.goto(`/invite/${token}`);
    await expect(page.getByText(`${club} invites you to Standblatt as Helper.`)).toBeVisible();
    await page.locator("input[name=password]").fill("my-chosen-password-1");
    await c(page, "button", "Accept and create the account").click();
    await expect(c(page, "a", club)).toBeVisible();
    await page.goto(`/invite/${token}`);
    await expect(page.getByText("accepted already")).toBeVisible();
  });

  test("accepts with an existing account after logging in", async ({ page }) => {
    await login(page, otherAdmin);
    const roles = await api(page, "GET", `/api/tenants/${ids.other}/roles`);
    await api(page, "POST", `/api/tenants/${ids.other}/invitations`, { emails: [guest], roleId: roles.find((r) => r.key === "member").id });
    await page.request.post("/api/auth/logout");
    const token = await mailToken(page, guest, other, invitePattern);
    await page.goto(`/invite/${token}`);
    await expect(page.getByText(`There is an account with ${guest}`)).toBeVisible();
    await c(page, "a", "Log in").click();
    await page.locator("input[name=name]").fill(guest);
    await page.locator("input[name=password]").fill("my-chosen-password-1");
    await c(page, "button", "Log in").click();
    await expect(page.getByText(`You join as ${guest}.`)).toBeVisible();
    await c(page, "button", "Accept the invitation").click();
    await expect(c(page, "a", other)).toBeVisible();
    await expect(c(page, "a", club)).toBeVisible();
  });

  test("withdraws an invitation", async ({ page }) => {
    await login(page, admin);
    const late = `late-${run}@example.org`;
    const roles = await api(page, "GET", `/api/tenants/${ids.club}/roles`);
    await api(page, "POST", `/api/tenants/${ids.club}/invitations`, { emails: [late], roleId: roles.find((r) => r.key === "staff").id });
    const token = await mailToken(page, late, "invites you", invitePattern);
    await page.goto(`/t/${ids.club}/users`);
    const row = c(page, "section[aria-label=Invitations] li", late);
    await c(row, "button", "Send again").click();
    await expect(c(page, "[data-sonner-toast]", `Invitation to ${late} sent again`)).toBeVisible();
    await c(row, "button", "Withdraw").click();
    await c(page, "[role=alertdialog] button", "Withdraw").click();
    await expect(c(page, "section[aria-label=Invitations] li", late)).toHaveCount(0);
    await page.goto(`/invite/${token}`);
    await expect(page.getByText("Invitation not found")).toBeVisible();
  });
});
