import type { SourceMetadata, ColumnConfig } from "./types"

// Source metadata from shared/sources.json
export const SOURCES: Record<string, SourceMetadata> = {
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

// Color mapping for source badges
export const SOURCE_COLORS: Record<string, string> = {
  slate: "bg-slate-500/15 text-slate-700 border-slate-300",
  blue: "bg-blue-500/15 text-blue-700 border-blue-300",
  red: "bg-red-500/15 text-red-700 border-red-300",
  purple: "bg-purple-500/15 text-purple-700 border-purple-300",
  sky: "bg-sky-500/15 text-sky-700 border-sky-300",
  green: "bg-green-500/15 text-green-700 border-green-300",
  orange: "bg-orange-500/15 text-orange-700 border-orange-300",
  gray: "bg-gray-500/15 text-gray-700 border-gray-300",
  pink: "bg-pink-500/15 text-pink-700 border-pink-300",
  cyan: "bg-cyan-500/15 text-cyan-700 border-cyan-300",
  yellow: "bg-yellow-500/15 text-yellow-700 border-yellow-300",
  indigo: "bg-indigo-500/15 text-indigo-700 border-indigo-300",
  teal: "bg-teal-500/15 text-teal-700 border-teal-300",
  black: "bg-gray-800/15 text-gray-800 border-gray-400",
}

// Column configuration
export const COLUMNS: ColumnConfig[] = [
  { id: "focus", name: "关注", sources: ["v2ex", "zhihu", "weibo"] },
  { id: "hottest", name: "最热", sources: ["zhihu", "weibo", "bilibili", "tieba"] },
  { id: "tech", name: "科技", sources: ["v2ex", "36kr", "ithome", "coolapk", "hackernews", "github", "producthunt"] },
  { id: "finance", name: "财经", sources: ["wallstreetcn", "cls", "xueqiu", "eastmoney"] },
  { id: "china", name: "国内", sources: ["zhihu", "weibo", "zaobao", "thepaper"] },
  { id: "world", name: "国际", sources: ["sputniknewscn", "cankaoxiaoxi"] },
]
