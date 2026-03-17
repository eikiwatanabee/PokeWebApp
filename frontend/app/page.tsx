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
    api.checkNeedsStarter().then(({ needs_starter }) => {
      if (needs_starter) { router.push('/starter'); return }
      return Promise.all([
        api.getBooks(),
        api.getPokedex(),
      ]).then(([booksData, pokedexData]) => {
        setBooks(booksData.books || [])
        setPokemon(pokedexData.pokemon || [])
      })
    }).catch(console.error).finally(() => setLoading(false))
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

  const unread = books.filter(b => b.status === 'unread').length
  const reading = books.filter(b => b.status === 'reading').length
  const finished = books.filter(b => b.status === 'finished').length

  return (
    <>
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">トレーナーダッシュボード</h1>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
          <StatCard label="総冊数" value={books.length} icon="📚" color="bg-gradient-to-br from-gray-700 to-gray-900" />
          <StatCard label="未読" value={unread} icon="📕" color="bg-gradient-to-br from-gray-400 to-gray-500" />
          <StatCard label="読書中" value={reading} icon="📖" color="bg-gradient-to-br from-blue-400 to-blue-600" />
          <StatCard label="読了" value={finished} icon="✅" color="bg-gradient-to-br from-green-400 to-green-600" />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white rounded-xl shadow-sm p-6 border-2 border-[#DC0A2D]/10">
            <div className="flex items-center gap-2 mb-4">
              <span className="text-xl">⚡</span>
              <h2 className="text-lg font-bold">ゲットしたポケモン</h2>
            </div>
            <div className="flex items-baseline gap-2">
              <p className="text-4xl font-bold text-[#DC0A2D]">{pokemon.length}</p>
              <p className="text-gray-500 text-sm">匹</p>
            </div>
            {pokemon.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-4">
                {pokemon.slice(0, 12).map(p => (
                  <img key={p.id} src={p.sprite_url} alt={p.pokemon_name} className="w-12 h-12 pixelated" />
                ))}
                {pokemon.length > 12 && (
                  <button
                    onClick={() => router.push('/pokedex')}
                    className="w-12 h-12 flex items-center justify-center text-sm text-[#DC0A2D] font-bold hover:bg-red-50 rounded-lg transition-colors"
                  >
                    +{pokemon.length - 12}
                  </button>
                )}
              </div>
            )}
            {pokemon.length === 0 && (
              <p className="text-gray-400 text-sm mt-3">本を読み終えるとポケモンをゲットできるよ！</p>
            )}
          </div>

          <div className="bg-white rounded-xl shadow-sm p-6 border-2 border-blue-100">
            <div className="flex items-center gap-2 mb-4">
              <span className="text-xl">📖</span>
              <h2 className="text-lg font-bold">いま読んでいる本</h2>
            </div>
            {books.filter(b => b.status === 'reading').length === 0 ? (
              <div className="text-center py-4">
                <p className="text-gray-400">読書中の本はありません</p>
                <button
                  onClick={() => router.push('/books')}
                  className="mt-2 text-sm text-[#DC0A2D] hover:underline"
                >
                  本棚を見る →
                </button>
              </div>
            ) : (
              <ul className="space-y-2">
                {books.filter(b => b.status === 'reading').slice(0, 5).map(b => (
                  <li key={b.id} className="flex items-center gap-3 cursor-pointer hover:bg-blue-50 p-2.5 rounded-lg transition-colors" onClick={() => router.push(`/books/${b.id}`)}>
                    <span className="w-2.5 h-2.5 bg-blue-500 rounded-full flex-shrink-0" />
                    <div className="min-w-0">
                      <span className="font-medium block truncate">{b.title}</span>
                      <span className="text-xs text-gray-400">{b.author}</span>
                    </div>
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

function StatCard({ label, value, icon, color }: { label: string; value: number; icon: string; color: string }) {
  return (
    <div className={`${color} rounded-xl shadow-sm p-5 text-white`}>
      <span className="text-2xl">{icon}</span>
      <p className="text-3xl font-bold mt-1">{value}</p>
      <p className="text-sm text-white/80">{label}</p>
    </div>
  )
}
