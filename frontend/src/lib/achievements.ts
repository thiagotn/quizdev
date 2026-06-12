// src/lib/achievements.ts
import type { Session } from '$lib/api/client';

export interface Achievement {
	icon: string;
	title: string;
	description: string;
}

export function detectAchievements(session: Session, isFirstSession: boolean): Achievement[] {
	const result: Achievement[] = [];
	const accuracy = session.total_questions > 0
		? Math.round((session.correct_answers / session.total_questions) * 100)
		: 0;

	if (isFirstSession) {
		result.push({
			icon: '🚀',
			title: 'Primeira arena!',
			description: 'Você sobreviveu ao primeiro round. A arena está aberta.'
		});
	}

	if (accuracy === 100) {
		result.push({
			icon: '💎',
			title: 'Round perfeito',
			description: 'Nenhum erro. Isso não acontece aqui com frequência.'
		});
	} else if (accuracy >= 90) {
		result.push({
			icon: '🏆',
			title: 'Quase perfeito',
			description: `${accuracy}% de precisão. Você domina ${session.level === 'beginner' ? 'o básico' : session.level === 'intermediate' ? 'o nível intermediário' : 'o avançado'} de Go!`
		});
	}

	if (session.score >= 150) {
		result.push({
			icon: '⭐',
			title: 'Pontuador de elite',
			description: `${session.score} pontos em uma sessão. Os streaks e bônus de velocidade somaram bastante!`
		});
	}

	return result;
}
