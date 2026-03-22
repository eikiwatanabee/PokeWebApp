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
  event_type: 'commit' | 'pr_merge' | 'pr_open' | 'issue_close' | 'review' | 'deploy'
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
  category: string
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

export interface DailyMission {
  id: string
  template_id: string
  title: string
  description: string
  icon: string
  progress: number
  required: number
  bonus_xp: number
  status: 'pending' | 'completed'
}

export interface DailyMissionsResult {
  missions: DailyMission[]
  all_completed: boolean
  completion_bonus: number
  date: string
}

export interface LoginBonusResult {
  claimed: boolean
  bonus_xp: number
  consecutive_days: number
  already_claimed: boolean
}

export interface UserRank {
  user_id: string
  name: string
  avatar_url: string
  github_username: string
  level: number
  total_xp: number
  pokemon_count: number
  current_streak: number
  max_streak: number
  rank: number
}

export interface UserRankingResult {
  rankings: UserRank[]
  sort_by: string
}

export interface FeedItem {
  id: string
  user_id: string
  user_name: string
  user_avatar_url: string
  github_username: string
  event_type: 'commit' | 'pr_merge' | 'pr_open' | 'issue_close' | 'review' | 'deploy'
  repo_name: string
  title: string
  url: string
  xp: number
  created_at: string
}

export interface TeamFeedResult {
  items: FeedItem[]
  total_count: number
}

export interface WeeklyEventTeamScore {
  team_id: string
  team_name: string
  score: number
  rank: number
}

export interface WeeklyEvent {
  event_type: string
  title: string
  description: string
  icon: string
  week_start: string
  week_end: string
  team_scores: WeeklyEventTeamScore[]
}

export interface LimitedEvent {
  id: string
  title: string
  description: string
  icon: string
  rarity_boost: number
  starts_at: string
  ends_at: string
  active: boolean
  hours_left: number
}

export interface LimitedEventsResult {
  events: LimitedEvent[]
}

export interface Trade {
  id: string
  offerer_id: string
  offerer_name: string
  offerer_avatar_url: string
  offered_pokemon_id: string
  offered_pokemon_name: string
  offered_pokemon_sprite: string
  offered_pokemon_rarity: string
  requested_pokemon_name: string
  status: 'open' | 'accepted' | 'cancelled'
  created_at: string
}

export interface TradesResult {
  trades: Trade[]
}

export interface TrainerCard {
  user_id: string
  name: string
  avatar_url: string
  github_username: string
  level: number
  total_xp: number
  current_streak: number
  max_streak: number
  pokemon_count: number
  featured_pokemon: { pokemon_name: string; sprite_url: string; rarity: string }[]
  achievements: { name: string; icon: string }[]
  achievement_count: number
}
