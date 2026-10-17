import { test, expect } from "@playwright/test";
import { c } from "../helpers.js";

test.describe("smoke", () => {
  test("serves the app and the API", async ({ page, request }) => {
    expect(await (await request.get("/healthz")).text()).toBe("ok");
    expect((await request.get("/readyz")).status()).toBe(200);
    expect(await (await request.get("/api/version")).json()).toHaveProperty("version");
    await page.goto("/");
    await expect(c(page, "h1", "Log in")).toBeVisible();
  });

  test("falls back to the SPA for client routes", async ({ page, request }) => {
    const res = await request.get("/some/deep/route");
    expect(res.status()).toBe(200);
    expect(await res.text()).toContain("/_app/immutable/");
    await page.goto("/some/deep/route");
    await expect(page.getByText("404").first()).toBeVisible();
  });

  test("returns 404 for unknown API paths and files", async ({ request }) => {
    expect((await request.get("/api/nope")).status()).toBe(404);
    expect((await request.get("/.htpasswd")).status()).toBe(404);
    expect((await request.get("/users.json")).status()).toBe(404);
    expect((await request.fetch("/", { method: "FOOBAR" })).status()).toBe(404);
  });
});
