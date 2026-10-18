<script lang="ts">
	import { t } from 'lingui-for-svelte/macro';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import Download from '@lucide/svelte/icons/download';
	import Pencil from '@lucide/svelte/icons/pencil';
	import Plus from '@lucide/svelte/icons/plus';
	import Smartphone from '@lucide/svelte/icons/smartphone';
	import Trash from '@lucide/svelte/icons/trash-2';
	import UserCog from '@lucide/svelte/icons/user-cog';
	import { api, call } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { ask } from '$lib/confirm.svelte';
	import FormField from '$lib/FormField.svelte';
	import { versionQuery } from '$lib/me';
	import { disciplineName, myKey, myQuery, type Discipline, type Weapon } from '$lib/my.svelte';
	import ThemeToggle from '$lib/ThemeToggle.svelte';

	const version = versionQuery();
	const queryClient = useQueryClient();
	const my = myQuery();
	const data = $derived(my.data!);

	let profile = $state({ firstname: '', lastname: '', birthday: '', licenseNumber: '' });
	$effect.pre(() => {
		const p = my.data?.profile;
		if (p) profile = { firstname: p.firstname, lastname: p.lastname, birthday: p.birthday ?? '', licenseNumber: p.licenseNumber ?? '' };
	});

	/** Saves the profile. */
	async function saveProfile(e: SubmitEvent) {
		e.preventDefault();
		try {
			await call(api.PUT('/api/my/profile', { body: profile }));
			await queryClient.invalidateQueries({ queryKey: myKey });
			toast.success($t`Profile saved`);
		} catch (err) {
			toast.error((err as Error).message);
		}
	}

	let disciplineOpen = $state(false);
	let editing = $state<Discipline>();
	let form = $state({ name: '', seriesCount: 6, shotsPerSeries: 10, maxShot: 10, decimal: false, target: '' as '' | 'air_pistol' | 'air_rifle' | 'pistol_precision' | 'rifle_50' | 'rifle_300' });
	/** Opens the discipline dialog for an existing discipline or a new one. */
	function openDiscipline(d?: Discipline) {
		editing = d;
		form = d
			? { name: d.name, seriesCount: d.seriesCount, shotsPerSeries: d.shotsPerSeries, maxShot: d.maxShot, decimal: d.decimal, target: d.target as typeof form.target }
			: { name: '', seriesCount: 6, shotsPerSeries: 10, maxShot: 10, decimal: false, target: '' };
		disciplineOpen = true;
	}
	/** Creates or updates the discipline from the dialog. */
	async function saveDiscipline(e: SubmitEvent) {
		e.preventDefault();
		try {
			const body = { ...form, maxShot: form.decimal ? 10 : form.maxShot };
			if (editing) await call(api.PUT('/api/my/disciplines/{id}', { params: { path: { id: editing.id } }, body }));
			else await call(api.POST('/api/my/disciplines', { body }));
			disciplineOpen = false;
			await queryClient.invalidateQueries({ queryKey: myKey });
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
	/** Deletes an own discipline after confirmation. */
	async function deleteDiscipline(d: Discipline) {
		if (!(await ask({ title: $t`Delete discipline "${d.name}"?`, action: $t`Delete`, destructive: true }))) return;
		try {
			await call(api.DELETE('/api/my/disciplines/{id}', { params: { path: { id: d.id } } }));
			await queryClient.invalidateQueries({ queryKey: myKey });
		} catch (err) {
			toast.error((err as Error).message);
		}
	}

	let weaponOpen = $state(false);
	let weapon = $state<Weapon>();
	let weaponName = $state('');
	/** Opens the weapon dialog for an existing weapon or a new one. */
	function openWeapon(w?: Weapon) {
		weapon = w;
		weaponName = w?.name ?? '';
		weaponOpen = true;
	}
	/** Creates or updates the weapon from the dialog. */
	async function saveWeapon(e: SubmitEvent) {
		e.preventDefault();
		try {
			const body = { name: weaponName };
			if (weapon) await call(api.PUT('/api/my/weapons/{id}', { params: { path: { id: weapon.id } }, body }));
			else await call(api.POST('/api/my/weapons', { body }));
			weaponOpen = false;
			await queryClient.invalidateQueries({ queryKey: myKey });
		} catch (err) {
			toast.error((err as Error).message);
		}
	}
	/** Deletes a weapon after confirmation; its sessions keep no weapon. */
	async function deleteWeapon(w: Weapon) {
		const ok = await ask({
			title: $t`Delete weapon "${w.name}"?`,
			description: $t`Its sessions stay in your log, without a weapon.`,
			action: $t`Delete`,
			destructive: true
		});
		if (!ok) return;
		try {
			await call(api.DELETE('/api/my/weapons/{id}', { params: { path: { id: w.id } } }));
			await queryClient.invalidateQueries({ queryKey: myKey });
		} catch (err) {
			toast.error((err as Error).message);
		}
	}

	let installPrompt = $state<Event & { prompt: () => Promise<void> }>();
	const ios = typeof navigator !== 'undefined' && /iphone|ipad/i.test(navigator.userAgent);
	const standalone = typeof window !== 'undefined' && matchMedia('(display-mode: standalone)').matches;
	onMount(() => {
		/** Keeps the browser's install prompt so the install button can show it later. */
		const handler = (e: Event) => {
			e.preventDefault();
			installPrompt = e as typeof installPrompt;
		};
		window.addEventListener('beforeinstallprompt', handler);
		return () => window.removeEventListener('beforeinstallprompt', handler);
	});

	const own = $derived(data.disciplines.filter((d) => !d.key));
</script>

<svelte:head><title>{$t`Settings`} · Standblatt</title></svelte:head>

<h1 class="text-lg font-semibold tracking-tight">{$t`Settings`}</h1>

<section class="rounded-xl border bg-card p-4 shadow-sm">
	<h2 class="mb-3 font-semibold">{$t`Profile`}</h2>
	<form class="grid gap-3 sm:grid-cols-2" onsubmit={saveProfile}>
		<FormField label={$t`First name`}><Input bind:value={profile.firstname} required name="firstname" /></FormField>
		<FormField label={$t`Last name`}><Input bind:value={profile.lastname} required name="lastname" /></FormField>
		<FormField label={$t`Birthday`}><Input type="date" bind:value={profile.birthday} name="birthday" /></FormField>
		<FormField label={$t`License number`}><Input bind:value={profile.licenseNumber} name="licenseNumber" inputmode="numeric" /></FormField>
		<p class="text-xs text-muted-foreground sm:col-span-2">
			{$t`Name, birthday and license number connect your log with your competition results: Standblatt finds the results of competitions run with it, so you can add them to your log with one tap.`}
		</p>
		<Button type="submit" class="w-fit">{$t`Save`}</Button>
	</form>
</section>

<section class="rounded-xl border bg-card p-4 shadow-sm">
	<div class="flex flex-wrap items-center gap-3">
		<div class="min-w-0 flex-1">
			<h2 class="font-semibold">{$t`Account`}</h2>
			<p class="text-sm text-muted-foreground">{$t`Email, password, language, second factor.`}</p>
		</div>
		<Button variant="outline" size="sm" href="/account"><UserCog />{$t`Open`}</Button>
	</div>
</section>

<section class="rounded-xl border bg-card p-4 shadow-sm">
	<div class="mb-2 flex items-center justify-between">
		<h2 class="font-semibold">{$t`Disciplines`}</h2>
		<Button variant="outline" size="sm" onclick={() => openDiscipline()}><Plus />{$t`Own discipline`}</Button>
	</div>
	<ul class="divide-y text-sm">
		{#each data.disciplines as d (d.id)}
			<li class="flex items-center gap-3 py-2">
				<span class="flex-1">{$t(disciplineName(d))}</span>
				<span class="text-xs text-muted-foreground">{d.seriesCount} × {d.shotsPerSeries} · {d.decimal ? '10.9' : d.maxShot}</span>
				{#if !d.key}
					<Button variant="ghost" size="icon-xs" onclick={() => openDiscipline(d)} aria-label={$t`Edit`}><Pencil /></Button>
					<Button variant="ghost" size="icon-xs" onclick={() => deleteDiscipline(d)} aria-label={$t`Delete`}><Trash /></Button>
				{/if}
			</li>
		{/each}
	</ul>
	{#if !own.length}<p class="mt-2 text-xs text-muted-foreground">{$t`Shoot something else, like a field programme? Add your own discipline with its format and scoring.`}</p>{/if}
</section>

<section class="rounded-xl border bg-card p-4 shadow-sm">
	<div class="mb-2 flex items-center justify-between">
		<h2 class="font-semibold">{$t`Weapons`}</h2>
		<Button variant="outline" size="sm" onclick={() => openWeapon()}><Plus />{$t`Add weapon`}</Button>
	</div>
	{#if data.weapons.length}
		<ul class="divide-y text-sm">
			{#each data.weapons as w (w.id)}
				<li class="flex items-center gap-3 py-2">
					<span class="flex-1">{w.name}</span>
					<Button variant="ghost" size="icon-xs" onclick={() => openWeapon(w)} aria-label={$t`Rename`}><Pencil /></Button>
					<Button variant="ghost" size="icon-xs" onclick={() => deleteWeapon(w)} aria-label={$t`Delete`}><Trash /></Button>
				</li>
			{/each}
		</ul>
	{:else}
		<p class="text-xs text-muted-foreground">{$t`Shoot more than one gun? Add them, choose one when you save a session and compare them on the progress page.`}</p>
	{/if}
</section>

<section class="grid gap-3 rounded-xl border bg-card p-4 text-sm shadow-sm">
	<h2 class="font-semibold">{$t`App`}</h2>
	<div class="flex items-center gap-3"><span class="flex-1">{$t`Theme`}</span><ThemeToggle /></div>
	{#if version.data}<div class="flex items-center gap-3"><span class="flex-1">{$t`Version`}</span><span class="text-muted-foreground" data-testid="version">{version.data.version || 'dev'}</span></div>{/if}
	{#if !standalone}
		<div class="flex flex-wrap items-center gap-3">
			<Smartphone class="size-4 text-muted-foreground" />
			<span class="min-w-0 flex-1">
				{#if installPrompt}{$t`Install Standblatt on this device to start a session with one tap.`}
				{:else if ios}{$t`On iPhone: tap Share and then "Add to Home Screen".`}
				{:else}{$t`Add Standblatt to your home screen from the browser menu.`}{/if}
			</span>
			{#if installPrompt}<Button size="sm" variant="outline" onclick={() => installPrompt?.prompt()}>{$t`Install`}</Button>{/if}
		</div>
	{/if}
</section>


<section class="grid gap-3 rounded-xl border bg-card p-4 text-sm shadow-sm">
	<h2 class="font-semibold">{$t`Your data`}</h2>
	<div class="flex flex-wrap gap-2">
		<Button variant="outline" size="sm" href="/api/my/export" download><Download />{$t`Export all sessions (CSV)`}</Button>
	</div>
</section>

<Dialog.Root bind:open={disciplineOpen}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header><Dialog.Title>{editing ? $t`Edit discipline` : $t`Own discipline`}</Dialog.Title></Dialog.Header>
		<form class="grid gap-3" onsubmit={saveDiscipline}>
			<FormField label={$t`Name`}><Input bind:value={form.name} required maxlength={60} name="name" /></FormField>
			<div class="grid grid-cols-2 gap-3">
				<FormField label={$t`Series`}><Input type="number" min="1" max="20" bind:value={form.seriesCount} required name="seriesCount" /></FormField>
				<FormField label={$t`Shots per series`}><Input type="number" min="1" max="100" bind:value={form.shotsPerSeries} required name="shotsPerSeries" /></FormField>
			</div>
			<label class="flex items-center gap-2 text-sm"><input type="checkbox" bind:checked={form.decimal} /> {$t`Decimal scoring (10.9)`}</label>
			{#if !form.decimal}
				<FormField label={$t`Highest value of a shot`}><Input type="number" min="1" max="100" bind:value={form.maxShot} required name="maxShot" /></FormField>
			{/if}
			<FormField label={$t`Target face for tapping hits`}>
				<select class="rounded-md border bg-background px-2 py-1.5 text-sm" bind:value={form.target}>
					<option value="">{$t`None`}</option>
					<option value="air_pistol">{$t`Air pistol 10 m`}</option>
					<option value="air_rifle">{$t`Air rifle 10 m`}</option>
					<option value="pistol_precision">{$t`Pistol precision 25/50 m`}</option>
					<option value="rifle_50">{$t`Rifle 50 m`}</option>
					<option value="rifle_300">{$t`Rifle 300 m`}</option>
				</select>
			</FormField>
			<Button type="submit">{$t`Save`}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={weaponOpen}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header><Dialog.Title>{weapon ? $t`Rename weapon` : $t`Add weapon`}</Dialog.Title></Dialog.Header>
		<form class="grid gap-3" onsubmit={saveWeapon}>
			<FormField label={$t`Name`}><Input bind:value={weaponName} required maxlength={60} name="weaponName" placeholder={$t`e.g. Steyr EVO 10`} /></FormField>
			<Button type="submit">{$t`Save`}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
