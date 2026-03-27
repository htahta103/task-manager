import { useEffect } from 'react'
import { useTasksStore } from '../../state/tasksStore'

export function TasksPage() {
  const { tasks, error, isLoading, refresh } = useTasksStore()

  useEffect(() => {
    void refresh()
  }, [refresh])

  return (
    <section className="space-y-6">
      <header className="space-y-2">
        <h1 className="text-2xl font-semibold tracking-tight text-white">
          Tasks
        </h1>
        <p className="text-sm text-[color:var(--muted)]">
          This is a scaffold page. It currently loads tasks from the API and
          renders a simple list.
        </p>
      </header>

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
                    {t.status} • {t.priority}
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

