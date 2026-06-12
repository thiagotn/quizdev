<!-- src/lib/components/OptionButton.svelte -->
<script lang="ts">
	import { scale } from 'svelte/transition';
	import { elasticOut } from 'svelte/easing';
	import CodeBlock from './CodeBlock.svelte';
	import type { Option } from '$lib/api/client';

	export let option: Option;
	export let state: 'idle' | 'selected' | 'correct' | 'wrong' | 'dimmed' = 'idle';
	export let onClick: () => void;
	export let disabled = false;

	const stateClasses: Record<typeof state, string> = {
		idle: 'border-zinc-700/80 bg-zinc-900 hover:border-emerald-500/50 hover:bg-zinc-800 cursor-pointer active:scale-[0.99]',
		selected: 'border-zinc-500 bg-zinc-800 cursor-pointer',
		correct: 'border-emerald-500 bg-emerald-950/50 ring-1 ring-emerald-500/30',
		wrong: 'border-red-500 bg-red-950/40 ring-1 ring-red-500/20',
		dimmed: 'border-zinc-800 bg-zinc-900/40 opacity-40',
	};

	// Bounce animation key on state change to correct
	let bounceKey = 0;
	$: if (state === 'correct') bounceKey++;
</script>

<button
	onclick={onClick}
	{disabled}
	class="w-full text-left rounded-xl border p-3 text-sm overflow-hidden
	       transition-all duration-200 {stateClasses[state]}"
>
	{#if state === 'correct'}
		<!-- re-mount to trigger animation -->
		{#key bounceKey}
			<div in:scale={{ duration: 300, easing: elasticOut, start: 0.95 }}>
				<CodeBlock code={option.code_snippet} lang="go" />
			</div>
		{/key}
	{:else}
		<CodeBlock code={option.code_snippet} lang="go" />
	{/if}
</button>
