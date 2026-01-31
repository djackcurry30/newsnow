// News item from API
export interface NewsItem {
  id: string
  title: string
  url: string
  extra?: {
    info?: string
    icon?: string
    date?: number | string
    diff?: number
  }
}

// API response for a source
export interface SourceResponse {
  id: string
  status: string
  items: NewsItem[]
  updatedTime?: number
}

// Source metadata
export interface SourceMetadata {
  name: string
  color: string
  home: string
  type?: "hottest" | "realtime"
  title?: string
  desc?: string
}

// Column configuration
export interface ColumnConfig {
  id: string
  name: string
  sources: string[]
}
