import { defineConfig } from "@playwright/test";
import { baseURL } from "./helpers.js";

// `e2e` is the test suite (`make e2e`); `docs` writes the help screenshots (`make help-screenshots`).
export default defineConfig({
  retries: 2,
  workers: process.env.CI ? 4 : undefined,
  reporter: process.env.CI ? [["list"], ["junit", { outputFile: "test-results/junit.xml" }]] : "list",
  expect: { timeout: 8000 },
  use: { baseURL, screenshot: "only-on-failure", trace: "retain-on-failure" },
  projects: [
    { name: "e2e", testDir: "tests", use: { viewport: { width: 1280, height: 800 } } },
    { name: "docs", testDir: "docs", retries: 0, use: { viewport: { width: 1440, height: 900 } } },
  ],
});
