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

export interface Pokemon {
  id: string
  pokemon_id: number
  pokemon_name: string
  sprite_url: string
  types: string[]
  book_id: string
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
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  user: User
}
