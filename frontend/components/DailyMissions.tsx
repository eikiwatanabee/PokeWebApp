'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import type { DailyMission, DailyMissionsResult, LoginBonusResult } from '@/lib/types'

export function DailyMissions() {
  const [data, setData] = useState<DailyMissionsResult | null>(null)
  const [loginBonus, setLoginBonus] = useState<LoginBonusResult | null>(null)
  const [claiming, setClaiming] = useState(false)

  useEffect(() => {
    api.getDailyMissions().then(setData).catch(console.error)
    api.claimLoginBonus().then(setLoginBonus).catch(console.error)
  }, [])

  const handleClaimBonus = async () => {
    if (claiming) return
    setClaiming(true)
    try {
      const result = await api.claimLoginBonus()
      setLoginBonus(result)
    } catch (e) {
      console.error(e)
    } finally {
      setClaiming(false)
    }
  }

  return (
    <div className="bg-white rounded-xl shadow-sm p-6 border-2 border-indigo-100">
      {/* Login Bonus */}
      {loginBonus && (
        <div className="mb-5 p-4 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-xl text-white">
          <div className="flex items-center justify-between">
            <div>
              <div className="flex items-center gap-2">
                <span className="text-xl">🎁</span>
                <h3 className="font-bold">ログインボーナス</h3>
              </div>
              {loginBonus.already_claimed ? (
                <p className="text-sm text-white/80 mt-1">
                  本日受取済み (+{loginBonus.bonus_xp} XP) | {loginBonus.consecutive_days}日連続
                </p>
              ) : loginBonus.claimed ? (
                <p className="text-sm text-white/80 mt-1">
                  +{loginBonus.bonus_xp} XP ゲット! | {loginBonus.consecutive_days}日連続
                </p>
              ) : null}
            </div>
            <div className="text-right">
              <span className="text-2xl font-bold">Day {loginBonus.consecutive_days}</span>
            </div>
          </div>
          {/* Streak dots */}
          <div className="flex gap-1 mt-3">
            {Array.from({ length: 7 }, (_, i) => (
              <div
                key={i}
                className={`h-2 flex-1 rounded-full ${
                  i < (loginBonus.consecutive_days % 7 || 7)
                    ? 'bg-white'
                    : 'bg-white/20'
                }`}
              />
            ))}
          </div>
        </div>
      )}

      {/* Daily Missions */}
      <div className="flex items-center gap-2 mb-4">
        <span className="text-xl">🎯</span>
        <h2 className="text-lg font-bold">デイリーミッション</h2>
        {data?.all_completed && (
          <span className="text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded-full ml-auto font-bold">
            ALL CLEAR! +{data.completion_bonus} XP
          </span>
        )}
      </div>

      {data?.missions && data.missions.length > 0 ? (
        <div className="space-y-3">
          {data.missions.map(mission => (
            <MissionCard key={mission.id} mission={mission} />
          ))}
        </div>
      ) : (
        <p className="text-gray-400 text-sm">ミッションを読み込み中...</p>
      )}
    </div>
  )
}

function MissionCard({ mission }: { mission: DailyMission }) {
  const isComplete = mission.status === 'completed'
  const progress = Math.min(mission.progress / mission.required * 100, 100)

  return (
    <div className={`p-3 rounded-lg border-2 transition-all ${
      isComplete
        ? 'border-green-300 bg-green-50'
        : 'border-gray-200 bg-gray-50'
    }`}>
      <div className="flex items-center gap-3">
        <span className="text-2xl">{mission.icon}</span>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className={`font-bold text-sm ${isComplete ? 'text-green-700' : 'text-gray-900'}`}>
              {mission.title}
            </span>
            {isComplete && <span className="text-xs">✅</span>}
          </div>
          <span className="text-xs text-gray-500">{mission.description}</span>
          {/* Progress bar */}
          <div className="mt-1.5 flex items-center gap-2">
            <div className="flex-1 bg-gray-200 rounded-full h-2">
              <div
                className={`h-2 rounded-full transition-all duration-300 ${
                  isComplete ? 'bg-green-500' : 'bg-indigo-500'
                }`}
                style={{ width: `${progress}%` }}
              />
            </div>
            <span className="text-xs font-mono text-gray-500 whitespace-nowrap">
              {mission.progress}/{mission.required}
            </span>
          </div>
        </div>
        <span className={`text-xs font-bold px-2 py-1 rounded-lg whitespace-nowrap ${
          isComplete
            ? 'bg-green-200 text-green-800'
            : 'bg-indigo-100 text-indigo-700'
        }`}>
          +{mission.bonus_xp} XP
        </span>
      </div>
    </div>
  )
}
