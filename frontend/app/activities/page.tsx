'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { FeedItem, GitHubActivity } from '@/lib/types'

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

type Tab = 'mine' | 'team'

export default function ActivitiesPage() {
  const router = useRouter()
  const [tab, setTab] = useState<Tab>('mine')
  const [activities, setActivities] = useState<GitHubActivity[]>([])
  const [feedItems, setFeedItems] = useState<FeedItem[]>([])
  const [totalCount, setTotalCount] = useState(0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
  }, [router])

  useEffect(() => {
    setLoading(true)
    if (tab === 'mine') {
      api.getActivities(100)
        .then(data => {
          setActivities(data.activities || [])
          setTotalCount(data.total_count)
        })
        .catch(console.error)
        .finally(() => setLoading(false))
    } else {
      api.getTeamFeed(100)
        .then(data => {
          setFeedItems(data.items || [])
          setTotalCount(data.total_count)
        })
        .catch(console.error)
        .finally(() => setLoading(false))
    }
  }, [tab])

  return (
    <>
      <Navbar />
      <main className="max-w-4xl mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">アクティビティ</h1>
            <p className="text-sm text-gray-500">{totalCount} 件</p>
          </div>
        </div>

        {/* Tab Switcher */}
        <div className="flex gap-2 mb-6">
          <button
            onClick={() => setTab('mine')}
            className={`px-5 py-2.5 rounded-lg text-sm font-medium transition-colors ${
              tab === 'mine' ? 'bg-gray-900 text-white shadow' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            👤 自分
          </button>
          <button
            onClick={() => setTab('team')}
            className={`px-5 py-2.5 rounded-lg text-sm font-medium transition-colors ${
              tab === 'team' ? 'bg-gray-900 text-white shadow' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            👥 チームフィード
          </button>
        </div>

        {loading ? (
          <p className="text-gray-400 text-center py-12">読み込み中...</p>
        ) : tab === 'mine' ? (
          /* My Activities */
          activities.length === 0 ? (
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
          )
        ) : (
          /* Team Feed */
          feedItems.length === 0 ? (
            <div className="text-center py-16 text-gray-400">
              <p className="text-4xl mb-4">👥</p>
              <p className="text-lg font-bold">まだチームのアクティビティがありません</p>
              <p className="text-sm mt-1">チームメンバーがGitHubで活動すると、ここに表示されます</p>
            </div>
          ) : (
            <div className="space-y-3">
              {feedItems.map(item => (
                <div
                  key={item.id}
                  className="bg-white rounded-xl shadow-sm p-4 border border-gray-100 hover:border-gray-200 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-2xl">{eventIcons[item.event_type] || '📌'}</span>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        {/* User info */}
                        <Link
                          href={`/trainers/${item.user_id}`}
                          className="flex items-center gap-1.5 hover:underline"
                        >
                          {item.user_avatar_url ? (
                            <img src={item.user_avatar_url} alt={item.user_name} className="w-5 h-5 rounded-full" />
                          ) : (
                            <span className="w-5 h-5 rounded-full bg-gray-200 flex items-center justify-center text-xs">👤</span>
                          )}
                          <span className="text-sm font-semibold text-gray-700">{item.user_name}</span>
                        </Link>
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${eventColors[item.event_type] || 'bg-gray-100'}`}>
                          {eventLabels[item.event_type] || item.event_type}
                        </span>
                      </div>
                      <p className="font-medium text-gray-900 truncate">{item.title}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <span className="text-xs text-gray-400">{item.repo_name}</span>
                        <span className="text-xs text-gray-300">·</span>
                        <span className="text-xs text-gray-400">
                          {new Date(item.created_at).toLocaleString('ja-JP')}
                        </span>
                      </div>
                    </div>
                    <div className="text-right">
                      <span className="text-lg font-bold text-green-600">+{item.xp}</span>
                      <span className="text-xs text-gray-400 block">XP</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )
        )}
      </main>
    </>
  )
}
