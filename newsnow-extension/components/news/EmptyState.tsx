export function EmptyState() {
  return (
    <div className="text-center py-8 text-gray-400 text-xs">
      <svg
        className="w-8 h-8 mx-auto mb-2 text-gray-300"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="1.5"
          d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
        />
      </svg>
      <p>暂无数据</p>
    </div>
  )
}
