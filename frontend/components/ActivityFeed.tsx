'use client'

import { useRouter } from 'next/navigation'
import type { GitHubActivity } from '@/lib/types'

const eventIcons: Record<string, string> = {
  commit: '💻',
  pr_merge: '🎉',
  pr_open: '📝',
  issue_close: '✅',
  review: '👀',
}

export function ActivityFeed({ activities }: { activities: GitHubActivity[] }) {
  const router = useRouter()

  return (
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
  )
}
