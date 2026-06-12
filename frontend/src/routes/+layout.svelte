<!-- src/routes/+layout.svelte -->
<script lang="ts">
  import '../app.css';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { isLoggedIn, auth } from '$lib/stores/auth';
  import { onMount } from 'svelte';

  const publicRoutes = ['/', '/login', '/register'];
  const quizRoutes   = ['/quiz'];

  $: isPublic = publicRoutes.includes($page.url.pathname);
  $: inQuiz   = quizRoutes.some(r => $page.url.pathname.startsWith(r));
  $: showNav  = $isLoggedIn && !inQuiz;

  onMount(() => {
    if (!$isLoggedIn && !isPublic) goto('/login');
  });
</script>

<div class="min-h-screen text-[#F1FAEE]" style="background: #111111; font-family: 'Inter', system-ui, sans-serif;">

  {#if showNav}
    <header class="border-b border-white/5 px-4 py-3 flex items-center justify-between sticky top-0 z-10"
      style="background: rgba(17,17,17,0.95); backdrop-filter: blur(8px);">
      <a href="/" class="flex items-center gap-1.5 font-black tracking-tight text-xs sm:text-base font-mono hover:opacity-80 transition-opacity">
        <img src="/logos/go-gopher.svg" alt="GoRound" class="h-6 w-6 sm:h-8 sm:w-8 shrink-0" />
        <span class="leading-tight">Go<span style="color:#E63946">Round</span></span>
      </a>
      <nav class="flex items-center gap-5 text-sm">
        <a href="/"            class="text-white/40 hover:text-white/80 transition-colors">Rounds</a>
        <a href="/leaderboard" class="text-white/40 hover:text-white/80 transition-colors">Ranking</a>
        <a href="/profile"     class="text-white/40 hover:text-white/80 transition-colors">Perfil</a>
        <button
          class="text-white/20 hover:text-red-400 transition-colors text-xs uppercase tracking-wider"
          onclick={() => { auth.logout(); goto('/login'); }}
        >Sair</button>
      </nav>
    </header>
  {/if}

  {#if inQuiz && $isLoggedIn}
    <header class="border-b border-white/5 px-4 py-3 flex items-center justify-between">
      <span class="font-black tracking-tight text-sm font-mono" style="color:#E63946">● GoRound</span>
      <a href="/" class="text-xs text-white/20 hover:text-white/40 transition-colors uppercase tracking-wider font-mono">
        ✕ Abandonar round
      </a>
    </header>
  {/if}

  <main class="w-full mx-auto px-3 py-6 sm:max-w-lg sm:mx-auto sm:px-4">
    <slot />
  </main>
</div>
