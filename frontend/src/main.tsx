import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider, createRouter, createRootRoute, createRoute } from '@tanstack/react-router'
import './index.css'
import App from './App.tsx'
import Signin from './pages/Signin.tsx'
import { useAuth } from './contexts/AuthContext.ts'
import { Outlet, Navigate } from '@tanstack/react-router'

const rootRoute = createRootRoute({
  component: App,
})

const Authenticated = () => {
  const { isLoggedIn, isLoading } = useAuth();

  if (isLoading) return <div>Loading...</div>;
  if (!isLoggedIn) return <Navigate to="/signin" replace />;

  return <Outlet />;
}

const authenticatedRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: '_authenticated',
  component: Authenticated,
})

const indexRoute = createRoute({
  getParentRoute: () => authenticatedRoute,
  path: '/',
  component: () => <Home />,
})

const chatRoute = createRoute({
  getParentRoute: () => authenticatedRoute,
  path: 'chat/$conversationId',
  component: () => <Home />,
})

const signinRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/signin',
  component: () => <Signin />,
})

const routeTree = rootRoute.addChildren([
  authenticatedRoute.addChildren([indexRoute, chatRoute]),
  signinRoute,
])

const router = createRouter({
  routeTree,
  defaultPreload: 'intent',
})

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Home from './pages/Home.tsx'

const queryClient = new QueryClient()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
