import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { api } from '../api'

const API_BASE = 'http://localhost:8080'

function mockFetchResponse(body: unknown, status = 200) {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  } as Response
}

describe('ApiClient', () => {
  let fetchSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    localStorage.clear()
    fetchSpy = vi.fn()
    vi.stubGlobal('fetch', fetchSpy)
  })

  afterEach(() => {
    vi.restoreAllMocks()
    vi.unstubAllGlobals()
  })

  describe('request adds Authorization header', () => {
    it('includes Bearer token when access_token exists', async () => {
      localStorage.setItem('access_token', 'my-token')
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({ books: [], total_count: 0, page: 1 }))

      await api.getBooks()

      expect(fetchSpy).toHaveBeenCalledWith(
        `${API_BASE}/api/books`,
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer my-token',
          }),
        }),
      )
    })

    it('does not include Authorization header when no token', async () => {
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({ url: 'https://google.com/auth' }))

      await api.getGoogleAuthURL()

      const headers = fetchSpy.mock.calls[0][1].headers
      expect(headers.Authorization).toBeUndefined()
    })
  })

  describe('getBooks', () => {
    it('sends GET request to /api/books', async () => {
      localStorage.setItem('access_token', 'token-1')
      const responseBody = { books: [{ id: '1', title: 'Test Book', author: 'Author', status: 'unread', tags: [] }], total_count: 1, page: 1 }
      fetchSpy.mockResolvedValueOnce(mockFetchResponse(responseBody))

      const result = await api.getBooks()

      expect(fetchSpy).toHaveBeenCalledWith(
        `${API_BASE}/api/books`,
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            Authorization: 'Bearer token-1',
          }),
        }),
      )
      expect(result).toEqual(responseBody)
    })

    it('sends query params when provided', async () => {
      localStorage.setItem('access_token', 'token-1')
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({ books: [], total_count: 0, page: 2 }))

      await api.getBooks({ status: 'reading', page: 2 })

      const url = fetchSpy.mock.calls[0][0] as string
      expect(url).toContain('status=reading')
      expect(url).toContain('page=2')
    })
  })

  describe('createBook', () => {
    it('sends POST request with correct body', async () => {
      localStorage.setItem('access_token', 'token-1')
      const bookData = { title: 'New Book', author: 'New Author', tag_ids: ['tag-1'] }
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({ BookID: 'book-123' }))

      const result = await api.createBook(bookData)

      expect(fetchSpy).toHaveBeenCalledWith(
        `${API_BASE}/api/books`,
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(bookData),
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            Authorization: 'Bearer token-1',
          }),
        }),
      )
      expect(result).toEqual({ BookID: 'book-123' })
    })
  })

  describe('401 triggers token refresh flow', () => {
    it('retries request after successful token refresh', async () => {
      localStorage.setItem('access_token', 'expired-token')
      localStorage.setItem('refresh_token', 'valid-refresh')

      // First call returns 401
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({}, 401))
      // Refresh call succeeds
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({ access_token: 'new-token' }))
      // Retry call succeeds
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({ books: [], total_count: 0, page: 1 }))

      const result = await api.getBooks()

      // Should have made 3 fetch calls: original, refresh, retry
      expect(fetchSpy).toHaveBeenCalledTimes(3)

      // Second call should be the refresh request
      expect(fetchSpy.mock.calls[1][0]).toBe(`${API_BASE}/api/auth/refresh`)
      expect(fetchSpy.mock.calls[1][1]).toEqual(
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ refresh_token: 'valid-refresh' }),
        }),
      )

      // New token should be saved
      expect(localStorage.getItem('access_token')).toBe('new-token')

      // Retry should use new token
      expect(fetchSpy.mock.calls[2][1].headers.Authorization).toBe('Bearer new-token')

      expect(result).toEqual({ books: [], total_count: 0, page: 1 })
    })

    it('redirects to /login when refresh fails', async () => {
      localStorage.setItem('access_token', 'expired-token')
      localStorage.setItem('refresh_token', 'invalid-refresh')

      // Mock window.location
      const originalLocation = window.location
      Object.defineProperty(window, 'location', {
        value: { ...originalLocation, href: '' },
        writable: true,
        configurable: true,
      })
      const hrefSetter = vi.fn()
      Object.defineProperty(window.location, 'href', {
        set: hrefSetter,
        get: () => '',
        configurable: true,
      })

      // First call returns 401
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({}, 401))
      // Refresh call fails
      fetchSpy.mockResolvedValueOnce(mockFetchResponse({}, 401))

      await expect(api.getBooks()).rejects.toThrow('Unauthorized')

      expect(hrefSetter).toHaveBeenCalledWith('/login')

      Object.defineProperty(window, 'location', {
        value: originalLocation,
        writable: true,
        configurable: true,
      })
    })

    it('throws Unauthorized when no refresh token available', async () => {
      localStorage.setItem('access_token', 'expired-token')
      // No refresh token set

      const originalLocation = window.location
      Object.defineProperty(window, 'location', {
        value: { ...originalLocation, href: '' },
        writable: true,
        configurable: true,
      })
      const hrefSetter = vi.fn()
      Object.defineProperty(window.location, 'href', {
        set: hrefSetter,
        get: () => '',
        configurable: true,
      })

      fetchSpy.mockResolvedValueOnce(mockFetchResponse({}, 401))

      await expect(api.getBooks()).rejects.toThrow('Unauthorized')

      Object.defineProperty(window, 'location', {
        value: originalLocation,
        writable: true,
        configurable: true,
      })
    })
  })

  describe('error handling', () => {
    it('throws error with message from response body', async () => {
      localStorage.setItem('access_token', 'token')
      fetchSpy.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: () => Promise.resolve({ error: 'Bad request data' }),
      } as unknown as Response)

      await expect(api.getBooks()).rejects.toThrow('Bad request data')
    })

    it('throws generic error when no error message in body', async () => {
      localStorage.setItem('access_token', 'token')
      fetchSpy.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () => Promise.resolve({}),
      } as unknown as Response)

      await expect(api.getBooks()).rejects.toThrow('API error: 500')
    })
  })
})
