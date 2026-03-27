import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <section className="space-y-4">
      <h1 className="text-2xl font-semibold tracking-tight text-white">
        Not found
      </h1>
      <p className="text-sm text-[color:var(--muted)]">
        That page doesn’t exist.
      </p>
      <div>
        <Link
          to="/"
          className="inline-flex items-center rounded-xl bg-white/10 px-3 py-2 text-sm text-white hover:bg-white/15"
        >
          Back to tasks
        </Link>
      </div>
    </section>
  )
}

