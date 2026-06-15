<!-- src/routes/leaderboard/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import { token } from '$lib/stores/auth';
	import { scoreApi, type LeaderboardEntry, type Level } from '$lib/api/client';

	type Tab = 'global' | Level;

	let entries: LeaderboardEntry[] = [];
	let loading = true;
	let activeTab: Tab = 'global';

	const tabs: { id: Tab; label: string; icon: string }[] = [
		{ id: 'global',       label: 'Global',        icon: '🌍' },
		{ id: 'beginner',     label: 'Iniciante',     icon: '🌱' },
		{ id: 'intermediate', label: 'Intermediário', icon: '⚡' },
		{ id: 'advanced',     label: 'Avançado',      icon: '🔥' },
	];

	onMount(() => loadTab('global'));

	async function loadTab(tab: Tab) {
		if (!$token) { goto('/login'); return; }
		activeTab = tab;
		loading = true;
		try {
			entries = await scoreApi.leaderboard($token, tab === 'global' ? undefined : tab);
		} finally {
			loading = false;
		}
	}

	const medalEmoji = (rank: number) =>
		rank === 1 ? '🥇' : rank === 2 ? '🥈' : rank === 3 ? '🥉' : null;
</script>

<svelte:head><title>Ranking — GoRound</title></svelte:head>

<div class="py-4">
	<!-- Header -->
	<div in:fly={{ y: 12, duration: 280, easing: cubicOut }} class="mb-6">
		<h1 class="text-2xl font-black text-zinc-100">Ranking</h1>
		<p class="text-zinc-500 text-sm mt-0.5">Sobreviventes classificados</p>
	</div>

	<!-- Tabs -->
	<div in:fly={{ y: 12, duration: 280, delay: 60, easing: cubicOut }}
		class="grid grid-cols-2 gap-1.5 mb-6 bg-zinc-900 border border-zinc-800 rounded-xl p-1 sm:grid-cols-4">
		{#each tabs as tab}
			<button
				onclick={() => loadTab(tab.id)}
				class="flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium whitespace-nowrap transition-all
					{activeTab === tab.id
						? 'bg-zinc-700 text-zinc-100'
						: 'text-zinc-500 hover:text-zinc-300'}"
			>
				<span>{tab.icon}</span>
				{tab.label}
			</button>
		{/each}
	</div>

	<!-- Table -->
	{#if loading}
		<div class="py-16 text-center">
			<div class="w-6 h-6 rounded-full border-2 border-emerald-500 border-t-transparent animate-spin mx-auto"></div>
		</div>

	{:else if entries.length === 0}
		<div class="py-16 text-center text-zinc-600">
			<p class="text-3xl mb-3">🏜️</p>
			<p class="text-sm">Nenhum jogador neste ranking ainda.</p>
			<a href="/" class="text-emerald-500 text-sm mt-3 block">Seja o primeiro →</a>
		</div>

	{:else}
		<div in:fly={{ y: 12, duration: 300, easing: cubicOut }} class="flex flex-col gap-2">
			{#each entries as entry}
				{@const medal = medalEmoji(entry.rank)}
				<div class="flex items-center gap-3 rounded-xl px-4 py-3 border transition-colors
					{entry.is_current_user
						? 'bg-emerald-950/30 border-emerald-800/50'
						: 'bg-zinc-900 border-zinc-800'}">

					<!-- Rank -->
					<div class="w-8 text-center shrink-0">
						{#if medal}
							<span class="text-lg">{medal}</span>
						{:else}
							<span class="text-zinc-600 font-bold text-sm tabular-nums">#{entry.rank}</span>
						{/if}
					</div>

					<!-- Username -->
					<div class="flex-1 min-w-0">
						<p class="font-bold text-sm truncate
							{entry.is_current_user ? 'text-emerald-400' : 'text-zinc-200'}">
							{entry.username}
							{#if entry.is_current_user}<span class="text-zinc-600 font-normal"> (você)</span>{/if}
						</p>
						<p class="text-xs text-zinc-600 tabular-nums">
							{entry.sessions_played} {entry.sessions_played === 1 ? 'sessão' : 'sessões'}
							· streak {entry.best_streak}🔥
						</p>
					</div>

					<!-- Points -->
					<div class="text-right shrink-0">
						<p class="font-black text-sm text-yellow-400 tabular-nums">
							{entry.total_points.toLocaleString('pt-BR')}
						</p>
						<p class="text-xs text-zinc-700">pts</p>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
