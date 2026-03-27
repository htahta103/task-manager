import { createBrowserRouter } from 'react-router-dom'
import { AppShell } from './shell/AppShell'
import { TasksPage } from '../routes/tasks/TasksPage'
import { SettingsPage } from '../routes/settings/SettingsPage'
import { NotFoundPage } from '../routes/not-found/NotFoundPage'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppShell />,
    children: [
      { index: true, element: <TasksPage /> },
      { path: 'settings', element: <SettingsPage /> },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
])

