'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { GitHubActivity } from '@/lib/types'

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

const eventColors: Record<string, string> = {
  commit: 'bg-gray-100 text-gray-700',
  pr_merge: 'bg-purple-100 text-purple-700',
  pr_open: 'bg-blue-100 text-blue-700',
  issue_close: 'bg-green-100 text-green-700',
  review: 'bg-yellow-100 text-yellow-700',
}

export default function ActivitiesPage() {
  const router = useRouter()
  const [activities, setActivities] = useState<GitHubActivity[]>([])
  const [totalCount, setTotalCount] = useState(0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    api.getActivities(100)
      .then(data => {
        setActivities(data.activities || [])
        setTotalCount(data.total_count)
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [router])

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
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">GitHub アクティビティ</h1>
            <p className="text-sm text-gray-500">{totalCount} 件のアクティビティ</p>
          </div>
        </div>

        {activities.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            <p className="text-4xl mb-4">💻</p>
            <p className="text-lg font-bold">まだアクティビティがありません</p>
            <p className="text-sm mt-1">GitHubにコミットやPRを作成すると、ここに表示されます</p>
          </div>
        ) : (
          <div className="space-y-3">
            {activities.map(activity => (
              <div
                key={activity.id}
                className="bg-white rounded-xl shadow-sm p-4 border border-gray-100 hover:border-gray-200 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <span className="text-2xl">{eventIcons[activity.event_type] || '📌'}</span>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${eventColors[activity.event_type] || 'bg-gray-100'}`}>
                        {eventLabels[activity.event_type] || activity.event_type}
                      </span>
                      <span className="text-xs text-gray-400">{activity.repo_name}</span>
                    </div>
                    <p className="font-medium text-gray-900 truncate">{activity.title}</p>
                    <p className="text-xs text-gray-400 mt-1">
                      {new Date(activity.created_at).toLocaleString('ja-JP')}
                    </p>
                  </div>
                  <div className="text-right">
                    <span className="text-lg font-bold text-green-600">+{activity.xp}</span>
                    <span className="text-xs text-gray-400 block">XP</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </>
  )
}
