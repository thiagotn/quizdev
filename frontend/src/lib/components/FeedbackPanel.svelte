<!-- src/lib/components/FeedbackPanel.svelte -->
<script lang="ts">
  import { fly } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import type { AnswerResult } from '$lib/api/client';
  import CodeBlock from './CodeBlock.svelte';
  import { feedbackCorrect, feedbackWrong } from '$lib/utils/labels';

  export let result: AnswerResult;
  export let timedOut = false;
  export let onNext: () => void;
  export let isLast = false;

  $: headline = timedOut ? 'Tempo esgotado.' : result.is_correct ? feedbackCorrect() : feedbackWrong();
</script>

<div
  in:fly={{ y: 20, duration: 300, easing: cubicOut }}
  class="rounded-2xl border overflow-hidden mb-5"
  style="
    border-color: {result.is_correct ? 'rgba(74,222,128,0.25)' : 'rgba(230,57,70,0.25)'};
    background:   {result.is_correct ? 'rgba(74,222,128,0.05)' : 'rgba(230,57,70,0.05)'};
  "
>
  <!-- Header -->
  <div class="px-4 py-3 flex items-center gap-3 border-b"
    style="border-color: {result.is_correct ? 'rgba(74,222,128,0.1)' : 'rgba(230,57,70,0.1)'}">
    <span class="text-xl">{timedOut ? '⏰' : result.is_correct ? '✅' : '❌'}</span>
    <div class="flex-1">
      <p class="font-bold text-sm font-mono"
        style="color: {result.is_correct ? '#4ADE80' : '#E63946'}">{headline}</p>
      {#if result.is_correct && result.points_earned > 0}
        <p class="text-xs font-mono" style="color: rgba(74,222,128,0.6)">+{result.points_earned} pts</p>
      {/if}
    </div>
    {#if result.current_streak >= 3}
      <div class="text-xs font-bold font-mono px-2 py-0.5 rounded"
        style="color:#F59E0B; background:rgba(245,158,11,0.1); border:1px solid rgba(245,158,11,0.2)">
        🔥 ×{result.current_streak}
      </div>
    {/if}
  </div>

  <!-- Body -->
  <div class="px-4 py-3">
    {#if !result.is_correct && result.correct_option}
      <p class="text-xs uppercase tracking-wider mb-2 font-mono" style="color:rgba(241,250,238,0.3)">Resposta correta</p>
      <div class="rounded-lg p-2 mb-3 border" style="background:#1A1A1A; border-color:rgba(74,222,128,0.15)">
        <CodeBlock code={result.correct_option.code_snippet} lang="go" />
      </div>
    {/if}
    <p class="text-sm leading-relaxed" style="color:rgba(241,250,238,0.7)">{result.explanation}</p>
  </div>
</div>

<button
  onclick={onNext}
  class="w-full font-bold font-mono rounded-xl py-3 text-sm transition-all active:scale-95 uppercase tracking-wider"
  style="background:{result.is_correct ? '#E63946' : '#2a2a2a'}; color:{result.is_correct ? '#111' : '#F1FAEE'}"
>
  {isLast ? 'Ver resultado →' : 'Próxima →'}
</button>
