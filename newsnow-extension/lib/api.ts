const STORAGE_KEYS = {
  API_URL: 'apiURL',
  TOKEN: 'token',
  USERNAME: 'username'
} as const

class ApiClient {
  private baseURL: string = 'http://localhost:8080'
  private token: string | null = null

  async init() {
    const result = await chrome.storage.local.get([
      STORAGE_KEYS.API_URL,
      STORAGE_KEYS.TOKEN
    ])
    this.baseURL = result[STORAGE_KEYS.API_URL] || 'http://localhost:8080'
    this.token = result[STORAGE_KEYS.TOKEN] || null
    console.log('API initialized:', { baseURL: this.baseURL, hasToken: !!this.token })
  }

  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>)
    }

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      return await response.json() as T
    } catch (error) {
      console.error('API request failed:', error)
      throw error
    }
  }

  async getLatest(sources?: string) {
    const defaultSources = 'v2ex,zhihu,weibo,36kr,ithome'
    const ids = sources || defaultSources
    return this.request(`/api/all?id=${ids}`)
  }

  async login(username: string, password: string) {
    const response = await this.request<{ token: string }>('/api/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    })
    
    if (response.token) {
      this.token = response.token
      await chrome.storage.local.set({
        [STORAGE_KEYS.TOKEN]: response.token,
        [STORAGE_KEYS.USERNAME]: username
      })
    }
    
    return response
  }

  async executePython(code: string) {
    return this.request<{ output: string; error?: string }>('/api/python/execute', {
      method: 'POST',
      body: JSON.stringify({ code })
    })
  }

  async cleanupEnvironment() {
    return this.request('/api/python/cleanup', {
      method: 'POST'
    })
  }

  async setApiUrl(url: string) {
    this.baseURL = url
    await chrome.storage.local.set({ [STORAGE_KEYS.API_URL]: url })
  }
}

export const api = new ApiClient()
