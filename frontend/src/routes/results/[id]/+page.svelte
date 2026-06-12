<!-- src/routes/results/[id]/+page.svelte -->
<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount, tick } from 'svelte';
	import { fly, scale } from 'svelte/transition';
	import { cubicOut, elasticOut } from 'svelte/easing';
	import { tweened } from 'svelte/motion';
	import { token } from '$lib/stores/auth';
	import { quiz } from '$lib/stores/quiz';
	import { quizApi, scoreApi, type Session } from '$lib/api/client';
	import { detectAchievements, type Achievement } from '$lib/achievements';
	import { gradeResult } from '$lib/utils/labels';
	import AchievementToast from '$lib/components/AchievementToast.svelte';

	let session: Session | null = null;
	let loading = true;
	let ready = false;
	let achievements: Achievement[] = [];
	let showAchievements = true;

	const tweenedScore    = tweened(0, { duration: 1200, easing: cubicOut });
	const tweenedAccuracy = tweened(0, { duration: 1000, easing: cubicOut });

	$: accuracy = session
		? Math.round((session.correct_answers / session.total_questions) * 100)
		: 0;

	$: grade = gradeResult(accuracy);

	$: barColor =
		accuracy >= 70 ? 'bg-emerald-500'
		: accuracy >= 50 ? 'bg-yellow-500'
		: 'bg-red-500';

	const levelLabel: Record<string, string> = {
		beginner: 'Round I', intermediate: 'Round II', advanced: 'Round III'
	};

	onMount(async () => {
		if (!$token) { goto('/login'); return; }
		try {
			session = await quizApi.getResult($page.params.id, $token);

			// Check if first ever session
			let isFirst = false;
			try {
				const profile = await scoreApi.myProfile($token);
				isFirst = profile.history.length <= 1;
			} catch {}

			achievements = detectAchievements(session, isFirst);

			await tick();
			ready = true;
			tweenedScore.set(session.score);
			tweenedAccuracy.set(accuracy);
		} finally {
			loading = false;
		}
	});

	function playAgain() {
		if (!session) return;
		quiz.reset();
		goto(`/quiz/${session.level}`);
	}
</script>

<svelte:head><title>Resultado — GoRound</title></svelte:head>

<!-- Achievement overlay -->
{#if showAchievements && achievements.length > 0}
	<AchievementToast
		{achievements}
		on:done={() => { showAchievements = false; }}
	/>
{/if}

{#if loading}
	<div class="flex flex-col items-center justify-center py-24 gap-3">
		<div class="w-8 h-8 rounded-full border-2 border-emerald-500 border-t-transparent animate-spin"></div>
		<p class="text-zinc-500 text-sm">Calculando resultado...</p>
	</div>

{:else if session}
	<div class="py-6">

		<!-- Hero -->
		<div in:scale={{ duration: 500, easing: elasticOut, start: 0.8 }} class="text-center mb-7">
			<div class="text-6xl mb-3">{grade.emoji}</div>
			<p class="text-zinc-600 text-xs uppercase tracking-widest mb-1">
				{levelLabel[session.level] ?? session.level} · Go
			</p>
			<h1 class="text-4xl font-black text-zinc-100 tabular-nums">
				{Math.round($tweenedScore)}<span class="text-emerald-500 text-2xl ml-1">pts</span>
			</h1>
			<p class="text-zinc-500 text-sm mt-1">{grade.label}</p>
		</div>

		<!-- Stats grid -->
		<div in:fly={{ y: 20, duration: 400, delay: 200, easing: cubicOut }}
			class="grid grid-cols-3 gap-2.5 mb-5">
			{#each [
				{ label: 'Acertos',  value: session.correct_answers,                           color: 'text-emerald-400' },
				{ label: 'Erros',    value: session.total_questions - session.correct_answers, color: 'text-red-400'     },
				{ label: 'Questões', value: session.total_questions,                           color: 'text-zinc-300'    },
			] as stat}
				<div class="bg-zinc-900 border border-zinc-800 rounded-xl p-4 text-center">
					<div class="text-2xl font-black {stat.color} tabular-nums">{stat.value}</div>
					<div class="text-xs text-zinc-600 mt-1">{stat.label}</div>
				</div>
			{/each}
		</div>

		<!-- Accuracy bar -->
		<div in:fly={{ y: 20, duration: 400, delay: 320, easing: cubicOut }}
			class="bg-zinc-900 border border-zinc-800 rounded-xl p-4 mb-5">
			<div class="flex justify-between text-xs text-zinc-500 mb-2.5">
				<span class="text-zinc-400 font-medium">Precisão</span>
				<span class="font-black text-zinc-200 tabular-nums">{Math.round($tweenedAccuracy)}%</span>
			</div>
			<div class="h-2.5 bg-zinc-800 rounded-full overflow-hidden">
				{#if ready}
					<div
						class="h-full rounded-full {barColor}"
						style="width: {accuracy}%; transition: width 1s cubic-bezier(0.16, 1, 0.3, 1)"
					></div>
				{/if}
			</div>
			<div class="flex justify-between mt-1.5 text-zinc-700 text-xs">
				<span>0%</span>
				<span>50%</span>
				<span>70%</span>
				<span>100%</span>
			</div>
		</div>

		<!-- Achievements teaser (if any but dismissed) -->
		{#if achievements.length > 0 && !showAchievements}
			<div in:fly={{ y: 12, duration: 300, easing: cubicOut }}
				class="flex items-center gap-3 bg-zinc-900 border border-zinc-700 rounded-xl px-4 py-3 mb-5">
				<span class="text-xl">{achievements[0].icon}</span>
				<div class="flex-1">
					<p class="text-xs font-bold text-emerald-400">Conquista desbloqueada</p>
					<p class="text-sm font-medium text-zinc-200">{achievements[0].title}</p>
				</div>
				{#if achievements.length > 1}
					<span class="text-xs text-zinc-600">+{achievements.length - 1}</span>
				{/if}
			</div>
		{/if}

		<!-- CTAs -->
		<div in:fly={{ y: 20, duration: 400, delay: 440, easing: cubicOut }} class="flex flex-col gap-3">
			<button
				onclick={playAgain}
				class="w-full bg-emerald-500 hover:bg-emerald-400 active:scale-95
				       text-zinc-950 font-bold rounded-xl py-3.5 text-sm transition-all"
			>
				Jogar novamente
			</button>

			<div class="grid grid-cols-2 gap-2.5">
				<a href="/profile"
					class="block text-center border border-zinc-700 hover:border-zinc-500
					       text-zinc-300 font-medium rounded-xl py-3 text-sm transition-all">
					Perfil
				</a>
				<a href="/leaderboard"
					class="block text-center border border-zinc-700 hover:border-zinc-500
					       text-zinc-300 font-medium rounded-xl py-3 text-sm transition-all">
					Ranking
				</a>
			</div>

			<a href="/" class="text-center text-zinc-600 hover:text-zinc-400 text-sm transition-colors py-1">
				Trocar nível
			</a>
		</div>

	</div>
{/if}
