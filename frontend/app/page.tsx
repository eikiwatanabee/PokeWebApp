'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn, getUser } from '@/lib/auth'
import type { Pokemon, GitHubActivity, UserStats } from '@/lib/types'

const eventIcons: Record<string, string> = {
  commit: '💻',
  pr_merge: '🎉',
  pr_open: '📝',
  issue_close: '✅',
  review: '👀',
}

const eventLabels: Record<string, string> = {
  commit: 'コミット',
  pr_merge: 'PRマージ',
  pr_open: 'PR作成',
  issue_close: 'Issue完了',
  review: 'レビュー',
}

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

        {/* Level & XP Bar */}
        {stats && (
          <div className="bg-gradient-to-r from-gray-800 to-gray-900 rounded-xl p-6 mb-6 text-white">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-3">
                <span className="text-3xl font-bold text-green-400">Lv.{stats.level}</span>
                <span className="text-gray-400">|</span>
                <span className="text-lg">{stats.total_xp} XP</span>
              </div>
              <span className="text-sm text-gray-400">次のレベルまで {stats.xp_to_next_level} XP</span>
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
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
          <StatCard label="総コミット" value={stats?.total_commits ?? 0} icon="💻" color="bg-gradient-to-br from-gray-700 to-gray-900" />
          <StatCard label="PRマージ" value={stats?.total_merges ?? 0} icon="🔀" color="bg-gradient-to-br from-purple-500 to-purple-700" />
          <StatCard label="ポケモン" value={stats?.pokemon_count ?? 0} icon="⚡" color="bg-gradient-to-br from-yellow-400 to-yellow-600" />
          <StatCard label="レベル" value={stats?.level ?? 1} icon="🌟" color="bg-gradient-to-br from-green-400 to-green-600" />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Pokemon Collection */}
          <div className="bg-white rounded-xl shadow-sm p-6 border-2 border-[#DC0A2D]/10">
            <div className="flex items-center gap-2 mb-4">
              <span className="text-xl">⚡</span>
              <h2 className="text-lg font-bold">ゲットしたポケモン</h2>
            </div>
            <div className="flex items-baseline gap-2">
              <p className="text-4xl font-bold text-[#DC0A2D]">{pokemon.length}</p>
              <p className="text-gray-500 text-sm">匹</p>
            </div>
            {pokemon.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-4">
                {pokemon.slice(0, 12).map(p => (
                  <img key={p.id} src={p.sprite_url} alt={p.pokemon_name} className="w-12 h-12 pixelated" />
                ))}
                {pokemon.length > 12 && (
                  <button
                    onClick={() => router.push('/pokedex')}
                    className="w-12 h-12 flex items-center justify-center text-sm text-[#DC0A2D] font-bold hover:bg-red-50 rounded-lg transition-colors"
                  >
                    +{pokemon.length - 12}
                  </button>
                )}
              </div>
            )}
            {pokemon.length === 0 && (
              <p className="text-gray-400 text-sm mt-3">PRをマージするとポケモンをゲットできるよ！</p>
            )}
          </div>

          {/* Recent Activity */}
          <div className="bg-white rounded-xl shadow-sm p-6 border-2 border-green-100">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <span className="text-xl">📊</span>
                <h2 className="text-lg font-bold">最近のアクティビティ</h2>
              </div>
              <button
                onClick={() => router.push('/activities')}
                className="text-xs text-green-600 hover:underline"
              >
                すべて見る →
              </button>
            </div>
            {activities.length === 0 ? (
              <div className="text-center py-4">
                <p className="text-gray-400">まだアクティビティがありません</p>
                <p className="text-sm text-gray-400 mt-1">GitHubにコミットして経験値をゲットしよう！</p>
              </div>
            ) : (
              <ul className="space-y-2">
                {activities.slice(0, 5).map(a => (
                  <li key={a.id} className="flex items-center gap-3 p-2.5 rounded-lg hover:bg-green-50 transition-colors">
                    <span className="text-lg">{eventIcons[a.event_type] || '📌'}</span>
                    <div className="min-w-0 flex-1">
                      <span className="font-medium block truncate text-sm">{a.title}</span>
                      <span className="text-xs text-gray-400">{a.repo_name}</span>
                    </div>
                    <span className="text-xs font-bold text-green-600 whitespace-nowrap">+{a.xp} XP</span>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

        {/* How it works */}
        <div className="mt-8 bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl p-6 border border-gray-200">
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
              <p className="text-xs text-gray-500">+50 XP + ポケモンゲット！</p>
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
