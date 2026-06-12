<!-- src/routes/login/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { authApi } from '$lib/api/client';

	let email = '';
	let password = '';
	let error = '';
	let loading = false;

	async function submit() {
		error = '';
		loading = true;
		try {
			const res = await authApi.login(email, password);
			auth.login(res.token, res.user);
			goto('/');
		} catch (e: any) {
			error = e.message ?? 'Erro ao fazer login';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head><title>Login — GoRound</title></svelte:head>

<div class="py-12 max-w-sm mx-auto">
	<div class="mb-8 text-center">
		<h1 class="text-2xl font-black font-mono" style="color:#F1FAEE">Entrar</h1>
		<p class="text-zinc-500 text-sm mt-1">Volte para a arena</p>
	</div>

	<div class="bg-zinc-900 border border-zinc-800 rounded-xl p-6 flex flex-col gap-4">
		{#if error}
			<div class="text-sm text-red-400 bg-red-950/50 border border-red-900 rounded-lg px-3 py-2">
				{error}
			</div>
		{/if}

		<div>
			<label for="email" class="block text-xs text-zinc-500 mb-1.5 uppercase tracking-wider">Email</label>
			<input
				id="email"
				type="email"
				bind:value={email}
				class="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-zinc-100
				       focus:outline-none focus:border-emerald-500 transition-colors"
				placeholder="you@example.com"
			/>
		</div>

		<div>
			<label for="password" class="block text-xs text-zinc-500 mb-1.5 uppercase tracking-wider">Senha</label>
			<input
				id="password"
				type="password"
				bind:value={password}
				onkeydown={(e) => e.key === 'Enter' && submit()}
				class="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-zinc-100
				       focus:outline-none focus:border-emerald-500 transition-colors"
				placeholder="••••••••"
			/>
		</div>

		<button
			onclick={submit}
			disabled={loading}
			class="w-full bg-emerald-500 hover:bg-emerald-400 disabled:opacity-50 disabled:cursor-not-allowed
			       text-zinc-950 font-bold rounded-lg py-2.5 text-sm transition-colors"
		>
			{loading ? 'Entrando...' : 'Entrar'}
		</button>
	</div>

	<p class="text-center text-zinc-500 text-sm mt-4">
		Não tem conta?
		<a href="/register" class="text-emerald-400 underline underline-offset-2">Criar conta</a>
	</p>
</div>
