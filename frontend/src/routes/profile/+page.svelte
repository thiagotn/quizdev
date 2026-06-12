<!-- src/routes/profile/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import { token } from '$lib/stores/auth';
	import { scoreApi, type ProfileResponse, type Level } from '$lib/api/client';
	import MiniChart from '$lib/components/MiniChart.svelte';
	import LevelCard from '$lib/components/LevelCard.svelte';

	let profile: ProfileResponse | null = null;
	let loading = true;

	const LEVEL_ORDER: Level[] = ['beginner', 'intermediate', 'advanced'];
	const levelMeta: Record<Level, { label: string; icon: string; color: string }> = {
		beginner:     { label: 'Iniciante',     icon: '🌱', color: 'emerald' },
		intermediate: { label: 'Intermediário', icon: '⚡', color: 'yellow'  },
		advanced:     { label: 'Avançado',      icon: '🔥', color: 'red'     },
	};

	onMount(async () => {
		if (!$token) { goto('/login'); return; }
		try {
			profile = await scoreApi.myProfile($token);
		} finally {
			loading = false;
		}
	});

	$: orderedScores = LEVEL_ORDER.map(lvl =>
		profile?.scores.find(s => s.level === lvl) ?? null
	);

	// Last 10 sessions reversed for chart (oldest first)
	$: chartHistory = [...(profile?.history ?? [])].slice(0, 10).reverse();

	function formatDate(iso: string) {
		return new Date(iso).toLocaleDateString('pt-BR', { day: '2-digit', month: 'short' });
	}
</script>

<svelte:head><title>Perfil — GoRound</title></svelte:head>

{#if loading}
	<div class="flex flex-col items-center justify-center py-24 gap-3">
		<div class="w-8 h-8 rounded-full border-2 border-emerald-500 border-t-transparent animate-spin"></div>
		<p class="text-zinc-500 text-sm">Carregando perfil...</p>
	</div>

{:else if profile}
	<!-- ── Avatar + rank ───────────────────────────────── -->
	<div in:fly={{ y: 16, duration: 300, easing: cubicOut }} class="flex items-center gap-4 mb-6 pt-2">
		<div class="relative shrink-0">
			<div class="w-16 h-16 rounded-2xl bg-zinc-800 border border-zinc-700 flex items-center justify-center text-2xl font-black text-emerald-400">
				{profile.user.username[0].toUpperCase()}
			</div>
			{#if profile.global_rank > 0}
				<div class="absolute -bottom-1.5 -right-1.5 bg-yellow-500 text-zinc-950 text-xs font-black px-1.5 py-0.5 rounded-full leading-none">
					#{profile.global_rank}
				</div>
			{/if}
		</div>
		<div class="flex-1 min-w-0">
			<h1 class="text-lg font-black text-zinc-100 truncate">{profile.user.username}</h1>
			<p class="text-zinc-500 text-xs truncate">{profile.user.email}</p>
		</div>
		<a href="/leaderboard" class="text-xs text-emerald-500 hover:text-emerald-400 transition-colors shrink-0">
			Ranking →
		</a>
	</div>

	<!-- ── Global stats ────────────────────────────────── -->
	<div in:fly={{ y: 16, duration: 300, delay: 60, easing: cubicOut }}
		class="grid grid-cols-3 gap-2.5 mb-6">
		{#each [
			{ label: 'Pontos',   value: profile.total_points.toLocaleString('pt-BR'), color: 'text-yellow-400' },
			{ label: 'Sessões',  value: profile.history.length,                       color: 'text-emerald-400' },
			{ label: 'Ranking',  value: profile.global_rank > 0 ? `#${profile.global_rank}` : '—', color: 'text-zinc-300' },
		] as stat}
			<div class="bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-center">
				<div class="text-xl font-black {stat.color} tabular-nums">{stat.value}</div>
				<div class="text-xs text-zinc-600 mt-0.5">{stat.label}</div>
			</div>
		{/each}
	</div>

	<!-- ── Score chart ─────────────────────────────────── -->
	{#if chartHistory.length >= 2}
		<div in:fly={{ y: 16, duration: 300, delay: 120, easing: cubicOut }}
			class="bg-zinc-900 border border-zinc-800 rounded-xl p-4 mb-6">
			<div class="flex justify-between items-center mb-3">
				<p class="text-xs font-medium text-zinc-400 uppercase tracking-wider">Evolução de pontos</p>
				<p class="text-xs text-zinc-600">{chartHistory.length} sessões</p>
			</div>
			<MiniChart data={chartHistory.map(s => s.score)} />
			<div class="flex justify-between mt-2 text-xs text-zinc-700">
				<span>{formatDate(chartHistory[0].finished_at)}</span>
				<span>{formatDate(chartHistory[chartHistory.length - 1].finished_at)}</span>
			</div>
		</div>
	{/if}

	<!-- ── Level cards ─────────────────────────────────── -->
	<div in:fly={{ y: 16, duration: 300, delay: 180, easing: cubicOut }} class="mb-6">
		<p class="text-xs uppercase tracking-widest text-zinc-600 mb-3">Progresso por nível</p>
		<div class="flex flex-col gap-2.5">
			{#each LEVEL_ORDER as lvl, i}
				{@const score = orderedScores[i]}
				{@const meta = levelMeta[lvl]}
				<LevelCard level={lvl} {meta} {score} />
			{/each}
		</div>
	</div>

	<!-- ── Session history ─────────────────────────────── -->
	{#if profile.history.length > 0}
		<div in:fly={{ y: 16, duration: 300, delay: 240, easing: cubicOut }} class="mb-6">
			<p class="text-xs uppercase tracking-widest text-zinc-600 mb-3">Histórico recente</p>
			<div class="flex flex-col gap-2">
				{#each profile.history.slice(0, 8) as session}
					{@const meta = levelMeta[session.level] ?? { icon: '📊', label: session.level, color: 'zinc' }}
					<div class="flex items-center gap-3 bg-zinc-900 border border-zinc-800 rounded-xl px-4 py-3">
						<span class="text-lg">{meta.icon}</span>
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium text-zinc-200">{meta.label}</p>
							<p class="text-xs text-zinc-600">
								{session.correct_answers}/{session.total_questions} corretas · {session.accuracy}%
							</p>
						</div>
						<div class="text-right shrink-0">
							<p class="text-sm font-black text-yellow-400">{session.score} pts</p>
							<p class="text-xs text-zinc-700">
								{session.finished_at ? formatDate(session.finished_at) : ''}
							</p>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{:else}
		<div class="text-center py-12">
			<p class="text-4xl mb-3">🎯</p>
			<p class="text-zinc-500 text-sm mb-4">Nenhum quiz concluído ainda.</p>
		</div>
	{/if}

	<a href="/" class="w-full block text-center bg-emerald-500 hover:bg-emerald-400
		text-zinc-950 font-bold rounded-xl py-3.5 text-sm transition-all active:scale-95">
		Jogar agora
	</a>
{/if}
