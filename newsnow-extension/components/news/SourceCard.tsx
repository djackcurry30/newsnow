import { cn } from "~/lib/utils"
import { EmptyState } from "./EmptyState"
import { NewsItem } from "./NewsItem"
import { getSourceColor, getSourceName, formatTime } from "../utils"
import { SOURCES } from "../constants"
import type { SourceResponse } from "../types"

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
          <div
            className="w-7 h-7 rounded-full bg-cover bg-center"
            style={{
              backgroundImage: `url(/icons/${sourceId.split("-")[0]}.png)`,
              backgroundColor: "#f3f4f6",
            }}
          />
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
