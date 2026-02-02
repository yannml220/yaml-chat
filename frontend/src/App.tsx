import { Link, Outlet } from '@tanstack/react-router'
import './App.css'
import { AuthProvider } from './components/providers/AuthProvider'
import { ThemeProvider } from '@emotion/react'
import { appTheme } from './style/theme'

function App() {
  return (
    <ThemeProvider theme={appTheme}>
      <AuthProvider>
        <div className="p-2 flex gap-2">
          <Link to="/" className="[&.active]:font-bold">
            Home
          </Link>
          <Link to="/signin" className="[&.active]:font-bold">
            Signin
          </Link>
        </div>
        <hr />
        <Outlet />
      </AuthProvider>
    </ThemeProvider>
  )
}

export default App
