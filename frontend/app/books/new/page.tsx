'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { Tag } from '@/lib/types'

export default function NewBookPage() {
  const router = useRouter()
  const [title, setTitle] = useState('')
  const [author, setAuthor] = useState('')
  const [selectedTags, setSelectedTags] = useState<string[]>([])
  const [tags, setTags] = useState<Tag[]>([])
  const [newTag, setNewTag] = useState('')
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    api.getTags().then(data => setTags(data.tags || [])).catch(console.error)
  }, [router])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    try {
      await api.createBook({ title, author, tag_ids: selectedTags })
      router.push('/books')
    } catch (err) {
      alert(err instanceof Error ? err.message : '本の作成に失敗しました')
    } finally {
      setSubmitting(false)
    }
  }

  const handleAddTag = async () => {
    if (!newTag.trim()) return
    try {
      const result = await api.createTag(newTag.trim())
      const tag: Tag = { id: result.TagID, name: newTag.trim() }
      setTags([...tags, tag])
      setSelectedTags([...selectedTags, tag.id])
      setNewTag('')
    } catch (err) {
      alert(err instanceof Error ? err.message : 'タグの作成に失敗しました')
    }
  }

  const toggleTag = (id: string) => {
    setSelectedTags(prev => prev.includes(id) ? prev.filter(t => t !== id) : [...prev, id])
  }

  return (
    <>
      <Navbar />
      <main className="max-w-2xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">📕 新しい本を追加</h1>
        <form onSubmit={handleSubmit} className="bg-white rounded-xl shadow-sm p-6 space-y-5">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">タイトル</label>
            <input
              type="text" required value={title} onChange={e => setTitle(e.target.value)}
              className="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-[#DC0A2D] focus:border-transparent outline-none"
              placeholder="本のタイトル"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">著者</label>
            <input
              type="text" required value={author} onChange={e => setAuthor(e.target.value)}
              className="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-[#DC0A2D] focus:border-transparent outline-none"
              placeholder="著者名"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">タグ</label>
            <div className="flex flex-wrap gap-2 mb-2">
              {tags.map(tag => (
                <button
                  key={tag.id} type="button" onClick={() => toggleTag(tag.id)}
                  className={`px-3 py-1 rounded-full text-sm transition-colors ${
                    selectedTags.includes(tag.id)
                      ? 'bg-[#DC0A2D] text-white'
                      : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                  }`}
                >
                  {tag.name}
                </button>
              ))}
            </div>
            <div className="flex gap-2">
              <input
                type="text" value={newTag} onChange={e => setNewTag(e.target.value)}
                className="flex-1 px-3 py-2 border border-gray-200 rounded-lg text-sm outline-none"
                placeholder="新しいタグ名"
                onKeyDown={e => e.key === 'Enter' && (e.preventDefault(), handleAddTag())}
              />
              <button type="button" onClick={handleAddTag} className="px-3 py-2 bg-gray-100 rounded-lg text-sm hover:bg-gray-200">
                追加
              </button>
            </div>
          </div>
          <div className="flex gap-3 pt-2">
            <button type="button" onClick={() => router.back()} className="px-4 py-2 border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-50">
              キャンセル
            </button>
            <button
              type="submit" disabled={submitting}
              className="px-6 py-2 bg-[#DC0A2D] text-white rounded-lg hover:bg-[#b8091f] transition-colors disabled:opacity-50"
            >
              {submitting ? '作成中...' : '本を登録する'}
            </button>
          </div>
        </form>
      </main>
    </>
  )
}
