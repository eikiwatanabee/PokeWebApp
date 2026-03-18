'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import type { LimitedEvent } from '@/lib/types'

export function LimitedEvents() {
  const [events, setEvents] = useState<LimitedEvent[]>([])

  useEffect(() => {
    api.getLimitedEvents()
      .then(res => setEvents(res.events || []))
      .catch(console.error)
  }, [])

  if (events.length === 0) return null

  return (
    <div className="space-y-3">
      {events.map(event => (
        <div
          key={event.id}
          className="bg-gradient-to-r from-amber-500 to-red-500 rounded-xl p-4 text-white shadow-lg"
        >
          <div className="flex items-center gap-3">
            <span className="text-3xl">{event.icon}</span>
            <div className="flex-1">
              <h3 className="font-bold">{event.title}</h3>
              <p className="text-sm text-white/80">{event.description}</p>
            </div>
            <div className="text-right">
              <p className="text-xs text-white/70">レアリティ</p>
              <p className="text-lg font-bold">x{event.rarity_boost}</p>
              <p className="text-xs text-white/70">残り {event.hours_left}h</p>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
