// Screenshots for the documentation (docs/content/images/screenshots). Not part of `make e2e`;
// `make docs-screenshots` starts an empty instance and writes the files into $SHOTS.
import { test, expect } from "@playwright/test";
import { adminLogin, api, c, createClub, login, mailToken, newPage } from "../helpers.js";

const shots = process.env.SHOTS || "screenshots";
const admin = "anna.admin";
const password = "docs-password-123";
const club = "PS Example";
const ids = {};

// Invented people and clubs only: screenshots end up in the public documentation.
const shooters = [
  ["Example", "Felix", "1986-09-02", "SG Hilltown", [94, 95, 96, 93]],
  ["Sample", "Martin", "1964-02-24", "PC Lakeside", [96, 93, 97, 88]],
  ["Muster", "Peter", "1978-04-28", "PC Lakeside", [93, 94, 93, 90]],
  ["Dummy", "Walter", "1949-11-05", "SV Riverbank", [91, 93, 87, 92]],
  ["Placeholder", "Rosa", "1948-03-12", "SV Riverbank", [86, 90, 88, 89]],
  ["Tester", "Lea", "2010-05-20", "SG Hilltown", [91, 90, 90, 89]],
  ["Demo", "Marco", "1968-07-01", "CP Valley", [88, 93, 92, 90]],
  ["Specimen", "Jonas", "1993-01-15", "PS Example", []],
  ["Mock", "Nina", "1987-06-30", "PS Example", []],
  ["Novice", "Hans", "1962-08-08", "PS Example", null],
];

async function shot(page, name) {
  await expect(page.locator("[data-sonner-toast]")).toHaveCount(0, { timeout: 15000 });
  await page.waitForTimeout(300);
  await page.screenshot({ path: `${shots}/${name}.png` });
}

async function visit(page, path, dark = false) {
  await page.emulateMedia({ colorScheme: dark ? "dark" : "light" });
  await page.goto(path);
}

test.describe.configure({ mode: "serial" });

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem("standblatt-locale", "en");
    localStorage.setItem("mode-watcher-mode", matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light");
  });
});

test.describe("documentation screenshots", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin, password);

    const base = `/api/tenants/${ids.tenant}`;
    await api(page, "POST", `${base}/users`, { name: "sam.staff", password, role: "staff" });
    const clubs = {};
    for (const [name, ssv] of [["PS Example", "1.99.0.00.001"], ["SG Hilltown", "1.99.0.00.002"], ["PC Lakeside", "1.99.0.00.003"], ["SV Riverbank", "1.99.0.00.004"], ["CP Valley", "1.99.0.00.005"]]) {
      clubs[name] = (await api(page, "POST", `${base}/clubs`, { name, ssvNumber: ssv })).id;
    }
    const range = await api(page, "POST", `${base}/ranges`, { name: "Air pistol 10m", laneCount: 12, disabledLanes: [12] });
    ids.range = range.id;
    await api(page, "POST", `${base}/ranges/${range.id}/print-agents`, { name: "Range PC" });
    const ev = await api(page, "POST", `${base}/events`, {
      rangeIds: [range.id], name: "42. Example Championship 2027", startsOn: "2027-01-08", organizer: "PS Example", locale: "en",
      settings: { showLaneStatus: true },
    });
    ids.event = ev.id;
    const evBase = `${base}/events/${ev.id}`;
    const prices = ["u21", "elite", "senior", "veteran", "senior_veteran"].map((cat) => ({ ageCategory: cat, price: cat === "u21" ? 10 : 20 }));
    prices.push({ ageCategory: "u17", price: 10 }, { ageCategory: "u19", price: 10 });
    const limits = (top) => ["u17", "u19", "u21", "elite", "senior", "veteran", "senior_veteran"].map((cat, i) => ({ ageCategory: cat, points: top - [20, 10, 10, 0, 0, 10, 20][i] }));
    const main = await api(page, "POST", `${evBase}/programmes`, {
      name: "Championship", variant: "Sport", barcode: "700005002334864165956154", countsForClubCompetition: true, showOnRankings: true,
      series: [1, 2, 3, 4].map(() => ({ shots: 10, maxPoints: 100 })), prices,
      awards: [
        { name: "Wreath card CHF 10", kind: "card", value: 10, limits: limits(370) },
        { name: "Wreath card CHF 8", kind: "card", value: 8, limits: limits(360) },
        { name: "Wreath card CHF 6", kind: "card", value: 6, limits: limits(350) },
      ],
    });
    ids.main = main.id;
    await api(page, "PUT", `${base}/programmes/${main.id}/special-prizes`, { ranks: 3, minCategory: 2, minClub: 2 });
    const bonus = await api(page, "POST", `${evBase}/programmes`, {
      name: "Bonus", variant: "Sport", showOnRankings: true, series: [1, 2].map(() => ({ shots: 10, maxPoints: 100 })), prices,
      awards: [{ name: "Prize in kind", kind: "natural", value: 0, limits: limits(175) }],
    });
    ids.bonus = bonus.id;
    await api(page, "PUT", `${base}/programmes/${bonus.id}/combination`, { combinableWithId: main.id });
    const start = new Date(Date.now() + 30 * 60000);
    await api(page, "POST", `${base}/ranges/${range.id}/reservations`, {
      lanes: "9-11", startsAt: start.toISOString(), endsAt: new Date(start.getTime() + 2 * 3600000).toISOString(), leadMinutes: 60, description: "Junior training",
    });
    ids.kiosk = (await api(page, "POST", `${evBase}/kiosk-tokens`, { name: "Hall screen" })).secret;

    for (const [i, [lastname, firstname, birthday, clubName, series]] of shooters.entries()) {
      const sh = await api(page, "POST", `${base}/shooters`, { lastname, firstname, birthday, licenseNumber: String(100100 + i), clubIds: [clubs[clubName]] });
      if (series === null) continue;
      const lane = series.length ? undefined : i - 6;
      const programmes = [{ programmeId: ids.main }];
      if (i % 2 === 0) programmes.push({ programmeId: ids.bonus, combinedSeries: [2, 3] });
      const reg = await api(page, "POST", `${evBase}/registrations`, { shooterId: sh.id, clubId: clubs[clubName], lane, programmes });
      if (!series.length) continue;
      if (i === 1) ids.registration = reg.id;
      const entry = reg.entries.find((e) => e.programmeId === ids.main);
      await api(page, "PUT", `${base}/entries/${entry.id}/results`, { series: series.map((points) => ({ points })) });
    }
    await page.close();
  });

  test.beforeEach(({ page }) => login(page, admin, password));

  test("login and events", async ({ page }) => {
    await page.context().clearCookies();
    await visit(page, "/login");
    await expect(c(page, "h1", "Log in")).toBeVisible();
    await shot(page, "login");
    await login(page, admin, password);
    await visit(page, "/");
    await expect(c(page, "h1", "Where to?")).toBeVisible();
    await shot(page, "home");
    await visit(page, "/account");
    await expect(c(page, "h1", "Account")).toBeVisible();
    await shot(page, "account");
    await visit(page, `/t/${ids.tenant}`);
    await expect(c(page, "td", "42. Example Championship 2027")).toBeVisible();
    await shot(page, "events");
  });

  test("dashboard and dialogs", async ({ page }) => {
    const event = `/t/${ids.tenant}/events/${ids.event}`;
    await visit(page, event);
    await expect(c(page, "[aria-label='Live updates']", "live")).toBeVisible();
    await expect(c(page, "h2", "Shooting (2)")).toBeVisible();
    await shot(page, "dashboard");
    await visit(page, event, true);
    await expect(c(page, "h2", "Shooting (2)")).toBeVisible();
    await shot(page, "dashboard-dark");

    await visit(page, event);
    await page.locator("input[aria-label='Search shooter']").fill("novice");
    await expect(c(page, "li", "Novice Hans")).toBeVisible();
    await shot(page, "search");
    await c(page, "button", "Register for PS Example").click();
    await c(page, "label", "Championship").locator("[role=checkbox]").click();
    await c(page, "label", "Bonus").locator("[role=checkbox]").click();
    await c(page, "label", "take 2 series from Championship (Sport)").locator("[role=checkbox]").click();
    const series = page.locator("[role=group][aria-label='Series of Championship (Sport)']");
    await c(series, "button", /^1$/).click();
    await c(series, "button", /^4$/).click();
    await c(page, "button", /^3$/).click();
    await shot(page, "registration");

    await visit(page, `${event}?results=${ids.registration}`);
    await expect(c(page, "[role=dialog]", "Results: Sample Martin")).toBeVisible();
    await shot(page, "results");
  });

  test("rankings, kiosk, statistics, reservations", async ({ page }) => {
    const event = `/t/${ids.tenant}/events/${ids.event}`;
    await visit(page, `${event}/rankings`);
    await expect(c(page, "section", "Championship (Sport)").locator("tbody tr")).toHaveCount(7);
    await shot(page, "rankings");
    await visit(page, `${event}/statistics`);
    await expect(c(page, "h1", "Statistics")).toBeVisible();
    await expect(page.getByText("Special prizes").first()).toBeVisible();
    await shot(page, "statistics");
    await visit(page, `${event}/reservations`);
    await expect(page.getByText("Junior training").first()).toBeVisible();
    await shot(page, "reservations");
    await page.context().clearCookies();
    await visit(page, `/kiosk/${ids.kiosk}?theme=dark`, true);
    await expect(c(page, "h1", "42. Example Championship 2027")).toBeVisible();
    await expect(page.locator("[title=live]")).toBeVisible();
    await shot(page, "kiosk");
  });

  test("helper link", async ({ page }) => {
    const helper = await api(page, "POST", `/api/tenants/${ids.tenant}/events/${ids.event}/helper-codes`, { name: "Scorer lanes 1-4", role: "scorer", rangeId: ids.range, lanes: "1-4" });
    await page.context().clearCookies();
    await page.setViewportSize({ width: 390, height: 844 });
    await visit(page, `/h/${helper.secret}`);
    await expect(c(page, "h1", "42. Example Championship 2027")).toBeVisible();
    await expect(page.locator("[aria-label^='Lane 1:']")).toBeVisible();
    await shot(page, "helper-scorer");
  });

  test("setup, programmes and special prizes", async ({ page }) => {
    await visit(page, `/t/${ids.tenant}/events/${ids.event}/config/programmes`);
    await expect(c(page, "td", "Championship")).toBeVisible();
    await shot(page, "setup");
    const actions = page.locator("button[aria-label='Actions for Championship (Sport)']");
    await actions.click();
    await c(page, "[role=menuitem]", "Special prizes").click();
    await expect(c(page, "[role=dialog]", "Special prizes: Championship (Sport)")).toBeVisible();
    await shot(page, "special-prizes");
    await page.keyboard.press("Escape");
    await actions.click();
    await c(page, "[role=menuitem]", "Edit").click();
    await expect(c(page, "[role=dialog]", "Edit programme")).toBeVisible();
    await shot(page, "programme-editor");
    await page.keyboard.press("Escape");
    await visit(page, `/t/${ids.tenant}/events/${ids.event}/config/close`);
    await c(c(page, "section", "Copy for next year"), "button", "Copy").click();
    await expect(page.locator("[role=dialog] input[name=name]")).toHaveValue("43. Example Championship 2028");
    await shot(page, "copy-event");
  });

  test("online booking", async ({ page }) => {
    const evBase = `/api/tenants/${ids.tenant}/events/${ids.event}`;
    const cfg = await api(page, "POST", `${evBase}/booking/code`);
    const { slots } = await api(page, "PUT", `${evBase}/booking`, { from: "08:00", to: "17:00", breakFrom: "12:00", breakTo: "13:00", slotMinutes: 60 });
    const code = new URL(cfg.link).pathname.split("/").pop();
    const ask = (club, contact, email, slot, lanes, shooters) =>
      api(page, "POST", `/api/booking/${code}/requests`, {
        club, contact, email, items: [{ rangeId: ids.range, start: slots[slot].start, end: slots[slot].end, lanes }], shooters,
      });
    await ask("SG Hilltown", "Lea Tester", "lea@example.org", 1, 4, [
      { lastname: "Hilltown", firstname: "Ben", birthday: "1990-03-03" },
      { lastname: "Example", firstname: "Felix", birthday: "1986-09-02", licenseNumber: "100100" },
    ]);
    await ask("CP Valley", "Marco Demo", "marco@example.org", 2, 3, [{ lastname: "Demo", firstname: "Marco", birthday: "1968-07-01" }]);
    const requests = await api(page, "GET", `${evBase}/booking-requests`);
    await api(page, "POST", `/api/tenants/${ids.tenant}/booking-requests/${requests.find((r) => r.club === "SG Hilltown").id}/confirm`);

    await visit(page, `/t/${ids.tenant}/events/${ids.event}/config/booking`);
    await expect(page.getByLabel("Booking link")).toHaveValue(cfg.link);
    await shot(page, "booking-setup");
    await visit(page, `/t/${ids.tenant}/events/${ids.event}/reservations`);
    await expect(page.locator("[aria-label='Request of CP Valley']")).toBeVisible();
    await shot(page, "booking-requests");
    await page.context().clearCookies();
    await visit(page, `/b/${code}`);
    await page.getByLabel("Lanes on Air pistol 10m at 09:00").fill("2");
    await c(page, "label", "Club").locator("input").fill("SV Riverbank");
    await c(page, "label", "Contact person").locator("input").fill("Rosa Placeholder");
    await c(page, "label", "Email").locator("input").fill("rosa@example.org");
    await page.evaluate(() => window.scrollTo(0, 0));
    await shot(page, "booking-public");
  });

  test("club pages and admin panel", async ({ page }) => {
    const pages = [
      ["shooters", "Example"],
      ["clubs", "SG Hilltown"],
      ["ranges", "Air pistol 10m"],
      [`ranges/${ids.range}/printing`, "Range PC", "printing"],
      ["users", admin],
    ];
    for (const [path, text, name = path] of pages) {
      await visit(page, `/t/${ids.tenant}/${path}`);
      await expect(page.getByText(text).first()).toBeVisible();
      await shot(page, name);
    }
    await page.context().clearCookies();
    await adminLogin(page);
    await expect(c(page, "h1", "Administration")).toBeVisible();
    await shot(page, "admin-panel");
    await page.goto("/admin/users");
    await c(page, "td a", admin).click();
    await expect(c(page, "h1", admin)).toBeVisible();
    await shot(page, "admin-user");
  });
});

// The personal training log, on a phone. Invented shooter, generated sessions.
test.describe("my training screenshots", () => {
  const email = "lea.tester@example.org";

  async function phone(page, path, name, ready) {
    await page.setViewportSize({ width: 390, height: 844 });
    await visit(page, path);
    await ready();
    await shot(page, name);
  }

  const scrollTo = (locator, offset) => locator.evaluate((el, top) => window.scrollTo(0, el.getBoundingClientRect().top + window.scrollY - top), offset);

  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    await api(page, "POST", "/api/auth/signup", { email, password, firstname: "Lea", lastname: "Tester", locale: "en" });
    await api(page, "POST", "/api/auth/verify", { token: await mailToken(page, email) });
    const weapons = [(await api(page, "POST", "/api/my/weapons", { name: "Steyr EVO 10" })).id];
    weapons.push((await api(page, "POST", "/api/my/weapons", { name: "Walther LP500" })).id);
    const my = await api(page, "GET", "/api/my");
    const pistol = my.disciplines.find((d) => d.key === "air_pistol").id;
    let seed = 7;
    const r = () => (seed = (seed * 16807) % 2147483647) / 2147483647;
    const now = Date.now();
    for (let i = 0; i < 26; i++) {
      const series = [];
      for (let k = 0; k < 5; k++) {
        if (i < 20) {
          series.push({ shots: 10, total: Math.round(90 + i * 0.12 + (r() - 0.5) * 5 - k * 0.5) });
          continue;
        }
        const values = [], x = [], y = [];
        for (let j = 0; j < 10; j++) {
          const dx = (r() - 0.42) * 20, dy = (r() - 0.6) * 20;
          values.push(Math.max(0, Math.min(10, Math.floor(10 + (8 - Math.hypot(dx, dy)) / 8))));
          x.push(Math.round(dx * 100) / 100);
          y.push(Math.round(dy * 100) / 100);
        }
        series.push({ shots: 10, total: 0, values, x, y });
      }
      const tags = i % 4 === 0 ? ["tired"] : i % 3 === 0 ? ["new grip"] : [];
      await api(page, "POST", "/api/my/trainings", {
        disciplineId: pistol, kind: i % 7 === 6 ? "competition" : "training", startedAt: new Date(now - (64 - i * 2.4) * 86400000).toISOString(),
        series, tags, rating: 2 + (i % 4), place: "Hilltown 10 m", weaponId: weapons[i % 3 === 0 ? 1 : 0],
      });
    }
    await api(page, "PUT", `/api/my/goals/${pistol}`, { targetTotal: 470 });
    ids.mySession = (await api(page, "GET", "/api/my/trainings?limit=1"))[0].id;
    await page.close();
  });

  test.beforeEach(({ page }) => login(page, email, password));

  test("today, progress, session", async ({ page }) => {
    await phone(page, "/my", "my-today", () => expect(c(page, "h2", "Recent sessions")).toBeVisible());
    await phone(page, "/my/progress", "my-progress", () => expect(c(page, "h2", "Trend")).toBeVisible());
    await phone(page, "/my/progress", "my-weapons", () => scrollTo(page.locator("[data-testid=weapons]"), 80));
    await phone(page, `/my/log/${ids.mySession}`, "my-session", () => scrollTo(c(page, "h2", "Your group"), 420));
  });

  test("logging a session", async ({ page }) => {
    await phone(page, "/my/train", "my-train-series", async () => {
      for (const k of ["9", "4", "9", "6", "9", "3", "9"]) await page.locator("[data-testid=keypad]").locator(`button[aria-label="${k}"]`).click();
    });
    await c(page, "[role=tab]", "Target").click();
    const target = page.locator("[aria-label='Tap where the shot hit']");
    const { width } = await target.boundingBox();
    for (const [fx, fy] of [[0.52, 0.47], [0.46, 0.55], [0.55, 0.56], [0.5, 0.6], [0.58, 0.5]]) await target.click({ position: { x: width * fx, y: width * fy } });
    await shot(page, "my-train-target");
    await c(page, "button", "Finish").click();
    await c(page.locator("[role=dialog]"), "button", "Save session").click();
    await expect(page.locator("[data-testid=result-total]")).toBeVisible();
    await shot(page, "my-result");
  });

  test("sign-up", async ({ page }) => {
    await page.context().clearCookies();
    await phone(page, "/signup", "my-signup", () => expect(c(page, "h1", "Create an account")).toBeVisible());
  });
});
