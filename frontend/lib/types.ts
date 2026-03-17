export interface Book {
  id: string
  title: string
  author: string
  status: 'unread' | 'reading' | 'finished'
  finished_at: string | null
  tags: Tag[]
  memos: Memo[]
  created_at: string
  updated_at: string
}

export interface BookListItem {
  id: string
  title: string
  author: string
  status: string
  tags: string[]
}

export interface Memo {
  id: string
  content: string
  created_at: string
  updated_at: string
}

export interface Tag {
  id: string
  name: string
}

export type PokemonRarity = 'common' | 'uncommon' | 'rare' | 'epic' | 'legendary'

export interface Pokemon {
  id: string
  pokemon_id: number
  pokemon_name: string
  sprite_url: string
  types: string[]
  rarity: PokemonRarity
  activity_id: string
  caught_at: string
}

export interface CaughtPokemon {
  pokemon_id: number
  pokemon_name: string
  sprite_url: string
}

export interface User {
  id: string
  name: string
  email: string
  role: string
  github_username?: string
  avatar_url?: string
  level?: number
  total_xp?: number
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  user: User
}

export interface GitHubActivity {
  id: string
  event_type: 'commit' | 'pr_merge' | 'pr_open' | 'issue_close' | 'review'
  repo_name: string
  title: string
  url: string
  xp: number
  created_at: string
}

export interface Achievement {
  type: string
  name: string
  description: string
  icon: string
  unlocked_at: string
}

export interface UserStats {
  total_xp: number
  level: number
  xp_to_next_level: number
  total_commits: number
  total_merges: number
  pokemon_count: number
  current_streak: number
  max_streak: number
  streak_multiplier: number
  achievements: Achievement[]
}
