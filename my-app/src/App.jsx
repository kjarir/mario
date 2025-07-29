import Index from './pages/index'
import NavbarComponent from './components/Navbar'
import { useLenis } from './hooks/useLenis'
import './output.css'

function App() {
  // Initialize Lenis smooth scrolling
  useLenis();

  return (
    <div className="min-h-screen">
      <NavbarComponent />
      <Index />
    </div>
  )
}

export default App
