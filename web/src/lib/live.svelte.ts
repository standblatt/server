import type { QueryClient, QueryKey } from '@tanstack/svelte-query';
import { keys } from './tenant';

export class Live {
	connected = $state(false);
	private socket: WebSocket | null = null;
	private retry = 1000;
	private timer: ReturnType<typeof setTimeout> | undefined;
	private closed = false;

	constructor(
		private path: string,
		private onTopic: (topic: string) => void
	) {
		this.connect();
	}

	/** Opens the WebSocket and reconnects with backoff when it closes. */
	private connect() {
		const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
		const socket = new WebSocket(`${proto}//${location.host}${this.path}`);
		this.socket = socket;
		socket.onopen = () => {
			this.connected = true;
			this.retry = 1000;
			this.onTopic('resync');
		};
		socket.onmessage = (m) => this.onTopic(JSON.parse(m.data).topic);
		socket.onclose = () => {
			this.connected = false;
			if (this.closed) return;
			this.timer = setTimeout(() => this.connect(), this.retry);
			this.retry = Math.min(this.retry * 2, 30000);
		};
	}

	/** Closes the socket for good; no reconnect. */
	close() {
		this.closed = true;
		clearTimeout(this.timer);
		this.socket?.close();
	}
}

/**
 * Reloads queries after a change elsewhere. A fetch already running may have started before the change,
 * and invalidating alone would join it instead of fetching again, so it is cancelled first.
 */
export function refresh(qc: QueryClient, queryKey: QueryKey) {
	qc.cancelQueries({ queryKey }).then(() => qc.invalidateQueries({ queryKey }));
}

/** Live connection of an event that refreshes the queries of each changed topic. */
export function eventLive(tenantId: number, eventId: number, qc: QueryClient, token?: string) {
	return new Live(`/api/tenants/${tenantId}/events/${eventId}/live${token ? `?token=${encodeURIComponent(token)}` : ''}`, (topic) => {
		const t = tenantId;
		const e = eventId;
		switch (topic) {
			case 'registrations':
				refresh(qc, keys.registrations(t, e));
				refresh(qc, keys.stats(t, e));
				refresh(qc, keys.reports(t, e));
				break;
			case 'results':
				refresh(qc, keys.registrations(t, e));
				refresh(qc, keys.rankings(t, e));
				refresh(qc, keys.reports(t, e));
				break;
			case 'lanes':
				refresh(qc, keys.lanes(t, e));
				refresh(qc, keys.reservations(t, e));
				break;
			case 'booking':
				refresh(qc, keys.booking(t, e));
				refresh(qc, keys.reservations(t, e));
				break;
			case 'programmes':
				refresh(qc, keys.programmes(t, e));
				refresh(qc, keys.rankings(t, e));
				break;
			default:
				refresh(qc, keys.event(t, e));
		}
	});
}
