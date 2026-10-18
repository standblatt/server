import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Print Club ${run}`;
const admin = `print-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("documents and printing", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 4 });
    ids.range = range.id;
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Indoor 2026", startsOn: "2026-11-20", settings: { printAwardValue: true } })).id;
    const main = await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Main", barcode: "123456", showOnRankings: true,
      series: [{ shots: 10, maxPoints: 100 }],
      prices: [{ ageCategory: "elite", price: 10 }],
      awards: [{ name: "Wreath", kind: "card", value: 25, limits: [{ ageCategory: "elite", points: 90 }] }],
    });
    ids.main = main.id;
    const sh = await api(page, "POST", `${base}/shooters`, { licenseNumber: "300001", firstname: "Anna", lastname: "Adult", birthday: "1990-01-01" });
    const reg = await api(page, "POST", `${base}/events/${ids.event}/registrations`, { shooterId: sh.id, lane: 1, programmes: [{ programmeId: main.id }] });
    ids.reg = reg.id;
    ids.entry = reg.entries[0].id;
    await page.close();
  });

  test("serves PDFs and opens them when no printer is set up", async ({ page }) => {
    await login(page, admin);
    await page.addInitScript(() => {
      window.opened = [];
      window.open = (url) => (window.opened.push(url), null);
    });
    await page.goto(`/t/${ids.tenant}/events/${ids.event}?registration=${ids.reg}`);
    const pdf = await page.request.get(await c(page, "a", "PDF").getAttribute("href"));
    expect(pdf.headers()["content-type"]).toBe("application/pdf");
    expect((await pdf.body()).subarray(0, 5).toString()).toBe("%PDF-");
    await c(page, "button", "Print lane sheet").click();
    await expect.poll(() => page.evaluate(() => window.opened)).toContain(`/api/tenants/${ids.tenant}/registrations/${ids.reg}/documents/lane_sheet`);
  });

  test("routes documents to a print agent", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/ranges`);
    await c(c(page, "tr", "10m"), "a", "Set up the printers").click();
    await page.locator("input[aria-label='Agent name']").fill("Pi hall 1");
    await c(page, "button", "Add print agent").click();
    const agent = { Authorization: `Bearer ${await page.locator("[data-testid=agent-token]").innerText()}` };
    await page.request.post("/api/agent/hello", {
      headers: agent,
      data: { printers: [{ name: "Office_A5", options: [{ name: "InputSlot", label: "Media Source", values: ["Auto", "Tray1", "Tray2"], default: "Auto" }] }], version: "v0.1.0", os: "linux" },
    });
    await page.reload();
    const entry = c(page, "li", "Pi hall 1");
    await expect(entry).toContainText("online");
    await expect(entry).toContainText("1 printer");
    await expect(entry).toContainText("v0.1.0 · linux");
    await page.locator("select[aria-label='Printer for Lane sheet (A5)']").selectOption("Office_A5");
    const row = c(page, "tr", "Lane sheet (A5)");
    await c(row, "label", "Media Source").locator("select").selectOption("Tray2");
    await page.locator("input[aria-label='More options for Lane sheet (A5)']").fill("fit-to-page");
    await c(row, "button", "Save").click();
    await expect(c(page, "[data-sonner-toast]", "Saved.")).toBeVisible();

    await page.goto(`/t/${ids.tenant}/events/${ids.event}?registration=${ids.reg}`);
    await c(page, "button", "Print lane sheet").click();
    await expect(c(page, "[data-sonner-toast]", "sent to printer")).toBeVisible();

    const next = await page.request.get("/api/agent/jobs/next", { headers: agent });
    expect(next.status()).toBe(200);
    const job = await next.json();
    expect(job.printer).toBe("Office_A5");
    expect(job.options).toEqual({ InputSlot: "Tray2", "fit-to-page": "" });
    await page.request.post(`/api/agent/jobs/${job.id}`, { headers: agent, data: { ok: true } });
    await page.goto(`/t/${ids.tenant}/ranges/${ids.range}/printing`);
    await expect(c(page, "tr", "Lane sheet Adult Anna")).toContainText("done");

    await c(c(page, "tr", "Lane sheet (A5)"), "button", "Test page").click();
    await expect(c(page, "[data-sonner-toast]", "Test page sent to Office_A5.")).toBeVisible();
    const test = await (await page.request.get("/api/agent/jobs/next", { headers: agent })).json();
    expect(test.title).toBe("Test page Office_A5");
    expect(test.options).toEqual({ InputSlot: "Tray2", "fit-to-page": "" });
    await page.request.post(`/api/agent/jobs/${test.id}`, { headers: agent, data: { ok: true } });
    await page.reload();
    await expect(c(page, "tr", "Test page Office_A5")).toContainText("done");
  });

  test("shows the options of a Windows agent in plain words", async ({ page }) => {
    await login(page, admin);
    const { secret } = await api(page, "POST", `/api/tenants/${ids.tenant}/ranges/${ids.range}/print-agents`, { name: "Laptop" });
    await page.request.post("/api/agent/hello", {
      headers: { Authorization: `Bearer ${secret}` },
      data: {
        version: "v0.1.0",
        os: "windows",
        printers: [{ name: "Brother QL-700", options: [
          { name: "paper", label: "Paper", values: ["29mm x 90mm", "62mm x 29mm"], default: "29mm x 90mm" },
          { name: "scale", label: "Scaling", values: ["noscale", "shrink", "fit"], default: "noscale" },
        ] }],
      },
    });
    await page.goto(`/t/${ids.tenant}/ranges/${ids.range}/printing`);
    const doc = "Control strip (label 90 × 29 mm)";
    await page.locator(`select[aria-label='Printer for ${doc}']`).selectOption("Brother QL-700");
    const row = c(page, "tr", doc);
    await expect(row.locator(`input[aria-label='More options for ${doc}']`)).toHaveCount(0);
    await c(row, "label", "Scaling").locator("select").selectOption({ label: "Shrink to fit" });
    await c(row, "button", "Save").click();
    await expect(c(page, "[data-sonner-toast]", "Saved.")).toBeVisible();
    const settings = await api(page, "GET", `/api/tenants/${ids.tenant}/ranges/${ids.range}/print-settings`);
    expect(settings.find((s) => s.document === "control_strip")).toMatchObject({ printer: "Brother QL-700", options: { scale: "shrink" } });
  });

  test("connects a computer by code on /connect", async ({ page }) => {
    const pairing = await (await page.request.post("/api/agent/pair", { data: { name: "DESKTOP-E2E", os: "windows", version: "v0.1.0" } })).json();
    await page.goto(`/connect?code=${pairing.code}`);
    await expect(c(page, "a", "Log in")).toHaveAttribute("href", `/login?next=${encodeURIComponent(`/connect?code=${pairing.code}`)}`);
    await login(page, admin);
    await page.goto(`/connect?code=${pairing.code}`);
    await expect(c(page, "[role=status]", "DESKTOP-E2E (Windows, v0.1.0) wants to connect.")).toBeVisible();
    await page.locator("select[name=range]").selectOption({ label: "10m" });
    await expect(page.locator("input[name=name]")).toHaveValue("DESKTOP-E2E");
    await c(page, "button", "Connect").click();
    await expect(c(page, "[role=status]", "Connected.")).toBeVisible();
    const status = await (await page.request.post("/api/agent/pair/status", { data: { secret: pairing.secret } })).json();
    expect(status).toMatchObject({ status: "connected", club, range: "10m" });
    await page.request.post("/api/agent/hello", { headers: { Authorization: `Bearer ${status.token}` }, data: { printers: [], version: "v0.1.0", os: "windows" } });
    await c(page, "a", "Set up the printers").click();
    await expect(c(page, "li", "DESKTOP-E2E")).toContainText("v0.1.0 · windows");
  });

  test("offers the wreath card and result sheet once results are in", async ({ page }) => {
    await login(page, admin);
    await api(page, "PUT", `/api/tenants/${ids.tenant}/entries/${ids.entry}/results`, { series: [{ points: 95 }] });
    await page.goto(`/t/${ids.tenant}/events/${ids.event}?results=${ids.reg}`);
    await expect(c(page, "legend", "Award: Wreath")).toBeVisible();
    await expect(c(page, "a", "PDF")).toBeVisible();
    await expect(c(page.locator("section[aria-label=Main]"), "button", "Print wreath card")).toBeVisible();
    await expect(c(c(page, "[role=status]", "All programmes have results"), "button", "Print result sheet")).toBeVisible();
    const type = async (url) => (await page.request.get(url)).headers()["content-type"];
    expect(await type(`/api/tenants/${ids.tenant}/entries/${ids.entry}/wreath-card`)).toBe("application/pdf");
    expect(await type(`/api/tenants/${ids.tenant}/programmes/${ids.main}/ranking`)).toBe("application/pdf");
  });
});
