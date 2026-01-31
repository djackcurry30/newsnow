import { useRef } from "react"
import { SourceCard } from "../news/SourceCard"
import { COLUMNS } from "../constants"
import type { SourceResponse } from "../types"

interface ContentAreaProps {
  activeColumn: string
  sourceData: Record<string, SourceResponse>
  loading: Record<string, boolean>
  error?: string
  onRefresh: (sourceId: string) => void
  onItemClick: (url: string) => void
}

export function ContentArea({
  activeColumn,
  sourceData,
  loading,
  error,
  onRefresh,
  onItemClick,
}: ContentAreaProps) {
  const scrollRef = useRef<HTMLDivElement>(null)
  const currentColumn = COLUMNS.find((c) => c.id === activeColumn)

  return (
    <div
      ref={scrollRef}
      className="flex-1 overflow-x-auto overflow-y-hidden scrollbar-thin"
    >
      <div className="flex gap-4 p-4 h-full">
        {currentColumn?.sources.map((sourceId) => (
          <SourceCard
            key={sourceId}
            sourceId={sourceId}
            data={sourceData[sourceId]}
            isLoading={loading[sourceId]}
            onRefresh={onRefresh}
            onItemClick={onItemClick}
          />
        ))}
      </div>
    </div>
  )
}
