"use client"
import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'

export default function Dashboard(){
  const [readings, setReadings] = useState<any[]>([])
  const router = useRouter()

  useEffect(()=>{
    const token = localStorage.getItem('token')
    if(!token){ router.push('/login'); return }
    fetch((process.env.NEXT_PUBLIC_API_URL||'http://localhost:8080') + '/api/v1/dashboard', {
      headers: { 'Authorization': `Bearer ${token}` }
    }).then(r=>{
      if(!r.ok) { router.push('/login'); return null }
      return r.json()
    }).then(j=>{ if(j) setReadings(j.readings||[]) }).catch(console.error)
  },[])

  return (
    <main style={{padding:20}}>
      <h1>Dashboard</h1>
      <table border={1} cellPadding={6}>
        <thead><tr><th>user_id</th><th>device</th><th>metric</th><th>value</th><th>timestamp</th></tr></thead>
        <tbody>
          {readings.map((r,i)=>(
            <tr key={i}><td>{r.user_id}</td><td>{r.device_id}</td><td>{r.metric}</td><td>{String(r.value)}</td><td>{r.timestamp}</td></tr>
          ))}
        </tbody>
      </table>
    </main>
  )
}
