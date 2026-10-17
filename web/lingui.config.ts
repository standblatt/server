import babelExtractor from '@lingui/cli/api/extractors/babel';
import { formatter } from '@lingui/format-po-gettext';
import { defineConfig } from 'lingui-for-svelte/config';
import { svelteExtractor } from 'lingui-for-svelte/extractor';

export default defineConfig({
	locales: ['en', 'de', 'fr', 'it'],
	sourceLocale: 'en',
	catalogs: [{ path: 'src/locales/{locale}', include: ['src'], exclude: ['src/lib/api.gen.ts'] }],
	format: formatter({ origins: true, lineNumbers: false }),
	extractors: [svelteExtractor, babelExtractor]
});
