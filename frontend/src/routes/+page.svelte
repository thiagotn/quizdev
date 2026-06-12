<!-- src/routes/+page.svelte — GoRound home -->
<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { fly } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import { isLoggedIn, token } from '$lib/stores/auth';
  import { scoreApi, type ProfileResponse } from '$lib/api/client';

  let profile: ProfileResponse | null = null;

  const rounds = [
    { id: 'beginner',     label: 'Round I',   desc: 'Sintaxe, tipos, loops, funções', icon: '🌱',
      border: 'border-emerald-900/60 hover:border-emerald-500/60' },
    { id: 'intermediate', label: 'Round II',  desc: 'Goroutines, interfaces, errors',  icon: '⚡',
      border: 'border-yellow-900/60 hover:border-yellow-500/60'  },
    { id: 'advanced',     label: 'Round III', desc: 'Generics, context, sync',         icon: '🔥',
      border: 'border-red-900/60 hover:border-red-500/60'        },
  ] as const;

  onMount(async () => {
    if ($isLoggedIn && $token) {
      try { profile = await scoreApi.myProfile($token); } catch {}
    }
  });

  function scoreFor(level: string) {
    return profile?.scores.find(s => s.level === level) ?? null;
  }

  function select(level: string) {
    if (!$isLoggedIn) { goto('/login'); return; }
    goto(`/quiz/${level}`);
  }
</script>

<svelte:head><title>GoRound — Survive the code</title></svelte:head>

<div class="py-6 sm:py-8">
  <!-- Hero -->
  <div in:fly={{ y: 12, duration: 300, easing: cubicOut }} class="text-center mb-12">
    <img src="/logos/go-gopher.svg" alt="GoRound Gopher" class="h-20 w-20 mx-auto mb-4" />

    <h1 class="text-5xl sm:text-6xl font-black tracking-tight mb-3 font-mono">
      Go<span style="color:#E63946">Round</span>
    </h1>
    <p class="text-xs uppercase tracking-[4px] font-mono" style="color: rgba(241,250,238,0.3)">
      · Survive the code ·
    </p>
  </div>

  <!-- Round selector -->
  <p class="text-xs uppercase tracking-widest mb-6 text-center font-mono" style="color: rgba(241,250,238,0.25)">
    Escolha o round
  </p>

  <div in:fly={{ y: 12, duration: 300, delay: 80, easing: cubicOut }} class="flex flex-col gap-4">
    {#each rounds as round}
      {@const score = scoreFor(round.id)}
      <button
        onclick={() => select(round.id)}
        class="w-full text-left p-6 rounded-xl border transition-all duration-200 group {round.border}"
        style="background: #1A1A1A;"
      >
        <div class="flex items-center gap-4">
          <span class="text-3xl">{round.icon}</span>
          <div class="flex-1">
            <div class="font-bold font-mono text-lg" style="color:#F1FAEE">{round.label}</div>
            <div class="text-sm mt-1" style="color: rgba(241,250,238,0.3)">{round.desc}</div>
            {#if score}
              <div class="mt-2 flex items-center gap-2">
                <div class="flex-1 h-1 rounded-full" style="background:#2a2a2a">
                  <div class="h-full rounded-full" style="background:#E63946; width:{Math.min(100, score.total_points / 10)}%"></div>
                </div>
                <span class="text-xs font-mono tabular-nums shrink-0" style="color:#E63946">
                  {score.total_points} pts
                </span>
              </div>
            {/if}
          </div>
          <span class="transition-colors text-lg group-hover:translate-x-1 transition-transform duration-200"
            style="color: rgba(241,250,238,0.2)">→</span>
        </div>
      </button>
    {/each}
  </div>

  {#if !$isLoggedIn}
    <p class="text-center text-sm mt-12 font-mono" style="color: rgba(241,250,238,0.3)">
      <a href="/login"    style="color:#E63946" class="underline underline-offset-2">Entre</a>
      ou
      <a href="/register" style="color:#E63946" class="underline underline-offset-2">crie uma conta</a>
      para salvar seu progresso.
    </p>
  {/if}
</div>
