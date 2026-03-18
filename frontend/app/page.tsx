'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { PokemonCollection } from '@/components/PokemonCollection'
import { ActivityFeed } from '@/components/ActivityFeed'
import { AchievementSection } from '@/components/AchievementSection'
import { DailyMissions } from '@/components/DailyMissions'
import { LimitedEvents } from '@/components/LimitedEvents'
import { api } from '@/lib/api'
import { isLoggedIn, getUser } from '@/lib/auth'
import type { Pokemon, GitHubActivity, UserStats } from '@/lib/types'

export default function DashboardPage() {
  const router = useRouter()
  const [pokemon, setPokemon] = useState<Pokemon[]>([])
  const [activities, setActivities] = useState<GitHubActivity[]>([])
  const [stats, setStats] = useState<UserStats | null>(null)
  const [loading, setLoading] = useState(true)
  const user = getUser()

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    api.checkNeedsStarter().then(({ needs_starter }) => {
      if (needs_starter) { router.push('/starter'); return }
      return Promise.all([
        api.getPokedex(),
        api.getActivities(10),
        api.getStats(),
      ]).then(([pokedexData, activitiesData, statsData]) => {
        setPokemon(pokedexData.pokemon || [])
        setActivities(activitiesData.activities || [])
        setStats(statsData)
      })
    }).catch(console.error).finally(() => setLoading(false))
  }, [router])

  if (loading) return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900">
      <div className="animate-pokeball-shake">
        <div className="w-16 h-16 bg-[#DC0A2D] rounded-full border-4 border-gray-600 relative">
          <div className="absolute top-1/2 left-0 right-0 h-1 bg-gray-600" />
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 bg-white rounded-full border-2 border-gray-600" />
        </div>
      </div>
    </div>
  )

  const xpProgress = stats ? Math.max(0, 100 - (stats.xp_to_next_level / (100 + (stats.level - 1) * 50) * 100)) : 0

  return (
    <>
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex items-center gap-3 mb-6">
          <h1 className="text-2xl font-bold text-gray-900">
            {user?.name || 'トレーナー'}のダッシュボード
          </h1>
          {user?.github_username && (
            <span className="text-sm text-gray-500 bg-gray-100 px-2 py-0.5 rounded">
              @{user.github_username}
            </span>
          )}
        </div>

        {/* Level & XP Bar + Streak */}
        {stats && (
          <div className="bg-gradient-to-r from-gray-800 to-gray-900 rounded-xl p-6 mb-6 text-white">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-3">
                <span className="text-3xl font-bold text-green-400">Lv.{stats.level}</span>
                <span className="text-gray-400">|</span>
                <span className="text-lg">{stats.total_xp} XP</span>
              </div>
              <div className="flex items-center gap-4">
                {stats.streak_multiplier > 1 && (
                  <span className="text-xs bg-orange-500/20 text-orange-300 px-2 py-1 rounded-full">
                    XP x{stats.streak_multiplier}
                  </span>
                )}
                <span className="text-sm text-gray-400">次のレベルまで {stats.xp_to_next_level} XP</span>
              </div>
            </div>
            <div className="w-full bg-gray-700 rounded-full h-3">
              <div
                className="bg-gradient-to-r from-green-400 to-green-500 h-3 rounded-full transition-all duration-500"
                style={{ width: `${xpProgress}%` }}
              />
            </div>
          </div>
        )}

        {/* Stats Grid */}
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-8">
          <StatCard label="総コミット" value={stats?.total_commits ?? 0} icon="💻" color="bg-gradient-to-br from-gray-700 to-gray-900" />
          <StatCard label="PRマージ" value={stats?.total_merges ?? 0} icon="🔀" color="bg-gradient-to-br from-purple-500 to-purple-700" />
          <StatCard label="ポケモン" value={stats?.pokemon_count ?? 0} icon="⚡" color="bg-gradient-to-br from-yellow-400 to-yellow-600" />
          <StreakCard currentStreak={stats?.current_streak ?? 0} maxStreak={stats?.max_streak ?? 0} />
          <StatCard label="レベル" value={stats?.level ?? 1} icon="🌟" color="bg-gradient-to-br from-green-400 to-green-600" />
        </div>

        <LimitedEvents />
        <DailyMissions />

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
          <PokemonCollection pokemon={pokemon} />
          <ActivityFeed activities={activities} />
        </div>

        <AchievementSection achievements={stats?.achievements ?? []} />

        {/* How it works */}
        <div className="mt-6 bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl p-6 border border-gray-200">
          <h2 className="text-lg font-bold text-gray-900 mb-4">仕組み</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
            <div>
              <div className="text-2xl mb-1">💻</div>
              <p className="font-bold text-sm">コミット</p>
              <p className="text-xs text-gray-500">+10 XP</p>
            </div>
            <div>
              <div className="text-2xl mb-1">📝</div>
              <p className="font-bold text-sm">PR作成</p>
              <p className="text-xs text-gray-500">+20 XP</p>
            </div>
            <div>
              <div className="text-2xl mb-1">👀</div>
              <p className="font-bold text-sm">レビュー</p>
              <p className="text-xs text-gray-500">+25 XP</p>
            </div>
            <div>
              <div className="text-2xl mb-1">🎉</div>
              <p className="font-bold text-sm">PRマージ</p>
              <p className="text-xs text-gray-500">+50 XP + ポケモン！</p>
            </div>
          </div>
          <div className="mt-4 pt-4 border-t border-gray-200">
            <p className="text-sm font-bold text-gray-700 mb-2">ストリークボーナス</p>
            <div className="grid grid-cols-3 gap-3 text-center text-sm">
              <div className="bg-white rounded-lg p-2 border">
                <p className="font-bold text-orange-500">7日連続</p>
                <p className="text-xs text-gray-500">XP x1.5</p>
              </div>
              <div className="bg-white rounded-lg p-2 border">
                <p className="font-bold text-orange-600">30日連続</p>
                <p className="text-xs text-gray-500">XP x2.0</p>
              </div>
              <div className="bg-white rounded-lg p-2 border">
                <p className="font-bold text-orange-700">100日連続</p>
                <p className="text-xs text-gray-500">XP x3.0</p>
              </div>
            </div>
          </div>
        </div>
      </main>
    </>
  )
}

function StatCard({ label, value, icon, color }: { label: string; value: number; icon: string; color: string }) {
  return (
    <div className={`${color} rounded-xl shadow-sm p-5 text-white`}>
      <span className="text-2xl">{icon}</span>
      <p className="text-3xl font-bold mt-1">{value}</p>
      <p className="text-sm text-white/80">{label}</p>
    </div>
  )
}

function StreakCard({ currentStreak, maxStreak }: { currentStreak: number; maxStreak: number }) {
  const streakColor = currentStreak >= 100
    ? 'from-red-500 to-orange-600'
    : currentStreak >= 30
    ? 'from-orange-500 to-yellow-500'
    : currentStreak >= 7
    ? 'from-orange-400 to-yellow-400'
    : 'from-orange-300 to-yellow-300'

  return (
    <div className={`bg-gradient-to-br ${streakColor} rounded-xl shadow-sm p-5 text-white`}>
      <span className="text-2xl">🔥</span>
      <p className="text-3xl font-bold mt-1">{currentStreak}</p>
      <p className="text-sm text-white/80">連続日数</p>
      {maxStreak > currentStreak && (
        <p className="text-xs text-white/60 mt-0.5">最高: {maxStreak}日</p>
      )}
    </div>
  )
}
