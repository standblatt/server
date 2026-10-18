import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const user = `mobile-${run}`;
const ids = {};
const statusBar = 47;

test.describe.configure({ mode: "serial" });
test.use({ viewport: { width: 390, height: 844 }, isMobile: true, hasTouch: true });

/** Top edge of a locator in CSS pixels. */
const top = async (locator) => (await locator.boundingBox()).y;

/** Replaces the share sheet with a recorder; with refuseFirst the first call fails like a share without a fresh tap. */
const recordShares = (page, refuseFirst = false) =>
  page.addInitScript((refuse) => {
    window.shared = [];
    let calls = 0;
    navigator.canShare = () => true;
    navigator.share = async ({ files }) => {
      if (refuse && calls++ === 0) throw new DOMException("needs a tap", "NotAllowedError");
      window.shared.push(`${files[0].name} ${files[0].type} ${files[0].size > 1000}`);
    };
  }, refuseFirst);

/** What scrolls sideways on the page: the page itself or a scroll container that is wider than its box, except rows marked data-swipe. */
const sideways = (page) =>
  page.evaluate(() => {
    const out = [];
    if (document.documentElement.scrollWidth > innerWidth) out.push(`page +${document.documentElement.scrollWidth - innerWidth}`);
    for (const el of document.querySelectorAll("body *")) {
      const x = getComputedStyle(el).overflowX;
      if ((x === "auto" || x === "scroll") && el.scrollWidth > el.clientWidth + 1 && !el.closest("[data-swipe]")) {
        out.push(`${el.tagName.toLowerCase()} +${el.scrollWidth - el.clientWidth}: ${(el.innerText || "").slice(0, 50).replace(/\s+/g, " ")}`);
      }
    }
    return out;
  });

test.describe("phone", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, `Mobile Club ${run}`, user);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "Pistol range 25 m", laneCount: 6 });
    const range50 = await api(page, "POST", `${base}/ranges`, { name: "Pistol range 50 m", laneCount: 6 });
    ids.range = range.id;
    const ev = await api(page, "POST", `${base}/events`, { rangeIds: [range.id, range50.id], name: "Autumn championship with a rather long name", startsOn: "2026-11-20" });
    ids.event = ev.id;
    const p = await api(page, "POST", `${base}/events/${ev.id}/programmes`, {
      name: "Supported elite programme", series: [{ shots: 10, maxPoints: 100 }, { shots: 10, maxPoints: 100 }, { shots: 10, maxPoints: 100 }],
      prices: [{ ageCategory: "elite", price: 25 }], awards: [],
    });
    ids.programme = p.id;
    ids.regs = [];
    for (let i = 0; i < 4; i++) {
      const sh = await api(page, "POST", `${base}/shooters`, {
        firstname: `Maximiliane${i}`, lastname: `Example-Longername${i}`, birthday: "1995-01-01", licenseNumber: String(810000 + i), clubNames: ["Shooting society Exampletown-Longvillage"],
      });
      ids.shooter = sh.id;
      const reg = await api(page, "POST", `${base}/events/${ev.id}/registrations`, { shooterId: sh.id, rangeId: i < 3 ? ids.range : undefined, lane: i < 3 ? i + 1 : undefined, programmes: [{ programmeId: p.id }] });
      ids.regs.push(reg.id);
      if (i < 2) await api(page, "PUT", `${base}/entries/${reg.entries[0].id}/results`, { series: [{ points: 90 + i }, { points: 88 }, { points: 87 }] });
    }
    await api(page, "POST", `${base}/ranges/${range.id}/reservations`, { lanes: "1-3", startsAt: "2026-11-20T09:00:00+01:00", endsAt: "2026-11-20T11:00:00+01:00", description: "Junior training" });
    ids.helper = (await api(page, "POST", `${base}/events/${ev.id}/helper-codes`, { name: "Results desk", role: "results" })).secret;
    ids.kiosk = (await api(page, "POST", `${base}/events/${ev.id}/kiosk-tokens`, { name: "Hall display" })).secret;
    ids.ammo = (await api(page, "POST", `${base}/events/${ev.id}/helper-codes`, { name: "Ammunition", role: "ammunition" })).secret;
    await page.close();
  });

  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem("standblatt-locale", "en"));
  });

  test("closes the sidebar after choosing a page", async ({ page }) => {
    await login(page, user);
    await page.goto(`/t/${ids.tenant}`);
    const sheet = page.locator("[data-mobile=true]");
    await page.locator("[data-slot=sidebar-trigger]").click();
    await expect(sheet).toBeVisible();
    await c(sheet, "a", "Shooters").click();
    await expect(page).toHaveURL((u) => u.pathname.endsWith("/shooters"));
    await expect(sheet).toHaveCount(0);
  });

  test("keeps content below the status bar of an installed iPhone app", async ({ page, context }) => {
    const cdp = await context.newCDPSession(page);
    await cdp.send("Emulation.setSafeAreaInsetsOverride", { insets: { top: statusBar } });
    await page.goto("/login");
    expect(await top(page.locator("input[name=name]"))).toBeGreaterThanOrEqual(statusBar);
    const strip = await page.evaluate(() => {
      const s = getComputedStyle(document.body, "::before");
      return { height: s.height, position: s.position, zIndex: s.zIndex };
    });
    expect(strip).toEqual({ height: `${statusBar}px`, position: "fixed", zIndex: "100" });

    await login(page, user);
    await page.goto(`/t/${ids.tenant}/shooters`);
    const trigger = page.locator("[data-slot=sidebar-trigger]");
    expect(await top(trigger)).toBeGreaterThanOrEqual(statusBar);
    await trigger.click();
    await expect.poll(() => top(page.locator("[data-mobile=true]"))).toBe(statusBar);
    await page.keyboard.press("Escape");

    await page.goto("/my");
    expect(await top(page.locator("header").first())).toBe(0);
    expect(await top(page.locator("header").first().locator("a, button").first())).toBeGreaterThanOrEqual(statusBar);
    await page.goto("/my/train");
    expect(await top(page.getByRole("button", { name: "Close" }))).toBeGreaterThanOrEqual(statusBar);
    await page.goto("/account");
    expect(await top(page.locator("main").locator("h1, h2").first())).toBeGreaterThanOrEqual(statusBar);
  });

  test("nothing scrolls sideways", async ({ page }) => {
    test.setTimeout(90000);
    await login(page, user);
    const t = `/t/${ids.tenant}`;
    const e = `${t}/events/${ids.event}`;
    const urls = [
      "/", "/account", t, `${t}/shooters`, `${t}/shooters?shooter=${ids.shooter}`, `${t}/clubs`, `${t}/ranges`, `${t}/ranges/${ids.range}/printing`, `${t}/users`, `${t}/forms`,
      e, `${e}?registration=${ids.regs[3]}`, `${e}?results=${ids.regs[0]}`, `${e}/entry`, `${e}/rankings`, `${e}/statistics`, `${e}/reservations`,
      ...["", "/modules", "/programmes", "/kiosk", "/helpers", "/booking", "/classifications", "/close"].map((x) => `${e}/config${x}`),
      "/my", "/my/log", "/my/progress", "/my/settings", "/my/train", "/help", "/help/user-guide/dashboard",
      `/h/${ids.helper}`, `/kiosk/${ids.kiosk}`,
    ];
    const found = {};
    for (const u of urls) {
      await page.goto(u);
      await page.waitForLoadState("networkidle");
      const s = await sideways(page);
      if (s.length) found[u] = s;
    }
    await page.goto(`${e}/config/programmes`);
    await page.getByRole("button", { name: /Actions for/ }).first().click();
    await page.getByRole("menuitem", { name: "Edit" }).click();
    await expect(page.locator("[role=dialog]")).toBeVisible();
    const s = await sideways(page);
    if (s.length) found["programme editor"] = s;
    expect(found).toEqual({});
  });

  test("reaches a configuration page from the sidebar", async ({ page }) => {
    await login(page, user);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}`);
    const sheet = page.locator("[data-mobile=true]");
    await page.locator("[data-slot=sidebar-trigger]").click();
    await c(sheet, "a", "Configuration").click();
    await expect(page).toHaveURL((u) => u.pathname.endsWith("/config"));
    await page.locator("[data-slot=sidebar-trigger]").click();
    await c(sheet, "a", "Helper links").click();
    await expect(c(page, "button", "Create helper link")).toBeVisible();
  });

  test("shares a PDF instead of opening a tab", async ({ page, context }) => {
    await recordShares(page);
    await login(page, user);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/statistics`);
    let tabs = 0;
    context.on("page", () => tabs++);
    await c(page.locator("section[aria-label=Exports]"), "a", "PDF").click();
    await expect.poll(() => page.evaluate(() => window.shared)).toEqual([expect.stringMatching(/\.pdf application\/pdf true$/)]);
    expect(tabs).toBe(0);
  });

  test("offers a share button when the phone wants another tap", async ({ page }) => {
    await recordShares(page, true);
    await login(page, user);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/statistics`);
    await c(page.locator("section[aria-label=Exports]"), "a", "PDF").click();
    const toast = page.locator("[data-sonner-toast]", { hasText: "The PDF is ready." });
    await c(toast, "button", "Share").click();
    await expect.poll(() => page.evaluate(() => window.shared)).toHaveLength(1);
  });

  test("shares a lane sheet from a helper link", async ({ page }) => {
    await recordShares(page);
    await page.goto(`/h/${ids.ammo}`);
    await c(page.locator("tr", { hasText: "Example-Longername3" }), "a", "PDF").click();
    await expect.poll(() => page.evaluate(() => window.shared)).toEqual([expect.stringMatching(/\.pdf application\/pdf true$/)]);
  });
});

test.describe("desktop helper link", () => {
  test.use({ viewport: { width: 1280, height: 800 }, isMobile: false, hasTouch: false });

  test("opens a lane sheet with the helper code", async ({ page }) => {
    test.skip(!ids.ammo, "needs the phone setup");
    await page.addInitScript(() => {
      window.opened = [];
      window.open = (url) => window.opened.push(url);
    });
    await page.goto(`/h/${ids.ammo}`);
    const pdf = page.waitForResponse((r) => r.url().includes("/documents/lane_sheet"));
    await c(page.locator("tr", { hasText: "Example-Longername3" }), "a", "PDF").click();
    expect((await pdf).status()).toBe(200);
    await expect.poll(() => page.evaluate(() => window.opened)).toEqual([expect.stringMatching(/^blob:/)]);
  });
});
