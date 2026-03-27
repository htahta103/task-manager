import { useEffect, useState } from 'react'
import { useTasksStore } from '../../state/tasksStore'

export function TasksPage() {
  const { tasks, error, isLoading, refresh, add } = useTasksStore()
  const [title, setTitle] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    void refresh()
  }, [refresh])

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    const trimmed = title.trim()
    if (!trimmed) return
    setIsSubmitting(true)
    try {
      await add(trimmed)
      setTitle('')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <section className="space-y-6">
      <header className="space-y-2">
        <h1 className="text-2xl font-semibold tracking-tight text-white">
          Tasks
        </h1>
        <p className="text-sm text-[color:var(--muted)]">
          Add a task and it will appear immediately in the list.
        </p>
      </header>

      <form
        onSubmit={onSubmit}
        className="flex flex-col gap-3 rounded-[var(--radius-lg)] border border-[var(--border)] bg-white/5 p-4 sm:flex-row sm:items-end"
      >
        <label className="flex-1 space-y-1">
          <span className="text-xs font-medium text-white/70">Title</span>
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="e.g. Buy milk"
            className="w-full rounded-md border border-white/10 bg-black/20 px-3 py-2 text-sm text-white placeholder:text-white/40 focus:outline-none focus:ring-2 focus:ring-white/20"
          />
        </label>
        <button
          type="submit"
          disabled={isSubmitting || title.trim().length === 0}
          className="rounded-md bg-white/90 px-4 py-2 text-sm font-medium text-black disabled:opacity-50"
        >
          {isSubmitting ? 'Adding…' : 'Add task'}
        </button>
      </form>

      <div className="rounded-[var(--radius-lg)] border border-[var(--border)] bg-white/5 p-4">
        {isLoading ? (
          <div className="text-sm text-white/70">Loading…</div>
        ) : error ? (
          <div className="text-sm text-[color:var(--danger)]">
            {error.message}
          </div>
        ) : (
          <ul className="divide-y divide-white/10">
            {tasks.map((t) => (
              <li key={t.id} className="flex items-center justify-between py-3">
                <div className="space-y-0.5">
                  <div className="text-sm font-medium text-white/90">
                    {t.title}
                  </div>
                  <div className="text-xs text-[color:var(--muted)]">
                    Created {new Date(t.created_at).toLocaleString()}
                  </div>
                </div>
              </li>
            ))}
            {tasks.length === 0 ? (
              <li className="py-6 text-sm text-white/60">No tasks yet.</li>
            ) : null}
          </ul>
        )}
      </div>
    </section>
  )
}

