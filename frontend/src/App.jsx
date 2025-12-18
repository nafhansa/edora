import { useEffect, useState } from 'react'
import axios from 'axios'

function App() {
  // --- STATE ---
  const [stats, setStats] = useState(null)
  const [patients, setPatients] = useState([])
  const [loading, setLoading] = useState(true)
  
  // State untuk Modal & Form
  const [showModal, setShowModal] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    nik: '',
    gender: 'M', // DEFAULT 'M' (Sesuai constraint DB)
    birth_date: '',
    address: ''
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  // --- FUNGSI ---

  // Ambil data dari Backend
  const refreshData = () => {
    const fetchStats = axios.get('/api/v1/dashboard/stats')
    const fetchPatients = axios.get('/api/v1/patients')

    Promise.all([fetchStats, fetchPatients])
      .then(([statsRes, patientsRes]) => {
        setStats(statsRes.data)
        // FIX: Paksa jadi array kosong jika null agar tidak error .map
        setPatients(patientsRes.data || []) 
        setLoading(false)
      })
      .catch(err => {
        console.error("Error fetching data:", err)
        setLoading(false)
        setPatients([]) // Safety jika error
      })
  }

  // Jalan pas pertama kali web dibuka
  useEffect(() => {
    refreshData()
  }, [])

  // Handle ketikan di form
  const handleInputChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  // Handle tombol Simpan
  const handleSubmit = async (e) => {
    e.preventDefault()
    setIsSubmitting(true)
    try {
      // Kirim data ke Backend
      await axios.post('/api/v1/patients', formData)
      
      alert("✅ Pasien berhasil ditambahkan!")
      setShowModal(false) // Tutup modal
      
      // FIX: Reset gender kembali ke 'M'
      setFormData({ name: '', nik: '', gender: 'M', birth_date: '', address: '' }) 
      
      refreshData() // Refresh tabel
    } catch (error) {
      alert("❌ Gagal menambah pasien. Cek console.")
      console.error(error)
    } finally {
      setIsSubmitting(false)
    }
  }

  if (loading) return <div className="p-10 text-center font-bold text-gray-500">Loading Dashboard...</div>

  return (
    <div className="min-h-screen bg-gray-100 p-8 font-sans relative">
      <div className="max-w-5xl mx-auto">
        
        {/* --- HEADER --- */}
        <header className="mb-8 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-blue-800">Edora Dashboard</h1>
            <p className="text-gray-500">Monitoring Kesehatan Tulang Real-time</p>
          </div>
          <div className="flex gap-3">
             <button onClick={refreshData} className="px-4 py-2 bg-white border rounded-lg hover:bg-gray-50 text-sm font-medium shadow-sm transition">
                🔄 Refresh Data
             </button>
             <div className="bg-green-100 px-4 py-2 rounded-lg text-sm font-bold text-green-700 flex items-center gap-2 shadow-sm">
               <span className="animate-pulse">●</span> Online
             </div>
          </div>
        </header>
        
        {/* --- STATS CARDS --- */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div className="bg-white p-6 rounded-xl shadow-sm border-l-4 border-blue-500">
            <h2 className="text-gray-500 text-sm uppercase font-bold tracking-wider">Total Scans Hari Ini</h2>
            <p className="text-5xl font-bold text-gray-800 mt-2">{stats?.totalPatientsToday || 0}</p>
          </div>
          <div className="bg-white p-6 rounded-xl shadow-sm border-l-4 border-red-500">
            <h2 className="text-gray-500 text-sm uppercase font-bold tracking-wider">Kasus Osteoporosis</h2>
            <p className="text-5xl font-bold text-red-600 mt-2">
              {stats?.osteoporosisCases || 0}
            </p>
          </div>
        </div>

        {/* --- TABEL PASIEN --- */}
        <div className="bg-white rounded-xl shadow-sm overflow-hidden border border-gray-200">
          <div className="p-6 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
            <h2 className="text-xl font-bold text-gray-800">Daftar Pasien</h2>
            <button 
              onClick={() => setShowModal(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-semibold transition shadow-md flex items-center gap-2"
            >
              <span>+</span> Input Manual
            </button>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead className="bg-gray-50 text-gray-600 uppercase text-xs tracking-wider border-b">
                <tr>
                  <th className="p-4">Nama Pasien</th>
                  <th className="p-4">NIK</th>
                  <th className="p-4">Tanggal Lahir</th>
                  <th className="p-4">Gender</th>
                  <th className="p-4">Alamat</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {/* FIX: Gunakan optional chaining (?.) agar tidak crash */}
                {patients?.map((p) => (
                  <tr key={p.id} className="hover:bg-blue-50 transition">
                    <td className="p-4 font-semibold text-gray-800">{p.name}</td>
                    <td className="p-4 text-gray-500 font-mono text-sm">{p.nik}</td>
                    <td className="p-4 text-gray-600 text-sm">
                      {p.birth_date ? p.birth_date.split('T')[0] : '-'}
                    </td>
                    <td className="p-4">
                      {/* FIX: Logika warna badge disesuaikan dengan M/F */}
                      <span className={`px-2 py-1 rounded-full text-xs font-bold ${p.gender === 'M' ? 'bg-blue-100 text-blue-700' : 'bg-pink-100 text-pink-700'}`}>
                        {p.gender === 'M' ? 'LAKI-LAKI' : 'PEREMPUAN'}
                      </span>
                    </td>
                    <td className="p-4 text-gray-500 text-sm">{p.address}</td>
                  </tr>
                ))}
                
                {/* Handling jika data kosong */}
                {(!patients || patients.length === 0) && (
                  <tr>
                    <td colSpan="5" className="p-8 text-center text-gray-400 italic">
                      Belum ada data pasien. Silakan input manual.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* --- MODAL / POPUP FORM --- */}
      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 backdrop-blur-sm p-4">
          <div className="bg-white w-full max-w-md rounded-2xl shadow-2xl p-6 transform transition-all scale-100 border border-gray-200">
            <h2 className="text-2xl font-bold mb-4 text-gray-800">Tambah Pasien Baru</h2>
            
            <form onSubmit={handleSubmit} className="space-y-4">
              {/* NIK */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">NIK</label>
                <input 
                  required
                  type="text" 
                  name="nik"
                  value={formData.nik}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
                  placeholder="Contoh: 3201..."
                />
              </div>

              {/* Nama */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap</label>
                <input 
                  required
                  type="text" 
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
                  placeholder="Nama Pasien"
                />
              </div>

              {/* Tgl Lahir & Gender */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Tgl Lahir</label>
                  <input 
                    required
                    type="date" 
                    name="birth_date"
                    value={formData.birth_date}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Gender</label>
                  <select 
                    name="gender"
                    value={formData.gender}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none bg-white"
                  >
                    {/* FIX: Value harus 'M' dan 'F' sesuai DB */}
                    <option value="M">Laki-laki</option>
                    <option value="F">Perempuan</option>
                  </select>
                </div>
              </div>

              {/* Alamat */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Alamat</label>
                <textarea 
                  name="address"
                  value={formData.address}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
                  rows="2"
                  placeholder="Alamat domisili..."
                ></textarea>
              </div>

              {/* Tombol Action */}
              <div className="flex justify-end gap-3 mt-6">
                <button 
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-lg font-medium transition"
                  disabled={isSubmitting}
                >
                  Batal
                </button>
                <button 
                  type="submit"
                  disabled={isSubmitting}
                  className="px-6 py-2 bg-blue-600 text-white rounded-lg font-bold hover:bg-blue-700 shadow-md transition disabled:opacity-50"
                >
                  {isSubmitting ? 'Menyimpan...' : 'Simpan Data'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  )
}

export default App