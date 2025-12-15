"use client"
import { useState } from 'react'
import { useRouter } from 'next/navigation'

export default function LoginPage(){
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [err, setErr] = useState('')
  const router = useRouter()

  async function submit(e:any){
    e.preventDefault()
    setErr('')
    try{
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/login`, {
        method: 'POST',
        headers: {'Content-Type':'application/json'},
        body: JSON.stringify({email, password})
      })
      const j = await res.json()
      if(!res.ok){ setErr(j.error||'login failed'); return }
      localStorage.setItem('token', j.token)
      router.push('/dashboard')
    }catch(err:any){ setErr(err.message) }
  }

  return (
    <main style={{padding:20}}>
      <h1>Login</h1>
      <form onSubmit={submit}>
        <div>
          <label>Email</label><br/>
          <input value={email} onChange={e=>setEmail(e.target.value)} />
        </div>
        <div>
          <label>Password</label><br/>
          <input type="password" value={password} onChange={e=>setPassword(e.target.value)} />
        </div>
        <div style={{marginTop:10}}>
          <button type="submit">Login</button>
        </div>
        {err && <div style={{color:'red'}}>{err}</div>}
      </form>
    </main>
  )
}
