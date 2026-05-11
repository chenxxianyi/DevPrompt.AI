/**
 * iOS26 frontend — constants
 * 仅服务于新版 ios26 视图，不影响旧版逻辑。
 */

export const IOS_ROOT = '/ios26'

export interface IosTabItem {
  /** vue-router 命名 */
  name: string
  /** 跳转路径 */
  path: string
  /** 显示标签 */
  label: string
  /** Lucide-style SVG path (24x24) */
  icon: string
  /** 用于判定 active：当 route.path 以此前缀开头视为激活 */
  match: string[]
}

/** 底部 TabBar 主入口 */
export const IOS_TABS: IosTabItem[] = [
  {
    name: 'ios26-home',
    path: '/ios26',
    label: '首页',
    icon: 'M3 11.5 12 4l9 7.5V20a1 1 0 0 1-1 1h-5v-6h-6v6H4a1 1 0 0 1-1-1Z',
    match: ['/ios26', '/ios26/'],
  },
  {
    name: 'ios26-prompts',
    path: '/ios26/prompts',
    label: '模板',
    icon: 'M4 6h16M4 12h16M4 18h10',
    match: ['/ios26/prompts'],
  },
  {
    name: 'ios26-generator',
    path: '/ios26/generator',
    label: '生成器',
    icon: 'M6.464 6.464A5 5 0 0 1 10 5.036V3a1 1 0 1 1 2 0v2.036a5 5 0 0 1 3.536 1.428l1.41 1.41-1.414 1.415-1.41-1.41A3 3 0 0 0 10 7.036v2.928l-3.536 3.536a3 3 0 0 0 4.242 4.242l1.41 1.41-1.414 1.415-1.41-1.41A5 5 0 0 1 6.464 6.464Z',
    match: ['/ios26/generator'],
  },
  {
    name: 'ios26-dashboard',
    path: '/ios26/dashboard',
    label: '我的',
    icon: 'M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8Zm0 2c-3.866 0-7 2.239-7 5v1h14v-1c0-2.761-3.134-5-7-5Z',
    match: ['/ios26/dashboard'],
  },
]

/** 顶部 NavBar 在不同路径下的标题映射；缺省时使用兜底 */
export const IOS_TITLES: Record<string, string> = {
  '/ios26': 'DevPrompt',
  '/ios26/prompts': '模板库',
  '/ios26/generator': '生成器',
  '/ios26/generator/history': '历史记录',
  '/ios26/dashboard': '我的',
  '/ios26/login': '登录',
  '/ios26/pricing': '会员方案',
}

/** 生成器类型，复用旧版逻辑的 GeneratorTab */
export const IOS_GENERATOR_TABS = [
  { id: 'project', label: '项目 Prompt' },
  { id: 'cursor-rules', label: 'Cursor Rules' },
  { id: 'claude-code', label: 'Claude Code' },
  { id: 'optimize', label: 'Prompt 优化' },
] as const

/** 主题存储 key */
export const IOS_THEME_STORAGE_KEY = 'ios26-theme'
export type IosTheme = 'system' | 'light' | 'dark'
