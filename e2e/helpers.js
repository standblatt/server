import { expect } from "@playwright/test";

export const baseURL = process.env.BASE_URL || "http://localhost:3002";
export const mailpit = process.env.MAILPIT || "http://localhost:8026";
export const adminToken = "e2e-admin-token-e2e-admin-token-e2e";
export const password = "e2e-password-123";

/** The first innermost element matching selector that contains text, like cy.contains(selector, text). */
export function c(scope, selector, text) {
  const root = typeof scope.page === "function" ? scope.page() : scope;
  return scope.locator(selector, { hasText: text }).filter({ hasNot: root.locator(selector, { hasText: text }) }).first();
}

/** A page in its own context for setup in beforeAll. */
export const newPage = (browser) => browser.newPage({ baseURL });

/** Calls the API with the cookies of the page and returns the body; any error status fails the test. */
export async function api(page, method, url, data) {
  const res = await page.request.fetch(url, { method, data });
  if (!res.ok()) throw new Error(`${method} ${url}: ${res.status()} ${await res.text()}`);
  return (res.headers()["content-type"] || "").includes("json") ? res.json() : res.text();
}

export const login = (page, name, pw = password) => api(page, "POST", "/api/auth/login", { name, password: pw });

export async function adminLogin(page) {
  await page.goto("/admin");
  await page.locator("input[name=token]").fill(adminToken);
  await c(page, "button", "Log in").click();
}

/** Creates a club and a user who administers it through the admin panel, and returns the club's ID. */
export async function createClub(page, club, user, pw = password) {
  await adminLogin(page);
  await page.goto("/admin/clubs");
  await page.getByPlaceholder("Club name").fill(club);
  await c(page, "button", "Create club").click();
  await expect(c(page, "h1", club)).toBeVisible();
  const clubUrl = page.url();
  await page.goto("/admin/users");
  await page.locator("form[action='/admin/users'] input[name=name]").fill(user);
  await page.locator("form[action='/admin/users'] input[name=password]").fill(pw);
  await c(page, "button", "Create user").click();
  await expect(c(page, "h1", user)).toBeVisible();
  await page.goto(clubUrl);
  await page.locator("input[name=user]").fill(user);
  await c(page, "button", "assign").click();
  await expect(c(page, "tr", user)).toBeVisible();
  await login(page, user, pw);
  const me = await api(page, "GET", "/api/me");
  return me.tenants.find((t) => t.tenantName === club).tenantId;
}

/** Adds a user with a password to a club with one of the default roles (admin, board, staff, member). */
export async function addUser(page, tenantId, name, role = "staff", pw = password) {
  const roles = await api(page, "GET", `/api/tenants/${tenantId}/roles`);
  await api(page, "POST", `/api/tenants/${tenantId}/users`, { name, roleId: roles.find((r) => r.key === role).id, password: pw });
}

/** The token from the newest mail to an address whose subject contains subject; pattern captures it. */
export async function mailToken(page, to, subject = "", pattern = /token=([A-Za-z0-9_-]+)/) {
  const query = encodeURIComponent(`to:${to}${subject ? ` subject:"${subject}"` : ""}`);
  let id;
  await expect
    .poll(async () => {
      const res = await page.request.get(`${mailpit}/api/v1/search?query=${query}`);
      id = (await res.json()).messages[0]?.ID;
      return id;
    })
    .toBeTruthy();
  const msg = await (await page.request.get(`${mailpit}/api/v1/message/${id}`)).json();
  return msg.Text.match(pattern)[1];
}
