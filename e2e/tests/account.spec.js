import { test, expect } from "@playwright/test";
import { createHmac } from "node:crypto";
import { c, createClub, login, mailpit, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Account Club ${run}`;
const admin = `account-${run}`;
const newPassword = "a-brand-new-password-1";
const email = `${admin}@example.org`;

/** RFC 6238 code (SHA-1, 30 s, six digits) for a base32 secret, like the authenticator app would show. */
function totp(secret, at = Date.now()) {
  const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567";
  let bits = "";
  for (const ch of secret.toUpperCase()) bits += alphabet.indexOf(ch).toString(2).padStart(5, "0");
  const key = Buffer.from(bits.match(/.{8}/g).map((b) => parseInt(b, 2)));
  const step = Buffer.alloc(8);
  step.writeBigUInt64BE(BigInt(Math.floor(at / 30000)));
  const mac = createHmac("sha1", key).update(step).digest();
  const off = mac[mac.length - 1] & 0x0f;
  return String((mac.readUInt32BE(off) & 0x7fffffff) % 1000000).padStart(6, "0");
}

test.describe.configure({ mode: "serial" });

test.describe("account page", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    await createClub(page, club, admin);
    await page.close();
  });

  test("opens My training without onboarding and shows the empty profile", async ({ page }) => {
    await login(page, admin);
    await page.goto("/my");
    await expect(c(page, "h1", admin)).toBeVisible();
    await page.goto("/my/settings");
    await expect(page.locator("input[name=firstname]")).toHaveValue("");
    await expect(c(page, "h2", "Account")).toBeVisible();
  });

  test("changes the password and sets an email address", async ({ page }) => {
    await login(page, admin);
    await page.goto("/account");
    await expect(c(page, "h1", "Account")).toBeVisible();
    await page.locator("input[name=currentPassword]").fill("wrong-password-1");
    await page.locator("input[name=newPassword]").fill(newPassword);
    await c(page, "button", "Change password").click();
    await expect(page.getByText("Wrong password.").first()).toBeVisible();
    await page.locator("input[name=currentPassword]").fill("e2e-password-123");
    await page.locator("input[name=newPassword]").fill(newPassword);
    await c(page, "button", "Change password").click();
    await expect(c(page, "[data-sonner-toast]", "Password changed")).toBeVisible();

    await c(page, "button", "Add").click();
    await page.locator("[role=dialog] input[name=email]").fill(email);
    await page.locator("[role=dialog] input[name=password]").fill(newPassword);
    await c(page, "[role=dialog] button", "Send confirmation link").click();
    await expect(c(page, "dd", email)).toContainText("not confirmed");
    const mails = await (await page.request.get(`${mailpit}/api/v1/search?query=${encodeURIComponent(`to:${email}`)}`)).json();
    expect(mails.messages.length).toBeGreaterThan(0);
    await expect(page.getByText("Your club manages this account").first()).toBeVisible();
    await expect(c(page, "button", "Delete account")).toHaveCount(0);
  });

  test("switches the second factor on and logs in with a code", async ({ page }) => {
    test.setTimeout(90000);
    await login(page, admin, newPassword);
    await page.goto("/account");
    await c(page, "button", "Set up").click();
    const secret = (await page.locator("[role=dialog] code").innerText()).trim();
    await page.locator("[role=dialog] input[name=code]").fill(totp(secret));
    await c(page, "[role=dialog] button", "Switch on").click();
    await expect(c(page, "[data-sonner-toast]", "Second factor switched on")).toBeVisible();
    await expect(c(page, "button", "Switch off")).toBeVisible();

    await page.context().clearCookies();
    await page.goto("/login");
    await page.locator("input[name=name]").fill(admin);
    await page.locator("input[name=password]").fill(newPassword);
    await c(page, "button", "Log in").click();
    const code = page.locator("input[name=code]");
    await expect(code).toBeVisible();
    await page.waitForTimeout(30000 - (Date.now() % 30000) + 300);
    await code.fill(totp(secret));
    await c(page, "button", "Log in").click();
    await expect(page).toHaveURL((u) => u.pathname === "/");
    await expect(c(page, "a", "My training")).toBeVisible();
  });

  test("logs in with the new password only", async ({ page }) => {
    const res = await page.request.post("/api/auth/login", { data: { name: admin, password: "e2e-password-123" } });
    expect(res.status()).toBe(401);
    const withoutCode = await page.request.post("/api/auth/login", { data: { name: admin, password: newPassword } });
    expect(withoutCode.status()).toBe(428);
  });
});
