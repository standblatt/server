import { test, expect } from "@playwright/test";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Programme Club ${run}`;
const admin = `prog-${run}`;
const ids = {};

const exportFile = {
  format: "standblatt.programme",
  version: 1,
  programme: {
    name: "Obligatory",
    variant: "",
    series: [{ shots: 20, maxPoints: 100 }],
    prices: [{ ageCategory: "elite", price: 0 }],
    awards: [],
  },
};

const json = (name, value) => ({ name, mimeType: "application/json", buffer: Buffer.from(JSON.stringify(value)) });

test.describe.configure({ mode: "serial" });

test.describe("programmes", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    await page.goto("/admin/library");
    await page
      .locator("input[type=file][name=file]")
      .setInputFiles(json("library.programme.json", { ...exportFile, programme: { ...exportFile.programme, name: `Library ${run}` } }));
    await c(page, "button", "Upload template").click();
    await expect(c(page, "td", `Library ${run}`)).toBeVisible();

    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "10m", laneCount: 10 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Indoor 2026", startsOn: "2026-11-20" })).id;
    ids.next = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Indoor 2027", startsOn: "2027-11-20" })).id;
    await page.close();
  });

  test("creates a programme in the editor", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/setup`);
    await c(page, "button", "New programme").click();
    const dialog = page.locator("[role=dialog]");
    await c(dialog, "label", "Name").locator("input").fill("Field match");
    await c(dialog, "label", "Variant").locator("input").fill("Stgw 90");
    await c(dialog, "button", "Add series").click();
    await expect(c(dialog, "h3", "Series (20 shots, max. 200 points)")).toBeVisible();
    await c(dialog, "label", "Elite").locator("[role=checkbox]").click();
    await dialog.locator("input[aria-label='Price Elite']").fill("20");
    await c(dialog, "button", "Add award").click();
    const award = c(dialog, "fieldset", "Award 1");
    await c(award, "label", "Name").locator("input").fill("Wreath");
    await c(award, "label", "Value (CHF)").locator("input").fill("25");
    const limit = c(award, "label", "Elite").locator("input");
    await limit.fill("250");
    expect(await limit.evaluate((el) => el.checkValidity()), "limit above the maximum is invalid").toBe(false);
    await limit.fill("180");
    await c(dialog, "button", "Create").click();
    await expect(dialog).toHaveCount(0);
    const row = c(page, "tr", "Field match (Stgw 90)");
    await expect(row).toContainText("20 / 200 pts");
    await expect(row).toContainText("CHF 20.00 (Elite)");
    await expect(row).toContainText("Wreath");
    await row.locator("[aria-label=Active]").click();
    await expect(row).toHaveClass(/opacity-50/);

    await c(page, "button", "New programme").click();
    await c(dialog, "label", "Name").locator("input").fill("Field match");
    const hint = page.locator("[data-testid=variant-hint]");
    await expect(hint).toContainText('Variant of "Field match"');
    await expect(hint).toContainText("Give each variant its own variant name");
    await c(dialog, "label", "Variant").locator("input").fill("Stgw 57");
    await expect(hint).not.toContainText("Give each variant");
    await page.keyboard.press("Escape");
  });

  test("imports, copies and adds library programmes", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/setup`);
    await c(page, "button", "Add from …").click();
    await c(page, "[role=menuitem]", "Programme file").click();
    await page.locator("input[aria-label='Programme file']").setInputFiles(json("obligatory.programme.json", exportFile));
    await c(page, "button", "Import programme file").click();
    await expect(c(page, "[data-sonner-toast]", "Imported Obligatory")).toBeVisible();
    await page.locator("[aria-label='Actions for Obligatory']").click();
    await expect(c(page, "[role=menuitem]", "Export")).toHaveAttribute("href", /\/export$/);
    await page.keyboard.press("Escape");

    await page.goto(`/t/${ids.tenant}/events/${ids.next}/setup`);
    await c(page, "button", "Add from …").click();
    await c(page, "[role=menuitem]", "Earlier event").click();
    await page.locator("select[aria-label='Copy from event']").selectOption({ label: "Indoor 2026 (2026-11-20)" });
    await c(page, "button", "Copy all programmes").click();
    await expect(c(page, "[data-sonner-toast]", "Copied 2 programmes.")).toBeVisible();
    await c(page, "button", "Add from …").click();
    await c(page, "[role=menuitem]", "Programme library").click();
    await page.locator("select[aria-label='Library template']").selectOption({ label: `Library ${run}` });
    await c(page, "button", "Add from library").click();
    await expect(c(page, "[data-sonner-toast]", `Added Library ${run}`)).toBeVisible();
    await expect(c(page, "tr", "Field match (Stgw 90)")).toBeVisible();
    await expect(c(page, "tr", `Library ${run}`)).toBeVisible();
  });
});
