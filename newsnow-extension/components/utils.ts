import { SOURCES, SOURCE_COLORS } from "./constants"

// Get source display color class
export function getSourceColor(source: string): string {
  const sourceKey = source.toLowerCase()
  const color = SOURCES[sourceKey]?.color || "gray"
  return SOURCE_COLORS[color] || SOURCE_COLORS.gray
}

// Get source display name
export function getSourceName(source: string): string {
  const sourceKey = source.toLowerCase()
  return SOURCES[sourceKey]?.name || source
}

// Format timestamp to relative time
export function formatTime(timeStr?: number | string): string {
  if (!timeStr) return ""
  try {
    const date = typeof timeStr === "number" ? new Date(timeStr) : new Date(timeStr)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (minutes < 1) return "刚刚"
    if (minutes < 60) return `${minutes}分钟前`
    if (hours < 24) return `${hours}小时前`
    if (days < 7) return `${days}天前`

    return date.toLocaleDateString("zh-CN", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  } catch {
    return String(timeStr)
  }
}
