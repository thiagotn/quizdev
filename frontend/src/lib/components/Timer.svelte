<!-- src/lib/components/Timer.svelte -->
<script lang="ts">
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';

	export let seconds = 20;
	export let paused = false;

	const dispatch = createEventDispatcher<{ timeout: void }>();

	let remaining = seconds;
	export let timeRemaining = seconds; // bindable for parent to read
	let interval: ReturnType<typeof setInterval>;

	const RADIUS = 22;
	const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

	$: dashoffset = CIRCUMFERENCE * (1 - remaining / seconds);
	$: pct = remaining / seconds;
	$: color = pct > 0.5 ? '#10b981' : pct > 0.25 ? '#eab308' : '#ef4444';
	$: urgent = pct <= 0.25;

	onMount(() => {
		interval = setInterval(() => {
			if (paused) return;
			remaining -= 1;
			timeRemaining = remaining;
			if (remaining <= 0) {
				remaining = 0;
				timeRemaining = 0;
				clearInterval(interval);
				dispatch('timeout');
			}
		}, 1000);
	});

	onDestroy(() => clearInterval(interval));

	export function stop() {
		clearInterval(interval);
	}
</script>

<div class="relative flex items-center justify-center w-14 h-14" class:animate-pulse={urgent}>
	<svg class="absolute inset-0 -rotate-90" width="56" height="56" viewBox="0 0 56 56">
		<!-- Track -->
		<circle cx="28" cy="28" r={RADIUS} fill="none" stroke="#27272a" stroke-width="3" />
		<!-- Progress -->
		<circle
			cx="28"
			cy="28"
			r={RADIUS}
			fill="none"
			stroke={color}
			stroke-width="3"
			stroke-linecap="round"
			stroke-dasharray={CIRCUMFERENCE}
			stroke-dashoffset={dashoffset}
			style="transition: stroke-dashoffset 0.9s linear, stroke 0.3s ease"
		/>
	</svg>
	<span
		class="text-sm font-black tabular-nums"
		style="color: {color}; transition: color 0.3s ease"
	>
		{remaining}
	</span>
</div>
