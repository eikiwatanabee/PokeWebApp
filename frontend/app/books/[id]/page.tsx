'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { StatusBadge } from '@/components/StatusBadge'
import { PokemonCatchModal } from '@/components/PokemonCatchModal'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { Book, CaughtPokemon } from '@/lib/types'

export default function BookDetailPage() {
  const params = useParams()
  const router = useRouter()
  const id = params.id as string
  const [book, setBook] = useState<Book | null>(null)
  const [newMemo, setNewMemo] = useState('')
  const [caughtPokemon, setCaughtPokemon] = useState<CaughtPokemon | null>(null)
  const [loading, setLoading] = useState(true)

  const fetchBook = () => {
    api.getBook(id).then(setBook).catch(console.error).finally(() => setLoading(false))
  }

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    fetchBook()
  }, [id, router])

  const handleStartReading = async () => {
    try {
      await api.startReading(id)
      fetchBook()
    } catch (err) { alert(err instanceof Error ? err.message : 'Error') }
  }

  const handleFinishReading = async () => {
    try {
      const result = await api.finishReading(id)
      setCaughtPokemon(result.pokemon)
      fetchBook()
    } catch (err) { alert(err instanceof Error ? err.message : 'Error') }
  }

  const handleAddMemo = async () => {
    if (!newMemo.trim()) return
    try {
      await api.addMemo(id, newMemo.trim())
      setNewMemo('')
      fetchBook()
    } catch (err) { alert(err instanceof Error ? err.message : 'Error') }
  }

  const handleDeleteMemo = async (memoId: string) => {
    if (!confirm('Delete this memo?')) return
    try {
      await api.deleteMemo(memoId)
      fetchBook()
    } catch (err) { alert(err instanceof Error ? err.message : 'Error') }
  }

  const handleDeleteBook = async () => {
    if (!confirm('Delete this book?')) return
    try {
      await api.deleteBook(id)
      router.push('/books')
    } catch (err) { alert(err instanceof Error ? err.message : 'Error') }
  }

  if (loading || !book) return <div className="min-h-screen flex items-center justify-center"><p>Loading...</p></div>

  return (
    <>
      <Navbar />
      <main className="max-w-4xl mx-auto px-4 py-8">
        <button onClick={() => router.push('/books')} className="text-sm text-gray-500 hover:text-gray-700 mb-4 inline-block">
          &larr; Back to Books
        </button>

        <div className="bg-white rounded-xl shadow-sm p-6 mb-6">
          <div className="flex items-start justify-between mb-4">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">{book.title}</h1>
              <p className="text-gray-500 mt-1">{book.author}</p>
            </div>
            <StatusBadge status={book.status} />
          </div>

          {book.tags.length > 0 && (
            <div className="flex flex-wrap gap-2 mb-4">
              {book.tags.map(tag => (
                <span key={tag.id} className="px-2 py-1 bg-[#FFDE00]/20 text-[#b8991a] rounded text-xs">
                  {tag.name}
                </span>
              ))}
            </div>
          )}

          <div className="flex gap-2 pt-2">
            {book.status === 'unread' && (
              <button onClick={handleStartReading} className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">
                Start Reading
              </button>
            )}
            {book.status === 'reading' && (
              <button onClick={handleFinishReading} className="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors">
                Finish Reading
              </button>
            )}
            <button onClick={handleDeleteBook} className="px-4 py-2 border border-red-200 text-red-500 rounded-lg hover:bg-red-50 transition-colors">
              Delete
            </button>
          </div>
        </div>

        {/* Memos */}
        <div className="bg-white rounded-xl shadow-sm p-6">
          <h2 className="text-lg font-bold mb-4">Memos</h2>

          <div className="flex gap-2 mb-4">
            <textarea
              value={newMemo} onChange={e => setNewMemo(e.target.value)}
              className="flex-1 px-3 py-2 border border-gray-200 rounded-lg text-sm outline-none resize-none"
              rows={2} placeholder="Write a memo..."
            />
            <button onClick={handleAddMemo} className="px-4 py-2 bg-[#DC0A2D] text-white rounded-lg hover:bg-[#b8091f] transition-colors self-end">
              Add
            </button>
          </div>

          {book.memos.length === 0 ? (
            <p className="text-gray-400 text-sm">No memos yet</p>
          ) : (
            <ul className="space-y-3">
              {book.memos.map(memo => (
                <li key={memo.id} className="bg-gray-50 rounded-lg p-4">
                  <p className="text-gray-800 whitespace-pre-wrap">{memo.content}</p>
                  <div className="flex items-center justify-between mt-2">
                    <span className="text-xs text-gray-400">
                      {new Date(memo.created_at).toLocaleDateString('ja-JP')}
                    </span>
                    <button onClick={() => handleDeleteMemo(memo.id)} className="text-xs text-red-400 hover:text-red-600">
                      Delete
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      </main>

      <PokemonCatchModal pokemon={caughtPokemon} onClose={() => setCaughtPokemon(null)} />
    </>
  )
}
