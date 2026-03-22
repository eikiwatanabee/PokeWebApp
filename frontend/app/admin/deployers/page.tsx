'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { UserRank } from '@/lib/types'

export default function DeployerManagementPage() {
  const router = useRouter()
  const [users, setUsers] = useState<UserRank[]>([])
  const [deployerIds, setDeployerIds] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [updating, setUpdating] = useState<string | null>(null)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    loadData()
  }, [router])

  const loadData = async () => {
    try {
      const [rankingData, deployerData] = await Promise.all([
        api.getUserRanking('xp'),
        api.getDeployers(),
      ])
      setUsers(rankingData.rankings || [])
      setDeployerIds(deployerData.deployer_ids || [])
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const toggleDeployer = async (userId: string, isCurrentlyDeployer: boolean) => {
    setUpdating(userId)
    try {
      if (isCurrentlyDeployer) {
        await api.removeDeployer(userId)
        setDeployerIds(prev => prev.filter(id => id !== userId))
      } else {
        await api.addDeployer(userId)
        setDeployerIds(prev => [...prev, userId])
      }
    } catch (err) {
      console.error(err)
    } finally {
      setUpdating(null)
    }
  }

  return (
    <>
      <Navbar />
      <main className="max-w-4xl mx-auto px-4 py-8">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900">デプロイヤー管理</h1>
          <p className="text-sm text-gray-500 mt-1">デプロイ報酬を受け取れるユーザーを管理します</p>
        </div>

        {loading ? (
          <p className="text-gray-400 text-center py-12">読み込み中...</p>
        ) : (
          <div className="space-y-2">
            {users.map(user => {
              const isDeployer = deployerIds.includes(user.user_id)
              return (
                <div
                  key={user.user_id}
                  className={`bg-white rounded-xl shadow-sm p-4 border transition-colors ${
                    isDeployer ? 'border-red-200 bg-red-50/30' : 'border-gray-100'
                  }`}
                >
                  <div className="flex items-center gap-3">
                    {user.avatar_url ? (
                      <img src={user.avatar_url} alt={user.name} className="w-10 h-10 rounded-full" />
                    ) : (
                      <span className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-lg">👤</span>
                    )}
                    <div className="flex-1 min-w-0">
                      <p className="font-semibold text-gray-900">{user.name}</p>
                      <p className="text-xs text-gray-500">@{user.github_username} · Lv.{user.level}</p>
                    </div>
                    <button
                      onClick={() => toggleDeployer(user.user_id, isDeployer)}
                      disabled={updating === user.user_id}
                      className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                        isDeployer
                          ? 'bg-red-100 text-red-700 hover:bg-red-200'
                          : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                      } disabled:opacity-50`}
                    >
                      {updating === user.user_id ? '...' : isDeployer ? '🚀 デプロイヤー' : '追加'}
                    </button>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </main>
    </>
  )
}
