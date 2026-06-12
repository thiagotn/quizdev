<!-- src/lib/components/CodeBlock.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import type { BundledLanguage } from 'shiki';

	export let code: string;
	export let lang: BundledLanguage = 'go';

	let html = '';
	let container: HTMLElement;

	// Simple fallback: display raw code while Shiki loads
	$: displayCode = code;

	onMount(async () => {
		try {
			const { codeToHtml } = await import('shiki');
			html = await codeToHtml(code, {
				lang,
				theme: 'github-dark-default'
			});
		} catch {
			// Fallback to plain display if shiki fails
			html = '';
		}
	});
</script>

{#if html}
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		bind:this={container}
		class="shiki-wrapper text-xs leading-relaxed overflow-x-auto rounded-lg"
	>
		{@html html}
	</div>
{:else}
	<pre class="text-xs text-zinc-300 leading-relaxed whitespace-pre-wrap break-words font-mono bg-transparent">{displayCode}</pre>
{/if}

<style>
	.shiki-wrapper :global(pre) {
		background: transparent !important;
		padding: 0;
		margin: 0;
		font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
		font-size: 0.75rem;
		line-height: 1.6;
		white-space: pre-wrap;
		word-break: break-word;
	}

	.shiki-wrapper :global(code) {
		background: transparent !important;
		font-family: inherit;
	}
</style>
