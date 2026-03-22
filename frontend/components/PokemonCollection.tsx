'use client'

import { useRouter } from 'next/navigation'
import type { Pokemon, PokemonRarity } from '@/lib/types'

const rarityColors: Record<PokemonRarity, string> = {
  common: 'border-gray-300 bg-gray-50',
  uncommon: 'border-green-400 bg-green-50',
  rare: 'border-blue-400 bg-blue-50',
  epic: 'border-purple-400 bg-purple-50',
  legendary: 'border-yellow-400 bg-yellow-50',
}

const rarityLabels: Record<PokemonRarity, string> = {
  common: 'Common',
  uncommon: 'Uncommon',
  rare: 'Rare',
  epic: 'Epic',
  legendary: 'Legendary',
}

const rarityBadgeColors: Record<PokemonRarity, string> = {
  common: 'bg-gray-200 text-gray-700',
  uncommon: 'bg-green-200 text-green-800',
  rare: 'bg-blue-200 text-blue-800',
  epic: 'bg-purple-200 text-purple-800',
  legendary: 'bg-yellow-200 text-yellow-800',
}

export function PokemonCollection({ pokemon }: { pokemon: Pokemon[] }) {
  const router = useRouter()

  return (
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
        <div className="flex flex-wrap gap-1.5 mt-4">
          {pokemon.slice(0, 12).map(p => (
            <div
              key={p.id}
              className={`relative rounded-lg border-2 ${rarityColors[p.rarity || 'common']} p-0.5`}
              title={`${p.pokemon_name} (${rarityLabels[p.rarity || 'common']})`}
            >
              <img src={p.sprite_url} alt={p.pokemon_name} className="w-12 h-12 pixelated" />
              {(p.rarity === 'epic' || p.rarity === 'legendary') && (
                <span className="absolute -top-1 -right-1 text-xs">
                  {p.rarity === 'legendary' ? '🌟' : '💎'}
                </span>
              )}
            </div>
          ))}
          {pokemon.length > 12 && (
            <button
              onClick={() => router.push('/pokedex')}
              className="w-14 h-14 flex items-center justify-center text-sm text-[#DC0A2D] font-bold hover:bg-red-50 rounded-lg transition-colors border-2 border-dashed border-red-200"
            >
              +{pokemon.length - 12}
            </button>
          )}
        </div>
      )}
      {pokemon.length === 0 && (
        <p className="text-gray-400 text-sm mt-3">PRをマージするとポケモンをゲットできるよ！</p>
      )}
      <div className="flex flex-wrap gap-2 mt-4 pt-3 border-t border-gray-100">
        {(['common', 'uncommon', 'rare', 'epic', 'legendary'] as PokemonRarity[]).map(r => (
          <span key={r} className={`text-xs px-2 py-0.5 rounded-full ${rarityBadgeColors[r]}`}>
            {rarityLabels[r]}
          </span>
        ))}
      </div>
    </div>
  )
}
