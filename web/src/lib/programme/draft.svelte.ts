import type { AgeCategory, ProgrammeDefinition } from '$lib/api';
import { ageCategories } from '$lib/categories';

/** An award tier as edited: limits are strings so a field can stay empty. */
export type AwardDraft = { name: string; kind: 'natural' | 'card' | 'cash'; value: number; limits: Record<AgeCategory, string> };

/** Award limits with an empty value for every age category. */
export function emptyLimits(): Record<AgeCategory, string> {
	return Object.fromEntries(ageCategories.map((c) => [c, ''])) as Record<AgeCategory, string>;
}

/** The programme being edited; every section of the programme page binds its part of it. */
export class ProgrammeDraft {
	name = $state('');
	variant = $state('');
	description = $state('');
	barcode = $state('');
	countsForClubCompetition = $state(false);
	showOnRankings = $state(true);
	attempts = $state(1);
	unlimited = $state(false);
	attemptRule = $state<'best' | 'first' | 'last'>('best');
	awardCountingOnly = $state(false);
	series = $state<{ shots: number; maxPoints: number; name: string; trackMouches: boolean; trackShots: boolean }[]>([]);
	prices = $state({} as Record<AgeCategory, { enabled: boolean; price: number; repeat: string }>);
	awards = $state<AwardDraft[]>([]);
	/** Age categories and number of awards that existed when editing started (they stay while shooters are registered). */
	readonly fixedPrices: Set<AgeCategory>;
	readonly fixedAwards: number;

	/** Starts from a programme, or from an empty one with a single series. */
	constructor(d?: ProgrammeDefinition) {
		const start = d ?? { name: '', variant: '', showOnRankings: true, series: [{ shots: 10, maxPoints: 100 }], prices: [], awards: [] };
		this.name = start.name;
		this.variant = start.variant ?? '';
		this.description = start.description ?? '';
		this.barcode = start.barcode ?? '';
		this.countsForClubCompetition = start.countsForClubCompetition ?? false;
		this.showOnRankings = start.showOnRankings ?? true;
		this.attempts = start.unlimitedAttempts ? 2 : Math.max(start.attempts ?? 1, 1);
		this.unlimited = start.unlimitedAttempts ?? false;
		this.attemptRule = start.attemptRule ?? 'best';
		this.awardCountingOnly = start.awardCountingOnly ?? false;
		this.series = start.series.map((s) => ({ shots: s.shots, maxPoints: s.maxPoints, name: s.name ?? '', trackMouches: s.trackMouches ?? false, trackShots: s.trackShots ?? false }));
		this.prices = Object.fromEntries(
			ageCategories.map((c) => {
				const p = start.prices.find((x) => x.ageCategory === c);
				return [c, { enabled: !!p, price: p?.price ?? 0, repeat: p?.repeatPrice == null ? '' : String(p.repeatPrice) }];
			})
		) as Record<AgeCategory, { enabled: boolean; price: number; repeat: string }>;
		this.awards = start.awards.map((a) => {
			const limits = emptyLimits();
			for (const l of a.limits) limits[l.ageCategory] = String(l.points);
			return { name: a.name, kind: a.kind, value: a.value, limits };
		});
		this.fixedPrices = new Set(start.prices.map((p) => p.ageCategory));
		this.fixedAwards = start.awards.length;
	}

	/** The programme can be shot more than once. */
	get repeatable() {
		return this.unlimited || Number(this.attempts) > 1;
	}

	/** Shots over all series. */
	get totalShots() {
		return this.series.reduce((n, s) => n + (Number(s.shots) || 0), 0);
	}

	/** Maximum points over all series. */
	get totalPoints() {
		return this.series.reduce((n, s) => n + (Number(s.maxPoints) || 0), 0);
	}

	/** The definition to save. */
	definition(): ProgrammeDefinition {
		const repeatable = this.repeatable;
		return {
			name: this.name,
			variant: this.variant,
			description: this.description,
			barcode: this.barcode,
			countsForClubCompetition: this.countsForClubCompetition,
			showOnRankings: this.showOnRankings,
			attempts: this.unlimited ? 0 : Number(this.attempts) || 1,
			unlimitedAttempts: this.unlimited,
			attemptRule: this.attemptRule,
			awardCountingOnly: repeatable && this.awardCountingOnly,
			series: this.series.map((s) => ({ shots: Number(s.shots), maxPoints: Number(s.maxPoints), name: s.name, trackMouches: s.trackMouches, trackShots: s.trackShots })),
			prices: ageCategories
				.filter((c) => this.prices[c].enabled)
				.map((c) => ({
					ageCategory: c,
					price: Number(this.prices[c].price) || 0,
					repeatPrice: repeatable && this.prices[c].repeat !== '' ? Number(this.prices[c].repeat) : undefined
				})),
			awards: this.awards.map((a) => ({
				name: a.name,
				kind: a.kind,
				value: Number(a.value) || 0,
				limits: ageCategories.filter((c) => a.limits[c] !== '').map((c) => ({ ageCategory: c, points: Number(a.limits[c]) }))
			}))
		};
	}
}
