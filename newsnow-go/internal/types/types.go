package types

type Color = string

type SourceID string

type AllSourceID string

type ColumnID string

type FixedColumnID string

type HiddenColumnID string

type SourceType string

const (
	SourceTypeNormal   SourceType = ""
	SourceTypeHottest  SourceType = "hottest"
	SourceTypeRealtime SourceType = "realtime"
)

type Metadata struct {
	UpdatedTime int64                  `json:"updatedTime"`
	Data        map[FixedColumnID][]SourceID `json:"data"`
	Action      string                 `json:"action"`
}

type OriginSource struct {
	Name     string                   `json:"name"`
	Sub      map[string]SubSource     `json:"sub,omitempty"`
	Interval int                      `json:"interval"`
	Color    Color                    `json:"color"`
	Title    string                   `json:"title,omitempty"`
	Desc     string                   `json:"desc,omitempty"`
	Type     SourceType               `json:"type,omitempty"`
	Column   HiddenColumnID           `json:"column,omitempty"`
	Home     string                   `json:"home,omitempty"`
	Disable  bool                     `json:"disable,omitempty"`
	Redirect SourceID                 `json:"redirect,omitempty"`
}

type SubSource struct {
	Title   string                   `json:"title"`
	Type    SourceType               `json:"type,omitempty"`
	Desc    string                   `json:"desc,omitempty"`
	Column  HiddenColumnID           `json:"column,omitempty"`
	Color   Color                    `json:"color,omitempty"`
	Home    string                   `json:"home,omitempty"`
	Disable bool                     `json:"disable,omitempty"`
	Interval int                     `json:"interval,omitempty"`
}

type Source struct {
	Name     string                   `json:"name"`
	Interval int                      `json:"interval"`
	Color    Color                    `json:"color"`
	Title    string                   `json:"title,omitempty"`
	Desc     string                   `json:"desc,omitempty"`
	Type     SourceType               `json:"type,omitempty"`
	Column   HiddenColumnID           `json:"column,omitempty"`
	Home     string                   `json:"home,omitempty"`
	Disable  bool                     `json:"disable,omitempty"`
	Redirect SourceID                 `json:"redirect,omitempty"`
}

type Column struct {
	Name    string      `json:"name"`
	Sources []SourceID  `json:"sources"`
}

type NewsItem struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	MobileURL string     `json:"mobileUrl,omitempty"`
	PubDate   interface{} `json:"pubDate,omitempty"`
	Extra     *NewsExtra `json:"extra,omitempty"`
}

type NewsExtra struct {
	Hover string            `json:"hover,omitempty"`
	Date  interface{}       `json:"date,omitempty"`
	Info  interface{}       `json:"info,omitempty"`
	Diff  int               `json:"diff,omitempty"`
	Icon  interface{}       `json:"icon,omitempty"`
}

type SourceResponse struct {
	Status      string      `json:"status"`
	ID          SourceID    `json:"id"`
	UpdatedTime interface{} `json:"updatedTime"`
	Items       []NewsItem  `json:"items"`
}

type CacheInfo struct {
	ID      SourceID    `json:"id"`
	Updated int64       `json:"updated"`
	Items   []NewsItem  `json:"items"`
}

type UserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Type    string `json:"type"`
	Data    string `json:"data"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}

type RSSInfo struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Image       string    `json:"image"`
	UpdatedTime string    `json:"updatedTime"`
	Items       []RSSItem `json:"items"`
}

type RSSItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Created     string `json:"created,omitempty"`
}

type RSSHubInfo struct {
	Title        string       `json:"title"`
	HomePageURL  string       `json:"home_page_url"`
	Description  string       `json:"description"`
	Items        []RSSHubItem `json:"items"`
}

type RSSHubItem struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	ContentHTML   string `json:"content_html"`
	DatePublished string `json:"date_published"`
}

type SourceGetter func() ([]NewsItem, error)

var SourceMap = map[SourceID]*Source{
	"v2ex":       {Name: "V2EX", Interval: 600, Color: "slate"},
	"v2ex-share": {Name: "V2EX", Interval: 600, Color: "slate"},
	"zhihu":        {Name: "知乎", Interval: 600, Color: "blue"},
	"weibo":        {Name: "微博", Interval: 120, Color: "red"},
	"weibo-hot":    {Name: "微博热搜", Interval: 120, Color: "red"},
	"zaobao":       {Name: "联合早报", Interval: 1800, Color: "red"},
	"coolapk":      {Name: "酷安", Interval: 600, Color: "green"},
	"mktnews":      {Name: "MKTNews", Interval: 120, Color: "indigo"},
	"mktnews-flash": {Name: "MKTNews", Interval: 120, Color: "indigo"},
	"wallstreetcn": {Name: "华尔街见闻", Interval: 300, Color: "blue"},
	"wallstreetcn-quick": {Name: "华尔街见闻", Interval: 300, Color: "blue"},
	"wallstreetcn-news": {Name: "华尔街见闻", Interval: 1800, Color: "blue"},
	"wallstreetcn-hot": {Name: "华尔街见闻", Interval: 1800, Color: "blue"},
	"36kr":         {Name: "36氪", Interval: 600, Color: "blue"},
	"36kr-quick":   {Name: "36氪", Interval: 600, Color: "blue"},
	"36kr-renqi":   {Name: "36氪", Interval: 600, Color: "blue"},
	"douyin":       {Name: "抖音", Interval: 600, Color: "gray"},
	"hupu":         {Name: "虎扑", Interval: 600, Color: "red"},
	"tieba":        {Name: "百度贴吧", Interval: 600, Color: "blue"},
	"toutiao":      {Name: "今日头条", Interval: 600, Color: "red"},
	"ithome":       {Name: "IT之家", Interval: 600, Color: "red"},
	"thepaper":     {Name: "澎湃新闻", Interval: 1800, Color: "gray"},
	"sputniknewscn": {Name: "卫星通讯社", Interval: 600, Color: "orange"},
	"cankaoxiaoxi": {Name: "参考消息", Interval: 1800, Color: "red"},
	"pcbeta":       {Name: "远景论坛", Interval: 300, Color: "blue"},
	"pcbeta-windows11": {Name: "远景论坛", Interval: 300, Color: "blue"},
	"cls":          {Name: "财联社", Interval: 300, Color: "red"},
	"cls-telegraph": {Name: "财联社", Interval: 300, Color: "red"},
	"cls-depth":    {Name: "财联社", Interval: 600, Color: "red"},
	"cls-hot":      {Name: "财联社", Interval: 600, Color: "red"},
	"xueqiu":       {Name: "雪球", Interval: 120, Color: "blue"},
	"xueqiu-hotstock": {Name: "雪球", Interval: 120, Color: "blue"},
	"gelonghui":    {Name: "格隆汇", Interval: 120, Color: "blue"},
	"fastbull":     {Name: "法布财经", Interval: 120, Color: "emerald"},
	"fastbull-express": {Name: "法布财经", Interval: 120, Color: "emerald"},
	"fastbull-news": {Name: "法布财经", Interval: 1800, Color: "emerald"},
	"solidot":      {Name: "Solidot", Interval: 3600, Color: "teal"},
	"hackernews":   {Name: "Hacker News", Interval: 600, Color: "orange"},
	"producthunt":  {Name: "Product Hunt", Interval: 600, Color: "red"},
	"github":                  {Name: "Github", Interval: 600, Color: "gray"},
	"github-trending-today": {Name: "Github", Interval: 600, Color: "gray"},
	"bilibili":     {Name: "哔哩哔哩", Interval: 600, Color: "blue"},
	"bilibili-hot-search": {Name: "哔哩哔哩", Interval: 600, Color: "blue"},
	"bilibili-hot-video": {Name: "哔哩哔哩", Interval: 600, Color: "blue"},
	"bilibili-ranking": {Name: "哔哩哔哩", Interval: 1800, Color: "blue"},
	"kuaishou":     {Name: "快手", Interval: 600, Color: "orange"},
	"kaopu":        {Name: "靠谱新闻", Interval: 1800, Color: "gray"},
	"jin10":        {Name: "金十数据", Interval: 600, Color: "blue"},
	"baidu":        {Name: "百度热搜", Interval: 600, Color: "blue"},
	"nowcoder":     {Name: "牛客", Interval: 600, Color: "blue"},
	"sspai":        {Name: "少数派", Interval: 600, Color: "red"},
	"juejin":       {Name: "稀土掘金", Interval: 600, Color: "blue"},
	"ifeng":        {Name: "凤凰网", Interval: 600, Color: "red"},
	"chongbuluo":   {Name: "虫部落", Interval: 1800, Color: "green"},
	"chongbuluo-latest": {Name: "虫部落", Interval: 1800, Color: "green"},
	"chongbuluo-hot": {Name: "虫部落", Interval: 1800, Color: "green"},
	"douban":       {Name: "豆瓣", Interval: 600, Color: "green"},
	"steam":        {Name: "Steam", Interval: 600, Color: "blue"},
	"tencent":      {Name: "腾讯新闻", Interval: 1800, Color: "blue"},
	"tencent-hot":  {Name: "腾讯新闻", Interval: 1800, Color: "blue"},
	"freebuf":      {Name: "Freebuf", Interval: 600, Color: "green"},
	"qqvideo":      {Name: "腾讯视频", Interval: 1800, Color: "blue"},
	"qqvideo-tv-hotsearch": {Name: "腾讯视频", Interval: 1800, Color: "blue"},
	"iqiyi":        {Name: "爱奇艺", Interval: 1800, Color: "green"},
	"iqiyi-hot-ranklist": {Name: "爱奇艺", Interval: 1800, Color: "green"},
}
