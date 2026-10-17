/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
import { build, files, version } from '$service-worker';

const sw = self as unknown as ServiceWorkerGlobalScope;
const cache = `standblatt-${version}`;
const assets = [...build, ...files.filter((f) => !f.endsWith('.webmanifest'))];

sw.addEventListener('install', (e) => {
	e.waitUntil(
		caches
			.open(cache)
			.then((c) => c.addAll([...assets, '/my']))
			.then(() => sw.skipWaiting())
	);
});

sw.addEventListener('activate', (e) => {
	e.waitUntil(
		caches
			.keys()
			.then((keys) => Promise.all(keys.filter((k) => k !== cache).map((k) => caches.delete(k))))
			.then(() => sw.clients.claim())
	);
});

sw.addEventListener('fetch', (e) => {
	const url = new URL(e.request.url);
	if (e.request.method !== 'GET' || url.origin !== location.origin || url.pathname.startsWith('/api/') || url.pathname.startsWith('/admin')) return;
	if (assets.includes(url.pathname)) {
		e.respondWith(caches.match(url.pathname).then((r) => r ?? fetch(e.request)));
		return;
	}
	if (e.request.mode === 'navigate') {
		e.respondWith(
			fetch(e.request).catch(async () => (await caches.match('/my')) ?? Response.error())
		);
	}
});
