import { cn } from "~/lib/utils"
import { COLUMNS } from "../constants"

interface TabBarProps {
  activeColumn: string
  onColumnChange: (columnId: string) => void
}

export function TabBar({ activeColumn, onColumnChange }: TabBarProps) {
  return (
    <div className="flex gap-1 bg-gray-100 p-1 rounded-lg">
      {COLUMNS.map((col) => (
        <button
          key={col.id}
          onClick={() => onColumnChange(col.id)}
          className={cn(
            "flex-1 text-xs py-1.5 px-2 rounded-md transition-all",
            activeColumn === col.id
              ? "bg-white text-gray-800 shadow-sm font-medium"
              : "text-gray-500 hover:text-gray-700"
          )}>
          {col.name}
        </button>
      ))}
    </div>
  )
}
