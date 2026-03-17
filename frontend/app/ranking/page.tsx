'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'

type TeamRanking = {
  team_id: string
  team_name: string
  member_count: number
  pokemon_count: number
  total_xp: number
}

const medals = ['🥇', '🥈', '🥉']
const rankColors = ['bg-yellow-50 border-yellow-400', 'bg-gray-50 border-gray-300', 'bg-orange-50 border-orange-300']

export default function RankingPage() {
  const router = useRouter()
  const [teams, setTeams] = useState<TeamRanking[]>([])
  const [loading, setLoading] = useState(true)
  const [newTeamName, setNewTeamName] = useState('')
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    loadRanking()
  }, [router])

  const loadRanking = () => {
    api.getTeamRanking()
      .then(data => setTeams(data.teams || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  const handleCreateTeam = async () => {
    if (!newTeamName.trim()) return
    setCreating(true)
    try {
      const result = await api.createTeam(newTeamName.trim())
      await api.joinTeam(result.TeamID)
      setNewTeamName('')
      loadRanking()
    } catch (e) {
      console.error(e)
    } finally {
      setCreating(false)
    }
  }

  const handleJoinTeam = async (teamId: string) => {
    try {
      await api.joinTeam(teamId)
      loadRanking()
    } catch (e) {
      console.error(e)
    }
  }

  if (loading) return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="animate-pokeball-shake">
        <div className="w-16 h-16 bg-[#DC0A2D] rounded-full border-4 border-gray-800 relative">
          <div className="absolute top-1/2 left-0 right-0 h-1 bg-gray-800" />
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 bg-white rounded-full border-2 border-gray-800" />
        </div>
      </div>
    </div>
  )

  return (
    <>
      <Navbar />
      <main className="max-w-4xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">チームランキング</h1>
        <p className="text-gray-500 mb-6">チーム対抗！XPとポケモン数で競おう</p>

        {teams.length === 0 ? (
          <div className="bg-white rounded-xl shadow-sm p-8 text-center">
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
      </main>
    </>
  )
}
