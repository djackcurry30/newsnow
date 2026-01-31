import { Button } from "~/components/ui/button"
import { TabBar } from "./TabBar"
import { LoadingSpinner } from "../common/LoadingSpinner"

interface HeaderProps {
  activeColumn: string
  isLoading: boolean
  onColumnChange: (columnId: string) => void
  onRefresh: () => void
}

export function Header({
  activeColumn,
  isLoading,
  onColumnChange,
  onRefresh,
}: HeaderProps) {
  return (
    <div className="px-4 py-3 bg-white border-b border-gray-200 shadow-sm flex-shrink-0">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-500 via-blue-600 to-indigo-600 rounded-lg flex items-center justify-center shadow-sm">
            <span className="text-white text-sm font-bold">N</span>
          </div>
          <div>
            <h1 className="text-sm font-bold text-gray-800">NewsNow</h1>
            <p className="text-[10px] text-gray-400">实时资讯聚合</p>
          </div>
        </div>
        <Button
          size="sm"
          variant="ghost"
          onClick={onRefresh}
          disabled={isLoading}
          className="h-8 w-8 p-0 text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg"
        >
          {isLoading ? (
            <LoadingSpinner />
          ) : (
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>
          )}
        </Button>
      </div>

      <TabBar activeColumn={activeColumn} onColumnChange={onColumnChange} />
    </div>
  )
}
