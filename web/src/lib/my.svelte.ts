import { msg } from '@lingui/core/macro';
import type { MessageDescriptor } from '@lingui/core';
import { createQuery, type QueryClient } from '@tanstack/svelte-query';
import { goto } from '$app/navigation';
import { api, call, ApiError, OfflineError, type Schemas } from './api';
import { intlLocale } from './locale.svelte';

export type MyOverview = Schemas['MyOverview'];
export type Discipline = Schemas['Discipline'];
export type Weapon = Schemas['Weapon'];
export type TrainingSession = Schemas['TrainingSession'];
export type TrainingInput = Schemas['TrainingInput'];
export type TrainingSeries = Schemas['TrainingSeries'];
export type TrainingSaved = Schemas['TrainingSaved'];
export type MyStats = Schemas['MyStats'];
export type CompetitionResult = Schemas['CompetitionResult'];

export const myKey = ['my'];
export const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;

/** The home data of the personal log; everything under /my invalidates the 'my' prefix after changes. */
export function myQuery() {
	return createQuery(() => ({
		queryKey: myKey,
		queryFn: async () => {
			const o = await call(api.GET('/api/my', { params: { query: { tz: timeZone } } }));
			cacheMy(o);
			return o;
		},
		placeholderData: cachedMy
	}));
}

/** Query of the statistics of a discipline, optionally filtered by weapon (0 = all). */
export function statsQuery(id: () => number | undefined, enabled: () => boolean = () => true, weapon: () => number = () => 0) {
	return createQuery(() => ({
		queryKey: [...myKey, 'stats', id(), weapon()],
		queryFn: () => call(api.GET('/api/my/stats/{id}', { params: { path: { id: id()! }, query: weapon() ? { weapon: weapon() } : {} } })),
		enabled: !!id() && enabled(),
		retry: false
	}));
}

const builtIn: Record<string, MessageDescriptor> = {
	air_pistol: msg`Air pistol 10 m`,
	air_rifle: msg`Air rifle 10 m`,
	pistol_25: msg`Pistol 25 m`,
	pistol_50: msg`Pistol 50 m`,
	rifle_50: msg`Rifle 50 m`,
	rifle_300: msg`Rifle 300 m`
};

/** Name of a discipline: built-in ones are translated, own ones keep their name. */
export function disciplineName(d: Discipline | undefined): MessageDescriptor {
	if (!d) return { id: '…', message: '…' };
	return builtIn[d.key] ?? { id: d.name, message: d.name };
}

/** Highest value of one shot in a discipline (10.9 with decimal scoring). */
export const maxShot = (d: Discipline) => (d.decimal ? 10.9 : d.maxShot);

/** Formats a number with the given decimals in the UI locale, a dash when missing. */
export function fmt(v: number | null | undefined, digits = 1): string {
	if (v == null) return '–';
	return new Intl.NumberFormat(intlLocale(), { minimumFractionDigits: digits, maximumFractionDigits: digits }).format(v);
}

/** A total: whole numbers without decimals, decimal scoring with one. */
export const fmtTotal = (v: number) => (Number.isInteger(v) ? fmt(v, 0) : fmt(v, 1));

/** Short weekday, day and month of a timestamp in the UI locale. */
export function fmtDay(iso: string): string {
	return new Intl.DateTimeFormat(intlLocale(), { weekday: 'short', day: 'numeric', month: 'short' }).format(new Date(iso));
}

/** Medium date and short time of a timestamp in the UI locale. */
export function fmtDateTime(iso: string): string {
	return new Intl.DateTimeFormat(intlLocale(), { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(iso));
}

export { localInput } from './format';

/** New random UUID. */
export const uuid = () => crypto.randomUUID();

const outboxKey = 'standblatt-outbox';
const myCacheKey = 'standblatt-my';

/** A session saved on this device for a user; error holds the server's answer when the upload was refused. */
type Queued = { userId: number; input: TrainingInput; error?: string };

/** Sessions in the offline outbox, empty when storage fails; entries of the old format (no user) belong to userId. */
function readOutbox(userId = 0): Queued[] {
	try {
		const list: (Queued | TrainingInput)[] = JSON.parse(localStorage.getItem(outboxKey) ?? localStorage.getItem('targetlog-outbox') ?? '[]');
		return list.map((x) => ('input' in x ? x : { userId, input: x }));
	} catch {
		return [];
	}
}

/** Stores the offline outbox and updates the counts of userId. */
function writeOutbox(list: Queued[], userId: number) {
	try {
		localStorage.setItem(outboxKey, JSON.stringify(list));
		localStorage.removeItem('targetlog-outbox');
	} catch {
		// nothing we can do without storage
	}
	countOutbox(list, userId);
}

/** Updates the counts of the outbox for userId. */
function countOutbox(list: Queued[], userId: number) {
	const mine = list.filter((x) => x.userId === userId);
	outbox.count = mine.filter((x) => !x.error).length;
	outbox.failed = mine.filter((x) => x.error).map((x) => x.error!);
}

/** Sessions saved on this device while offline, uploaded by flushOutbox; failed holds the errors of refused uploads. */
export const outbox = $state<{ count: number; failed: string[] }>({ count: 0, failed: [] });

/** Puts a session of userId into the offline outbox, replacing one with the same client ID. */
export function queueOffline(input: TrainingInput, userId: number) {
	writeOutbox([...readOutbox(userId).filter((x) => x.input.clientId !== input.clientId), { userId, input }], userId);
}

/** Removes the sessions of userId whose upload was refused. */
export function discardFailed(userId: number) {
	writeOutbox(readOutbox(userId).filter((x) => x.userId !== userId || !x.error), userId);
}

let flushing = false;

/** Uploads the sessions userId saved offline; the client ID makes a retry safe, refused ones stay with their error. */
export async function flushOutbox(queryClient: QueryClient, userId: number) {
	const list = readOutbox(userId);
	countOutbox(list, userId);
	if (flushing || !navigator.onLine) return;
	flushing = true;
	let sent = 0;
	try {
		for (const q of list.filter((x) => x.userId === userId && !x.error)) {
			let error: string | undefined;
			try {
				await call(api.POST('/api/my/trainings', { body: q.input }));
				sent++;
			} catch (err) {
				if (!(err instanceof ApiError) || err.status === 401 || err.status >= 500) break;
				error = err.message;
			}
			const now = readOutbox(userId);
			writeOutbox(error ? now.map((x) => (x.input.clientId === q.input.clientId ? { ...x, error } : x)) : now.filter((x) => x.input.clientId !== q.input.clientId), userId);
		}
	} finally {
		flushing = false;
	}
	if (sent) await queryClient.invalidateQueries({ queryKey: myKey });
}

/** The last home data kept on this device, so the app opens offline. */
function cachedMy(): MyOverview | undefined {
	try {
		return JSON.parse(localStorage.getItem(myCacheKey) ?? 'null') ?? undefined;
	} catch {
		return undefined;
	}
}

/** Keeps the home data on this device. */
function cacheMy(o: MyOverview) {
	try {
		localStorage.setItem(myCacheKey, JSON.stringify(o));
	} catch {
		// the app then only opens online
	}
}

/** True for errors where the request never reached the server. */
export const isOffline = (err: unknown) => err instanceof OfflineError;

const kinds: Record<string, MessageDescriptor> = { training: msg`Training`, competition: msg`Competition` };
/** Label descriptor of a session kind (training, competition), the raw kind otherwise. */
export const kindLabel = (k: string) => kinds[k] ?? { id: k, message: k };

/** Drops the home data kept on this device, so the next user does not see it while theirs loads. */
export function forgetMy() {
	try {
		localStorage.removeItem(myCacheKey);
	} catch {
		// nothing kept
	}
}

/** Ends the session and returns to the login page. */
export async function logout(queryClient: QueryClient) {
	await call(api.POST('/api/auth/logout'));
	forgetMy();
	queryClient.clear();
	await goto('/login');
}

/** The total a points-per-shot average gives over the standard match of a discipline, rounded like its scores. */
export function projected(perShot: number, d: Discipline): number {
	const v = perShot * d.seriesCount * d.shotsPerSeries;
	return d.decimal ? Math.round(v * 10) / 10 : Math.round(v);
}
