// src/lib/api/client.ts
const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

async function request<T>(
	path: string,
	options: RequestInit = {},
	token?: string
): Promise<T> {
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(options.headers as Record<string, string>)
	};
	if (token) headers['Authorization'] = `Bearer ${token}`;

	const res = await fetch(`${BASE_URL}${path}`, { ...options, headers });
	if (!res.ok) {
		const error = await res.json().catch(() => ({ error: 'Unknown error' }));
		throw new Error(error.error ?? 'Request failed');
	}
	return res.json();
}

// ─── Auth ─────────────────────────────────────────────────────────────────────

export interface User {
	id: string; email: string; username: string; provider: string; created_at: string;
}
export interface AuthResponse { token: string; user: User; }

export const authApi = {
	register: (email: string, username: string, password: string) =>
		request<AuthResponse>('/auth/register', { method: 'POST', body: JSON.stringify({ email, username, password }) }),
	login: (email: string, password: string) =>
		request<AuthResponse>('/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) }),
	me: (token: string) => request<User>('/auth/me', {}, token)
};

// ─── Quiz ─────────────────────────────────────────────────────────────────────

export type Level = 'beginner' | 'intermediate' | 'advanced';

export interface Option   { id: string; code_snippet: string; display_order: number; }
export interface Question { id: string; language: string; level: Level; title: string; explanation?: string; options: Option[]; }
export interface Session  {
	id: string; user_id: string; language: string; level: Level;
	total_questions: number; correct_answers: number; score: number;
	started_at: string; finished_at?: string; questions: Question[];
}
export interface AnswerResult {
	is_correct: boolean; explanation: string; correct_option: Option;
	points_earned: number; current_streak: number; total_score: number;
}

export const quizApi = {
	createSession: (language: string, level: Level, token: string) =>
		request<Session>('/sessions', { method: 'POST', body: JSON.stringify({ language, level }) }, token),
	submitAnswer: (sessionId: string, questionId: string, optionId: string, token: string, timeRemaining = 0) =>
		request<AnswerResult>(`/sessions/${sessionId}/answer`, {
			method: 'POST',
			body: JSON.stringify({ question_id: questionId, option_id: optionId, time_remaining: timeRemaining })
		}, token),
	getResult: (sessionId: string, token: string) =>
		request<Session>(`/sessions/${sessionId}/result`, {}, token)
};

// ─── Scores ───────────────────────────────────────────────────────────────────

export interface Score {
	id: string; language: string; level: Level;
	total_points: number; best_streak: number; best_accuracy: number;
	sessions_played: number; updated_at: string;
}

export interface SessionSummary {
	id: string; level: Level; score: number;
	correct_answers: number; total_questions: number;
	accuracy: number; finished_at: string;
}

export interface ProfileResponse {
	user: { id: string; username: string; email: string };
	scores: Score[];
	history: SessionSummary[];
	global_rank: number;
	total_points: number;
}

export interface LeaderboardEntry {
	rank: number; username: string; total_points: number;
	sessions_played: number; best_streak: number; is_current_user: boolean;
}

export const scoreApi = {
	myProfile: (token: string) =>
		request<ProfileResponse>('/scores/me', {}, token),
	leaderboard: (token: string, level?: Level) =>
		request<LeaderboardEntry[]>(
			`/scores/leaderboard${level ? `?level=${level}` : ''}`,
			{}, token
		)
};
