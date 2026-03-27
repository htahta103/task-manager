import { NavLink, Outlet } from 'react-router-dom'

export function AppShell() {
  return (
    <div className="min-h-screen">
      <header className="sticky top-0 z-10 border-b border-[var(--border)] bg-[color:var(--bg)]/60 backdrop-blur">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-gradient-to-br from-[color:var(--primary)] to-[color:var(--primary-2)] shadow-[var(--shadow)]" />
            <div className="leading-tight">
              <div className="text-sm font-semibold tracking-wide text-white/90">
                Task Manager
              </div>
              <div className="text-xs text-[color:var(--muted)]">v1 scaffold</div>
            </div>
          </div>

          <nav className="flex items-center gap-2 text-sm">
            <NavLink
              to="/"
              className={({ isActive }) =>
                [
                  'rounded-xl px-3 py-2 transition',
                  isActive
                    ? 'bg-white/10 text-white'
                    : 'text-white/70 hover:bg-white/5 hover:text-white',
                ].join(' ')
              }
              end
            >
              Tasks
            </NavLink>
            <NavLink
              to="/settings"
              className={({ isActive }) =>
                [
                  'rounded-xl px-3 py-2 transition',
                  isActive
                    ? 'bg-white/10 text-white'
                    : 'text-white/70 hover:bg-white/5 hover:text-white',
                ].join(' ')
              }
            >
              Settings
            </NavLink>
          </nav>
        </div>
      </header>

      <main className="mx-auto w-full max-w-5xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}

