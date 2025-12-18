import { useEffect, useState } from 'react'
import axios from 'axios'

function App() {
  // --- STATE DATA UTAMA ---
  const [stats, setStats] = useState(null)
  const [patients, setPatients] = useState([])
  const [loading, setLoading] = useState(true)
  
  // --- STATE MODAL PASIEN (INPUT/EDIT) ---
  const [showModal, setShowModal] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false) 
  const [editId, setEditId] = useState(null)          
  const [formData, setFormData] = useState({ name: '', nik: '', gender: 'M', birth_date: '', address: '' })

  // --- STATE REKAM MEDIS (IOT SYNC) ---
  const [selectedPatient, setSelectedPatient] = useState(null) 
  const [records, setRecords] = useState([]) 
  const [isSyncing, setIsSyncing] = useState(false) // Indikator sensor aktif

  // --- STATE UI (TOAST & DELETE) ---
  const [toast, setToast] = useState({ show: false, message: '', type: 'success' })
  const [deleteModal, setDeleteModal] = useState({ show: false, id: null, name: '' })
  const [isSubmitting, setIsSubmitting] = useState(false)

  // --- HELPER: TAMPILKAN TOAST ---
  const showToast = (message, type = 'success') => {
    setToast({ show: true, message, type })
    setTimeout(() => setToast((prev) => ({ ...prev, show: false })), 3000)
  }

  // --- 1. AMBIL DATA AWAL (DASHBOARD & LIST) ---
  const refreshData = () => {
    const fetchStats = axios.get('/api/v1/dashboard/stats')
    const fetchPatients = axios.get('/api/v1/patients')

    Promise.all([fetchStats, fetchPatients])
      .then(([statsRes, patientsRes]) => {
        setStats(statsRes.data)
        setPatients(patientsRes.data || []) 
        setLoading(false)
      })
      .catch(() => {
        setLoading(false)
        showToast("Gagal sinkronisasi server", "error")
      })
  }

  useEffect(() => {
    refreshData()
  }, [])

  // --- 2. LOGIKA AUTO-SYNC IOT (POLLING SETIAP 5 DETIK) ---
  useEffect(() => {
    let interval;
    if (selectedPatient) {
      setIsSyncing(true);
      // Fungsi untuk tarik data terbaru
      const fetchLatestRecords = () => {
        axios.get(`/api/v1/patients/${selectedPatient.id}/scan`)
          .then(res => {
            setRecords(res.data || []);
            // Update stats dashboard juga siapa tahu ada status Osteoporosis baru
            axios.get('/api/v1/dashboard/stats').then(s => setStats(s.data));
          })
          .catch(err => console.error("Sync error:", err));
      };

      fetchLatestRecords(); // Jalankan sekali saat klik
      interval = setInterval(fetchLatestRecords, 5000); // Lalu cek tiap 5 detik
    } else {
      setIsSyncing(false);
    }

    return () => clearInterval(interval); // Bersihkan saat ganti pasien atau tutup
  }, [selectedPatient]);

  // --- 3. CRUD PASIEN ---
  const handleInputChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const openModalAdd = () => {
    setIsEditMode(false)
    setEditId(null)
    setFormData({ name: '', nik: '', gender: 'M', birth_date: '', address: '' })
    setShowModal(true)
  }

  const openModalEdit = (p) => {
    setIsEditMode(true)
    setEditId(p.id)
    setFormData({ ...p, birth_date: p.birth_date ? p.birth_date.split('T')[0] : '' })
    setShowModal(true)
  }

  const handleSubmitPatient = async (e) => {
    e.preventDefault()
    setIsSubmitting(true)
    try {
      if (isEditMode) {
        await axios.put(`/api/v1/patients/${editId}`, formData)
        showToast("✏️ Data pasien diperbarui!", "success")
      } else {
        await axios.post('/api/v1/patients', formData)
        showToast("✅ Pasien berhasil didaftarkan!", "success")
      }
      setShowModal(false)
      refreshData()
    } catch (error) {
      showToast("Gagal menyimpan data", "error")
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleDeleteExecute = async () => {
    try {
      await axios.delete(`/api/v1/patients/${deleteModal.id}`)
      showToast("🗑️ Data berhasil dihapus.", "success")
      setDeleteModal({ show: false })
      setSelectedPatient(null)
      refreshData()
    } catch (err) { showToast("Gagal menghapus", "error") }
  }

  if (loading) return <div className="flex h-screen items-center justify-center font-bold text-gray-400 animate-pulse">Inisialisasi Edora IoT...</div>

  return (
    <div className="min-h-screen bg-slate-50 p-4 md:p-8 font-sans">
      
      {/* TOAST NOTIFICATION */}
      <div className={`fixed top-5 right-5 z-[100] transition-all duration-500 transform ${toast.show ? 'translate-x-0 opacity-100' : 'translate-x-10 opacity-0 pointer-events-none'}`}>
        <div className={`px-6 py-4 rounded-2xl shadow-2xl border-l-4 bg-white ${toast.type === 'success' ? 'border-green-500 text-green-800' : 'border-red-500 text-red-800'}`}>
          <p className="font-bold text-sm">{toast.message}</p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        {/* KOLOM KIRI: DASHBOARD & TABEL */}
        <div className="lg:col-span-2 space-y-8">
          
          <header className="flex flex-col md:flex-row justify-between items-center gap-4">
            <div>
              <h1 className="text-3xl font-black text-slate-800 tracking-tight italic">EDORA <span className="text-blue-600 not-italic">v2.0</span></h1>
              <p className="text-slate-500 font-medium">Bone Health IoT Ecosystem</p>
            </div>
            <div className="flex gap-2">
              <button onClick={openModalAdd} className="bg-slate-900 hover:bg-black text-white px-6 py-3 rounded-2xl font-bold text-sm shadow-xl transition active:scale-95 flex items-center gap-2">
                <span>+</span> Daftarkan Pasien
              </button>
            </div>
          </header>

          {/* STATS */}
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-white p-6 rounded-[2rem] shadow-sm border border-slate-100">
              <h2 className="text-slate-400 text-[10px] font-black uppercase tracking-widest mb-1">Total Monitoring Hari Ini</h2>
              <p className="text-4xl font-black text-slate-800">{stats?.totalPatientsToday || 0}</p>
            </div>
            <div className="bg-white p-6 rounded-[2rem] shadow-sm border border-slate-100">
              <h2 className="text-slate-400 text-[10px] font-black uppercase tracking-widest mb-1">Terdeteksi Osteoporosis</h2>
              <p className="text-4xl font-black text-red-500">{stats?.osteoporosisCases || 0}</p>
            </div>
          </div>

          {/* TABEL PASIEN */}
          <div className="bg-white rounded-[2rem] shadow-sm border border-slate-200 overflow-hidden">
            <div className="p-6 bg-slate-50/50 border-b border-slate-100 flex justify-between items-center">
              <h2 className="font-bold text-slate-700">📋 Database Pasien</h2>
              <span className="text-[10px] font-bold text-slate-400">TOTAL: {patients.length}</span>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead className="bg-slate-50/30 text-slate-400 text-[9px] font-black uppercase tracking-[0.2em]">
                  <tr>
                    <th className="p-5">Informasi Pasien</th>
                    <th className="p-5">Gender</th>
                    <th className="p-5 text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-50">
                  {patients.map(p => (
                    <tr 
                      key={p.id} 
                      onClick={() => setSelectedPatient(p)}
                      className={`hover:bg-blue-50/40 transition cursor-pointer group ${selectedPatient?.id === p.id ? 'bg-blue-50' : ''}`}
                    >
                      <td className="p-5">
                        <div className="font-bold text-slate-700 group-hover:text-blue-600 transition">{p.name}</div>
                        <div className="text-[10px] text-slate-400 font-mono">NIK: {p.nik}</div>
                      </td>
                      <td className="p-5">
                        <span className={`px-3 py-1 rounded-lg text-[9px] font-black ${p.gender === 'M' ? 'bg-blue-100 text-blue-700' : 'bg-pink-100 text-pink-700'}`}>
                          {p.gender === 'M' ? 'MALE' : 'FEMALE'}
                        </span>
                      </td>
                      <td className="p-5 flex justify-center gap-3">
                        <button onClick={(e) => { e.stopPropagation(); openModalEdit(p); }} className="opacity-40 hover:opacity-100 transition">✏️</button>
                        <button onClick={(e) => { e.stopPropagation(); setDeleteModal({ show: true, id: p.id, name: p.name }); }} className="opacity-40 hover:opacity-100 transition hover:text-red-500">🗑️</button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        {/* KOLOM KANAN: PANEL MONITORING IOT */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-[2.5rem] shadow-2xl border border-slate-200 h-full min-h-[600px] sticky top-8 flex flex-col overflow-hidden">
            {!selectedPatient ? (
              <div className="flex-1 flex flex-col items-center justify-center p-12 text-center space-y-4">
                <div className="w-24 h-24 bg-slate-50 rounded-full flex items-center justify-center text-4xl animate-bounce">📡</div>
                <h3 className="font-black text-slate-800">Menunggu Koneksi...</h3>
                <p className="text-xs text-slate-400 leading-relaxed">Pilih pasien dari daftar untuk mulai menerima data dari sensor IoT Edora.</p>
              </div>
            ) : (
              <>
                <div className="p-8 border-b bg-slate-900 text-white relative overflow-hidden">
                  <div className="relative z-10">
                    <div className="flex items-center gap-2 mb-4">
                      <div className="h-2 w-2 rounded-full bg-emerald-400 animate-ping"></div>
                      <span className="text-[9px] font-black tracking-widest text-emerald-400 uppercase">Live Monitoring</span>
                    </div>
                    <h3 className="text-2xl font-black italic">{selectedPatient.name}</h3>
                    <p className="text-[10px] opacity-60 mt-1 font-mono">{selectedPatient.id}</p>
                  </div>
                  <div className="absolute -right-4 -bottom-4 text-7xl opacity-10">🦴</div>
                </div>

                <div className="flex-1 p-6 overflow-y-auto space-y-4 bg-slate-50/30">
                  <div className="flex justify-between items-center px-2">
                    <h4 className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Aktivitas Sensor</h4>
                    {isSyncing && <span className="text-[9px] font-bold text-blue-500 animate-pulse">SINKRONISASI...</span>}
                  </div>
                  
                  {records.length === 0 ? (
                    <div className="text-center py-20 border-2 border-dashed border-slate-200 rounded-[2rem] text-slate-300 text-xs italic font-medium p-4">
                      Belum ada data masuk.<br/>Silakan aktifkan alat scan pada pasien.
                    </div>
                  ) : (
                    records.map(r => (
                      <div key={r.id} className="p-6 rounded-[2rem] bg-white shadow-sm border border-slate-100 hover:shadow-md transition group">
                        <div className="flex justify-between items-start mb-3">
                          <span className="text-[10px] font-bold text-slate-400">
                            {new Date(r.scan_date).toLocaleTimeString('id-ID')} - {new Date(r.scan_date).toLocaleDateString('id-ID')}
                          </span>
                          <span className={`text-[9px] font-black px-3 py-1 rounded-full ${r.diagnosis === 'Osteoporosis' ? 'bg-red-500 text-white shadow-lg shadow-red-100' : r.diagnosis === 'Osteopenia' ? 'bg-yellow-400 text-white shadow-lg shadow-yellow-100' : 'bg-emerald-500 text-white shadow-lg shadow-emerald-100'}`}>
                            {r.diagnosis.toUpperCase()}
                          </span>
                        </div>
                        <div className="text-4xl font-black text-slate-800 tracking-tighter">
                          {r.t_score} <span className="text-sm font-bold text-slate-300 ml-1">T-Score</span>
                        </div>
                        {r.notes && (
                          <div className="mt-4 p-3 bg-slate-50 rounded-xl text-[11px] text-slate-500 italic border-l-2 border-slate-200">
                            "{r.notes}"
                          </div>
                        )}
                      </div>
                    ))
                  )}
                </div>
              </>
            )}
          </div>
        </div>
      </div>

      {/* --- MODAL FORM PASIEN --- */}
      {showModal && (
        <div className="fixed inset-0 bg-slate-900/70 z-[110] backdrop-blur-md flex items-center justify-center p-4">
          <div className="bg-white rounded-[2.5rem] p-10 w-full max-w-md shadow-2xl animate-in zoom-in duration-300">
            <h2 className="text-2xl font-black text-slate-800 mb-8 tracking-tight">
              {isEditMode ? 'Edit Profil Pasien' : 'Registrasi Pasien Baru'}
            </h2>
            <form onSubmit={handleSubmitPatient} className="space-y-5">
              <div className="space-y-1">
                <label className="text-[10px] font-black text-slate-400 uppercase ml-2">NIK Identitas</label>
                <input required name="nik" value={formData.nik} onChange={handleInputChange} placeholder="Masukkan 16 digit NIK" className="w-full p-4 bg-slate-50 border-none rounded-2xl outline-none focus:ring-2 focus:ring-blue-500 transition font-medium" />
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black text-slate-400 uppercase ml-2">Nama Lengkap</label>
                <input required name="name" value={formData.name} onChange={handleInputChange} placeholder="Nama sesuai KTP" className="w-full p-4 bg-slate-50 border-none rounded-2xl outline-none focus:ring-2 focus:ring-blue-500 transition font-medium" />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label className="text-[10px] font-black text-slate-400 uppercase ml-2">Tgl Lahir</label>
                  <input required type="date" name="birth_date" value={formData.birth_date} onChange={handleInputChange} className="w-full p-4 bg-slate-50 border-none rounded-2xl outline-none focus:ring-2 focus:ring-blue-500 transition" />
                </div>
                <div className="space-y-1">
                  <label className="text-[10px] font-black text-slate-400 uppercase ml-2">Gender</label>
                  <select name="gender" value={formData.gender} onChange={handleInputChange} className="w-full p-4 bg-slate-50 border-none rounded-2xl outline-none bg-white focus:ring-2 focus:ring-blue-500 transition">
                    <option value="M">Pria</option>
                    <option value="F">Wanita</option>
                  </select>
                </div>
              </div>
              <div className="space-y-1">
                <label className="text-[10px] font-black text-slate-400 uppercase ml-2">Alamat Domisili</label>
                <textarea name="address" value={formData.address} onChange={handleInputChange} placeholder="Alamat lengkap..." className="w-full p-4 bg-slate-50 border-none rounded-2xl outline-none focus:ring-2 focus:ring-blue-500 transition h-24" />
              </div>
              <div className="flex gap-4 pt-6">
                <button type="button" onClick={() => setShowModal(false)} className="flex-1 py-4 font-bold text-slate-400 transition hover:text-slate-600">Batal</button>
                <button type="submit" disabled={isSubmitting} className="flex-[2] py-4 bg-slate-900 text-white rounded-[1.25rem] font-black shadow-2xl shadow-slate-200 hover:bg-black active:scale-95 transition">
                  {isSubmitting ? 'MEMPROSES...' : 'KONFIRMASI DATA'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- MODAL DELETE --- */}
      {deleteModal.show && (
        <div className="fixed inset-0 bg-slate-900/60 z-[130] flex items-center justify-center p-4">
          <div className="bg-white rounded-[2rem] p-8 w-full max-w-xs text-center shadow-2xl">
            <div className="w-16 h-16 bg-red-50 text-red-500 rounded-full flex items-center justify-center text-3xl mx-auto mb-4">⚠️</div>
            <h3 className="font-extrabold text-slate-800 text-lg italic uppercase tracking-tighter">Hapus Pasien?</h3>
            <p className="text-xs text-slate-400 mt-2 mb-8 font-medium">Seluruh data riwayat scan pasien <span className="text-slate-800 font-bold">{deleteModal.name}</span> akan hilang permanen.</p>
            <div className="flex gap-3">
              <button onClick={() => setDeleteModal({ show: false })} className="flex-1 py-3 text-slate-400 font-bold text-xs">Batal</button>
              <button onClick={handleDeleteExecute} className="flex-1 py-3 bg-red-500 text-white rounded-xl font-bold text-xs shadow-lg shadow-red-100">Ya, Hapus</button>
            </div>
          </div>
        </div>
      )}

    </div>
  )
}

export default App