import { test, expect } from "@playwright/test";
import { c, createClub, login, newPage } from "../helpers.js";

const run = Date.now().toString(36);
const club = `Master Data Club ${run}`;
const admin = `md-${run}`;

async function open(page) {
  await login(page, admin);
  await page.goto("/");
  await c(page, "a", club).click();
  await expect(c(page, "header", club)).toBeVisible();
  await c(page, "a", "Shooters").click();
}

const csv = (name, text) => ({ name, mimeType: "text/csv", buffer: Buffer.from(text) });

test.describe.configure({ mode: "serial" });

test.describe("shooters and clubs", () => {
  test.beforeAll(async ({ browser }) => {
    const page = await newPage(browser);
    await createClub(page, club, admin);
    await page.close();
  });

  test("creates, finds and edits a shooter", async ({ page }) => {
    await open(page);
    await c(page, "button", "New shooter").click();
    await page.getByRole("dialog").getByLabel("Club", { exact: true }).fill("Sargans Pistol");
    await c(page, "label", "License number").locator("input").fill("13756");
    await c(page, "label", "First name").locator("input").fill("Sandra");
    await c(page, "label", "Last name").locator("input").fill("Schwitter");
    await c(page, "label", "Birthday").locator("input").fill("1969-06-02");
    await expect(c(page, "[role=status]", "“Sargans Pistol” is not in the club list yet.")).toBeVisible();
    await page.getByRole("button", { name: "Create", exact: true }).click();
    await expect(c(page, "[role=dialog]", "Pick an existing club or confirm the new one.")).toBeVisible();
    await c(page, "button", "Create new club “Sargans Pistol”").click();
    await expect(c(page, "[role=status]", "“Sargans Pistol” will be added as a new club.")).toBeVisible();
    await page.getByRole("button", { name: "Create", exact: true }).click();
    const dialog = c(page, "[role=dialog]", "Schwitter Sandra");
    await expect(dialog).toContainText("License 013756 · born 02.06.1969");
    await expect(dialog).toContainText("Sargans Pistol");
    await dialog.getByLabel("Club", { exact: true }).fill("Second Club");
    await c(dialog, "button", "Create new club “Second Club”").click();
    await c(dialog, "button", "Add club").click();
    await expect(dialog.locator("[data-slot=badge]", { hasText: "Second Club" })).toBeVisible();
    await c(page, "[role=dialog] button", "Close").click();
    await expect(page.locator("[role=dialog]")).toHaveCount(0);
    await page.locator("input[type=search]").fill("schwitter sandra");
    await expect(c(page, "tr", "013756")).toContainText("Sargans Pistol");
    await expect(c(page, "tr", "013756")).toContainText("Second Club");
    await page.locator("input[type=search]").fill("13756");
    await expect(c(page, "tr", "Schwitter Sandra")).toBeVisible();
  });

  test("suggests a similar club instead of creating a duplicate", async ({ page }) => {
    await open(page);
    await c(page, "button", "New shooter").click();
    const dialog = page.getByRole("dialog");
    await dialog.getByLabel("Club", { exact: true }).fill("sargans pistl");
    await c(dialog, "[role=status] button", "Sargans Pistol").click();
    await expect(dialog.getByLabel("Club", { exact: true })).toHaveValue("Sargans Pistol");
    await expect(dialog.locator("[role=status]")).toHaveCount(0);
    await c(page, "label", "First name").locator("input").fill("Tina");
    await c(page, "label", "Last name").locator("input").fill("Typo");
    await page.getByRole("button", { name: "Create", exact: true }).click();
    await expect(c(page, "[role=dialog]", "Typo Tina")).toContainText("Sargans Pistol");
    await c(page, "[role=dialog] button", "Close").click();
    await c(page, "a", "Clubs").click();
    await expect(c(page, "tr", "Sargans Pistol")).toContainText("2");
    await expect(page.locator("tbody tr")).toHaveCount(2);
  });

  test("imports shooters from CSV and reports bad lines", async ({ page }) => {
    await open(page);
    await c(page, "button", "Import CSV").click();
    const dialog = page.locator("[role=dialog]");
    await page.locator("input[type=file]").setInputFiles(csv("bad.csv", "license;firstname;lastname;birthday;club\n1;;Nofirst;;\n"));
    await c(dialog, "button", /^Import$/).click();
    await expect(c(page, "[role=alert]", "Nothing imported")).toBeVisible();
    await expect(c(page, "li", "Line 2")).toBeVisible();

    await page.locator("input[type=file]").setInputFiles(
      csv(
        "good.csv",
        "license;firstname;lastname;birthday;postalcode;city;club\n" +
          "117804;Max;Abderhalden;17.04.1949;8450;Andelfingen;Winterthur-Veltheim\n" +
          ";Guest;Person;;;;Winterthur-Veltheim\n" +
          ";Paula;Pistole;;;;Sargans Pistole\n" +
          ";Sam;Second;;;;second club\n",
      ),
    );
    await c(dialog, "button", /^Import$/).click();
    await expect(c(dialog, "[role=status]", "Nothing imported yet")).toBeVisible();
    await expect(dialog.getByLabel("Club for Winterthur-Veltheim")).toHaveValue("new");
    await expect(dialog.getByLabel("Club for second club")).toHaveCount(0);
    await dialog.getByLabel("Club for Sargans Pistole").selectOption({ label: "Sargans Pistol" });
    await c(dialog, "button", /^Import$/).click();
    await expect(c(page, "[role=status]", "4 new, 0 updated, 1 new clubs, 4 club assignments")).toBeVisible();
    await c(page, "[role=dialog] button", "Close").click();

    await c(page, "a", "Clubs").click();
    await expect(c(page, "tr", "Winterthur-Veltheim")).toContainText("2");
    await expect(c(page, "tr", "Sargans Pistol")).toContainText("3");
    await expect(c(page, "tr", "Second Club")).toContainText("2");
  });

  test("merges a duplicate club", async ({ page }) => {
    await open(page);
    await c(page, "a", "Clubs").click();
    await page.getByLabel("Actions for Second Club").click();
    await c(page, "[role=menuitem]", "Merge into…").click();
    const dialog = page.getByRole("dialog");
    await dialog.getByLabel("Club that stays").selectOption({ label: "Sargans Pistol" });
    await c(dialog, "button", "Merge").click();
    await c(page, "[role=alertdialog] button", "Merge").click();
    await expect(c(page, "[data-sonner-toast]", "Clubs merged")).toBeVisible();
    await expect(c(page, "tr", "Second Club")).toHaveCount(0);
    await expect(c(page, "tr", "Sargans Pistol")).toContainText("4");
  });
});
