import { useState } from 'react';
import axios from 'axios';

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await axios.post('/api/v1/login', { username, password });
      const { token, role } = res.data;
      localStorage.setItem('edora_token', token);
      if (role) localStorage.setItem('edora_role', role);
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      onLogin && onLogin(role);
    } catch (err) {
      const msg = err?.response?.data?.error || 'Login Gagal! Periksa username dan password.';
      alert(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-900 flex items-center justify-center p-6">
      <div className="bg-white/10 backdrop-blur-xl border border-white/20 p-10 rounded-[2.5rem] w-full max-w-md shadow-2xl">
        <div className="text-center mb-10">
          <h1 className="text-4xl font-black text-white tracking-tighter italic">EDORA</h1>
          <p className="text-slate-400 text-sm font-medium uppercase tracking-widest">Admin Portal</p>
        </div>
        <form onSubmit={handleLogin} className="space-y-6">
          <input 
            type="text" placeholder="Username" 
            className="w-full p-4 bg-white/5 border border-white/10 rounded-2xl text-white outline-none focus:ring-2 focus:ring-blue-500"
            value={username} onChange={(e) => setUsername(e.target.value)}
          />
          <input 
            type="password" placeholder="Password" 
            className="w-full p-4 bg-white/5 border border-white/10 rounded-2xl text-white outline-none focus:ring-2 focus:ring-blue-500"
            value={password} onChange={(e) => setPassword(e.target.value)}
          />
          <button 
            type="submit" disabled={loading}
            className="w-full py-4 bg-blue-600 hover:bg-blue-700 text-white rounded-2xl font-black transition active:scale-95 disabled:opacity-50"
          >
            {loading ? 'AUTHENTICATING...' : 'ACCESS DASHBOARD'}
          </button>
        </form>
      </div>
    </div>
  );
}

export default Login;