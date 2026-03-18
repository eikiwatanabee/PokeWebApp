'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { UserRank } from '@/lib/types'

type TeamRanking = {
  team_id: string
  team_name: string
  member_count: number
  pokemon_count: number
  total_xp: number
}

type Tab = 'users' | 'teams'

const medals = ['🥇', '🥈', '🥉']
const rankColors = ['bg-yellow-50 border-yellow-400', 'bg-gray-50 border-gray-300', 'bg-orange-50 border-orange-300']

const userSortOptions = [
  { value: 'xp', label: 'XP', icon: '⭐' },
  { value: 'level', label: 'レベル', icon: '📊' },
  { value: 'pokemon', label: 'ポケモン', icon: '⚡' },
  { value: 'streak', label: 'ストリーク', icon: '🔥' },
]

export default function RankingPage() {
  const router = useRouter()
  const [tab, setTab] = useState<Tab>('users')
  const [teams, setTeams] = useState<TeamRanking[]>([])
  const [userRankings, setUserRankings] = useState<UserRank[]>([])
  const [userSortBy, setUserSortBy] = useState('xp')
  const [loading, setLoading] = useState(true)
  const [newTeamName, setNewTeamName] = useState('')
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
  }, [router])

  useEffect(() => {
    setLoading(true)
    if (tab === 'teams') {
      api.getTeamRanking()
        .then(data => setTeams(data.teams || []))
        .catch(console.error)
        .finally(() => setLoading(false))
    } else {
      api.getUserRanking(userSortBy)
        .then(data => setUserRankings(data.rankings || []))
        .catch(console.error)
        .finally(() => setLoading(false))
    }
  }, [tab, userSortBy])

  const loadTeamRanking = () => {
    api.getTeamRanking()
      .then(data => setTeams(data.teams || []))
      .catch(console.error)
  }

  const handleCreateTeam = async () => {
    if (!newTeamName.trim()) return
    setCreating(true)
    try {
      const result = await api.createTeam(newTeamName.trim())
      await api.joinTeam(result.TeamID)
      setNewTeamName('')
      loadTeamRanking()
    } catch (e) {
      console.error(e)
    } finally {
      setCreating(false)
    }
  }

  const handleJoinTeam = async (teamId: string) => {
    try {
      await api.joinTeam(teamId)
      loadTeamRanking()
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <>
      <Navbar />
      <main className="max-w-4xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">ランキング</h1>
        <p className="text-gray-500 mb-6">トレーナーとチームの競争を楽しもう</p>

        {/* Tab Switcher */}
        <div className="flex gap-2 mb-6">
          <button
            onClick={() => setTab('users')}
            className={`px-5 py-2.5 rounded-lg text-sm font-medium transition-colors ${
              tab === 'users' ? 'bg-gray-900 text-white shadow' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            👤 トレーナー
          </button>
          <button
            onClick={() => setTab('teams')}
            className={`px-5 py-2.5 rounded-lg text-sm font-medium transition-colors ${
              tab === 'teams' ? 'bg-gray-900 text-white shadow' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            🏟️ チーム
          </button>
        </div>

        {/* User Rankings Tab */}
        {tab === 'users' && (
          <>
            {/* Sort Options */}
            <div className="flex gap-2 mb-4">
              {userSortOptions.map(opt => (
                <button
                  key={opt.value}
                  onClick={() => setUserSortBy(opt.value)}
                  className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors flex items-center gap-1 ${
                    userSortBy === opt.value
                      ? 'bg-green-600 text-white'
                      : 'bg-gray-100 text-gray-500 hover:bg-gray-200'
                  }`}
                >
                  <span>{opt.icon}</span>
                  {opt.label}
                </button>
              ))}
            </div>

            {loading ? (
              <p className="text-gray-400 text-center py-12">読み込み中...</p>
            ) : userRankings.length === 0 ? (
              <p className="text-gray-500 text-center py-12">ランキングデータがありません</p>
            ) : (
              <div className="space-y-2">
                {userRankings.map((user) => (
                  <Link
                    key={user.user_id}
                    href={`/trainers/${user.user_id}`}
                    className="flex items-center gap-4 p-4 bg-white rounded-xl border border-gray-200 hover:border-gray-300 hover:shadow-sm transition-all"
                  >
                    {/* Rank */}
                    <div className="w-10 text-center flex-shrink-0">
                      {user.rank <= 3 ? (
                        <span className="text-2xl">{medals[user.rank - 1]}</span>
                      ) : (
                        <span className="text-lg font-bold text-gray-400">{user.rank}</span>
                      )}
                    </div>

                    {/* Avatar */}
                    {user.avatar_url ? (
                      <img src={user.avatar_url} alt={user.name} className="w-10 h-10 rounded-full" />
                    ) : (
                      <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-lg">
                        👤
                      </div>
                    )}

                    {/* Name */}
                    <div className="flex-1 min-w-0">
                      <p className="font-semibold text-gray-900 truncate">{user.name}</p>
                      {user.github_username && (
                        <p className="text-sm text-gray-500 truncate">@{user.github_username}</p>
                      )}
                    </div>

                    {/* Stats */}
                    <div className="flex items-center gap-4 text-sm flex-shrink-0">
                      <div className="text-center">
                        <p className="font-bold text-green-600">Lv.{user.level}</p>
                      </div>
                      <div className="text-center">
                        <p className="font-bold text-yellow-600">{user.total_xp.toLocaleString()}</p>
                        <p className="text-xs text-gray-400">XP</p>
                      </div>
                      <div className="text-center">
                        <p className="font-bold text-blue-600">{user.pokemon_count}</p>
                        <p className="text-xs text-gray-400">匹</p>
                      </div>
                      <div className="text-center">
                        <p className="font-bold text-orange-600">{user.current_streak}</p>
                        <p className="text-xs text-gray-400">🔥</p>
                      </div>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </>
        )}

        {/* Team Rankings Tab */}
        {tab === 'teams' && (
          <>
            {loading ? (
              <p className="text-gray-400 text-center py-12">読み込み中...</p>
            ) : teams.length === 0 ? (
              <div className="bg-white rounded-xl shadow-sm p-8 text-center mb-8">
                <p className="text-4xl mb-4">🏟️</p>
                <p className="text-gray-400 text-lg mb-2">まだチームがありません</p>
                <p className="text-gray-500 text-sm">最初のチームを作って仲間を集めよう！</p>
              </div>
            ) : (
              <div className="space-y-4 mb-8">
                {teams.map((team, idx) => (
                  <div
                    key={team.team_id}
                    className={`rounded-xl border-2 p-5 transition-all ${
                      idx < 3 ? rankColors[idx] : 'bg-white border-gray-200'
                    }`}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <span className="text-3xl font-bold text-gray-300 w-10 text-center">
                          {idx < 3 ? medals[idx] : idx + 1}
                        </span>
                        <div>
                          <h2 className="text-lg font-bold text-gray-900">{team.team_name}</h2>
                          <p className="text-sm text-gray-500">{team.member_count} 人のメンバー</p>
                        </div>
                      </div>
                      <div className="flex gap-6 items-center">
                        <div className="text-center">
                          <p className="text-2xl font-bold text-green-600">{team.total_xp.toLocaleString()}</p>
                          <p className="text-xs text-gray-500">XP</p>
                        </div>
                        <div className="text-center">
                          <p className="text-2xl font-bold text-[#DC0A2D]">{team.pokemon_count}</p>
                          <p className="text-xs text-gray-500">ポケモン</p>
                        </div>
                        <button
                          onClick={() => handleJoinTeam(team.team_id)}
                          className="bg-blue-500 hover:bg-blue-600 text-white text-sm font-medium py-2 px-4 rounded-lg transition-colors"
                        >
                          参加する
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            <div className="bg-white rounded-xl shadow-sm p-6">
              <h2 className="text-lg font-bold mb-4">新しいチームを作る</h2>
              <div className="flex gap-3">
                <input
                  type="text"
                  value={newTeamName}
                  onChange={e => setNewTeamName(e.target.value)}
                  placeholder="チーム名を入力..."
                  className="flex-1 border border-gray-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-gray-800 focus:border-transparent"
                  onKeyDown={e => e.key === 'Enter' && handleCreateTeam()}
                />
                <button
                  onClick={handleCreateTeam}
                  disabled={creating || !newTeamName.trim()}
                  className="bg-gray-900 hover:bg-gray-700 disabled:bg-gray-300 text-white font-medium py-2 px-6 rounded-lg transition-colors"
                >
                  {creating ? '作成中...' : 'チーム作成'}
                </button>
              </div>
            </div>
          </>
        )}
      </main>
    </>
  )
}
