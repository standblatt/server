<script lang="ts">
	import { plural, t } from 'lingui-for-svelte/macro';
	import { page } from '$app/state';
	import BarChart from '@lucide/svelte/icons/chart-column';
	import Building from '@lucide/svelte/icons/building-2';
	import CalendarClock from '@lucide/svelte/icons/calendar-clock';
	import CalendarDays from '@lucide/svelte/icons/calendar-days';
	import Check from '@lucide/svelte/icons/check';
	import CircleHelp from '@lucide/svelte/icons/circle-help';
	import ChevronsUpDown from '@lucide/svelte/icons/chevrons-up-down';
	import Crosshair from '@lucide/svelte/icons/crosshair';
	import LayoutDashboard from '@lucide/svelte/icons/layout-dashboard';
	import ListOrdered from '@lucide/svelte/icons/list-ordered';
	import LogOut from '@lucide/svelte/icons/log-out';
	import Repeat from '@lucide/svelte/icons/repeat';
	import Settings from '@lucide/svelte/icons/settings-2';
	import UserCog from '@lucide/svelte/icons/user-cog';
	import BadgeCheck from '@lucide/svelte/icons/badge-check';
	import Dumbbell from '@lucide/svelte/icons/dumbbell';
	import ShieldCheck from '@lucide/svelte/icons/shield-check';
	import { authOptionsQuery, versionQuery } from '$lib/me';
	import { createQuery } from '@tanstack/svelte-query';
	import { api, call } from '$lib/api';
	import { keys } from '$lib/tenant';
	import Users from '@lucide/svelte/icons/users';
	import type { Event, Me } from '$lib/api';
	import { liveLabel, roleLabel } from '$lib/labels';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { formatDays, rangeNames } from '$lib/format';
	import Languages from '@lucide/svelte/icons/languages';
	import { currentLocale, localeNames } from '$lib/locale.svelte';
	import Logo from '$lib/Logo.svelte';
	import Wordmark from '$lib/Wordmark.svelte';

	let {
		me,
		tenantId,
		event,
		live,
		onlogout,
		onlocale
	}: {
		me: Me;
		tenantId: number;
		event?: Event;
		live?: boolean;
		onlogout: () => void;
		onlocale: (locale: string) => void;
	} = $props();

	const options = authOptionsQuery();
	const version = versionQuery();
	const membership = $derived(me.tenants.find((t) => t.tenantId === tenantId)!);
	const isAdmin = $derived(membership.role === 'admin');
	const base = $derived(`/t/${tenantId}`);
	const eventBase = $derived(`${base}/events/${event?.id}`);
	const path = $derived(page.url.pathname);

	const bookings = createQuery(() => ({
		queryKey: [...keys.booking(tenantId, event?.id ?? 0), 'requests'],
		queryFn: () => call(api.GET('/api/tenants/{tenantId}/events/{eventId}/booking-requests', { params: { path: { tenantId, eventId: event!.id } } })),
		enabled: !!event
	}));
	const pending = $derived((bookings.data ?? []).filter((r) => r.status === 'pending').length);

	const eventItems = $derived([
		{ href: eventBase, label: $t`Dashboard`, icon: LayoutDashboard, active: path === eventBase },
		{ href: `${eventBase}/reservations`, label: $t`Reservations`, icon: CalendarClock, badge: pending },
		{ href: `${eventBase}/rankings`, label: $t`Rankings`, icon: ListOrdered },
		{ href: `${eventBase}/statistics`, label: $t`Statistics`, icon: BarChart },
		{ href: `${eventBase}/setup`, label: $t`Setup`, icon: Settings }
	]);
	const clubItems = $derived([
		{ href: base, label: $t`Events`, icon: CalendarDays, active: path === base },
		{ href: `${base}/shooters`, label: $t`Shooters`, icon: Users },
		{ href: `${base}/clubs`, label: $t`Clubs`, icon: Building },
		...(isAdmin
			? [
					{ href: `${base}/ranges`, label: $t`Ranges`, icon: Crosshair },
					{ href: `${base}/users`, label: $t`Users`, icon: UserCog }
				]
			: [])
	]);
	const initials = $derived(me.user.name.slice(0, 2).toUpperCase());
	const sidebar = Sidebar.useSidebar();
</script>

{#snippet nav(items: { href: string; label: string; icon: typeof Users; active?: boolean; badge?: number }[])}
	<Sidebar.Menu>
		{#each items as item (item.href)}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton isActive={item.active ?? path.startsWith(item.href)} tooltipContent={item.label}>
					{#snippet child({ props })}
						<a href={item.href} {...props}><item.icon /><span>{item.label}</span></a>
					{/snippet}
				</Sidebar.MenuButton>
				{#if item.badge}<Sidebar.MenuBadge aria-label={$plural(item.badge, { one: '# open request', other: '# open requests' })}>{item.badge}</Sidebar.MenuBadge>{/if}
			</Sidebar.MenuItem>
		{/each}
	</Sidebar.Menu>
{/snippet}

<Sidebar.Root collapsible="icon">
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuButton size="lg" {...props}>
								<div class="flex size-8 shrink-0 items-center justify-center text-sidebar-primary"><Logo class="size-7" /></div>
								<div class="grid flex-1 text-left leading-tight">
									<Wordmark class="truncate text-base text-sidebar-accent-foreground" />
									<span class="truncate text-xs">{membership.tenantName}</span>
								</div>
								{#if me.tenants.length > 1}<ChevronsUpDown class="ml-auto" />{/if}
							</Sidebar.MenuButton>
						{/snippet}
					</DropdownMenu.Trigger>
					{#if me.tenants.length > 1}
						<DropdownMenu.Content class="w-64" align="start">
							<DropdownMenu.Label class="text-xs text-muted-foreground">{$t`Clubs`}</DropdownMenu.Label>
							{#each me.tenants as t (t.tenantId)}
								<DropdownMenu.Item>
									{#snippet child({ props })}
										<a href="/t/{t.tenantId}" {...props}>
											{t.tenantName}
											{#if t.tenantId === tenantId}<Check class="ml-auto" />{/if}
										</a>
									{/snippet}
								</DropdownMenu.Item>
							{/each}
						</DropdownMenu.Content>
					{/if}
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>

	<Sidebar.Content>
		{#if event}
			<Sidebar.Group>
				<Sidebar.GroupLabel>{$t`Event`}</Sidebar.GroupLabel>
				<div data-sidebar="event" class="rounded-md bg-sidebar-accent/40 px-2 py-1.5 text-sidebar-foreground group-data-[collapsible=icon]:hidden">
					<div class="flex items-center gap-2">
						{#if live !== undefined}
							<span class="size-2 shrink-0 rounded-full {live ? 'bg-emerald-400' : 'bg-red-400'}" title={$t(liveLabel(live))}></span>
						{/if}
						<span class="truncate font-medium text-sidebar-accent-foreground">{event.name}</span>
					</div>
					<span class="block truncate text-xs">{formatDays(event.days, event.startsOn)} · {rangeNames(event.ranges)}{event.archivedAt ? ` · ${$t`closed`}` : ''}</span>
					<a href={base} class="mt-1 inline-flex items-center gap-1 text-xs text-sidebar-foreground/80 underline-offset-2 hover:underline"><Repeat class="size-3" />{$t`Change event`}</a>
				</div>
				<div class="mt-1">{@render nav(eventItems)}</div>
			</Sidebar.Group>
		{/if}

		<Sidebar.Group>
			<Sidebar.GroupLabel>{$t`Club`}</Sidebar.GroupLabel>
			{@render nav(clubItems)}
		</Sidebar.Group>
	</Sidebar.Content>

	<Sidebar.Footer>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton tooltipContent={$t`Help`}>
					{#snippet child({ props })}<a href="/help" target="_blank" rel="noopener" {...props}><CircleHelp /><span>{$t`Help`}</span></a>{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
			<Sidebar.MenuItem>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuButton size="lg" aria-label={$t`User menu`} {...props}>
								<Avatar.Root class="size-8 rounded-md">
									<Avatar.Fallback class="rounded-md bg-sidebar-accent text-xs text-sidebar-accent-foreground">{initials}</Avatar.Fallback>
								</Avatar.Root>
								<div class="grid flex-1 text-left leading-tight">
									<span class="truncate font-medium text-sidebar-accent-foreground">{me.user.name}</span>
									<span class="truncate text-xs">{$t(roleLabel(membership.role))}</span>
								</div>
								<ChevronsUpDown class="ml-auto" />
							</Sidebar.MenuButton>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content class="w-56" side={sidebar.isMobile ? 'top' : 'right'} align="end">
						<DropdownMenu.Label>
							<div class="grid leading-tight">
								<span>{me.user.name}</span>
								<span class="text-xs font-normal text-muted-foreground">{membership.tenantName} · {$t(roleLabel(membership.role))}</span>
							</div>
						</DropdownMenu.Label>
						<DropdownMenu.Separator />
						<DropdownMenu.Sub>
							<DropdownMenu.SubTrigger><Languages />{$t`Language`}</DropdownMenu.SubTrigger>
							<DropdownMenu.SubContent>
								<DropdownMenu.RadioGroup value={currentLocale()} onValueChange={onlocale}>
									{#each Object.entries(localeNames) as [value, label] (value)}<DropdownMenu.RadioItem {value}>{label}</DropdownMenu.RadioItem>{/each}
								</DropdownMenu.RadioGroup>
							</DropdownMenu.SubContent>
						</DropdownMenu.Sub>
						{#if me.tenants.length > 1}
							<DropdownMenu.Item>
								{#snippet child({ props })}<a href="/" {...props}><Repeat />{$t`Switch club`}</a>{/snippet}
							</DropdownMenu.Item>
						{/if}
						<DropdownMenu.Item>
							{#snippet child({ props })}<a href="/my" {...props}><Dumbbell />{$t`My training`}</a>{/snippet}
						</DropdownMenu.Item>
						<DropdownMenu.Item>
							{#snippet child({ props })}<a href="/account" {...props}><ShieldCheck />{$t`Account`}</a>{/snippet}
						</DropdownMenu.Item>
						<DropdownMenu.Item onSelect={onlogout}><LogOut />{$t`Log out`}</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
		{#if version.data}
			<p class="truncate px-2 text-[11px] text-sidebar-foreground/60 group-data-[collapsible=icon]:hidden" title={$t`Version`}>Standblatt {version.data.version || 'dev'}</p>
		{/if}
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
