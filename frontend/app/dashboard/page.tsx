"use client"
import React, { useMemo } from "react"
import useSWR from "swr"
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  ResponsiveContainer,
} from "recharts"
import type { DashboardStats, SignalPoint, RecentScan } from "../../types/dashboard"

const API = `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/v1/dashboard/stats`
const fetcher = async (url: string) => {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null
  const res = await fetch(url, { headers: token ? { Authorization: `Bearer ${token}` } : undefined })
  if (!res.ok) throw new Error("failed to fetch")
  return res.json()
}

export default function DashboardPage() {
  const { data, error } = useSWR<DashboardStats>(API, fetcher, { refreshInterval: 5000 })

  const signalMock: SignalPoint[] = useMemo(() => {
    const freqs = [50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000, 50000]
    return freqs.map((f, i) => ({ freq: f, impedance: Math.round((100 + Math.sin(i / 2) * 30 + Math.random() * 10) * 10) / 10 }))
  }, [])

  const stats = data ?? {
    totalPatientsToday: 0,
    activeDevices: 0,
    osteoporosisCases: 0,
    recentScans: [] as RecentScan[],
  }

  return (
    <main className="min-h-screen bg-sky-50 p-6">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-2xl font-semibold text-slate-800 mb-4">Monitoring Dashboard</h1>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
          <div className="bg-white rounded shadow p-4">
            <div className="text-sm text-slate-500">Patients Today</div>
            <div className="text-2xl font-bold text-slate-800">{stats.totalPatientsToday}</div>
          </div>
          <div className="bg-white rounded shadow p-4">
            <div className="text-sm text-slate-500">Active Devices</div>
            <div className="text-2xl font-bold text-slate-800">{stats.activeDevices}</div>
          </div>
          <div className="bg-white rounded shadow p-4">
            <div className="text-sm text-slate-500">Osteoporosis Cases</div>
            <div className={`text-2xl font-bold ${stats.osteoporosisCases > 0 ? "text-red-600" : "text-slate-800"}`}>
              {stats.osteoporosisCases}
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-white rounded shadow p-4">
            <div className="flex items-center justify-between mb-2">
              <h2 className="text-lg font-medium text-slate-800">Impedance Signal (mock)</h2>
              <div className="text-sm text-slate-500">Frequency vs Impedance (Ω)</div>
            </div>
            <div style={{ width: "100%", height: 320 }}>
              <ResponsiveContainer>
                <LineChart data={signalMock as any}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="freq" tickFormatter={(v) => `${v}`} />
                  <YAxis />
                  <Tooltip formatter={(v: any) => [`${v} Ω`, "Impedance"]} labelFormatter={(l) => `Freq: ${l} Hz`} />
                  <Line type="monotone" dataKey="impedance" stroke="#2563EB" strokeWidth={2} dot={{ r: 3 }} />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </div>

          <div className="bg-white rounded shadow p-4">
            <h3 className="text-lg font-medium text-slate-800 mb-3">Recent Scans</h3>
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead>
                  <tr className="text-slate-500">
                    <th className="py-2 px-2">Patient</th>
                    <th className="py-2 px-2">Device</th>
                    <th className="py-2 px-2">Metric</th>
                    <th className="py-2 px-2">Value</th>
                    <th className="py-2 px-2">Time</th>
                  </tr>
                </thead>
                <tbody>
                  {(stats.recentScans.length ? stats.recentScans : [
                    { id: "m-1", patientName: "Demo Patient", deviceId: "ESP32-001", metric: "Impedance", value: 120.5, timestamp: new Date().toISOString(), status: "Normal" }
                  ]).map((r: any) => (
                    <tr key={r.id} className="border-t">
                      <td className="py-2 px-2 text-slate-700">{r.patientName}</td>
                      <td className="py-2 px-2 text-slate-700">{r.deviceId}</td>
                      <td className="py-2 px-2 text-slate-700">{r.metric}</td>
                      <td className={`py-2 px-2 font-medium ${r.status === "Osteoporosis" ? "text-red-600" : "text-slate-700"}`}>{r.value}</td>
                      <td className="py-2 px-2 text-slate-500">{new Date(r.timestamp).toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </main>
  )
}
