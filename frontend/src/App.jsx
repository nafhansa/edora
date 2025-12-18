import { useEffect, useState } from 'react'
import axios from 'axios'

function App() {
  // --- STATE DATA UTAMA ---
  const [stats, setStats] = useState(null)
  const [patients, setPatients] = useState([])
  const [loading, setLoading] = useState(true)
  
  // --- STATE MODAL FORM (INPUT/EDIT) ---
  const [showModal, setShowModal] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false) 
  const [editId, setEditId] = useState(null)          
  
  const [formData, setFormData] = useState({
    name: '',
    nik: '',
    gender: 'M',
    birth_date: '',
    address: ''
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  // --- STATE UI BARU (TOAST & DELETE MODAL) ---
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' }) // Type: success | error
  const [deleteModal, setDeleteModal] = useState({ show: false, id: null, name: '' })

  // --- HELPER: TAMPILKAN TOAST ---
  const showToast = (message, type = 'success') => {
    setToast({ show: true, message, type })
    // Hilang otomatis setelah 3 detik
    setTimeout(() => {
      setToast((prev) => ({ ...prev, show: false }))
    }, 3000)
  }

  // --- 1. FUNGSI AMBIL DATA ---
  const refreshData = () => {
    const fetchStats = axios.get('/api/v1/dashboard/stats')
    const fetchPatients = axios.get('/api/v1/patients')

    Promise.all([fetchStats, fetchPatients])
      .then(([statsRes, patientsRes]) => {
        setStats(statsRes.data)
        setPatients(patientsRes.data || []) 
        setLoading(false)
      })
      .catch(err => {
        console.error("Error fetching data:", err)
        setLoading(false)
        setPatients([]) 
        showToast("Gagal mengambil data dari server", "error")
      })
  }

  useEffect(() => {
    refreshData()
  }, [])

  // --- 2. HANDLE FORM INPUT ---
  const handleInputChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  // --- 3. BUKA MODAL ---
  const openModalAdd = () => {
    setIsEditMode(false)
    setEditId(null)
    setFormData({ name: '', nik: '', gender: 'M', birth_date: '', address: '' })
    setShowModal(true)
  }

  const openModalEdit = (patient) => {
    setIsEditMode(true)
    setEditId(patient.id)
    setFormData({
      name: patient.name,
      nik: patient.nik,
      gender: patient.gender,
      birth_date: patient.birth_date ? patient.birth_date.split('T')[0] : '',
      address: patient.address
    })
    setShowModal(true)
  }

  // --- 4. SIMPAN DATA (CREATE / UPDATE) ---
  const handleSubmit = async (e) => {
    e.preventDefault()
    setIsSubmitting(true)
    
    try {
      if (isEditMode) {
        await axios.put(`/api/v1/patients/${editId}`, formData)
        showToast("✏️ Data pasien berhasil diperbarui!", "success")
      } else {
        await axios.post('/api/v1/patients', formData)
        showToast("✅ Pasien baru berhasil ditambahkan!", "success")
      }

      setShowModal(false)
      refreshData()
      
    } catch (error) {
      console.error(error)
      
      // VALIDASI ERROR DARI BACKEND
      const errMsg = error.response?.data?.error || JSON.stringify(error.response?.data) || error.message;
      
      if (errMsg && (errMsg.includes("duplicate") || errMsg.includes("unique"))) {
        showToast("⚠️ NIK sudah terdaftar! Gunakan NIK lain.", "error")
      } else {
        showToast("❌ Terjadi kesalahan server.", "error")
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  // --- 5. LOGIKA HAPUS (DELETE) ---
  // Tahap 1: Buka Modal Konfirmasi
  const confirmDelete = (id, name) => {
    setDeleteModal({ show: true, id, name })
  }

  // Tahap 2: Eksekusi Hapus ke Backend
  const handleDeleteExecute = async () => {
    try {
      await axios.delete(`/api/v1/patients/${deleteModal.id}`)
      showToast("🗑️ Data berhasil dihapus permanen.", "success")
      setDeleteModal({ show: false, id: null, name: '' }) // Tutup modal
      refreshData()
    } catch (error) {
      console.error(error)
      showToast("❌ Gagal menghapus data.", "error")
    }
  }

  if (loading) return <div className="flex h-screen items-center justify-center bg-gray-100 font-bold text-gray-500">Loading Dashboard...</div>

  return (
    <div className="min-h-screen bg-gray-100 p-8 font-sans relative overflow-x-hidden">
      
      {/* --- KOMPONEN TOAST NOTIFICATION --- */}
      <div className={`fixed top-5 right-5 z-[100] transition-all duration-500 transform ${toast.show ? 'translate-x-0 opacity-100' : 'translate-x-10 opacity-0 pointer-events-none'}`}>
        <div className={`flex items-center gap-3 px-6 py-4 rounded-xl shadow-2xl border-l-4 ${toast.type === 'success' ? 'bg-white border-green-500 text-green-800' : 'bg-white border-red-500 text-red-800'}`}>
          <div className={`p-2 rounded-full ${toast.type === 'success' ? 'bg-green-100' : 'bg-red-100'}`}>
            {toast.type === 'success' ? '✔' : '✖'}
          </div>
          <div>
            <h4 className="font-bold text-sm">{toast.type === 'success' ? 'Berhasil' : 'Gagal'}</h4>
            <p className="text-sm font-medium opacity-90">{toast.message}</p>
          </div>
        </div>
      </div>

      <div className="max-w-6xl mx-auto">
        
        {/* HEADER */}
        <header className="mb-8 flex flex-col md:flex-row justify-between items-center gap-4">
          <div>
            <h1 className="text-3xl font-extrabold text-slate-800 tracking-tight">Edora Dashboard</h1>
            <p className="text-slate-500 font-medium">Monitoring Kesehatan Tulang Real-time</p>
          </div>
          <div className="flex gap-3">
             <button onClick={refreshData} className="px-4 py-2 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 text-sm font-bold text-gray-600 shadow-sm transition active:scale-95">
                🔄 Refresh
             </button>
             <div className="bg-emerald-100 px-4 py-2 rounded-xl text-sm font-bold text-emerald-700 flex items-center gap-2 shadow-sm border border-emerald-200">
               <span className="relative flex h-3 w-3">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
                </span>
               System Online
             </div>
          </div>
        </header>
        
        {/* STATS CARDS */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition">
            <div className="flex items-center gap-4 mb-2">
              <div className="p-3 bg-blue-100 text-blue-600 rounded-lg">📊</div>
              <h2 className="text-gray-500 text-sm uppercase font-bold tracking-wider">Total Pasien Hari Ini</h2>
            </div>
            <p className="text-5xl font-extrabold text-slate-800">{stats?.totalPatientsToday || 0}</p>
          </div>
          <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition">
            <div className="flex items-center gap-4 mb-2">
              <div className="p-3 bg-red-100 text-red-600 rounded-lg">⚠️</div>
              <h2 className="text-gray-500 text-sm uppercase font-bold tracking-wider">Risiko Osteoporosis</h2>
            </div>
            <p className="text-5xl font-extrabold text-red-500">
              {stats?.osteoporosisCases || 0}
            </p>
          </div>
        </div>

        {/* TABEL DATA */}
        <div className="bg-white rounded-2xl shadow-sm overflow-hidden border border-gray-200">
          <div className="p-6 border-b border-gray-100 bg-gray-50/50 flex flex-col sm:flex-row justify-between items-center gap-4">
            <h2 className="text-xl font-bold text-slate-800">📋 Daftar Pasien Terdaftar</h2>
            <button 
              onClick={openModalAdd} 
              className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-xl text-sm font-bold transition shadow-lg shadow-blue-200 active:scale-95 flex items-center gap-2"
            >
              <span>+</span> Input Pasien Baru
            </button>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead className="bg-slate-50 text-slate-500 uppercase text-xs font-bold tracking-wider border-b border-gray-100">
                <tr>
                  <th className="p-4 pl-6">Nama Pasien</th>
                  <th className="p-4">NIK</th>
                  <th className="p-4">Tgl Lahir</th>
                  <th className="p-4">Gender</th>
                  <th className="p-4">Alamat</th>
                  <th className="p-4 text-center">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {patients?.map((p) => (
                  <tr key={p.id} className="hover:bg-slate-50 transition group">
                    <td className="p-4 pl-6">
                      <div className="font-bold text-slate-700">{p.name}</div>
                    </td>
                    <td className="p-4">
                      <span className="bg-gray-100 text-gray-600 px-2 py-1 rounded text-xs font-mono font-medium">{p.nik}</span>
                    </td>
                    <td className="p-4 text-gray-500 text-sm font-medium">
                      {p.birth_date ? new Date(p.birth_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '-'}
                    </td>
                    <td className="p-4">
                      <span className={`px-3 py-1 rounded-full text-[10px] font-extrabold tracking-wide ${p.gender === 'M' ? 'bg-blue-50 text-blue-600 border border-blue-100' : 'bg-pink-50 text-pink-600 border border-pink-100'}`}>
                        {p.gender === 'M' ? 'PRIA' : 'WANITA'}
                      </span>
                    </td>
                    <td className="p-4 text-gray-500 text-sm max-w-[200px] truncate" title={p.address}>{p.address}</td>
                    
                    <td className="p-4 flex justify-center gap-2">
                      <button 
                        onClick={() => openModalEdit(p)}
                        className="p-2 bg-white border border-gray-200 text-gray-600 rounded-lg hover:bg-yellow-50 hover:text-yellow-600 hover:border-yellow-200 transition shadow-sm"
                        title="Edit Data"
                      >
                        ✏️
                      </button>
                      <button 
                        onClick={() => confirmDelete(p.id, p.name)}
                        className="p-2 bg-white border border-gray-200 text-gray-600 rounded-lg hover:bg-red-50 hover:text-red-600 hover:border-red-200 transition shadow-sm"
                        title="Hapus Data"
                      >
                        🗑️
                      </button>
                    </td>
                  </tr>
                ))}
                
                {(!patients || patients.length === 0) && (
                  <tr>
                    <td colSpan="6" className="p-12 text-center">
                      <div className="flex flex-col items-center justify-center text-gray-400">
                        <span className="text-4xl mb-2">📂</span>
                        <p className="font-medium">Belum ada data pasien.</p>
                      </div>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* --- MODAL FORM INPUT/EDIT --- */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-900/60 flex items-center justify-center z-50 backdrop-blur-sm p-4 animate-in fade-in duration-200">
          <div className="bg-white w-full max-w-lg rounded-2xl shadow-2xl p-8 transform transition-all scale-100 border border-gray-200">
            <div className="flex justify-between items-center mb-6">
               <h2 className="text-2xl font-bold text-slate-800">
                {isEditMode ? '✏️ Edit Informasi' : '➕ Tambah Pasien'}
              </h2>
              <button onClick={() => setShowModal(false)} className="text-gray-400 hover:text-gray-600 text-2xl">&times;</button>
            </div>
            
            <form onSubmit={handleSubmit} className="space-y-5">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                <div>
                  <label className="block text-xs font-bold text-gray-500 uppercase mb-2">NIK (Identitas)</label>
                  <input 
                    required
                    type="text" 
                    name="nik"
                    value={formData.nik}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:outline-none transition"
                    placeholder="3201xxxxxx"
                  />
                </div>
                <div>
                   <label className="block text-xs font-bold text-gray-500 uppercase mb-2">Nama Lengkap</label>
                  <input 
                    required
                    type="text" 
                    name="name"
                    value={formData.name}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:outline-none transition"
                    placeholder="Sesuai KTP"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-5">
                <div>
                   <label className="block text-xs font-bold text-gray-500 uppercase mb-2">Tanggal Lahir</label>
                  <input 
                    required
                    type="date" 
                    name="birth_date"
                    value={formData.birth_date}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:outline-none transition text-gray-700"
                  />
                </div>
                <div>
                   <label className="block text-xs font-bold text-gray-500 uppercase mb-2">Gender</label>
                  <select 
                    name="gender"
                    value={formData.gender}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:outline-none bg-white transition"
                  >
                    <option value="M">Laki-laki</option>
                    <option value="F">Perempuan</option>
                  </select>
                </div>
              </div>

              <div>
                 <label className="block text-xs font-bold text-gray-500 uppercase mb-2">Alamat Domisili</label>
                <textarea 
                  name="address"
                  value={formData.address}
                  onChange={handleInputChange}
                  className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:outline-none transition"
                  rows="3"
                  placeholder="Jalan, No Rumah, RT/RW..."
                ></textarea>
              </div>

              <div className="flex justify-end gap-3 mt-8 pt-4 border-t border-gray-100">
                <button 
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="px-6 py-2.5 text-gray-600 hover:bg-gray-100 rounded-xl font-bold text-sm transition"
                  disabled={isSubmitting}
                >
                  Batal
                </button>
                <button 
                  type="submit"
                  disabled={isSubmitting}
                  className={`px-8 py-2.5 text-white rounded-xl font-bold text-sm shadow-lg shadow-blue-200 transition transform hover:-translate-y-1 disabled:opacity-50 disabled:transform-none ${isEditMode ? 'bg-yellow-500 hover:bg-yellow-600' : 'bg-blue-600 hover:bg-blue-700'}`}
                >
                  {isSubmitting ? 'Memproses...' : (isEditMode ? 'Simpan Perubahan' : 'Simpan Data')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- MODAL KONFIRMASI DELETE (BARU) --- */}
      {deleteModal.show && (
        <div className="fixed inset-0 bg-slate-900/60 flex items-center justify-center z-[60] backdrop-blur-sm p-4 animate-in fade-in duration-200">
          <div className="bg-white w-full max-w-sm rounded-2xl shadow-2xl p-6 text-center border border-gray-200">
            <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <span className="text-3xl">🗑️</span>
            </div>
            <h3 className="text-xl font-bold text-gray-800 mb-2">Hapus Data?</h3>
            <p className="text-gray-500 text-sm mb-6">
              Anda yakin ingin menghapus data pasien <br/>
              <span className="font-bold text-gray-800">"{deleteModal.name}"</span>?
              <br/>Tindakan ini tidak bisa dibatalkan.
            </p>
            <div className="flex gap-3 justify-center">
              <button 
                onClick={() => setDeleteModal({ show: false, id: null, name: '' })}
                className="px-5 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-xl font-bold text-sm transition"
              >
                Batal
              </button>
              <button 
                onClick={handleDeleteExecute}
                className="px-5 py-2.5 bg-red-500 hover:bg-red-600 text-white rounded-xl font-bold text-sm shadow-lg shadow-red-200 transition"
              >
                Ya, Hapus
              </button>
            </div>
          </div>
        </div>
      )}

    </div>
  )
}

export default App