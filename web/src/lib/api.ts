import createClient from 'openapi-fetch';
import { goto } from '$app/navigation';
import type { components, paths } from './api.gen';
import { msg } from '@lingui/core/macro';
import { currentLocale, i18n } from './locale.svelte';

export type Schemas = components['schemas'];
export type Me = Schemas['Me'];
export type Range = Schemas['Range'];
export type RangeInput = Schemas['RangeInput'];
export type Event = Schemas['Event'];
export type EventInput = Schemas['EventInput'];
export type Schedule = Schemas['Schedule'];
export type ScheduleDay = Schemas['ScheduleDay'];
export type WeeklyTimes = Schemas['WeeklyTimes'];
export type ShootingTimes = Schemas['BookingHours'];
export type RangeGroup = Schemas['RangeGroup'];
export type Member = Schemas['Member'];
export type Club = Schemas['Club'];
export type Shooter = Schemas['Shooter'];
export type ClubRef = Schemas['ClubRef'];
export type ShooterInput = Schemas['ShooterInput'];
export type ImportReport = Schemas['ImportReport'];
export type Programme = Schemas['Programme'];
export type ProgrammeDefinition = Schemas['ProgrammeDefinition'];
export type ProgrammeTemplate = Schemas['ProgrammeTemplate'];
export type ProgrammeExport = Schemas['ProgrammeExport'];
export type AgeCategory = Schemas['Price']['ageCategory'];
export type Registration = Schemas['Registration'];
export type RegistrationInput = Schemas['RegistrationInput'];
export type ProgrammeChoice = Schemas['ProgrammeChoice'];
export type LaneState = Schemas['LaneState'];
export type RangeLanes = Schemas['RangeLanes'];
export type ResultSheet = Schemas['ResultSheet'];
export type EntryResult = Schemas['EntryResult'];
export type SeriesResult = Schemas['SeriesResult'];
export type ProgrammeRanking = Schemas['ProgrammeRanking'];
export type KioskToken = Schemas['KioskToken'];
export type HelperCode = Schemas['HelperCode'];
export type Helper = Schemas['Helper'];
export type PrintAgent = Schemas['PrintAgent'];
export type PrintSetting = Schemas['PrintSetting'];
export type Printer = Schemas['Printer'];
export type ReservationView = Schemas['ReservationView'];
export type ReservationInput = Schemas['ReservationInput'];
export type Role = Member['role'];
export type BookingConfig = Schemas['BookingConfig'];
export type BookingSettings = Schemas['BookingSettings'];
export type BookingInfo = Schemas['BookingInfo'];
export type BookingItem = Schemas['BookingItem'];
export type BookingShooter = Schemas['BookingShooter'];
export type BookingRequest = Schemas['BookingRequest'];
export type BookingRequestInput = Schemas['BookingRequestInput'];

export const api = createClient<paths>({ baseUrl: '' });
/** The helper link code when the app runs under /h/<code>; it authenticates every request instead of a session. */
export const helperCode = () => (typeof location === 'undefined' ? undefined : /^\/h\/([^/?#]+)/.exec(location.pathname)?.[1]);

api.use({
	/** Adds the UI language and, under /h/<code>, the helper link as Bearer token to every request. */
	onRequest({ request }) {
		request.headers.set('Accept-Language', currentLocale());
		const code = helperCode();
		if (code) request.headers.set('Authorization', `Bearer ${code}`);
	}
});

/** OfflineError is thrown by call() when the request never reached the server. */
export class OfflineError extends Error {}

const noConnection = msg`No connection. Nothing was saved.`;

export class ApiError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
	}
}

type Result<T> = { data?: T; error?: unknown; response: Response };

/** Awaits an openapi-fetch call and returns its data; throws ApiError otherwise and handles 401 (login redirect, or an event on helper pages). */
export async function call<T>(p: Promise<Result<T>>): Promise<T> {
	const { data, error, response } = await p.catch((err) => {
		throw err instanceof TypeError ? new OfflineError(i18n._(noConnection)) : err;
	});
	if (response.ok) return data as T;
	if (response.status === 401 && location.pathname !== '/login' && !location.pathname.startsWith('/kiosk/') && !location.pathname.startsWith('/b/')) {
		if (location.pathname.startsWith('/h/')) window.dispatchEvent(new Event('standblatt:unauthorized'));
		else await goto('/login');
	}
	const detail = (error as { detail?: string } | undefined)?.detail;
	throw new ApiError(response.status, detail || response.statusText);
}

/** Fetches a file the API serves (a PDF, a CSV) with the helper code if any, named after its Content-Disposition; throws like call(). */
export async function fetchFile(url: string): Promise<File> {
	const code = helperCode();
	const res = await fetch(url, { headers: { 'Accept-Language': currentLocale(), ...(code ? { Authorization: `Bearer ${code}` } : {}) } }).catch(() => {
		throw new OfflineError(i18n._(noConnection));
	});
	if (!res.ok) throw new ApiError(res.status, ((await res.json().catch(() => ({}))) as { detail?: string }).detail || res.statusText);
	const name = /filename="?([^";]+)"?/.exec(res.headers.get('Content-Disposition') ?? '')?.[1] ?? 'document.pdf';
	return new File([await res.blob()], name, { type: res.headers.get('Content-Type') ?? 'application/octet-stream' });
}
