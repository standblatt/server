/** Geometry of a target face in mm: radius of the 10 ring, ring width, radius of the projectile and of the black aiming mark. */
export type Target = { ten: number; ring: number; bullet: number; black: number };

export const targets: Record<string, Target> = {
	air_pistol: { ten: 5.75, ring: 8, bullet: 2.25, black: 29.75 },
	air_rifle: { ten: 0.25, ring: 2.5, bullet: 2.25, black: 15.25 },
	pistol_precision: { ten: 25, ring: 25, bullet: 2.8, black: 100 },
	rifle_50: { ten: 5.2, ring: 8, bullet: 2.8, black: 56.2 },
	rifle_300: { ten: 50, ring: 50, bullet: 4, black: 300 }
};

/** Outer radius of ring n (1 to 10). */
export const ringRadius = (t: Target, n: number) => t.ten + (10 - n) * t.ring;

/** Score of a shot at x/y mm from the centre: decimal values are truncated to tenths and capped at 10.9. */
export function score(t: Target, x: number, y: number, decimal: boolean): number {
	const d = Math.hypot(x, y);
	const s = Math.min(10.9, 10 + (t.ten + t.bullet - d) / t.ring);
	if (s < 1) return 0;
	return decimal ? Math.floor(s * 10 + 1e-9) / 10 : Math.min(10, Math.floor(s));
}

/** Clock-like description of an offset, e.g. right and low, used to explain the group centre. */
export function direction(x: number, y: number, tolerance: number): { horizontal: 'left' | 'right' | ''; vertical: 'high' | 'low' | '' } {
	return {
		horizontal: Math.abs(x) < tolerance ? '' : x > 0 ? 'right' : 'left',
		vertical: Math.abs(y) < tolerance ? '' : y > 0 ? 'high' : 'low'
	};
}
