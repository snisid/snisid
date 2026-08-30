const API_BASE = '/intelligence/v1'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...init?.headers },
    ...init,
  })
  if (!res.ok) throw new Error(`API ${res.status}: ${res.statusText}`)
  return res.json()
}

export const api = {
  search: (q: string) =>
    request<any>(`/search/unified?q=${encodeURIComponent(q)}`),

  getAlerts: (severity?: string) =>
    request<any[]>(`/alerts${severity ? `?severity=${severity}` : ''}`),

  getAlert: (id: string) =>
    request<any>(`/alerts/${id}`),

  getReport: (type: string, id: string) =>
    request<any>(`/ai/report/${type}/${id}`),

  getDashboard: () =>
    request<any>('/dashboard/stats'),

  getWantedPersons: (params?: string) =>
    request<any[]>(`/ncid/wanted-persons${params ? `?${params}` : ''}`),

  getEntityGraph: (type: string, id: string) =>
    request<any>(`/search/graph/${type}/${id}`),

  generateIntelligence: (body: any) =>
    request<any>('/ai/report', {
      method: 'POST',
      body: JSON.stringify(body),
    }),

  searchIOC: (ioc: string) =>
    request<any>(`/cyber/iocs?q=${encodeURIComponent(ioc)}`),
}
