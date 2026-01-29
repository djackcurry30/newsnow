import { useCallback, useEffect, useState, useRef } from "react"

import { Button } from "~/components/ui/button"
import { cn } from "~/lib/utils"
import { api } from "~/lib/api"

import "./globals.css"

// Source metadata from shared/sources.json
const SOURCES: Record<string, {
  name: string
  color: string
  home: string
  type?: "hottest" | "realtime"
  title?: string
  desc?: string
}> = {
  "v2ex": { name: "V2EX", color: "slate", home: "https://v2ex.com/", title: "最新分享" },
  "v2ex-share": { name: "V2EX", color: "slate", home: "https://v2ex.com/", title: "最新分享" },
  "zhihu": { name: "知乎", color: "blue", home: "https://www.zhihu.com", type: "hottest" },
  "weibo": { name: "微博", color: "red", home: "https://weibo.com", type: "hottest", title: "实时热搜" },
  "zaobao": { name: "联合早报", color: "red", home: "https://www.zaobao.com", type: "realtime" },
  "coolapk": { name: "酷安", color: "green", home: "https://coolapk.com", type: "hottest", title: "今日最热" },
  "wallstreetcn": { name: "华尔街见闻", color: "blue", home: "https://wallstreetcn.com/", type: "realtime", title: "快讯" },
  "wallstreetcn-quick": { name: "华尔街见闻", color: "blue", home: "https://wallstreetcn.com/", type: "realtime", title: "快讯" },
  "wallstreetcn-news": { name: "华尔街见闻", color: "blue", home: "https://wallstreetcn.com/", title: "最新" },
  "wallstreetcn-hot": { name: "华尔街见闻", color: "blue", home: "https://wallstreetcn.com/", type: "hottest", title: "最热" },
  "36kr": { name: "36氪", color: "blue", home: "https://36kr.com", type: "realtime", title: "快讯" },
  "36kr-quick": { name: "36氪", color: "blue", home: "https://36kr.com", type: "realtime", title: "快讯" },
  "36kr-renqi": { name: "36氪", color: "blue", home: "https://36kr.com", type: "hottest", title: "人气榜" },
  "douyin": { name: "抖音", color: "gray", home: "https://www.douyin.com", type: "hottest" },
  "hupu": { name: "虎扑", color: "red", home: "https://hupu.com", type: "hottest", title: "主干道热帖" },
  "tieba": { name: "百度贴吧", color: "blue", home: "https://tieba.baidu.com", type: "hottest", title: "热议" },
  "toutiao": { name: "今日头条", color: "red", home: "https://www.toutiao.com", type: "hottest" },
  "ithome": { name: "IT之家", color: "red", home: "https://www.ithome.com", type: "realtime" },
  "thepaper": { name: "澎湃新闻", color: "gray", home: "https://www.thepaper.cn", type: "hottest", title: "热榜" },
  "sputniknewscn": { name: "卫星通讯社", color: "orange", home: "https://sputniknews.cn" },
  "cankaoxiaoxi": { name: "参考消息", color: "red", home: "https://china.cankaoxiaoxi.com" },
  "hackernews": { name: "Hacker News", color: "orange", home: "https://news.ycombinator.com", type: "hottest" },
  "github": { name: "GitHub", color: "gray", home: "https://github.com", type: "hottest" },
  "producthunt": { name: "Product Hunt", color: "orange", home: "https://www.producthunt.com", type: "hottest" },
  "bilibili": { name: "哔哩哔哩", color: "pink", home: "https://www.bilibili.com", type: "hottest" },
  "cls": { name: "财联社", color: "red", home: "https://www.cls.cn", type: "realtime", title: "快讯" },
  "xueqiu": { name: "雪球", color: "blue", home: "https://xueqiu.com", type: "hottest" },
  "eastmoney": { name: "东方财富", color: "red", home: "https://eastmoney.com", type: "hottest" },
}

interface NewsItem {
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

interface SourceResponse {
  id: string
  status: string
  items: NewsItem[]
  updatedTime?: number
}

// Columns configuration based on metadata
const COLUMNS = [
  { id: "focus", name: "关注", sources: ["v2ex", "zhihu", "weibo"] },
  { id: "hottest", name: "最热", sources: ["zhihu", "weibo", "bilibili", "tieba"] },
  { id: "tech", name: "科技", sources: ["v2ex", "36kr", "ithome", "coolapk", "hackernews", "github", "producthunt"] },
  { id: "finance", name: "财经", sources: ["wallstreetcn", "cls", "xueqiu", "eastmoney"] },
  { id: "china", name: "国内", sources: ["zhihu", "weibo", "zaobao", "thepaper"] },
  { id: "world", name: "国际", sources: ["sputniknewscn", "cankaoxiaoxi"] },
] as const

function getSourceColor(source: string): string {
  const sourceKey = source.toLowerCase()
  const color = SOURCES[sourceKey]?.color || "gray"
  const colors: Record<string, string> = {
    "slate": "bg-slate-500/15 text-slate-700 border-slate-300",
    "blue": "bg-blue-500/15 text-blue-700 border-blue-300",
    "red": "bg-red-500/15 text-red-700 border-red-300",
    "purple": "bg-purple-500/15 text-purple-700 border-purple-300",
    "sky": "bg-sky-500/15 text-sky-700 border-sky-300",
    "green": "bg-green-500/15 text-green-700 border-green-300",
    "orange": "bg-orange-500/15 text-orange-700 border-orange-300",
    "gray": "bg-gray-500/15 text-gray-700 border-gray-300",
    "pink": "bg-pink-500/15 text-pink-700 border-pink-300",
    "cyan": "bg-cyan-500/15 text-cyan-700 border-cyan-300",
    "yellow": "bg-yellow-500/15 text-yellow-700 border-yellow-300",
    "indigo": "bg-indigo-500/15 text-indigo-700 border-indigo-300",
    "teal": "bg-teal-500/15 text-teal-700 border-teal-300",
    "black": "bg-gray-800/15 text-gray-800 border-gray-400",
  }
  return colors[color] || colors["gray"]
}

function getSourceName(source: string): string {
  const sourceKey = source.toLowerCase()
  return SOURCES[sourceKey]?.name || source
}

function getSourceTitle(source: string): string {
  const sourceKey = source.toLowerCase()
  return SOURCES[sourceKey]?.title || ""
}

function formatTime(timeStr?: number | string): string {
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

function IndexPopup() {
  const [activeColumn, setActiveColumn] = useState("focus")
  const [sourceData, setSourceData] = useState<Record<string, SourceResponse>>({})
  const [loading, setLoading] = useState<Record<string, boolean>>({})
  const [error, setError] = useState("")
  const scrollRef = useRef<HTMLDivElement>(null)

  // Initialize API client
  useEffect(() => {
    api.init()
  }, [])

  // Load news for a specific source
  const loadSource = useCallback(async (sourceId: string) => {
    setLoading(prev => ({ ...prev, [sourceId]: true }))
    setError("")

    try {
      const data = await api.request<SourceResponse>(`/api/s/${sourceId}`)
      setSourceData(prev => ({ ...prev, [sourceId]: data }))
    } catch (err: any) {
      console.error(`Failed to load ${sourceId}:`, err)
      setError(`加载 ${getSourceName(sourceId)} 失败`)
    } finally {
      setLoading(prev => ({ ...prev, [sourceId]: false }))
    }
  }, [])

  // Load all sources in active column
  const loadColumnNews = useCallback(async (columnId: string) => {
    const column = COLUMNS.find(c => c.id === columnId)
    if (!column) return

    // Load each source in parallel
    await Promise.all(column.sources.map(sourceId => loadSource(sourceId)))
  }, [loadSource])

  // Load initial data
  useEffect(() => {
    loadColumnNews("focus")
  }, [loadColumnNews])

  // Refresh all sources in active column
  const refreshAll = useCallback(async () => {
    setError("")
    await loadColumnNews(activeColumn)
  }, [activeColumn, loadColumnNews])

  // Handle column change
  const handleColumnChange = (columnId: string) => {
    setActiveColumn(columnId)
    loadColumnNews(columnId)
  }

  const handleNewsItemClick = (url: string) => {
    if (url && url !== "#") {
      window.open(url, "_blank")
    }
  }

  const LoadingSpinner = () => (
    <svg className="animate-spin h-4 w-4 text-blue-600" viewBox="0 0 24 24">
      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
    </svg>
  )

  const EmptyState = () => (
    <div className="text-center py-8 text-gray-400 text-xs">
      <svg className="w-8 h-8 mx-auto mb-2 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
      </svg>
      <p>暂无数据</p>
    </div>
  )

  const NewsItemComponent = ({ item, index }: { item: NewsItem; index: number }) => (
    <div
      className={cn(
        "group flex items-start gap-2 p-2 rounded-lg",
        "hover:bg-gray-100/50 transition-colors cursor-pointer"
      )}
      onClick={() => handleNewsItemClick(item.url)}
    >
      {item.extra?.diff !== undefined && (
        <span className={cn(
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

  const SourceCard = ({ sourceId }: { sourceId: string }) => {
    const data = sourceData[sourceId]
    const isLoading = loading[sourceId]
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
                backgroundColor: "#f3f4f6"
              }}
            />
            <div className="flex flex-col">
              <div className="flex items-center gap-1.5">
                <span className="text-sm font-semibold text-gray-800">{getSourceName(sourceId)}</span>
                {sourceInfo?.title && (
                  <span className={cn("text-[10px] px-1 py-0.5 rounded", getSourceColor(sourceId))}>
                    {sourceInfo.title}
                  </span>
                )}
              </div>
              <span className="text-[10px] text-gray-400">
                {isLoading ? "加载中..." : data?.updatedTime ? formatTime(data.updatedTime) : "获取失败"}
              </span>
            </div>
          </div>
          <div className="flex items-center gap-1">
            <button
              type="button"
              onClick={() => loadSource(sourceId)}
              className={cn(
                "p-1.5 rounded-lg text-gray-500 hover:text-blue-600 hover:bg-blue-50 transition-colors",
                isLoading && "animate-spin"
              )}
              disabled={isLoading}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
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
              <NewsItemComponent key={item.id || index} item={item} index={index} />
            ))
          )}
        </div>
      </div>
    )
  }

  const currentColumn = COLUMNS.find(c => c.id === activeColumn)
  const isAnyLoading = Object.values(loading).some(Boolean)

  return (
    <div className="w-[600px] h-[500px] bg-gray-100 flex flex-col overflow-hidden">
      {/* Header */}
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
            onClick={refreshAll}
            disabled={isAnyLoading}
            className="h-8 w-8 p-0 text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg"
          >
            {isAnyLoading ? <LoadingSpinner /> : (
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            )}
          </Button>
        </div>

        {/* Column Tabs */}
        <div className="flex gap-1 bg-gray-100 p-1 rounded-lg">
          {COLUMNS.map(col => (
            <button
              key={col.id}
              onClick={() => handleColumnChange(col.id)}
              className={cn(
                "flex-1 text-xs py-1.5 px-2 rounded-md transition-all",
                activeColumn === col.id
                  ? "bg-white text-gray-800 shadow-sm font-medium"
                  : "text-gray-500 hover:text-gray-700"
              )}
            >
              {col.name}
            </button>
          ))}
        </div>
      </div>

      {/* Error Banner */}
      {error && (
        <div className="mx-4 mt-3 p-2.5 bg-red-50 border border-red-100 rounded-lg text-red-600 text-xs flex items-center gap-2 flex-shrink-0">
          <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {error}
        </div>
      )}

      {/* Content - Horizontal Scroll */}
      <div
        ref={scrollRef}
        className="flex-1 overflow-x-auto overflow-y-hidden scrollbar-thin"
      >
        <div className="flex gap-4 p-4 h-full">
          {currentColumn?.sources.map(sourceId => (
            <SourceCard key={sourceId} sourceId={sourceId} />
          ))}
        </div>
      </div>

      {/* Footer */}
      <div className="px-4 py-2 bg-white border-t border-gray-200 text-center flex-shrink-0">
        <span className="text-[10px] text-gray-400">NewsNow Extension</span>
      </div>
    </div>
  )
}

export default IndexPopup
