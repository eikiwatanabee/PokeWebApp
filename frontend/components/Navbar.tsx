'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { getUser, logout } from '@/lib/auth'

const navItems = [
  { href: '/', label: 'Dashboard' },
  { href: '/books', label: 'Books' },
  { href: '/pokedex', label: 'Pokedex' },
  { href: '/ranking', label: 'Ranking' },
  { href: '/settings', label: 'Settings' },
]

export function Navbar() {
  const pathname = usePathname()
  const user = getUser()

  return (
    <nav className="bg-[#DC0A2D] text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <Link href="/" className="text-xl font-bold tracking-tight">
            PokeBookManager
          </Link>
          <div className="flex items-center gap-1">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  pathname === item.href
                    ? 'bg-white/20'
                    : 'hover:bg-white/10'
                }`}
              >
                {item.label}
              </Link>
            ))}
            {user && (
              <div className="flex items-center gap-3 ml-4 pl-4 border-l border-white/30">
                <span className="text-sm">{user.name}</span>
                <button
                  onClick={logout}
                  className="text-sm px-3 py-1 rounded bg-white/10 hover:bg-white/20 transition-colors"
                >
                  Logout
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}
