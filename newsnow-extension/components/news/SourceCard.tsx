import { cn } from "~/lib/utils"
import { EmptyState } from "./EmptyState"
import { NewsItem } from "./NewsItem"
import { getSourceColor, getSourceName, formatTime } from "../utils"
import { SOURCES, SOURCE_COLORS } from "../constants"
import type { SourceResponse, SourceMetadata } from "../types"
import { useState } from "react"

// Source icon component with image fallback to letter
function SourceIcon({ sourceId, sourceInfo }: { sourceId: string; sourceInfo?: SourceMetadata }) {
  const [imgError, setImgError] = useState(false)
  const baseSourceId = sourceId.split("-")[0]
  const iconUrl = `/icons/${baseSourceId}.png`

  // Get color for background
  const colorKey = sourceInfo?.color || "gray"
  const colorClass = SOURCE_COLORS[colorKey] || SOURCE_COLORS.gray

  // Extract color from tailwind class for background
  const bgColors: Record<string, string> = {
    slate: "#64748b",
    blue: "#3b82f6",
    red: "#ef4444",
    purple: "#a855f7",
    sky: "#0ea5e9",
    green: "#22c55e",
    orange: "#f97316",
    gray: "#6b7280",
    pink: "#ec4899",
    cyan: "#06b6d4",
    yellow: "#eab308",
    indigo: "#6366f1",
    teal: "#14b8a6",
    black: "#1f2937",
  }
  const bgColor = bgColors[colorKey] || bgColors.gray

  // Get first letter of source name
  const letter = (sourceInfo?.name || sourceId).charAt(0).toUpperCase()

  if (imgError) {
    return (
      <div
        className={cn("w-7 h-7 rounded-full flex items-center justify-center text-white text-xs font-bold", colorClass)}
        style={{ backgroundColor: bgColor }}
        title={sourceInfo?.name || sourceId}
      >
        {letter}
      </div>
    )
  }

  return (
    <img
      src={iconUrl}
      alt={sourceInfo?.name || sourceId}
      className="w-7 h-7 rounded-full object-cover"
      onError={() => setImgError(true)}
      style={{ backgroundColor: "#f3f4f6" }}
    />
  )
}

interface SourceCardProps {
  sourceId: string
  data?: SourceResponse
  isLoading?: boolean
  onRefresh: (sourceId: string) => void
  onItemClick: (url: string) => void
}

export function SourceCard({
  sourceId,
  data,
  isLoading,
  onRefresh,
  onItemClick,
}: SourceCardProps) {
  const sourceInfo = SOURCES[sourceId.toLowerCase()]

  return (
    <div className="flex-shrink-0 w-[280px] bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
      {/* Source Header */}
      <div className="px-3 py-2.5 bg-gray-50 border-b border-gray-200 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <SourceIcon sourceId={sourceId} sourceInfo={sourceInfo} />
          <div className="flex flex-col">
            <div className="flex items-center gap-1.5">
              <span className="text-sm font-semibold text-gray-800">
                {getSourceName(sourceId)}
              </span>
              {sourceInfo?.title && (
                <span
                  className={cn(
                    "text-[10px] px-1 py-0.5 rounded",
                    getSourceColor(sourceId)
                  )}>
                  {sourceInfo.title}
                </span>
              )}
            </div>
            <span className="text-[10px] text-gray-400">
              {isLoading
                ? "加载中..."
                : data?.updatedTime
                  ? formatTime(data.updatedTime)
                  : "获取失败"}
            </span>
          </div>
        </div>
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => onRefresh(sourceId)}
            className={cn(
              "p-1.5 rounded-lg text-gray-500 hover:text-blue-600 hover:bg-blue-50 transition-colors",
              isLoading && "animate-spin"
            )}
            disabled={isLoading}>
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>
          </button>
        </div>
      </div>

      {/* News List */}
      <div className="p-2 space-y-0.5 max-h-[320px] overflow-y-auto scrollbar-thin">
        {!data?.items?.length && !isLoading ? (
          <EmptyState />
        ) : (
          data?.items?.slice(0, 15).map((item, index) => (
            <NewsItem
              key={item.id || index}
              item={item}
              index={index}
              onClick={onItemClick}
            />
          ))
        )}
      </div>
    </div>
  )
}
