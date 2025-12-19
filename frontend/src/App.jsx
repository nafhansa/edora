import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import axios from 'axios';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [checkingAuth, setCheckingAuth] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('edora_token');
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      setIsAuthenticated(true);
    }
    setCheckingAuth(false);
  }, []);

  if (checkingAuth) return <div className="flex h-screen items-center justify-center font-bold text-gray-400">Loading Auth...</div>;

  return (
    <Router>
      <Routes>
        {/* Jika belum login, ke halaman Login */}
        <Route 
          path="/login" 
          element={!isAuthenticated ? <Login onLogin={() => setIsAuthenticated(true)} /> : <Navigate to="/" />} 
        />
        
        {/* Jika sudah login, ke Dashboard. Jika belum, tendang ke Login */}
        <Route 
          path="/" 
          element={isAuthenticated ? <Dashboard /> : <Navigate to="/login" />} 
        />
        
        {/* Tombol Logout bisa ditaruh di dalam Dashboard atau Layout */}
      </Routes>
    </Router>
  );
}

export default App;