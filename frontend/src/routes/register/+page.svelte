<!-- src/routes/register/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { authApi } from '$lib/api/client';

	let email = '';
	let username = '';
	let password = '';
	let error = '';
	let loading = false;

	async function submit() {
		error = '';
		if (password.length < 8) {
			error = 'A senha deve ter no mínimo 8 caracteres';
			return;
		}
		loading = true;
		try {
			const res = await authApi.register(email, username, password);
			auth.login(res.token, res.user);
			goto('/');
		} catch (e: any) {
			error = e.message ?? 'Erro ao criar conta';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head><title>Criar conta — GoRound</title></svelte:head>

<div class="py-12 max-w-sm mx-auto">
	<div class="mb-8 text-center">
		<h1 class="text-2xl font-black text-zinc-100">Criar conta</h1>
		<p class="text-zinc-500 text-sm mt-1">Entre na arena</p>
	</div>

	<div class="bg-zinc-900 border border-zinc-800 rounded-xl p-6 flex flex-col gap-4">
		{#if error}
			<div class="text-sm text-red-400 bg-red-950/50 border border-red-900 rounded-lg px-3 py-2">
				{error}
			</div>
		{/if}

		<div>
			<label for="username" class="block text-xs text-zinc-500 mb-1.5 uppercase tracking-wider">Usuário</label>
			<input
				id="username"
				type="text"
				bind:value={username}
				class="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-zinc-100
				       focus:outline-none focus:border-emerald-500 transition-colors"
				placeholder="gopher42"
			/>
		</div>

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
				placeholder="mín. 8 caracteres"
			/>
		</div>

		<button
			onclick={submit}
			disabled={loading}
			class="w-full bg-emerald-500 hover:bg-emerald-400 disabled:opacity-50 disabled:cursor-not-allowed
			       text-zinc-950 font-bold rounded-lg py-2.5 text-sm transition-colors"
		>
			{loading ? 'Criando...' : 'Criar conta'}
		</button>
	</div>

	<p class="text-center text-zinc-500 text-sm mt-4">
		Já tem conta?
		<a href="/login" class="text-emerald-400 underline underline-offset-2">Entrar</a>
	</p>
</div>
