<!-- src/lib/components/AchievementToast.svelte -->
<script lang="ts">
	import { fly, scale } from 'svelte/transition';
	import { elasticOut, cubicOut } from 'svelte/easing';
	import { onMount } from 'svelte';
	import { createEventDispatcher } from 'svelte';

	export let achievements: { icon: string; title: string; description: string }[] = [];

	const dispatch = createEventDispatcher();
	let visible = false;
	let currentIdx = 0;

	$: current = achievements[currentIdx];

	onMount(() => {
		if (achievements.length > 0) {
			setTimeout(() => { visible = true; }, 400);
		}
	});

	function dismiss() {
		visible = false;
		setTimeout(() => {
			currentIdx++;
			if (currentIdx < achievements.length) {
				setTimeout(() => { visible = true; }, 200);
			} else {
				dispatch('done');
			}
		}, 300);
	}
</script>

{#if visible && current}
	<!-- Backdrop -->
	<div
		in:fly={{ duration: 200 }}
		out:fly={{ duration: 200 }}
		class="fixed inset-0 bg-zinc-950/80 backdrop-blur-sm z-50 flex items-end justify-center pb-8 px-4"
		onclick={dismiss}
		role="button"
		tabindex="-1"
	>
		<!-- Card -->
		<div
			in:scale={{ duration: 500, easing: elasticOut, start: 0.7 }}
			out:fly={{ y: 30, duration: 250, easing: cubicOut }}
			class="w-full max-w-sm bg-zinc-900 border border-zinc-700 rounded-2xl p-6 text-center shadow-2xl"
			onclick|stopPropagation={() => {}}
			role="dialog"
		>
			<!-- Burst icon -->
			<div class="text-5xl mb-3 animate-bounce">{current.icon}</div>

			<p class="text-xs uppercase tracking-widest text-emerald-500 mb-1 font-bold">Conquista desbloqueada</p>
			<h2 class="text-xl font-black text-zinc-100 mb-2">{current.title}</h2>
			<p class="text-zinc-400 text-sm leading-relaxed mb-6">{current.description}</p>

			<button
				onclick={dismiss}
				class="w-full bg-emerald-500 hover:bg-emerald-400 text-zinc-950 font-bold
				       rounded-xl py-3 text-sm transition-all active:scale-95"
			>
				{currentIdx < achievements.length - 1 ? 'Próxima →' : 'Continuar'}
			</button>

			{#if achievements.length > 1}
				<div class="flex justify-center gap-1.5 mt-4">
					{#each achievements as _, i}
						<div class="w-1.5 h-1.5 rounded-full {i === currentIdx ? 'bg-emerald-400' : 'bg-zinc-700'}"></div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}
