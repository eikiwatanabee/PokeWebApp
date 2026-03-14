'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { BookListItem, Pokemon } from '@/lib/types'

export default function DashboardPage() {
  const router = useRouter()
  const [books, setBooks] = useState<BookListItem[]>([])
  const [pokemon, setPokemon] = useState<Pokemon[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    Promise.all([
      api.getBooks(),
      api.getPokedex(),
    ]).then(([booksData, pokedexData]) => {
      setBooks(booksData.books || [])
      setPokemon(pokedexData.pokemon || [])
    }).catch(console.error).finally(() => setLoading(false))
  }, [router])

  if (loading) return <div className="min-h-screen flex items-center justify-center"><p>Loading...</p></div>

  const unread = books.filter(b => b.status === 'unread').length
  const reading = books.filter(b => b.status === 'reading').length
  const finished = books.filter(b => b.status === 'finished').length

  return (
    <>
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Dashboard</h1>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <StatCard label="Total Books" value={books.length} color="bg-gray-800" />
          <StatCard label="Unread" value={unread} color="bg-gray-400" />
          <StatCard label="Reading" value={reading} color="bg-blue-500" />
          <StatCard label="Finished" value={finished} color="bg-green-500" />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white rounded-xl shadow-sm p-6">
            <h2 className="text-lg font-bold mb-4">Pokemon Collection</h2>
            <p className="text-4xl font-bold text-[#DC0A2D]">{pokemon.length}</p>
            <p className="text-gray-500 text-sm">Pokemon caught</p>
            {pokemon.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-4">
                {pokemon.slice(0, 12).map(p => (
                  <img key={p.id} src={p.sprite_url} alt={p.pokemon_name} className="w-12 h-12" />
                ))}
                {pokemon.length > 12 && (
                  <div className="w-12 h-12 flex items-center justify-center text-sm text-gray-400">
                    +{pokemon.length - 12}
                  </div>
                )}
              </div>
            )}
          </div>

          <div className="bg-white rounded-xl shadow-sm p-6">
            <h2 className="text-lg font-bold mb-4">Currently Reading</h2>
            {books.filter(b => b.status === 'reading').length === 0 ? (
              <p className="text-gray-400">No books in progress</p>
            ) : (
              <ul className="space-y-2">
                {books.filter(b => b.status === 'reading').slice(0, 5).map(b => (
                  <li key={b.id} className="flex items-center gap-2 cursor-pointer hover:bg-gray-50 p-2 rounded" onClick={() => router.push(`/books/${b.id}`)}>
                    <span className="w-2 h-2 bg-blue-500 rounded-full" />
                    <span className="font-medium">{b.title}</span>
                    <span className="text-sm text-gray-400">{b.author}</span>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      </main>
    </>
  )
}

function StatCard({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <div className="bg-white rounded-xl shadow-sm p-6">
      <div className={`w-3 h-3 rounded-full ${color} mb-2`} />
      <p className="text-3xl font-bold">{value}</p>
      <p className="text-sm text-gray-500">{label}</p>
    </div>
  )
}
