// src/lib/stores/quiz.ts
import { writable, derived } from 'svelte/store';
import type { Session, Question, AnswerResult, Option } from '$lib/api/client';

export interface QuizState {
	session: Session | null;
	currentIndex: number;
	lastAnswer: AnswerResult | null;
	answered: boolean;
	finished: boolean;
	streak: number;
	shuffledOptions: Option[][];   // pre-shuffled options per question index
	timedOut: boolean;
}

function shuffle<T>(arr: T[]): T[] {
	const a = [...arr];
	for (let i = a.length - 1; i > 0; i--) {
		const j = Math.floor(Math.random() * (i + 1));
		[a[i], a[j]] = [a[j], a[i]];
	}
	return a;
}

function createQuizStore() {
	const initial: QuizState = {
		session: null,
		currentIndex: 0,
		lastAnswer: null,
		answered: false,
		finished: false,
		streak: 0,
		shuffledOptions: [],
		timedOut: false,
	};

	const { subscribe, set, update } = writable<QuizState>(initial);

	return {
		subscribe,
		startSession(session: Session) {
			// Pre-shuffle all option arrays once so they stay stable during re-renders
			const shuffledOptions = session.questions.map((q) => shuffle(q.options));
			set({ ...initial, session, shuffledOptions });
		},
		recordAnswer(result: AnswerResult, timedOut = false) {
			update((s) => ({
				...s,
				lastAnswer: result,
				answered: true,
				timedOut,
				streak: result.is_correct ? result.current_streak : 0,
			}));
		},
		nextQuestion() {
			update((s) => {
				const nextIndex = s.currentIndex + 1;
				const finished = nextIndex >= (s.session?.total_questions ?? 0);
				return { ...s, currentIndex: nextIndex, answered: false, lastAnswer: null, finished, timedOut: false };
			});
		},
		reset() {
			set(initial);
		},
	};
}

export const quiz = createQuizStore();

export const currentQuestion = derived(quiz, ($quiz): Question | null => {
	if (!$quiz.session?.questions) return null;
	return $quiz.session.questions[$quiz.currentIndex] ?? null;
});

export const currentOptions = derived(quiz, ($quiz): Option[] => {
	return $quiz.shuffledOptions[$quiz.currentIndex] ?? [];
});

export const progress = derived(quiz, ($quiz) => {
	if (!$quiz.session) return { current: 0, total: 0, percent: 0 };
	const total = $quiz.session.total_questions;
	const current = $quiz.currentIndex + 1;
	return { current, total, percent: Math.round((current / total) * 100) };
});
