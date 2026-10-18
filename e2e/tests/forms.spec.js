import { test, expect } from "@playwright/test";
import { readFile } from "node:fs/promises";
import { api, c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Forms Club ${run}`;
const admin = `forms-${run}`;
const ids = {};

test.describe.configure({ mode: "serial" });

test.describe("print forms", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    ids.tenant = await createClub(page, club, admin);
    const base = `/api/tenants/${ids.tenant}`;
    const range = await api(page, "POST", `${base}/ranges`, { name: "300m", laneCount: 4 });
    ids.event = (await api(page, "POST", `${base}/events`, { rangeIds: [range.id], name: "Obligatory 2027", startsOn: "2027-05-12" })).id;
    ids.programme = (await api(page, "POST", `${base}/events/${ids.event}/programmes`, {
      name: "Obligatory", showOnRankings: true, series: [{ shots: 3, maxPoints: 15, trackShots: true }], prices: [{ ageCategory: "elite", price: 0 }], awards: [],
    })).id;
    const sh = await api(page, "POST", `${base}/shooters`, { firstname: "Fritz", lastname: "Form", birthday: "1990-01-01" });
    const reg = await api(page, "POST", `${base}/events/${ids.event}/registrations`, { shooterId: sh.id, lane: 1, programmes: [{ programmeId: ids.programme }] });
    ids.registration = reg.id;
    ids.entry = reg.entries[0].id;
    await api(page, "PUT", `${base}/entries/${ids.entry}/results`, { series: [{ points: 0, shots: [5, 4, 5] }] });
    await page.close();
  });

  test("builds a print form from the fields palette", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}`);
    await page.locator("[data-sidebar=sidebar]").getByRole("link", { name: "Print forms" }).click();
    await c(page, "button", "New print form").click();
    const dialog = page.getByRole("dialog");
    await dialog.getByLabel("Name").fill("Army sheet");
    await dialog.getByLabel("Paper size").selectOption("a5");
    await c(dialog, "button", "Create").click();
    await page.waitForURL(/\/forms\/\d+$/);
    ids.layout = Number(page.url().split("/").pop());

    const fields = page.getByRole("complementary", { name: "Fields" });
    const sheet = page.getByRole("application", { name: "Sheet" });
    await c(fields, "button", /^Name/).click();
    await expect(sheet).toContainText("Example Anna");
    await c(fields, "button", "Shots of series 1").click();
    await expect(sheet.locator("[data-copy]")).toHaveCount(11);
    await page.getByRole("switch").click();
    await expect(sheet).toContainText("{{shooter.name}}");
    await page.getByRole("switch").click();

    await page.getByRole("button", { name: /^Save/ }).click();
    await expect(c(page, "[data-sonner-toast]", "Print form saved")).toBeVisible();
    const layout = await api(page, "GET", `/api/tenants/${ids.tenant}/layouts/${ids.layout}`);
    expect(layout.elements).toMatchObject([{ value: "{{shooter.name}}", x: 10, y: 10 }, { value: "{{s1.shot#}}", repeat: 10, dx: 8, y: 18 }]);
    const preview = await page.request.get(`/api/tenants/${ids.tenant}/layouts/${ids.layout}/preview`);
    expect(preview.headers()["content-type"]).toBe("application/pdf");
  });

  test("drops, snaps, selects, aligns and undoes on the sheet", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/forms/${ids.layout}`);
    const sheet = page.getByRole("application", { name: "Sheet" });
    const props = page.getByRole("complementary", { name: "Properties" });
    await expect(sheet).toContainText("Example Anna");
    let box = await sheet.boundingBox();
    const mm = box.width / 148;

    await c(page.getByRole("complementary", { name: "Fields" }), "button", /^Club/).dragTo(sheet, { targetPosition: { x: 60 * mm, y: 60 * mm } });
    await expect(props.getByLabel("x (mm)", { exact: true })).toHaveValue("60");
    await expect(props.getByLabel("y (mm)", { exact: true })).toHaveValue("60");

    const club = sheet.locator("[data-id]").nth(2);
    const from = await club.boundingBox();
    await page.mouse.move(from.x + 2, from.y + 2);
    await page.mouse.down();
    await page.mouse.move(from.x + 2 - 49.6 * mm, from.y + 2, { steps: 6 });
    await expect(sheet.locator("line[stroke='#db2777']").first()).toBeAttached();
    await page.mouse.up();
    await expect(props.getByLabel("x (mm)", { exact: true })).toHaveValue("10");

    box = await sheet.boundingBox();
    await page.mouse.move(box.x + 100 * mm, box.y + 5 * mm);
    await page.mouse.down();
    await page.mouse.move(box.x + 3 * mm, box.y + 70 * mm, { steps: 6 });
    await page.mouse.up();
    await expect(props).toContainText("3 elements selected");

    await page.keyboard.press("Escape");
    await sheet.locator("[data-id]").nth(0).click();
    await sheet.locator("[data-id]").nth(2).click({ modifiers: ["Shift"] });
    await expect(props).toContainText("2 elements selected");
    await sheet.locator("[data-id]").nth(2).click({ button: "right" });
    await page.getByRole("menuitem", { name: "Align" }).hover();
    await page.getByRole("menuitem", { name: "Align top edges" }).click();
    await sheet.locator("[data-id]").nth(2).click();
    await expect(props.getByLabel("y (mm)", { exact: true })).toHaveValue("10");
    await page.keyboard.press("Control+z");
    await expect(props.getByLabel("y (mm)", { exact: true })).toHaveValue("60");

    await page.locator("[data-sidebar=sidebar]").getByRole("link", { name: "Events" }).click();
    await expect(c(page, "[role=alertdialog]", "Leave without saving?")).toBeVisible();
    await c(page, "[role=alertdialog] button", "Cancel").click();
    await expect(page).toHaveURL(/\/forms\/\d+$/);
    await page.keyboard.press("Control+s");
    await expect(c(page, "[data-sonner-toast]", "Print form saved")).toBeVisible();
    const layout = await api(page, "GET", `/api/tenants/${ids.tenant}/layouts/${ids.layout}`);
    expect(layout.elements).toHaveLength(3);
    expect(layout.elements[2]).toMatchObject({ value: "{{club}}", x: 10, y: 60 });
  });

  test("saves changes every minute unless autosave is off", async ({ page }) => {
    await login(page, admin);
    await page.clock.install();
    await page.goto(`/t/${ids.tenant}/forms/${ids.layout}`);
    const sheet = page.getByRole("application", { name: "Sheet" });
    const x = async () => (await api(page, "GET", `/api/tenants/${ids.tenant}/layouts/${ids.layout}`)).elements[0].x;
    await expect(sheet).toContainText("Example Anna");
    await sheet.locator("[data-id]").nth(0).click();
    await page.keyboard.press("ArrowRight");
    await page.clock.fastForward(61_000);
    await expect.poll(x).toBe(10.5);
    await expect(page.getByRole("button", { name: "Save", exact: true })).toBeVisible();

    await page.getByRole("checkbox", { name: "Save automatically" }).click();
    await sheet.locator("[data-id]").nth(0).click();
    await page.keyboard.press("ArrowRight");
    await page.clock.fastForward(61_000);
    await expect(page.getByRole("button", { name: "Save •" })).toBeVisible();
    expect(await x()).toBe(10.5);
    await page.keyboard.press("Control+s");
    await expect.poll(x).toBe(11);
    await page.reload();
    await expect(page.getByRole("checkbox", { name: "Save automatically" })).not.toBeChecked();
  });

  test("prints the results onto the form of the programme", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/config/programmes`);
    await page.getByLabel("Print form").selectOption({ label: "Army sheet" });
    await expect.poll(async () => (await api(page, "GET", `/api/tenants/${ids.tenant}/events/${ids.event}/programmes`))[0].layoutId).toBe(ids.layout);

    await page.goto(`/t/${ids.tenant}/events/${ids.event}?results=${ids.registration}`);
    const print = c(page.locator("section[aria-label=Obligatory]"), "button", "Print onto form");
    await expect(print).toBeVisible();
    const form = await page.request.get(`/api/tenants/${ids.tenant}/entries/${ids.entry}/form`);
    expect(form.headers()["content-type"]).toBe("application/pdf");
  });

  test("exports a print form and imports it again", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/forms`);
    const download = page.waitForEvent("download");
    await page.getByLabel("Export Army sheet").click();
    const file = JSON.parse(await readFile(await (await download).path(), "utf8"));
    expect(file).toMatchObject({ name: "Army sheet", width: 148, height: 210 });
    file.name = "Army sheet copy";
    await page.getByLabel("Import a print form").setInputFiles({ name: "form.json", mimeType: "application/json", buffer: Buffer.from(JSON.stringify(file)) });
    await page.waitForURL(/\/forms\/\d+$/);
    await expect(page.getByRole("application", { name: "Sheet" }).locator("[data-copy]")).toHaveCount(12);
  });
});
