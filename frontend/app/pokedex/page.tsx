'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { Pokemon } from '@/lib/types'

const typeColors: Record<string, string> = {
  normal: 'bg-gray-400', fire: 'bg-red-500', water: 'bg-blue-500',
  electric: 'bg-yellow-400', grass: 'bg-green-500', ice: 'bg-blue-200',
  fighting: 'bg-red-700', poison: 'bg-purple-500', ground: 'bg-yellow-600',
  flying: 'bg-indigo-300', psychic: 'bg-pink-500', bug: 'bg-lime-500',
  rock: 'bg-yellow-700', ghost: 'bg-purple-700', dragon: 'bg-indigo-600',
  dark: 'bg-gray-700', steel: 'bg-gray-400', fairy: 'bg-pink-300',
}

export default function PokedexPage() {
  const router = useRouter()
  const [pokemon, setPokemon] = useState<Pokemon[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    api.getPokedex()
      .then(data => setPokemon(data.pokemon || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [router])

  if (loading) return <div className="min-h-screen flex items-center justify-center"><p>Loading...</p></div>

  return (
    <>
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Pokedex</h1>
          <span className="text-gray-500">{pokemon.length} caught</span>
        </div>

        {pokemon.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            <div className="w-24 h-24 mx-auto mb-4 bg-gray-200 rounded-full border-4 border-gray-300 relative">
              <div className="absolute top-1/2 left-0 right-0 h-1 bg-gray-300" />
              <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-5 h-5 bg-white rounded-full border-2 border-gray-300" />
            </div>
            <p className="text-lg">No Pokemon yet</p>
            <p className="text-sm mt-1">Finish reading a book to catch your first Pokemon!</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {pokemon.map(p => (
              <div
                key={p.id}
                className="bg-white rounded-xl shadow-sm p-4 text-center hover:shadow-md transition-shadow cursor-pointer"
                onClick={() => router.push(`/books/${p.book_id}`)}
              >
                <img src={p.sprite_url} alt={p.pokemon_name} className="w-20 h-20 mx-auto" />
                <p className="font-bold capitalize text-sm mt-1">{p.pokemon_name}</p>
                <p className="text-xs text-gray-400">#{p.pokemon_id}</p>
                <div className="flex justify-center gap-1 mt-2">
                  {p.types.map(type => (
                    <span key={type} className={`px-2 py-0.5 rounded text-white text-[10px] ${typeColors[type] || 'bg-gray-400'}`}>
                      {type}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </>
  )
}
