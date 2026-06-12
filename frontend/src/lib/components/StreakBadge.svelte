<!-- src/lib/components/StreakBadge.svelte -->
<script lang="ts">
	import { fly, scale } from 'svelte/transition';
	import { elasticOut } from 'svelte/easing';

	export let streak = 0;

	$: show = streak >= 2;
	$: label =
		streak >= 7 ? '🔥🔥 INSANO'
		: streak >= 5 ? '🔥 Em chamas'
		: streak >= 3 ? '⚡ Streak'
		: `${streak}x`;
</script>

{#if show}
	<div
		in:scale={{ duration: 350, easing: elasticOut, start: 0.5 }}
		out:fly={{ y: -8, duration: 200 }}
		class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold
		       bg-yellow-500/10 border border-yellow-500/30 text-yellow-400"
	>
		{label}
		{#if streak >= 3}
			<span class="text-yellow-300 font-black">{streak}</span>
		{/if}
	</div>
{/if}
