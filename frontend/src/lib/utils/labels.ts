// src/lib/utils/labels.ts — GoRound naming system
import type { Level } from '$lib/api/client';

export const ROUND_LABEL: Record<Level, string> = {
  beginner:     'Round I',
  intermediate: 'Round II',
  advanced:     'Round III',
};

export const ROUND_ICON: Record<Level, string> = {
  beginner:     '🌱',
  intermediate: '⚡',
  advanced:     '🔥',
};

export function feedbackCorrect(): string {
  const phrases = ['Sobreviveu.', 'Correto — avance.', 'Eliminação evitada.', 'Passou para o próximo.'];
  return phrases[Math.floor(Math.random() * phrases.length)];
}

export function feedbackWrong(): string {
  const phrases = ['Eliminado nesta questão.', 'Errou — sem pontos.', 'A arena não perdoa.', 'Incorreto.'];
  return phrases[Math.floor(Math.random() * phrases.length)];
}

export function gradeResult(accuracy: number): { emoji: string; label: string } {
  if (accuracy === 100) return { emoji: '💎', label: 'Imaculado. Round dominado.' };
  if (accuracy >= 90)   return { emoji: '🏆', label: 'Sobreviveu com distinção.' };
  if (accuracy >= 70)   return { emoji: '⭐', label: 'Sobreviveu ao round.' };
  if (accuracy >= 50)   return { emoji: '🩹', label: 'Saiu ferido, mas vivo.' };
  return                       { emoji: '💀', label: 'Eliminado. Tente novamente.' };
}
