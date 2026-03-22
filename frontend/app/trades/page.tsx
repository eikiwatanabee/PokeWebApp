'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Navbar } from '@/components/Navbar'
import { api } from '@/lib/api'
import { isLoggedIn, getUser } from '@/lib/auth'
import type { Trade, Pokemon } from '@/lib/types'

const rarityColors: Record<string, string> = {
  legendary: 'border-yellow-400 bg-yellow-50',
  epic: 'border-purple-400 bg-purple-50',
  rare: 'border-blue-400 bg-blue-50',
  uncommon: 'border-green-400 bg-green-50',
  common: 'border-gray-300 bg-gray-50',
}

export default function TradesPage() {
  const router = useRouter()
  const user = getUser()
  const [trades, setTrades] = useState<Trade[]>([])
  const [myPokemon, setMyPokemon] = useState<Pokemon[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreate, setShowCreate] = useState(false)
  const [selectedOffer, setSelectedOffer] = useState('')
  const [requestedName, setRequestedName] = useState('')
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    if (!isLoggedIn()) { router.push('/login'); return }
    Promise.all([
      api.getTrades().then(res => setTrades(res.trades || [])),
      api.getPokedex().then(res => setMyPokemon(res.pokemon || [])),
    ])
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [router])

  const loadTrades = () => {
    api.getTrades()
      .then(res => setTrades(res.trades || []))
      .catch(console.error)
  }

  const handleCreate = async () => {
    if (!selectedOffer) return
    setCreating(true)
    try {
      await api.createTrade(selectedOffer, requestedName)
      setShowCreate(false)
      setSelectedOffer('')
      setRequestedName('')
      loadTrades()
    } catch (e) {
      console.error(e)
    } finally {
      setCreating(false)
    }
  }

  const handleCancel = async (tradeId: string) => {
    try {
      await api.cancelTrade(tradeId)
      loadTrades()
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <>
      <Navbar />
      <main className="max-w-4xl mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">ポケモン交換</h1>
            <p className="text-sm text-gray-500">仲間とポケモンを交換しよう</p>
          </div>
          <button
            onClick={() => setShowCreate(!showCreate)}
            className="bg-gray-900 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-700 transition-colors"
          >
            {showCreate ? 'キャンセル' : '交換を出す'}
          </button>
        </div>

        {/* Create Trade Form */}
        {showCreate && (
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 mb-6">
            <h2 className="font-bold mb-4">新しい交換を作成</h2>

            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">出すポケモン</label>
              <div className="grid grid-cols-6 gap-2 max-h-48 overflow-y-auto">
                {myPokemon.map(p => (
                  <button
                    key={p.id}
                    onClick={() => setSelectedOffer(p.id)}
                    className={`rounded-lg border-2 p-1 transition-all ${
                      selectedOffer === p.id
                        ? 'border-gray-900 shadow-md scale-105'
                        : `${rarityColors[p.rarity] || 'border-gray-200'} hover:scale-105`
                    }`}
                    title={p.pokemon_name}
                  >
                    <img src={p.sprite_url} alt={p.pokemon_name} className="w-full pixelated" />
                  </button>
                ))}
              </div>
            </div>

            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                欲しいポケモン名（空欄＝なんでもOK）
              </label>
              <input
                type="text"
                value={requestedName}
                onChange={e => setRequestedName(e.target.value)}
                placeholder="例: pikachu"
                className="w-full border border-gray-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-gray-800"
              />
            </div>

            <button
              onClick={handleCreate}
              disabled={!selectedOffer || creating}
              className="bg-green-600 text-white px-6 py-2 rounded-lg font-medium hover:bg-green-700 disabled:bg-gray-300 transition-colors"
            >
              {creating ? '作成中...' : '交換を公開する'}
            </button>
          </div>
        )}

        {/* Trade List */}
        {loading ? (
          <p className="text-gray-400 text-center py-12">読み込み中...</p>
        ) : trades.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            <p className="text-4xl mb-4">🔄</p>
            <p className="text-lg font-bold">まだ交換がありません</p>
            <p className="text-sm mt-1">最初の交換を出してみよう！</p>
          </div>
        ) : (
          <div className="space-y-3">
            {trades.map(trade => (
              <div
                key={trade.id}
                className="bg-white rounded-xl shadow-sm border border-gray-200 p-4"
              >
                <div className="flex items-center gap-4">
                  {/* Offerer */}
                  <Link
                    href={`/trainers/${trade.offerer_id}`}
                    className="flex items-center gap-2 hover:underline flex-shrink-0"
                  >
                    {trade.offerer_avatar_url ? (
                      <img src={trade.offerer_avatar_url} alt={trade.offerer_name} className="w-8 h-8 rounded-full" />
                    ) : (
                      <span className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-sm">👤</span>
                    )}
                    <span className="text-sm font-medium text-gray-700">{trade.offerer_name}</span>
                  </Link>

                  {/* Offered Pokemon */}
                  <div className={`rounded-lg border-2 p-1 ${rarityColors[trade.offered_pokemon_rarity] || 'border-gray-200'}`}>
                    <img
                      src={trade.offered_pokemon_sprite}
                      alt={trade.offered_pokemon_name}
                      className="w-12 h-12 pixelated"
                    />
                  </div>
                  <div className="text-center">
                    <p className="text-sm font-bold">{trade.offered_pokemon_name}</p>
                    <p className="text-xs text-gray-400">{trade.offered_pokemon_rarity}</p>
                  </div>

                  {/* Arrow */}
                  <span className="text-2xl text-gray-300">⇄</span>

                  {/* Requested */}
                  <div className="flex-1">
                    {trade.requested_pokemon_name ? (
                      <p className="text-sm font-medium text-gray-600">
                        欲しい: <span className="font-bold">{trade.requested_pokemon_name}</span>
                      </p>
                    ) : (
                      <p className="text-sm text-gray-400">なんでもOK</p>
                    )}
                  </div>

                  {/* Actions */}
                  <div className="flex-shrink-0">
                    {user && trade.offerer_id === user.id ? (
                      <button
                        onClick={() => handleCancel(trade.id)}
                        className="text-sm text-red-500 hover:text-red-700 font-medium"
                      >
                        取り消す
                      </button>
                    ) : (
                      <span className="text-xs text-gray-400">
                        {new Date(trade.created_at).toLocaleDateString('ja-JP')}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </>
  )
}
