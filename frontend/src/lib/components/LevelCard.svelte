<!-- src/lib/components/LevelCard.svelte -->
<script lang="ts">
	import type { Score, Level } from '$lib/api/client';

	export let level: Level;
	export let meta: { label: string; icon: string; color: string };
	export let score: Score | null;

	// Tier thresholds for each level
	const TIERS: Record<Level, { label: string; min: number; max: number }[]> = {
		beginner:     [{ label: 'Bronze', min: 0, max: 100 }, { label: 'Prata', min: 100, max: 300 }, { label: 'Ouro', min: 300, max: 600 }, { label: 'Platina', min: 600, max: Infinity }],
		intermediate: [{ label: 'Bronze', min: 0, max: 150 }, { label: 'Prata', min: 150, max: 400 }, { label: 'Ouro', min: 400, max: 800 }, { label: 'Platina', min: 800, max: Infinity }],
		advanced:     [{ label: 'Bronze', min: 0, max: 200 }, { label: 'Prata', min: 200, max: 500 }, { label: 'Ouro', min: 500, max: 1000 }, { label: 'Platina', min: 1000, max: Infinity }],
	};

	const colorMap: Record<string, string> = {
		emerald: 'bg-emerald-500',
		yellow:  'bg-yellow-500',
		red:     'bg-red-500',
	};

	$: tiers = TIERS[level];
	$: pts = score?.total_points ?? 0;
	$: currentTier = tiers.findLast(t => pts >= t.min) ?? tiers[0];
	$: nextTier = tiers[tiers.indexOf(currentTier) + 1] ?? null;
	$: pct = nextTier
		? Math.min(100, Math.round(((pts - currentTier.min) / (nextTier.min - currentTier.min)) * 100))
		: 100;

	$: tierEmoji = currentTier.label === 'Platina' ? '💎'
		: currentTier.label === 'Ouro' ? '🥇'
		: currentTier.label === 'Prata' ? '🥈'
		: '🥉';
</script>

<div class="bg-zinc-900 border border-zinc-800 rounded-xl p-4">
	<div class="flex items-center gap-3 mb-3">
		<span class="text-xl">{meta.icon}</span>
		<div class="flex-1">
			<p class="text-sm font-bold text-zinc-200">{meta.label}</p>
			{#if score}
				<p class="text-xs text-zinc-600">{score.sessions_played} sessão{score.sessions_played !== 1 ? 'ões' : ''}</p>
			{:else}
				<p class="text-xs text-zinc-700">Não jogado ainda</p>
			{/if}
		</div>
		<div class="text-right">
			{#if score}
				<p class="text-sm font-black text-yellow-400 tabular-nums">{pts.toLocaleString('pt-BR')} pts</p>
				<p class="text-xs text-zinc-600">{tierEmoji} {currentTier.label}</p>
			{:else}
				<a href="/quiz/{level}" class="text-xs text-emerald-500 hover:text-emerald-400">Começar →</a>
			{/if}
		</div>
	</div>

	{#if score}
		<!-- Tier progress bar -->
		<div class="h-1.5 bg-zinc-800 rounded-full overflow-hidden mb-2">
			<div
				class="h-full rounded-full {colorMap[meta.color] ?? 'bg-emerald-500'}"
				style="width: {pct}%; transition: width 0.8s cubic-bezier(0.16,1,0.3,1)"
			></div>
		</div>
		<div class="flex justify-between text-xs text-zinc-700">
			<span>{currentTier.label}</span>
			{#if nextTier}
				<span class="text-zinc-600">{pct}% → {nextTier.label}</span>
			{:else}
				<span class="text-yellow-600">Nível máximo 💎</span>
			{/if}
		</div>

		<!-- Mini stats row -->
		<div class="flex gap-4 mt-3 pt-3 border-t border-zinc-800/60 text-xs text-zinc-600">
			<span>🔥 streak <span class="text-zinc-400 font-bold">{score.best_streak}</span></span>
			<span>🎯 precisão <span class="text-zinc-400 font-bold">{score.best_accuracy}%</span></span>
		</div>
	{/if}
</div>
