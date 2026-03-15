'use client'

import { useState, useEffect } from 'react'
import type { CaughtPokemon } from '@/lib/types'

interface Props {
  pokemon: CaughtPokemon | null
  onClose: () => void
}

export function PokemonCatchModal({ pokemon, onClose }: Props) {
  const [phase, setPhase] = useState<'shake' | 'caught' | 'reveal'>('shake')

  useEffect(() => {
    if (!pokemon) return
    setPhase('shake')
    const t1 = setTimeout(() => setPhase('caught'), 1500)
    const t2 = setTimeout(() => setPhase('reveal'), 2500)
    return () => { clearTimeout(t1); clearTimeout(t2) }
  }, [pokemon])

  if (!pokemon) return null

  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50" onClick={onClose}>
      <div
        className="bg-white rounded-2xl p-8 max-w-sm w-full mx-4 text-center"
        onClick={(e) => e.stopPropagation()}
      >
        {phase === 'shake' && (
          <div className="animate-bounce">
            <div className="w-24 h-24 mx-auto bg-[#DC0A2D] rounded-full border-4 border-gray-800 relative">
              <div className="absolute top-1/2 left-0 right-0 h-1 bg-gray-800" />
              <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-6 h-6 bg-white rounded-full border-2 border-gray-800" />
            </div>
            <p className="mt-4 text-lg font-bold text-gray-700">...</p>
          </div>
        )}

        {phase === 'caught' && (
          <div>
            <div className="text-6xl mb-4">*</div>
            <p className="text-xl font-bold text-[#DC0A2D]">Gotcha!</p>
          </div>
        )}

        {phase === 'reveal' && (
          <div>
            <img
              src={pokemon.sprite_url}
              alt={pokemon.pokemon_name}
              className="w-32 h-32 mx-auto pixelated"
            />
            <h3 className="text-2xl font-bold mt-2 capitalize">{pokemon.pokemon_name}</h3>
            <p className="text-gray-500 mt-1">#{pokemon.pokemon_id}</p>
            <p className="text-lg font-bold text-[#DC0A2D] mt-4">
              Pokemon GET!
            </p>
            <button
              onClick={onClose}
              className="mt-6 px-6 py-2 bg-[#DC0A2D] text-white rounded-lg hover:bg-[#b8091f] transition-colors"
            >
              OK
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
