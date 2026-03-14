'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Navbar } from '@/components/Navbar'
import { StatusBadge } from '@/components/StatusBadge'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { BookListItem, Tag } from '@/lib/types'

export default function BooksPage() {
  const router = useRouter()
  const [books, setBooks] = useState<BookListItem[]>([])
  const [tags, setTags] = useState<Tag[]>([])
  const [statusFilter, setStatusFilter] = useState('')
  const [tagFilter, setTagFilter] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    Promise.all([
      api.getBooks(),
      api.getTags(),
    ]).then(([booksData, tagsData]) => {
      setBooks(booksData.books || [])
      setTags(tagsData.tags || [])
    }).catch(console.error).finally(() => setLoading(false))
  }, [router])

  const refetch = () => {
    api.getBooks({ status: statusFilter || undefined, tag_id: tagFilter || undefined })
      .then(data => setBooks(data.books || []))
      .catch(console.error)
  }

  useEffect(() => { if (!loading) refetch() }, [statusFilter, tagFilter])

  if (loading) return <div className="min-h-screen flex items-center justify-center"><p>Loading...</p></div>

  return (
    <>
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Books</h1>
          <Link
            href="/books/new"
            className="px-4 py-2 bg-[#DC0A2D] text-white rounded-lg hover:bg-[#b8091f] transition-colors font-medium"
          >
            + Add Book
          </Link>
        </div>

        <div className="flex gap-3 mb-6">
          <select
            value={statusFilter}
            onChange={e => setStatusFilter(e.target.value)}
            className="px-3 py-2 border border-gray-200 rounded-lg text-sm"
          >
            <option value="">All Status</option>
            <option value="unread">未読</option>
            <option value="reading">読書中</option>
            <option value="finished">読了</option>
          </select>
          <select
            value={tagFilter}
            onChange={e => setTagFilter(e.target.value)}
            className="px-3 py-2 border border-gray-200 rounded-lg text-sm"
          >
            <option value="">All Tags</option>
            {tags.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
          </select>
        </div>

        {books.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            <p className="text-lg">No books yet</p>
            <p className="text-sm mt-1">Add your first book to start catching Pokemon!</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {books.map(book => (
              <Link
                key={book.id}
                href={`/books/${book.id}`}
                className="bg-white rounded-xl shadow-sm p-5 hover:shadow-md transition-shadow"
              >
                <div className="flex items-start justify-between mb-2">
                  <h3 className="font-bold text-gray-900 line-clamp-1">{book.title}</h3>
                  <StatusBadge status={book.status} />
                </div>
                <p className="text-sm text-gray-500 mb-3">{book.author}</p>
                {book.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1">
                    {book.tags.map(tag => (
                      <span key={tag} className="px-2 py-0.5 bg-[#FFDE00]/20 text-[#b8991a] rounded text-xs">
                        {tag}
                      </span>
                    ))}
                  </div>
                )}
              </Link>
            ))}
          </div>
        )}
      </main>
    </>
  )
}
