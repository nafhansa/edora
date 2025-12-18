export interface SignalPoint {
  freq: number
  impedance: number
}

export interface RecentScan {
  id: string
  patientName: string
  deviceId: string
  metric: string
  value: number
  timestamp: string
  status?: "Normal" | "Osteoporosis" | "Alert" | string
}

export interface DashboardStats {
  totalPatientsToday: number
  activeDevices: number
  osteoporosisCases: number
  recentScans: RecentScan[]
  signal?: SignalPoint[]
}
