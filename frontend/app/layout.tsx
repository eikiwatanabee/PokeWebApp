import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'PokeGitHub - GitHubで活動してポケモンをゲットしよう！',
  description: 'GitHubのコミットで経験値UP、PRマージでポケモンゲット！',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body className="bg-gray-50 min-h-screen">{children}</body>
    </html>
  )
}
