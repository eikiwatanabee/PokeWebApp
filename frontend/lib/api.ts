import type { AuthTokens, Book, BookListItem, CaughtPokemon, DailyMissionsResult, GitHubActivity, LimitedEventsResult, LoginBonusResult, Memo, Pokemon, Tag, TeamFeedResult, TradesResult, TrainerCard, UserRankingResult, UserStats, WeeklyEvent } from './types'

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

class ApiClient {
  private getToken(): string | null {
    if (typeof window === 'undefined') return null
    return localStorage.getItem('access_token')
  }

  private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const token = this.getToken()
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...((options.headers as Record<string, string>) || {}),
    }
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const res = await fetch(`${API_BASE}${path}`, { ...options, headers })

    if (res.status === 401) {
      // Try refresh
      const refreshed = await this.refreshToken()
      if (refreshed) {
        headers['Authorization'] = `Bearer ${this.getToken()}`
        const retry = await fetch(`${API_BASE}${path}`, { ...options, headers })
        if (!retry.ok) throw new Error(`API error: ${retry.status}`)
        return retry.json()
      }
      // Redirect to login
      if (typeof window !== 'undefined') {
        window.location.href = '/login'
      }
      throw new Error('Unauthorized')
    }

    if (!res.ok) {
      const body = await res.json().catch(() => ({}))
      throw new Error(body.error || `API error: ${res.status}`)
    }

    return res.json()
  }

  private async refreshToken(): Promise<boolean> {
    const refreshToken = localStorage.getItem('refresh_token')
    if (!refreshToken) return false
    try {
      const res = await fetch(`${API_BASE}/api/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      })
      if (!res.ok) return false
      const data = await res.json()
      localStorage.setItem('access_token', data.access_token)
      return true
    } catch {
      return false
    }
  }

  // Auth
  async getGoogleAuthURL(): Promise<{ url: string }> {
    return this.request('/api/auth/google')
  }

  async googleCallback(code: string): Promise<AuthTokens> {
    return this.request(`/api/auth/google/callback?code=${code}`)
  }

  async getGitHubAuthURL(): Promise<{ url: string }> {
    return this.request('/api/auth/github')
  }

  async githubCallback(code: string): Promise<AuthTokens> {
    return this.request(`/api/auth/github/callback?code=${code}`)
  }

  async devLogin(): Promise<AuthTokens> {
    return this.request('/api/auth/dev-login', { method: 'POST' })
  }

  // GitHub Activities
  async getActivities(limit?: number): Promise<{ activities: GitHubActivity[]; total_count: number }> {
    const qs = limit ? `?limit=${limit}` : ''
    return this.request(`/api/activities${qs}`)
  }

  // User Stats
  async getStats(): Promise<UserStats> {
    return this.request('/api/stats')
  }

  // Books
  async getBooks(params?: { status?: string; tag_id?: string; page?: number }): Promise<{ books: BookListItem[]; total_count: number; page: number }> {
    const searchParams = new URLSearchParams()
    if (params?.status) searchParams.set('status', params.status)
    if (params?.tag_id) searchParams.set('tag_id', params.tag_id)
    if (params?.page) searchParams.set('page', params.page.toString())
    const qs = searchParams.toString()
    return this.request(`/api/books${qs ? `?${qs}` : ''}`)
  }

  async getBook(id: string): Promise<Book> {
    return this.request(`/api/books/${id}`)
  }

  async createBook(data: { title: string; author: string; tag_ids?: string[] }): Promise<{ BookID: string }> {
    return this.request('/api/books', { method: 'POST', body: JSON.stringify(data) })
  }

  async updateBook(id: string, data: { title?: string; author?: string; tag_ids?: string[] }): Promise<void> {
    await this.request(`/api/books/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  }

  async deleteBook(id: string): Promise<void> {
    await this.request(`/api/books/${id}`, { method: 'DELETE' })
  }

  async startReading(id: string): Promise<void> {
    await this.request(`/api/books/${id}/start`, { method: 'POST' })
  }

  async finishReading(id: string): Promise<{ message: string; pokemon: CaughtPokemon }> {
    return this.request(`/api/books/${id}/finish`, { method: 'POST' })
  }

  // Memos
  async getMemos(bookId: string): Promise<{ memos: Memo[] }> {
    return this.request(`/api/books/${bookId}/memos`)
  }

  async addMemo(bookId: string, content: string): Promise<{ MemoID: string }> {
    return this.request(`/api/books/${bookId}/memos`, { method: 'POST', body: JSON.stringify({ content }) })
  }

  async updateMemo(memoId: string, content: string): Promise<void> {
    await this.request(`/api/memos/${memoId}`, { method: 'PUT', body: JSON.stringify({ content }) })
  }

  async deleteMemo(memoId: string): Promise<void> {
    await this.request(`/api/memos/${memoId}`, { method: 'DELETE' })
  }

  // Tags
  async getTags(): Promise<{ tags: Tag[] }> {
    return this.request('/api/tags')
  }

  async createTag(name: string): Promise<{ TagID: string }> {
    return this.request('/api/tags', { method: 'POST', body: JSON.stringify({ name }) })
  }

  async deleteTag(id: string): Promise<void> {
    await this.request(`/api/tags/${id}`, { method: 'DELETE' })
  }

  // Pokedex
  async getPokedex(): Promise<{ pokemon: Pokemon[]; total: number }> {
    return this.request('/api/pokedex')
  }

  // Starter Pokemon
  async checkNeedsStarter(): Promise<{ needs_starter: boolean }> {
    return this.request('/api/starter/check')
  }

  async chooseStarter(pokemonId: number): Promise<{ pokemon_id: number; pokemon_name: string; sprite_url: string }> {
    return this.request('/api/starter/choose', { method: 'POST', body: JSON.stringify({ pokemon_id: pokemonId }) })
  }

  // Teams
  async createTeam(name: string): Promise<{ TeamID: string }> {
    return this.request('/api/teams', { method: 'POST', body: JSON.stringify({ name }) })
  }

  async joinTeam(teamId: string): Promise<void> {
    await this.request(`/api/teams/${teamId}/join`, { method: 'POST' })
  }

  async getTeamRanking(): Promise<{ teams: { team_id: string; team_name: string; member_count: number; pokemon_count: number; total_xp: number }[] }> {
    return this.request('/api/teams/ranking')
  }

  // Daily Missions
  async getDailyMissions(): Promise<DailyMissionsResult> {
    return this.request('/api/daily/missions')
  }

  // Login Bonus
  async claimLoginBonus(): Promise<LoginBonusResult> {
    return this.request('/api/daily/login-bonus', { method: 'POST' })
  }

  // User Ranking
  async getUserRanking(sortBy?: string): Promise<UserRankingResult> {
    const qs = sortBy ? `?sort_by=${sortBy}` : ''
    return this.request(`/api/ranking/users${qs}`)
  }

  // Team Feed
  async getTeamFeed(limit?: number): Promise<TeamFeedResult> {
    const qs = limit ? `?limit=${limit}` : ''
    return this.request(`/api/feed${qs}`)
  }

  // Limited Events
  async getLimitedEvents(): Promise<LimitedEventsResult> {
    return this.request('/api/events/limited')
  }

  // Weekly Event
  async getWeeklyEvent(): Promise<WeeklyEvent> {
    return this.request('/api/events/weekly')
  }

  // Trades
  async getTrades(): Promise<TradesResult> {
    return this.request('/api/trades')
  }

  async createTrade(offeredPokemonId: string, requestedPokemonName?: string): Promise<{ trade_id: string }> {
    return this.request('/api/trades', {
      method: 'POST',
      body: JSON.stringify({ offered_pokemon_id: offeredPokemonId, requested_pokemon_name: requestedPokemonName || '' }),
    })
  }

  async acceptTrade(tradeId: string, offeredPokemonId: string): Promise<void> {
    await this.request(`/api/trades/${tradeId}/accept`, {
      method: 'POST',
      body: JSON.stringify({ offered_pokemon_id: offeredPokemonId }),
    })
  }

  async cancelTrade(tradeId: string): Promise<void> {
    await this.request(`/api/trades/${tradeId}/cancel`, { method: 'POST' })
  }

  // Admin: Deployers
  async getDeployers(): Promise<{ deployer_ids: string[] }> {
    return this.request('/api/admin/deployers')
  }

  async addDeployer(userId: string): Promise<void> {
    await this.request('/api/admin/deployers', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    })
  }

  async removeDeployer(userId: string): Promise<void> {
    await this.request(`/api/admin/deployers/${userId}`, { method: 'DELETE' })
  }

  // Trainer Card
  async getTrainerCard(userId?: string): Promise<TrainerCard> {
    const id = userId || 'me'
    return this.request(`/api/trainers/${id}`)
  }
}

export const api = new ApiClient()
