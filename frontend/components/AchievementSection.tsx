'use client'

import type { Achievement } from '@/lib/types'

const categoryLabels: Record<string, string> = {
  commit: '💻 コミット',
  streak: '🔥 ストリーク',
  level: '⬆️ レベル',
  xp: '✨ 経験値',
  pokemon: '⚡ ポケモン',
  rarity: '💎 レアリティ',
  pr: '🔀 プルリクエスト',
  review: '👀 レビュー',
  issue: '🐛 Issue',
  special: '🎯 スペシャル',
}

function groupByCategory(achievements: Achievement[]): Record<string, Achievement[]> {
  const groups: Record<string, Achievement[]> = {}
  for (const a of achievements) {
    const cat = a.category || 'other'
    if (!groups[cat]) groups[cat] = []
    groups[cat].push(a)
  }
  return groups
}

export function AchievementSection({ achievements }: { achievements: Achievement[] }) {
  if (!achievements || achievements.length === 0) return null

  return (
    <div className="mt-6 bg-white rounded-xl shadow-sm p-6 border-2 border-yellow-100">
      <div className="flex items-center gap-2 mb-4">
        <span className="text-xl">🏅</span>
        <h2 className="text-lg font-bold">アチーブメント</h2>
        <span className="text-sm text-gray-400 ml-auto">{achievements.length}個 達成</span>
      </div>
      {Object.entries(groupByCategory(achievements)).map(([category, items]) => (
        <div key={category} className="mb-4 last:mb-0">
          <p className="text-xs font-bold text-gray-500 uppercase mb-2">{categoryLabels[category] || category}</p>
          <div className="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2">
            {items.map(a => (
              <div
                key={a.type}
                className="flex flex-col items-center p-2.5 bg-gradient-to-b from-yellow-50 to-white rounded-xl border border-yellow-200 hover:shadow-md transition-shadow"
                title={a.description}
              >
                <span className="text-2xl mb-0.5">{a.icon}</span>
                <span className="text-[11px] font-bold text-center leading-tight">{a.name}</span>
                <span className="text-[9px] text-gray-400 text-center mt-0.5">{a.description}</span>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}
