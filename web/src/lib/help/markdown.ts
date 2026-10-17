import { Marked, type Tokens } from 'marked';
import changelog from '../../../../CHANGELOG.md?raw';

/** Every help page of every language as a lazy raw import, keyed like /src/help/de/user-guide/results.md. */
const files = import.meta.glob('/src/help/**/*.md', { query: '?raw', import: 'default' }) as Record<string, () => Promise<string>>;

/** A help page: its source file, the language it was found in and its rendered parts. */
export type HelpPage = { slug: string; lang: string; title: string; html: string; sections: Section[] };
/** A part of a page under one heading, for the search. */
export type Section = { id: string; heading: string; text: string };

/** The source file of a page in a language, falling back to English; undefined when the page does not exist. */
function source(lang: string, slug: string): { key: string; lang: string; file: string } | undefined {
	const names = slug ? [`${slug}.md`, `${slug}/index.md`] : ['index.md'];
	for (const l of [lang, 'en']) {
		for (const file of names) {
			const key = `/src/help/${l}/${file}`;
			if (key in files) return { key, lang: l, file };
		}
	}
}

/** Whether a help page with this slug exists (in English at least). */
export const pageExists = (slug: string) => slug === 'changelog' || !!source('en', slug);

/** The slug of a page file: user-guide/results.md → user-guide/results, my-training/index.md → my-training. */
const slugOf = (file: string) => file.replace(/\.md$/, '').replace(/(^|\/)index$/, '');

/** The id of a heading: an explicit {#id}, else the text in lower case with dashes. */
function headingId(text: string): string {
	const m = /\{#([\w-]+)\}\s*$/.exec(text);
	if (m) return m[1];
	return text
		.toLowerCase()
		.replace(/[^\p{L}\p{N}\s-]/gu, '')
		.trim()
		.replace(/\s+/g, '-');
}

/** Removes markdown syntax for the search text. */
const plain = (md: string) =>
	md
		.replace(/!\[[^\]]*\]\([^)]*\)(\{[^}]*\})?/g, '')
		.replace(/\[([^\]]*)\]\([^)]*\)/g, '$1')
		.replace(/<[^>]+>/g, '')
		.replace(/^\s*(?:[-*+]|\d+\.)\s+/gm, '')
		.replace(/[*_`]+/g, '')
		.replace(/[|>#]+|^-{3,}$/gm, ' ')
		.replace(/\s+/g, ' ')
		.trim();

/** Splits a page into the parts under its headings. */
function sections(md: string): Section[] {
	const out: Section[] = [];
	let current: Section = { id: '', heading: '', text: '' };
	for (const line of md.split('\n')) {
		const h = /^#{1,6}\s+(.*)$/.exec(line);
		if (h) {
			if (current.heading || current.text) out.push({ ...current, text: plain(current.text) });
			const heading = h[1].replace(/\s*\{#[\w-]+\}\s*$/, '');
			current = { id: headingId(h[1]), heading, text: '' };
		} else current.text += line + '\n';
	}
	out.push({ ...current, text: plain(current.text) });
	return out;
}

/**
 * Turns the Zensical-flavoured markdown of the help pages into what marked understands: images with
 * { .screenshot } become <img> tags, `!!! kind "Title"` blocks become boxes, and `markdown` attributes of
 * raw <div> blocks are dropped and their classes get a help- prefix (marked parses markdown between block tags that
 * stand on their own lines; `grid` would otherwise be Tailwind's utility).
 */
function prepare(md: string): string {
	return md
		.replace(/!\[([^\]]*)\]\(([^)\s]+)\)\{([^}]*)\}/g, (_, alt: string, src: string, attrs: string) => {
			const classes = attrs.match(/\.[\w-]+/g)?.map((c) => c.slice(1)) ?? [];
			return `<img src="${imageSrc(src)}" alt="${alt.replace(/"/g, '&quot;')}" class="${classes.join(' ')}" loading="lazy">`;
		})
		.replace(/^!!! (\w+)(?: "([^"]*)")?\n((?:\n| {4}.*\n?)*)/gm, (_, kind: string, title: string | undefined, body: string) => {
			const text = body.replace(/^ {4}/gm, '').trim();
			return `<div class="admonition ${kind}">\n\n${title ? `**${title}**\n\n` : ''}${text}\n\n</div>\n\n`;
		})
		.replace(/<div class="([^"]*)"\s+markdown(?:="1")?>/g, (_, classes: string) => `<div class="${classes.split(/\s+/).map((c) => `help-${c}`).join(' ')}">`);
}

/** The URL of an image of the pages: everything lives in /help/screenshots/. */
const imageSrc = (src: string) => (src.startsWith('http') ? src : `/help/screenshots/${src.split('/').pop()}`);

/** A marked instance that resolves links relative to a page file and gives headings their ids. */
function renderer(file: string, newTab: boolean) {
	const target = newTab ? ' target="_blank" rel="noopener"' : '';
	return new Marked({
		renderer: {
			/** Heading with the id of {#id} or of its text; the {#id} itself is not shown. */
			heading(this: { parser: { parseInline: (t: Tokens.Generic[]) => string } }, token: Tokens.Heading) {
				const inner = this.parser.parseInline(token.tokens).replace(/\s*\{#[\w-]+\}\s*$/, '');
				return `<h${token.depth} id="${headingId(token.text)}">${inner}</h${token.depth}>\n`;
			},
			/** Links to other pages go to /help/…; links to pages that are not part of the help become plain text. */
			link(this: { parser: { parseInline: (t: Tokens.Generic[]) => string } }, token: Tokens.Link) {
				const text = this.parser.parseInline(token.tokens);
				const href = token.href;
				if (/^https?:|^mailto:/.test(href)) return `<a href="${href}" target="_blank" rel="noopener">${text}</a>`;
				if (href.startsWith('#')) return `<a href="${href}">${text}</a>`;
				const [path, hash] = href.split('#');
				if (!path.endsWith('.md')) return text;
				const resolved = new URL(path, `http://help/${file}`).pathname.slice(1);
				const slug = slugOf(resolved);
				if (!pageExists(slug)) return text;
				return `<a href="/help/${slug}${hash ? `#${hash}` : ''}"${target}>${text}</a>`;
			},
			/** Images without attributes, e.g. inside tables. */
			image(token: Tokens.Image) {
				return `<img src="${imageSrc(token.href)}" alt="${token.text.replace(/"/g, '&quot;')}" loading="lazy">`;
			}
		}
	});
}

/** Loads and renders a page in a language (English when it is not translated). */
export async function loadPage(lang: string, slug: string, newTab = false): Promise<HelpPage | undefined> {
	let md: string;
	let found: { lang: string; file: string };
	if (slug === 'changelog') {
		md = changelog;
		found = { lang: 'en', file: 'changelog.md' };
	} else {
		const src = source(lang, slug);
		if (!src) return undefined;
		md = await files[src.key]();
		found = src;
	}
	const title = /^#\s+(.*)$/m.exec(md)?.[1].replace(/\s*\{#[\w-]+\}\s*$/, '') ?? slug;
	const html = renderer(found.file, newTab).parse(prepare(md)) as string;
	return { slug, lang: found.lang, title, html, sections: sections(md) };
}
