'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import type { WeeklyEvent } from '@/lib/types'

const medals = ['🥇', '🥈', '🥉']

export function WeeklyEventCard() {
  const [event, setEvent] = useState<WeeklyEvent | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getWeeklyEvent()
      .then(setEvent)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return (
    <div className="bg-gradient-to-r from-indigo-500 to-purple-600 rounded-xl p-6 text-white animate-pulse">
      <div className="h-6 bg-white/20 rounded w-48 mb-2" />
      <div className="h-4 bg-white/20 rounded w-64" />
    </div>
  )

  if (!event) return null

  const weekStartDate = new Date(event.week_start)
  const weekEndDate = new Date(event.week_end)
  const formatDate = (d: Date) => `${d.getMonth() + 1}/${d.getDate()}`

  return (
    <div className="bg-gradient-to-r from-indigo-500 to-purple-600 rounded-xl p-6 text-white shadow-lg mb-6">
      {/* Header */}
      <div className="flex items-center gap-3 mb-4">
        <span className="text-3xl">{event.icon}</span>
        <div>
          <h2 className="text-lg font-bold">週間イベント: {event.title}</h2>
          <p className="text-sm text-white/80">{event.description}</p>
          <p className="text-xs text-white/60 mt-0.5">
            {formatDate(weekStartDate)} ~ {formatDate(weekEndDate)}
          </p>
        </div>
      </div>

      {/* Team Scores */}
      {event.team_scores.length === 0 ? (
        <p className="text-sm text-white/60 text-center py-4">チームを作成して参加しよう！</p>
      ) : (
        <div className="space-y-2">
          {event.team_scores.map((team) => {
            const maxScore = event.team_scores[0]?.score || 1
            const percentage = maxScore > 0 ? (team.score / maxScore) * 100 : 0

            return (
              <div key={team.team_id} className="flex items-center gap-3">
                <span className="w-8 text-center text-lg">
                  {team.rank <= 3 ? medals[team.rank - 1] : team.rank}
                </span>
                <div className="flex-1">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-sm font-medium">{team.team_name}</span>
                    <span className="text-sm font-bold">{team.score.toLocaleString()}</span>
                  </div>
                  <div className="h-2 bg-white/20 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-white/80 rounded-full transition-all duration-500"
                      style={{ width: `${Math.max(percentage, 2)}%` }}
                    />
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
