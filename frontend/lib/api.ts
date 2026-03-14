import type { AuthTokens, Book, BookListItem, CaughtPokemon, Memo, Pokemon, Tag } from './types'

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
}

export const api = new ApiClient()
