'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { api } from '@/lib/api'

const starters = [
  {
    id: 1,
    name: 'フシギダネ',
    nameEn: 'Bulbasaur',
    type: 'くさ / どく',
    color: 'from-green-400 to-green-600',
    bgColor: 'bg-green-50',
    borderColor: 'border-green-400',
    hoverColor: 'hover:border-green-500 hover:shadow-green-200',
    sprite: 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/1.png',
    description: '背中のタネが日に日に大きくなる。知的で落ち着いた性格。',
  },
  {
    id: 4,
    name: 'ヒトカゲ',
    nameEn: 'Charmander',
    type: 'ほのお',
    color: 'from-orange-400 to-red-500',
    bgColor: 'bg-orange-50',
    borderColor: 'border-orange-400',
    hoverColor: 'hover:border-orange-500 hover:shadow-orange-200',
    sprite: 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/4.png',
    description: 'しっぽの炎は生命力の証。情熱的で冒険好き。',
  },
  {
    id: 7,
    name: 'ゼニガメ',
    nameEn: 'Squirtle',
    type: 'みず',
    color: 'from-blue-400 to-blue-600',
    bgColor: 'bg-blue-50',
    borderColor: 'border-blue-400',
    hoverColor: 'hover:border-blue-500 hover:shadow-blue-200',
    sprite: 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/7.png',
    description: '甲羅は身を守るだけでなく水の抵抗を減らす。忠実で粘り強い。',
  },
]

export default function StarterPage() {
  const router = useRouter()
  const [selected, setSelected] = useState<number | null>(null)
  const [confirming, setConfirming] = useState(false)
  const [result, setResult] = useState<{ pokemon_name: string; sprite_url: string } | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleChoose = async () => {
    if (selected === null) return
    setLoading(true)
    setError('')
    try {
      const res = await api.chooseStarter(selected)
      setResult(res)
    } catch (e) {
      setError(e instanceof Error ? e.message : '選択に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  if (result) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-gray-900 to-gray-800 flex items-center justify-center p-4">
        <div className="text-center">
          <div className="animate-bounce mb-6">
            <img
              src={result.sprite_url}
              alt={result.pokemon_name}
              className="w-48 h-48 mx-auto drop-shadow-2xl"
            />
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">
            おめでとう！
          </h1>
          <p className="text-xl text-gray-300 mb-8">
            <span className="text-yellow-400 font-bold">{starters.find(s => s.id === selected)?.name}</span>
            をパートナーに選んだ！
          </p>
          <p className="text-gray-400 mb-8">
            本を読了してポケモンをたくさんゲットしよう！
          </p>
          <button
            onClick={() => router.push('/')}
            className="bg-red-500 hover:bg-red-600 text-white font-bold py-3 px-8 rounded-full transition-colors shadow-lg"
          >
            冒険をはじめる
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-900 to-gray-800 p-4">
      <div className="max-w-4xl mx-auto pt-12">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-white mb-4">
            パートナーを選ぼう！
          </h1>
          <p className="text-gray-400 text-lg">
            最初のポケモンを1匹選んでください。あなたの読書の旅のパートナーになります。
          </p>
        </div>

        {error && (
          <div className="bg-red-900/50 border border-red-500 text-red-300 px-4 py-3 rounded-lg mb-6 text-center">
            {error}
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          {starters.map((starter) => (
            <button
              key={starter.id}
              onClick={() => { setSelected(starter.id); setConfirming(false) }}
              className={`relative rounded-2xl border-4 p-6 transition-all duration-300 cursor-pointer ${
                selected === starter.id
                  ? `${starter.borderColor} ${starter.bgColor} shadow-xl scale-105`
                  : `border-gray-600 bg-gray-800 ${starter.hoverColor} hover:shadow-lg`
              }`}
            >
              <div className={`absolute inset-0 rounded-2xl bg-gradient-to-b ${starter.color} opacity-${selected === starter.id ? '10' : '0'} transition-opacity`} />
              <div className="relative">
                <img
                  src={starter.sprite}
                  alt={starter.name}
                  className={`w-32 h-32 mx-auto mb-4 drop-shadow-lg transition-transform ${
                    selected === starter.id ? 'scale-110' : ''
                  }`}
                />
                <h2 className={`text-2xl font-bold mb-1 ${selected === starter.id ? 'text-gray-900' : 'text-white'}`}>
                  {starter.name}
                </h2>
                <p className={`text-sm mb-2 ${selected === starter.id ? 'text-gray-600' : 'text-gray-400'}`}>
                  {starter.nameEn}
                </p>
                <span className={`inline-block px-3 py-1 rounded-full text-xs font-semibold bg-gradient-to-r ${starter.color} text-white mb-3`}>
                  {starter.type}
                </span>
                <p className={`text-sm ${selected === starter.id ? 'text-gray-700' : 'text-gray-400'}`}>
                  {starter.description}
                </p>
              </div>
            </button>
          ))}
        </div>

        {selected !== null && !confirming && (
          <div className="text-center">
            <button
              onClick={() => setConfirming(true)}
              className="bg-gradient-to-r from-yellow-400 to-yellow-500 hover:from-yellow-500 hover:to-yellow-600 text-gray-900 font-bold py-3 px-10 rounded-full transition-all shadow-lg text-lg"
            >
              {starters.find(s => s.id === selected)?.name}に決めた！
            </button>
          </div>
        )}

        {confirming && (
          <div className="text-center space-y-4">
            <p className="text-white text-lg">
              <span className="text-yellow-400 font-bold">{starters.find(s => s.id === selected)?.name}</span>
              でいいですか？
            </p>
            <div className="flex justify-center gap-4">
              <button
                onClick={handleChoose}
                disabled={loading}
                className="bg-red-500 hover:bg-red-600 disabled:bg-gray-500 text-white font-bold py-3 px-8 rounded-full transition-colors shadow-lg"
              >
                {loading ? 'ゲット中...' : 'はい！'}
              </button>
              <button
                onClick={() => setConfirming(false)}
                className="bg-gray-600 hover:bg-gray-500 text-white font-bold py-3 px-8 rounded-full transition-colors"
              >
                もどる
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
