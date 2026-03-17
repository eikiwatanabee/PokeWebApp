'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { getUser, logout } from '@/lib/auth'

const navItems = [
  { href: '/', label: 'ホーム', icon: '🏠' },
  { href: '/books', label: '本棚', icon: '📚' },
  { href: '/pokedex', label: '図鑑', icon: '📖' },
  { href: '/ranking', label: 'ランキング', icon: '🏆' },
  { href: '/settings', label: '設定', icon: '⚙️' },
]

export function Navbar() {
  const pathname = usePathname()
  const user = getUser()

  return (
    <nav className="bg-gradient-to-r from-[#DC0A2D] to-[#b8091f] text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <Link href="/" className="flex items-center gap-2">
            <div className="w-8 h-8 bg-white rounded-full border-2 border-gray-800 relative flex-shrink-0">
              <div className="absolute top-1/2 left-0 right-0 h-0.5 bg-gray-800" />
              <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-2.5 h-2.5 bg-white rounded-full border-[1.5px] border-gray-800" />
            </div>
            <span className="text-xl font-bold tracking-tight">PokeBookManager</span>
          </Link>
          <div className="flex items-center gap-1">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-1.5 ${
                  pathname === item.href
                    ? 'bg-white/25 shadow-inner'
                    : 'hover:bg-white/10'
                }`}
              >
                <span className="text-base">{item.icon}</span>
                <span className="hidden sm:inline">{item.label}</span>
              </Link>
            ))}
            {user && (
              <div className="flex items-center gap-3 ml-4 pl-4 border-l border-white/30">
                <span className="text-sm hidden md:inline">{user.name}</span>
                <button
                  onClick={logout}
                  className="text-sm px-3 py-1 rounded bg-white/10 hover:bg-white/20 transition-colors"
                >
                  ログアウト
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}
