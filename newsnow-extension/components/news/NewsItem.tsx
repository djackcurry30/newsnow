import { cn } from "~/lib/utils"
import type { NewsItem as NewsItemType } from "../types"

interface NewsItemProps {
  item: NewsItemType
  index: number
  onClick: (url: string) => void
}

export function NewsItem({ item, index, onClick }: NewsItemProps) {
  return (
    <div
      className={cn(
        "group flex items-start gap-2 p-2 rounded-lg",
        "hover:bg-gray-100/50 transition-colors cursor-pointer"
      )}
      onClick={() => onClick(item.url)}>
      {item.extra?.diff !== undefined && (
        <span
          className={cn(
            "text-[10px] font-medium mt-0.5",
            item.extra.diff < 0 ? "text-green-600" : "text-red-500"
          )}>
          {item.extra.diff > 0 ? `+${item.extra.diff}` : item.extra.diff}
        </span>
      )}
      <div className="flex-1 min-w-0">
        <span className="text-xs text-gray-700 leading-relaxed line-clamp-2 group-hover:text-gray-900">
          {index + 1}. {item.title}
        </span>
        {item.extra?.info && (
          <span className="text-[10px] text-gray-400 mt-0.5 block">{item.extra.info}</span>
        )}
      </div>
    </div>
  )
}
