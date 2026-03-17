'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn } from '@/lib/auth'
import type { Pokemon } from '@/lib/types'

const typeColors: Record<string, string> = {
  normal: 'bg-gray-400', fire: 'bg-red-500', water: 'bg-blue-500',
  electric: 'bg-yellow-400', grass: 'bg-green-500', ice: 'bg-cyan-300',
  fighting: 'bg-red-700', poison: 'bg-purple-500', ground: 'bg-yellow-600',
  flying: 'bg-indigo-300', psychic: 'bg-pink-500', bug: 'bg-lime-500',
  rock: 'bg-yellow-700', ghost: 'bg-purple-700', dragon: 'bg-indigo-600',
  dark: 'bg-gray-700', steel: 'bg-gray-400', fairy: 'bg-pink-300',
}

const typeLabels: Record<string, string> = {
  normal: 'ノーマル', fire: 'ほのお', water: 'みず',
  electric: 'でんき', grass: 'くさ', ice: 'こおり',
  fighting: 'かくとう', poison: 'どく', ground: 'じめん',
  flying: 'ひこう', psychic: 'エスパー', bug: 'むし',
  rock: 'いわ', ghost: 'ゴースト', dragon: 'ドラゴン',
  dark: 'あく', steel: 'はがね', fairy: 'フェアリー',
}

export default function PokedexPage() {
  const router = useRouter()
  const [pokemon, setPokemon] = useState<Pokemon[]>([])
  const [selected, setSelected] = useState<Pokemon | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    api.getPokedex()
      .then(data => {
        const sorted = (data.pokemon || []).sort((a: Pokemon, b: Pokemon) => a.pokemon_id - b.pokemon_id)
        setPokemon(sorted)
      })
      .catch(console.error)
      .finally(() => setLoading(false))
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

  return (
    <>
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 py-8">
        {/* Pokedex Device Frame */}
        <div className="pokedex-frame mb-8 pt-8">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-2xl font-bold text-white">ポケモン図鑑</h1>
              <p className="text-white/70 text-sm">PRマージで ゲットした ポケモンの データ</p>
            </div>
            <div className="bg-white/20 rounded-full px-4 py-2">
              <span className="text-white font-bold text-lg">{pokemon.length}</span>
              <span className="text-white/70 text-sm ml-1">匹 発見</span>
            </div>
          </div>

          {/* Detail Screen */}
          {selected ? (
            <div className="pokedex-screen p-6 mb-4">
              <div className="flex flex-col sm:flex-row items-center gap-6">
                <div className="w-40 h-40 bg-gradient-to-br from-green-900/30 to-green-700/10 rounded-2xl flex items-center justify-center">
                  <img src={selected.sprite_url} alt={selected.pokemon_name} className="w-32 h-32 pixelated drop-shadow-lg" />
                </div>
                <div className="text-center sm:text-left flex-1">
                  <p className="text-green-400 font-mono text-sm">No.{String(selected.pokemon_id).padStart(4, '0')}</p>
                  <h2 className="text-white text-2xl font-bold capitalize mt-1">{selected.pokemon_name}</h2>
                  <div className="flex gap-2 mt-3 justify-center sm:justify-start">
                    {selected.types.map(type => (
                      <span key={type} className={`px-3 py-1 rounded-full text-white text-xs font-bold ${typeColors[type] || 'bg-gray-500'}`}>
                        {typeLabels[type] || type}
                      </span>
                    ))}
                  </div>
                  <p className="text-gray-400 text-xs mt-3">
                    捕獲日: {new Date(selected.caught_at).toLocaleDateString('ja-JP')}
                  </p>
                </div>
              </div>
              <button
                onClick={() => setSelected(null)}
                className="mt-4 text-gray-500 text-xs hover:text-gray-300 w-full text-center"
              >
                ▼ 閉じる
              </button>
            </div>
          ) : (
            <div className="pokedex-screen p-8 mb-4 flex items-center justify-center">
              <p className="text-green-400/50 font-mono text-sm">ポケモンを選択してデータを表示...</p>
            </div>
          )}
        </div>

        {/* Pokemon Grid */}
        {pokemon.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            <div className="w-24 h-24 mx-auto mb-4 bg-[#DC0A2D] rounded-full border-4 border-gray-300 relative opacity-30">
              <div className="absolute top-1/2 left-0 right-0 h-1 bg-gray-300" />
              <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-5 h-5 bg-white rounded-full border-2 border-gray-300" />
            </div>
            <p className="text-lg font-bold">まだポケモンがいません</p>
            <p className="text-sm mt-1">PRをマージして、最初のポケモンをゲットしよう！</p>
          </div>
        ) : (
          <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6 gap-3">
            {pokemon.map(p => (
              <button
                key={p.id}
                onClick={() => setSelected(p)}
                className={`pokemon-card bg-white rounded-xl p-3 text-center cursor-pointer border-2 transition-all ${
                  selected?.id === p.id
                    ? 'border-[#DC0A2D] shadow-lg shadow-red-200'
                    : 'border-transparent hover:border-gray-200'
                }`}
              >
                <div className="bg-gray-50 rounded-lg p-2 mb-2">
                  <img src={p.sprite_url} alt={p.pokemon_name} className="w-16 h-16 mx-auto pixelated" />
                </div>
                <p className="text-[10px] text-gray-400 font-mono">No.{String(p.pokemon_id).padStart(4, '0')}</p>
                <p className="font-bold capitalize text-xs mt-0.5 truncate">{p.pokemon_name}</p>
                <div className="flex justify-center gap-0.5 mt-1.5">
                  {p.types.map(type => (
                    <span key={type} className={`w-2 h-2 rounded-full ${typeColors[type] || 'bg-gray-400'}`} title={typeLabels[type] || type} />
                  ))}
                </div>
              </button>
            ))}
          </div>
        )}
      </main>
    </>
  )
}
