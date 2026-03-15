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
  book_count: number
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

  if (loading) return <div className="min-h-screen flex items-center justify-center"><p>Loading...</p></div>

  return (
    <>
      <Navbar />
      <main className="max-w-4xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Team Ranking</h1>
        <p className="text-gray-500 mb-6">チーム対抗！ポケモンゲット数で競おう</p>

        {teams.length === 0 ? (
          <div className="bg-white rounded-xl shadow-sm p-8 text-center">
            <p className="text-gray-400 text-lg mb-4">まだチームがありません</p>
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
                      <p className="text-sm text-gray-500">{team.member_count} members</p>
                    </div>
                  </div>
                  <div className="flex gap-6 items-center">
                    <div className="text-center">
                      <p className="text-2xl font-bold text-[#DC0A2D]">{team.pokemon_count}</p>
                      <p className="text-xs text-gray-500">Pokemon</p>
                    </div>
                    <div className="text-center">
                      <p className="text-2xl font-bold text-green-600">{team.book_count}</p>
                      <p className="text-xs text-gray-500">Books</p>
                    </div>
                    <button
                      onClick={() => handleJoinTeam(team.team_id)}
                      className="bg-blue-500 hover:bg-blue-600 text-white text-sm font-medium py-2 px-4 rounded-lg transition-colors"
                    >
                      Join
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        <div className="bg-white rounded-xl shadow-sm p-6">
          <h2 className="text-lg font-bold mb-4">Create New Team</h2>
          <div className="flex gap-3">
            <input
              type="text"
              value={newTeamName}
              onChange={e => setNewTeamName(e.target.value)}
              placeholder="チーム名を入力..."
              className="flex-1 border border-gray-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-[#DC0A2D] focus:border-transparent"
              onKeyDown={e => e.key === 'Enter' && handleCreateTeam()}
            />
            <button
              onClick={handleCreateTeam}
              disabled={creating || !newTeamName.trim()}
              className="bg-[#DC0A2D] hover:bg-red-700 disabled:bg-gray-300 text-white font-medium py-2 px-6 rounded-lg transition-colors"
            >
              {creating ? '作成中...' : 'Create'}
            </button>
          </div>
        </div>
      </main>
    </>
  )
}
