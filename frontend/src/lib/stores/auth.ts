// src/lib/stores/auth.ts
import { writable, derived } from 'svelte/store';
import type { User } from '$lib/api/client';

interface AuthState {
	token: string | null;
	user: User | null;
}

function createAuthStore() {
	// Hydrate from localStorage on init (browser only)
	const initial: AuthState =
		typeof localStorage !== 'undefined'
			? {
					token: localStorage.getItem('token'),
					user: JSON.parse(localStorage.getItem('user') ?? 'null')
				}
			: { token: null, user: null };

	const { subscribe, set, update } = writable<AuthState>(initial);

	return {
		subscribe,
		login(token: string, user: User) {
			localStorage.setItem('token', token);
			localStorage.setItem('user', JSON.stringify(user));
			set({ token, user });
		},
		logout() {
			localStorage.removeItem('token');
			localStorage.removeItem('user');
			set({ token: null, user: null });
		}
	};
}

export const auth = createAuthStore();
export const isLoggedIn = derived(auth, ($auth) => !!$auth.token);
export const currentUser = derived(auth, ($auth) => $auth.user);
export const token = derived(auth, ($auth) => $auth.token);
