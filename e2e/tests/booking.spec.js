import { test, expect } from "@playwright/test";
import { api, baseURL, c, createClub, login, mailpit, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Booking Host ${run}`;
const admin = `book-${run}`;
const guest = `SV Gast ${run}`;
const email = `gast-${run}@example.org`;
const ids = {};
const day = new Date(Date.now() + 30 * 86400000).toISOString().slice(0, 10);

test.describe.configure({ mode: "serial" });

/** The manage link of the newest booking mail to an address. */
async function manageLink(page, to) {
  let id;
  await expect
    .poll(async () => {
      const res = await page.request.get(`${mailpit}/api/v1/search?query=${encodeURIComponent(`to:${to}`)}`);
      id = (await res.json()).messages.find((m) => m.Subject.includes("lane request"))?.ID;
      return id;
    })
    .toBeTruthy();
  const msg = await (await page.request.get(`${mailpit}/api/v1/message/${id}`)).json();
  return msg.Text.match(/\/b\/manage\/(bk_[A-Za-z0-9_-]+)/)[0];
}

test.describe("online booking", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 6 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Lupi Cup", startsOn: day, locale: "en" })).id;
    await page.close();
  });

  test("admin switches online booking on", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/config/booking`);
    await c(page, "button", "Switch on").click();
    const link = page.getByLabel("Booking link");
    await expect(link).toHaveValue(/\/b\/[A-Za-z0-9_-]+$/);
    ids.link = new URL(await link.inputValue()).pathname;
    await expect(c(page, "section", "Slots")).toContainText("18 slots");
  });

  test("a club asks for lanes and announces a shooter", async ({ browser }) => {
    const page = await browser.newPage({ baseURL });
    await page.goto(ids.link);
    await expect(c(page, "h1", "Lupi Cup")).toBeVisible();
    await page.getByLabel("Lanes on 10m at 08:00").fill("2");
    await c(page, "label", "Club").locator("input").fill(guest);
    await c(page, "label", "Contact person").locator("input").fill("Gina Gast");
    await c(page, "label", "Email").locator("input").fill(email);
    await c(page, "button", "Add shooter").click();
    await page.getByLabel("Last name of shooter 1").fill("Gast");
    await page.getByLabel("First name of shooter 1").fill("Gino");
    await page.getByLabel("Birthday of shooter 1").fill("1985-04-03");
    await c(page, "button", "Send request").click();
    await expect(c(page, "[role=status]", "Request received")).toBeVisible();
    await page.close();
  });

  test("admin confirms the request", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/reservations`);
    await expect(page.getByLabel("1 open request")).toBeVisible();
    const request = page.locator(`[aria-label='Request of ${guest}']`);
    await expect(request).toContainText("waiting for confirmation");
    await c(request, "summary", "1 shooter announced").click();
    await expect(request).toContainText("new");
    await api(page, "POST", `/api/tenants/${ids.tenant}/clubs`, { name: `${guest} 1920` });
    await c(request, "button", "Confirm").click();
    const dialog = page.getByRole("dialog");
    await dialog.getByLabel("Club").selectOption({ label: `${guest} 1920` });
    await c(dialog, "button", "Confirm").click();
    await expect(request).toContainText("confirmed");
    await expect(request).toContainText(`Club in the master data: ${guest} 1920`);
    await expect(request).toContainText("lanes 1, 2");
    await expect(page.locator(`[aria-label='Reservation ${guest}']`)).toContainText("lanes 1, 2");
    await page.goto(`/t/${ids.tenant}/shooters`);
    await page.getByPlaceholder("License number or name").fill("Gino Gast");
    await expect(c(page, "tr", "Gast")).toContainText(`${guest} 1920`);
  });

  test("the club cancels through its manage link", async ({ browser }) => {
    const page = await browser.newPage({ baseURL });
    await page.goto(await manageLink(page, email));
    await expect(c(page, "section", "Your request")).toContainText("confirmed");
    await c(page, "button", "Cancel request").click();
    await c(page, "[role=alertdialog] button", "Cancel request").click();
    await expect(c(page, "section", "Your request")).toContainText("cancelled");
    await page.close();
  });
});
