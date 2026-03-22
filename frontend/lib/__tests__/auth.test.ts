import { describe, it, expect, beforeEach, vi } from 'vitest'
import { getUser, isLoggedIn, saveAuth, logout } from '../auth'
import type { User } from '../types'

const mockUser: User = {
  id: 'user-1',
  name: 'Ash Ketchum',
  email: 'ash@pokemon.com',
  role: 'trainer',
  github_username: 'ash-ketchum',
  avatar_url: 'https://example.com/ash.png',
}

describe('auth', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
  })

  describe('getUser', () => {
    it('returns null when no user in localStorage', () => {
      expect(getUser()).toBeNull()
    })

    it('returns parsed user when user exists in localStorage', () => {
      localStorage.setItem('user', JSON.stringify(mockUser))
      expect(getUser()).toEqual(mockUser)
    })

    it('returns user without optional fields', () => {
      const minimalUser: User = { id: '1', name: 'Misty', email: 'misty@pokemon.com', role: 'trainer' }
      localStorage.setItem('user', JSON.stringify(minimalUser))
      expect(getUser()).toEqual(minimalUser)
    })
  })

  describe('isLoggedIn', () => {
    it('returns false when no access_token', () => {
      expect(isLoggedIn()).toBe(false)
    })

    it('returns true when access_token exists', () => {
      localStorage.setItem('access_token', 'some-token')
      expect(isLoggedIn()).toBe(true)
    })
  })

  describe('saveAuth', () => {
    it('saves access_token, refresh_token, and user to localStorage', () => {
      saveAuth('access-123', 'refresh-456', mockUser)

      expect(localStorage.getItem('access_token')).toBe('access-123')
      expect(localStorage.getItem('refresh_token')).toBe('refresh-456')
      expect(JSON.parse(localStorage.getItem('user')!)).toEqual(mockUser)
    })
  })

  describe('logout', () => {
    it('removes all auth data from localStorage', () => {
      saveAuth('access-123', 'refresh-456', mockUser)

      // Mock window.location.href to prevent navigation error
      const locationSpy = vi.spyOn(window, 'location', 'get').mockReturnValue({
        ...window.location,
        href: '',
        set href(_url: string) {},
      } as Location)

      // Use Object.defineProperty to make href settable
      const hrefSetter = vi.fn()
      Object.defineProperty(window.location, 'href', {
        set: hrefSetter,
        configurable: true,
      })

      logout()

      expect(localStorage.getItem('access_token')).toBeNull()
      expect(localStorage.getItem('refresh_token')).toBeNull()
      expect(localStorage.getItem('user')).toBeNull()

      locationSpy.mockRestore()
    })

    it('redirects to /login', () => {
      saveAuth('access-123', 'refresh-456', mockUser)

      const hrefSetter = vi.fn()
      const originalLocation = window.location

      Object.defineProperty(window, 'location', {
        value: { ...originalLocation, href: '' },
        writable: true,
        configurable: true,
      })

      Object.defineProperty(window.location, 'href', {
        set: hrefSetter,
        get: () => '',
        configurable: true,
      })

      logout()

      expect(hrefSetter).toHaveBeenCalledWith('/login')

      Object.defineProperty(window, 'location', {
        value: originalLocation,
        writable: true,
        configurable: true,
      })
    })
  })
})
