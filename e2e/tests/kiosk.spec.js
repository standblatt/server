import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Kiosk Club ${run}`;
const admin = `kiosk-${run}`;
const ids = {};

async function shoot(page, name, birthday, points, lane) {
  const base = `/api/tenants/${ids.tenant}`;
  const sh = await api(page, "POST", `${base}/shooters`, { firstname: name, lastname: `${name}son`, birthday });
  const reg = await api(page, "POST", `${base}/events/${ids.event}/registrations`, { shooterId: sh.id, lane, programmes: [{ programmeId: ids.main }] });
  if (points) await api(page, "PUT", `${base}/entries/${reg.entries[0].id}/results`, { series: points.map((p) => ({ points: p })) });
}

test.describe.configure({ mode: "serial" });

test.describe("rankings and kiosk", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 3 });
    ids.event = (await api(page, "POST", `${base}/events`, {
      rangeIds: [range.id], name: "Indoor 2026", startsOn: "2026-11-20", settings: { showLaneStatus: true },
    })).id;
    ids.main = (await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Main", showOnRankings: true,
      series: [{ shots: 10, maxPoints: 100 }, { shots: 10, maxPoints: 100 }],
      prices: [{ ageCategory: "elite", price: 10 }],
      awards: [],
    })).id;
    await shoot(page, "Anna", "1990-01-01", [90, 95]);
    await shoot(page, "Bert", "1991-01-01", [95, 95]);
    await shoot(page, "Carl", "1992-01-01", null, 2);
    await page.close();
  });

  test("shows the rankings to staff", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/rankings`);
    const first = page.locator("section[aria-label='Ranking Main'] tbody tr").first();
    await expect(first).toContainText("Bertson");
    await expect(first).toContainText("190");

    const section = page.locator("section[aria-label='Ranking Main']");
    await expect(section.locator("button", { hasText: "Print" })).toHaveCount(0);
    const pdf = c(section, "a", "PDF");
    const csv = c(section, "a", "CSV");
    await expect(pdf).toHaveAttribute("href", `/api/tenants/${ids.tenant}/programmes/${ids.main}/ranking`);
    const type = async (link) => (await page.request.get(await link.getAttribute("href"))).headers()["content-type"];
    expect(await type(pdf)).toBe("application/pdf");
    expect(await type(csv)).toContain("text/csv");

    const group = (label) => page.locator("[data-sidebar=group]").filter({ has: page.locator("[data-sidebar=group-label]", { hasText: new RegExp(`^${label}$`) }) });
    await expect(group("Master data").locator("a")).toHaveText(["Shooters", "Clubs"]);
    await expect(group("Club").locator("a")).toHaveText(["Events", "Ranges", "Users", "Print forms"]);
  });

  test("serves a live kiosk display without login", async ({ page, browser }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/setup?tab=kiosk`);
    await page.locator("input[aria-label='Display name']").fill("Hall 1");
    await c(page, "button", "Create kiosk link").click();
    const url = new URL(await page.locator("[data-testid=kiosk-url]").innerText());
    expect(url.searchParams.get("theme")).toBe("dark");
    ids.kiosk = url.pathname;

    const staff = await newPage(browser);
    await login(staff, admin);
    await page.context().clearCookies();
    await page.goto(`${ids.kiosk}?theme=dark`);
    await expect(c(page, "h1", "Indoor 2026")).toBeVisible();
    await expect(page.locator("html")).toHaveClass(/dark/);
    const rows = page.locator("section[aria-label='Ranking Main'] tbody tr");
    await expect(rows).toHaveCount(2);
    await expect(rows.first()).toContainText("1");
    await expect(rows.first()).toContainText("Bertson");
    await expect(page.locator("[aria-label='Lane 2: Carlson Carl']")).toBeVisible();
    await expect(page.locator("[title=live]")).toBeVisible();

    const regs = await api(staff, "GET", `/api/tenants/${ids.tenant}/events/${ids.event}/registrations`);
    const carl = regs.find((r) => r.shooter.firstname === "Carl");
    await api(staff, "PUT", `/api/tenants/${ids.tenant}/entries/${carl.entries[0].id}/results`, { series: [{ points: 100 }, { points: 100 }] });
    await staff.close();
    await expect(rows).toHaveCount(3);
    await expect(rows.first()).toContainText("Carlson");
    await expect(rows.first()).toContainText("200");
    await expect(page.locator("[aria-label='Lane 2: free']")).toBeVisible();
  });

  test("stops working when the kiosk link is revoked", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/setup?tab=kiosk`);
    const hall = c(page, "li", "Hall 1");
    await expect(hall).toContainText("last seen");
    await c(hall, "button", "revoke").click();
    await c(c(page, "[role=alertdialog]", "Revoke kiosk display?"), "button", "Revoke").click();
    await expect(c(page, "li", "No kiosk displays")).toBeVisible();
    await page.context().clearCookies();
    await page.goto(ids.kiosk);
    await expect(c(page, "[role=alert]", "not valid")).toBeVisible();
  });
});
