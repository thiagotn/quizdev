<!-- src/routes/quiz/[level]/+page.svelte -->
<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';

	import { quiz, currentQuestion, currentOptions, progress } from '$lib/stores/quiz';
	import { token } from '$lib/stores/auth';
	import { quizApi, type Level } from '$lib/api/client';

	import Timer from '$lib/components/Timer.svelte';
	import StreakBadge from '$lib/components/StreakBadge.svelte';
	import OptionButton from '$lib/components/OptionButton.svelte';
	import FeedbackPanel from '$lib/components/FeedbackPanel.svelte';

	const TIMER_SECONDS = 20;
	const TIME_BONUS_THRESHOLD = 10; // seconds remaining to earn bonus

	let loading = true;
	let error = '';
	let selectedOptionId: string | null = null;
	let submitting = false;
	let timerRef: Timer;
	let timerPaused = false;
	let timerSecondsLeft = TIMER_SECONDS;
	let questionKey = 0; // force re-mount of timer on question change

	$: level = $page.params.level as Level;
	$: isLast = $quiz.currentIndex >= ($quiz.session?.total_questions ?? 1) - 1;

	onMount(async () => {
		if (!$token) { goto('/login'); return; }
		try {
			const session = await quizApi.createSession('go', level, $token);
			quiz.startSession(session);
		} catch (e: any) {
			error = e.message ?? 'Falha ao carregar quiz';
		} finally {
			loading = false;
		}
	});

	function optionState(optId: string): 'idle' | 'selected' | 'correct' | 'wrong' | 'dimmed' {
		if (!$quiz.answered) {
			return selectedOptionId === optId ? 'selected' : 'idle';
		}
		const correctId = $quiz.lastAnswer?.correct_option.id;
		if (optId === correctId) return 'correct';
		if (optId === selectedOptionId) return 'wrong';
		return 'dimmed';
	}

	async function handleSelect(optionId: string) {
		if ($quiz.answered || submitting || !$quiz.session || !$token) return;
		selectedOptionId = optionId;
		submitting = true;
		timerPaused = true;
		timerRef?.stop();

		try {
			const result = await quizApi.submitAnswer(
				$quiz.session.id,
				$currentQuestion!.id,
				optionId,
				$token,
				timerSecondsLeft
			);
			quiz.recordAnswer(result, false);
		} catch (e: any) {
			error = e.message;
		} finally {
			submitting = false;
		}
	}

	async function handleTimeout() {
		if ($quiz.answered || !$quiz.session || !$token) return;
		timerPaused = true;
		submitting = true;

		// Auto-submit wrong answer using first option as placeholder
		const dummyOptionId = $currentOptions[0]?.id ?? '';
		if (!dummyOptionId) { submitting = false; return; }

		selectedOptionId = null;
		try {
			const result = await quizApi.submitAnswer(
				$quiz.session.id,
				$currentQuestion!.id,
				dummyOptionId,
				$token
			);
			// Override: mark as incorrect regardless since it timed out
			quiz.recordAnswer({ ...result, is_correct: false, points_earned: 0 }, true);
		} catch {
			quiz.recordAnswer(
				{
					is_correct: false,
					explanation: 'Tempo esgotado.',
					correct_option: $currentOptions.find(o => o) ?? $currentOptions[0],
					points_earned: 0,
					current_streak: 0,
					total_score: $quiz.session?.score ?? 0,
				},
				true
			);
		} finally {
			submitting = false;
		}
	}

	async function handleNext() {
		if ($quiz.finished) {
			await quizApi.getResult($quiz.session!.id, $token!);
			goto(`/results/${$quiz.session!.id}`);
			return;
		}
		selectedOptionId = null;
		timerPaused = false;
		timerSecondsLeft = TIMER_SECONDS;
		questionKey++;
		quiz.nextQuestion();
	}
</script>

<svelte:head><title>{level} — GoRound</title></svelte:head>

<!-- ─── Loading ─────────────────────────────────────────── -->
{#if loading}
	<div class="flex flex-col items-center justify-center py-24 gap-3">
		<div class="w-8 h-8 rounded-full border-2 border-emerald-500 border-t-transparent animate-spin"></div>
		<p class="text-zinc-500 text-sm">Preparando questões...</p>
	</div>

<!-- ─── Error ────────────────────────────────────────────── -->
{:else if error}
	<div class="text-center py-20">
		<p class="text-4xl mb-4">💥</p>
		<p class="text-red-400 mb-6 text-sm">{error}</p>
		<a href="/" class="text-emerald-400 underline text-sm">Voltar ao início</a>
	</div>

<!-- ─── Quiz ─────────────────────────────────────────────── -->
{:else if $currentQuestion}
	<div class="py-2">

		<!-- Top bar: progress + score + streak -->
		<div class="flex items-center justify-between mb-3">
			<div class="flex items-center gap-2">
				<span class="text-xs text-zinc-500">
					{$progress.current}<span class="text-zinc-700">/</span>{$progress.total}
				</span>
				<StreakBadge streak={$quiz.streak} />
			</div>
			<span class="text-sm font-black text-emerald-400 tabular-nums">
				{$quiz.session?.score ?? 0} <span class="text-zinc-600 font-normal text-xs">pts</span>
			</span>
		</div>

		<!-- Progress bar -->
		<div class="h-1 bg-zinc-800 rounded-full overflow-hidden mb-6">
			<div
				class="h-full bg-emerald-500 rounded-full transition-all duration-500"
				style="width: {$progress.percent}%"
			></div>
		</div>

		<!-- Question card -->
		{#key questionKey}
			<div in:fly={{ y: 16, duration: 280, easing: cubicOut }} class="mb-5">
				<div class="flex items-start gap-4 mb-5">
					<!-- Timer -->
					{#if !$quiz.answered}
						<div class="shrink-0 mt-0.5">
							<Timer
								bind:this={timerRef}
								bind:timeRemaining={timerSecondsLeft}
								seconds={TIMER_SECONDS}
								paused={timerPaused}
								on:timeout={handleTimeout}
							/>
						</div>
					{/if}

					<p class="text-zinc-100 font-medium leading-relaxed text-sm flex-1">
						{$currentQuestion.title}
					</p>
				</div>

				<!-- Options -->
				<div class="flex flex-col gap-2.5">
					{#each $currentOptions as option (option.id)}
						<OptionButton
							{option}
							state={optionState(option.id)}
							onClick={() => handleSelect(option.id)}
							disabled={$quiz.answered || submitting}
						/>
					{/each}
				</div>
			</div>

			<!-- Feedback -->
			{#if $quiz.answered && $quiz.lastAnswer}
				<FeedbackPanel
					result={$quiz.lastAnswer}
					timedOut={$quiz.timedOut}
					onNext={handleNext}
					{isLast}
				/>
			{/if}
		{/key}

	</div>
{/if}
