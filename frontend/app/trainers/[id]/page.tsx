'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import type { TrainerCard } from '@/lib/types'

const rarityBorder: Record<string, string> = {
  legendary: 'border-yellow-400 shadow-yellow-200/50',
  epic: 'border-purple-400 shadow-purple-200/50',
  rare: 'border-blue-400',
  uncommon: 'border-green-400',
  common: 'border-gray-300',
}

export default function TrainerCardPage() {
  const params = useParams()
  const [card, setCard] = useState<TrainerCard | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const id = params.id as string
    api.getTrainerCard(id === 'me' ? undefined : id)
      .then(setCard)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [params.id])

  if (loading) return (
    <>
      <Navbar />
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-400">読み込み中...</p>
      </div>
    </>
  )

  if (!card) return (
    <>
      <Navbar />
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-500">トレーナーが見つかりません</p>
      </div>
    </>
  )

  return (
    <>
      <Navbar />
      <main className="max-w-2xl mx-auto px-4 py-8">
        {/* Trainer Card */}
        <div className="bg-gradient-to-br from-gray-800 via-gray-900 to-black rounded-2xl p-6 text-white shadow-2xl border border-gray-700">
          {/* Header */}
          <div className="flex items-center gap-4 mb-6">
            {card.avatar_url ? (
              <img src={card.avatar_url} alt={card.name} className="w-16 h-16 rounded-full border-2 border-green-400" />
            ) : (
              <div className="w-16 h-16 rounded-full bg-gray-700 border-2 border-green-400 flex items-center justify-center text-2xl">
                🎮
              </div>
            )}
            <div>
              <h1 className="text-xl font-bold">{card.name}</h1>
              {card.github_username && (
                <p className="text-sm text-gray-400">@{card.github_username}</p>
              )}
            </div>
            <div className="ml-auto text-right">
              <p className="text-3xl font-bold text-green-400">Lv.{card.level}</p>
              <p className="text-xs text-gray-400">{card.total_xp.toLocaleString()} XP</p>
            </div>
          </div>

          {/* Stats Row */}
          <div className="grid grid-cols-3 gap-3 mb-6">
            <div className="bg-white/5 rounded-xl p-3 text-center">
              <p className="text-2xl font-bold text-yellow-400">⚡ {card.pokemon_count}</p>
              <p className="text-xs text-gray-400">ポケモン</p>
            </div>
            <div className="bg-white/5 rounded-xl p-3 text-center">
              <p className="text-2xl font-bold text-orange-400">🔥 {card.current_streak}</p>
              <p className="text-xs text-gray-400">ストリーク</p>
            </div>
            <div className="bg-white/5 rounded-xl p-3 text-center">
              <p className="text-2xl font-bold text-blue-400">🏅 {card.achievement_count}</p>
              <p className="text-xs text-gray-400">アチーブメント</p>
            </div>
          </div>

          {/* Featured Pokemon */}
          {card.featured_pokemon.length > 0 && (
            <div className="mb-6">
              <p className="text-xs font-bold text-gray-400 uppercase mb-3">お気に入りポケモン</p>
              <div className="flex gap-2 flex-wrap">
                {card.featured_pokemon.map((p, i) => (
                  <div
                    key={i}
                    className={`rounded-xl border-2 bg-white/5 p-1 ${rarityBorder[p.rarity] || 'border-gray-600'}`}
                    title={p.pokemon_name}
                  >
                    <img src={p.sprite_url} alt={p.pokemon_name} className="w-14 h-14 pixelated" />
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Achievements */}
          {card.achievements.length > 0 && (
            <div>
              <p className="text-xs font-bold text-gray-400 uppercase mb-2">アチーブメント</p>
              <div className="flex flex-wrap gap-1.5">
                {card.achievements.slice(0, 12).map((a, i) => (
                  <span key={i} className="text-lg" title={a.name}>{a.icon}</span>
                ))}
                {card.achievements.length > 12 && (
                  <span className="text-xs text-gray-500 self-center">+{card.achievements.length - 12}</span>
                )}
              </div>
            </div>
          )}
        </div>
      </main>
    </>
  )
}
