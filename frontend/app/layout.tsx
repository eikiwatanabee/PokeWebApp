import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'PokeBookManager',
  description: 'Read books, catch Pokemon!',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body className="bg-gray-50 min-h-screen">{children}</body>
    </html>
  )
}
