'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { getUser, isLoggedIn } from '@/lib/auth'
import type { Tag } from '@/lib/types'

export default function SettingsPage() {
  const router = useRouter()
  const user = getUser()
  const [tags, setTags] = useState<Tag[]>([])
  const [newTag, setNewTag] = useState('')

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    api.getTags().then(data => setTags(data.tags || [])).catch(console.error)
  }, [router])

  const handleAddTag = async () => {
    if (!newTag.trim()) return
    try {
      const result = await api.createTag(newTag.trim())
      setTags([...tags, { id: result.TagID, name: newTag.trim() }])
      setNewTag('')
    } catch (err) { alert(err instanceof Error ? err.message : 'エラーが発生しました') }
  }

  const handleDeleteTag = async (id: string) => {
    if (!confirm('このタグを削除しますか？')) return
    try {
      await api.deleteTag(id)
      setTags(tags.filter(t => t.id !== id))
    } catch (err) { alert(err instanceof Error ? err.message : 'エラーが発生しました') }
  }

  const roleLabels: Record<string, string> = {
    admin: '管理者',
    member: 'メンバー',
  }

  return (
    <>
      <Navbar />
      <main className="max-w-2xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">設定</h1>

        {/* User Info */}
        <div className="bg-white rounded-xl shadow-sm p-6 mb-6">
          <h2 className="text-lg font-bold mb-4">トレーナー情報</h2>
          <div className="space-y-3 text-sm">
            <div className="flex justify-between items-center py-2 border-b border-gray-100">
              <span className="text-gray-500">名前</span>
              <span className="font-medium">{user?.name}</span>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-gray-100">
              <span className="text-gray-500">メール</span>
              <span className="font-medium">{user?.email}</span>
            </div>
            {user?.github_username && (
              <div className="flex justify-between items-center py-2 border-b border-gray-100">
                <span className="text-gray-500">GitHub</span>
                <span className="font-medium">@{user.github_username}</span>
              </div>
            )}
            <div className="flex justify-between items-center py-2 border-b border-gray-100">
              <span className="text-gray-500">レベル</span>
              <span className="font-medium text-green-600">Lv.{user?.level || 1}</span>
            </div>
            <div className="flex justify-between items-center py-2">
              <span className="text-gray-500">ロール</span>
              <span className="font-medium">{user?.role ? roleLabels[user.role] || user.role : ''}</span>
            </div>
          </div>
        </div>

        {/* Webhook Setup */}
        <div className="bg-white rounded-xl shadow-sm p-6 mb-6">
          <h2 className="text-lg font-bold mb-2">GitHub Webhook 設定</h2>
          <p className="text-sm text-gray-500 mb-4">
            GitHubリポジトリにWebhookを設定すると、コミットやPRマージが自動で記録されます。
          </p>
          <div className="bg-gray-50 rounded-lg p-4 space-y-3 text-sm">
            <div>
              <span className="text-gray-500 block mb-1">Payload URL:</span>
              <code className="bg-gray-200 px-2 py-1 rounded text-xs">
                {typeof window !== 'undefined' ? `${window.location.origin.replace(':3000', ':8080')}/api/webhook/github` : '/api/webhook/github'}
              </code>
            </div>
            <div>
              <span className="text-gray-500 block mb-1">Content type:</span>
              <code className="bg-gray-200 px-2 py-1 rounded text-xs">application/json</code>
            </div>
            <div>
              <span className="text-gray-500 block mb-1">Events:</span>
              <span className="text-gray-700">Push, Pull requests, Pull request reviews, Issues, Deployments</span>
            </div>
          </div>
        </div>

        {/* Admin: Deployer Management */}
        {user?.role === 'admin' && (
          <div className="bg-white rounded-xl shadow-sm p-6 mb-6">
            <h2 className="text-lg font-bold mb-2">管理者メニュー</h2>
            <Link
              href="/admin/deployers"
              className="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <span className="text-2xl">🚀</span>
              <div>
                <p className="font-medium text-gray-900">デプロイヤー管理</p>
                <p className="text-xs text-gray-500">デプロイ報酬を受け取れるユーザーを管理</p>
              </div>
              <span className="ml-auto text-gray-400">→</span>
            </Link>
          </div>
        )}

        {/* Tag Management */}
        <div className="bg-white rounded-xl shadow-sm p-6">
          <h2 className="text-lg font-bold mb-2">タグ管理</h2>
          <p className="text-sm text-gray-500 mb-4">タグはチーム内で共有されます。</p>

          <div className="flex gap-2 mb-4">
            <input
              type="text" value={newTag} onChange={e => setNewTag(e.target.value)}
              className="flex-1 px-3 py-2 border border-gray-200 rounded-lg text-sm outline-none"
              placeholder="新しいタグ名"
              onKeyDown={e => e.key === 'Enter' && handleAddTag()}
            />
            <button onClick={handleAddTag} className="px-4 py-2 bg-gray-900 text-white rounded-lg text-sm hover:bg-gray-700">
              追加
            </button>
          </div>

          {tags.length === 0 ? (
            <p className="text-gray-400 text-sm">まだタグはありません</p>
          ) : (
            <div className="flex flex-wrap gap-2">
              {tags.map(tag => (
                <div key={tag.id} className="flex items-center gap-1 px-3 py-1 bg-gray-100 rounded-full">
                  <span className="text-sm">{tag.name}</span>
                  <button onClick={() => handleDeleteTag(tag.id)} className="text-gray-400 hover:text-red-500 text-xs ml-1">
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </>
  )
}
