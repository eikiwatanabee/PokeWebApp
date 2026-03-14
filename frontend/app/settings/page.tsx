'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
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
    } catch (err) { alert(err instanceof Error ? err.message : 'Error') }
  }

  const handleDeleteTag = async (id: string) => {
    if (!confirm('Delete this tag?')) return
    try {
      await api.deleteTag(id)
      setTags(tags.filter(t => t.id !== id))
    } catch (err) { alert(err instanceof Error ? err.message : 'Error') }
  }

  return (
    <>
      <Navbar />
      <main className="max-w-2xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Settings</h1>

        {/* User Info */}
        <div className="bg-white rounded-xl shadow-sm p-6 mb-6">
          <h2 className="text-lg font-bold mb-4">User Info</h2>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-500">Name</span>
              <span className="font-medium">{user?.name}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">Email</span>
              <span className="font-medium">{user?.email}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">Role</span>
              <span className="font-medium capitalize">{user?.role}</span>
            </div>
          </div>
        </div>

        {/* Tag Management */}
        <div className="bg-white rounded-xl shadow-sm p-6">
          <h2 className="text-lg font-bold mb-4">Tag Management</h2>
          <p className="text-sm text-gray-500 mb-4">Tags are shared across your team.</p>

          <div className="flex gap-2 mb-4">
            <input
              type="text" value={newTag} onChange={e => setNewTag(e.target.value)}
              className="flex-1 px-3 py-2 border border-gray-200 rounded-lg text-sm outline-none"
              placeholder="New tag name"
              onKeyDown={e => e.key === 'Enter' && handleAddTag()}
            />
            <button onClick={handleAddTag} className="px-4 py-2 bg-[#DC0A2D] text-white rounded-lg text-sm hover:bg-[#b8091f]">
              Add
            </button>
          </div>

          {tags.length === 0 ? (
            <p className="text-gray-400 text-sm">No tags yet</p>
          ) : (
            <div className="flex flex-wrap gap-2">
              {tags.map(tag => (
                <div key={tag.id} className="flex items-center gap-1 px-3 py-1 bg-gray-100 rounded-full">
                  <span className="text-sm">{tag.name}</span>
                  <button onClick={() => handleDeleteTag(tag.id)} className="text-gray-400 hover:text-red-500 text-xs ml-1">
                    x
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
