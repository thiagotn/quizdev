<!-- src/lib/components/MiniChart.svelte -->
<script lang="ts">
	export let data: number[] = [];

	const W = 320;
	const H = 64;
	const PAD = 4;

	$: min = Math.min(...data);
	$: max = Math.max(...data);
	$: range = max - min || 1;

	function x(i: number) {
		return PAD + (i / (data.length - 1)) * (W - PAD * 2);
	}

	function y(v: number) {
		return H - PAD - ((v - min) / range) * (H - PAD * 2);
	}

	$: points = data.map((v, i) => `${x(i)},${y(v)}`).join(' ');

	// Area fill path
	$: area = data.length > 1
		? `M${x(0)},${H} ` +
		  data.map((v, i) => `L${x(i)},${y(v)}`).join(' ') +
		  ` L${x(data.length - 1)},${H} Z`
		: '';

	// Line path
	$: line = data.length > 1
		? `M` + data.map((v, i) => `${x(i)},${y(v)}`).join(' L')
		: '';

	$: lastX = x(data.length - 1);
	$: lastY = y(data[data.length - 1]);
</script>

{#if data.length >= 2}
	<svg
		viewBox="0 0 {W} {H}"
		width="100%"
		height={H}
		class="overflow-visible"
		aria-hidden="true"
	>
		<defs>
			<linearGradient id="chartGrad" x1="0" y1="0" x2="0" y2="1">
				<stop offset="0%" stop-color="#10b981" stop-opacity="0.25" />
				<stop offset="100%" stop-color="#10b981" stop-opacity="0" />
			</linearGradient>
		</defs>

		<!-- Area fill -->
		<path d={area} fill="url(#chartGrad)" />

		<!-- Line -->
		<polyline
			{points}
			fill="none"
			stroke="#10b981"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		/>

		<!-- Dots -->
		{#each data as v, i}
			<circle cx={x(i)} cy={y(v)} r="3" fill="#10b981" opacity={i === data.length - 1 ? 1 : 0.4} />
		{/each}

		<!-- Last value label -->
		<text
			x={lastX}
			y={lastY - 8}
			text-anchor="middle"
			font-size="10"
			font-weight="700"
			fill="#10b981"
			font-family="monospace"
		>
			{data[data.length - 1]}
		</text>
	</svg>
{/if}
