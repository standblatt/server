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

  test("builds a print form in the editor", async ({ page }) => {
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

    await c(page, "button", "Add element").click();
    await page.getByLabel("Value of element 1").click();
    await c(page, "section[aria-label=Placeholders] button", /^shooter\.name$/).click();
    await expect(page.getByLabel("Value of element 1")).toHaveValue("{{shooter.name}}");
    await c(page, "button", "Add element").click();
    await page.getByLabel("Repeat element 2").fill("3");
    await page.getByLabel("dx of element 2").fill("9");
    await page.getByLabel("Value of element 2").click();
    await c(page, "section[aria-label=Placeholders] button", /^s1\.shot1$/).click();
    await expect(page.getByLabel("Value of element 2")).toHaveValue("{{s1.shot#}}");
    await page.getByRole("button", { name: "Save", exact: true }).click();
    await expect(c(page, "[data-sonner-toast]", "Print form saved")).toBeVisible();
    await page.getByRole("tab", { name: "PDF preview" }).click();
    await expect(page.locator("iframe[title=Preview]")).toHaveAttribute("src", /preview\?v=1$/);

    const layout = await api(page, "GET", `/api/tenants/${ids.tenant}/layouts/${ids.layout}`);
    expect(layout.elements).toMatchObject([{ value: "{{shooter.name}}", y: 13 }, { value: "{{s1.shot#}}", repeat: 3, dx: 9 }]);
    const preview = await page.request.get(`/api/tenants/${ids.tenant}/layouts/${ids.layout}/preview`);
    expect(preview.headers()["content-type"]).toBe("application/pdf");
    expect(preview.headers()["x-frame-options"]).toBe("SAMEORIGIN");
  });

  test("moves elements on the sheet with the mouse and the keyboard", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/forms/${ids.layout}`);
    const sheet = page.getByRole("application", { name: "Sheet" });
    const first = sheet.locator("[data-element='0']");
    await expect(first).toContainText("Example Anna");
    await expect(sheet.locator("[data-element='1']")).toHaveCount(3);

    const x = page.getByLabel("x of element 1", { exact: true });
    const y = page.getByLabel("y of element 1", { exact: true });
    const before = Number(await x.inputValue());
    const box = await first.boundingBox();
    const sheetBox = await sheet.boundingBox();
    const mm = sheetBox.width / 148;
    await page.mouse.move(box.x + 5, box.y + box.height / 2);
    await page.mouse.down();
    await page.mouse.move(box.x + 5 + 20 * mm, box.y + box.height / 2 + 10 * mm, { steps: 5 });
    await page.mouse.up();
    await expect(x).toHaveValue(String(before + 20));
    await expect(y).toHaveValue("23");

    await page.keyboard.press("ArrowRight");
    await page.keyboard.press("Shift+ArrowDown");
    await expect(x).toHaveValue(String(before + 20.5));
    await expect(y).toHaveValue("28");

    const now = await sheet.boundingBox();
    await page.mouse.dblclick(now.x + 100 * mm, now.y + 150 * mm);
    await expect(page.getByLabel("x of element 3", { exact: true })).toHaveValue("100");
    await expect(page.getByLabel("y of element 3", { exact: true })).toHaveValue("150");
    await page.getByLabel("Remove element 3").click();
    await page.getByRole("button", { name: "Save", exact: true }).click();
    await expect(c(page, "[data-sonner-toast]", "Print form saved")).toBeVisible();
    const layout = await api(page, "GET", `/api/tenants/${ids.tenant}/layouts/${ids.layout}`);
    expect(layout.elements[0]).toMatchObject({ x: before + 20.5, y: 28 });
  });

  test("prints the results onto the form of the programme", async ({ page }) => {
    await login(page, admin);
    await page.goto(`/t/${ids.tenant}/events/${ids.event}/setup`);
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
    await expect(page.getByLabel("Value of element 2")).toHaveValue("{{s1.shot#}}");
  });
});
