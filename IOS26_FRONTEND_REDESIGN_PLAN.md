# iOS 26 Liquid Glass Frontend Redesign Plan

## 目标

在不改变现有前端风格和文件结构的前提下，为 `web` 前端新增一版基于 iOS 26 / Liquid Glass 设计语言的新界面。

新版前端作为独立风格层开发，通过 `/ios26` 路由访问。旧版页面继续保留在当前路径，避免影响现有功能、样式和上线稳定性。

## 当前项目判断

当前项目包含两个 Vite 前端入口：

- `web`: 面向用户侧的 Vue 3 + Vite + Tailwind 前端
- `admin`: 管理端 Vite 前端

本次改造建议优先作用于 `web`，因为当前用户侧页面集中在：

```text
web/src/views/
web/src/components/
web/src/router/
web/src/assets/main.css
```

旧版风格主要由 `web/src/assets/main.css`、`web/tailwind.config.js` 和现有 Vue 组件共同构成。新版不修改这些旧样式，而是在 `web/src/ios26` 下新建独立目录。

## 总体方案

采用“同一 Vue/Vite 应用内新增 iOS26 独立前端层”的方式：

- 旧版首页继续使用 `/`
- 新版首页使用 `/ios26`
- 旧版生成器继续使用 `/generator`
- 新版生成器使用 `/ios26/generator`
- 旧版组件、旧版 CSS、旧版 Tailwind 配置不做破坏性修改
- 新版 UI 组件、样式 token、页面视图全部放入新目录

推荐新增目录：

```text
web/src/ios26/
  styles/
    tokens.css
    materials.css
    typography.css
    animations.css
    index.css
  components/
    IosAppShell.vue
    IosNavBar.vue
    IosTabBar.vue
    IosButton.vue
    IosGlassPanel.vue
    IosListRow.vue
    IosSheet.vue
    IosSegmentedControl.vue
    IosToast.vue
  views/
    IosHomeView.vue
    IosPromptsView.vue
    IosPromptDetailView.vue
    IosGeneratorView.vue
    IosHistoryView.vue
    IosDashboardView.vue
    IosLoginView.vue
    IosPricingView.vue
  router.ts
  constants.ts
```

## 路由设计

在 `web/src/router/index.ts` 中新增 `/ios26` 路由组。旧路由保持不变。

示例结构：

```ts
{
  path: '/ios26',
  component: () => import('@/ios26/components/IosAppShell.vue'),
  children: [
    {
      path: '',
      name: 'ios26-home',
      component: () => import('@/ios26/views/IosHomeView.vue'),
    },
    {
      path: 'prompts',
      name: 'ios26-prompts',
      component: () => import('@/ios26/views/IosPromptsView.vue'),
    },
    {
      path: 'prompts/:slug',
      name: 'ios26-prompt-detail',
      component: () => import('@/ios26/views/IosPromptDetailView.vue'),
    },
    {
      path: 'generator',
      name: 'ios26-generator',
      component: () => import('@/ios26/views/IosGeneratorView.vue'),
    },
    {
      path: 'generator/history',
      name: 'ios26-history',
      component: () => import('@/ios26/views/IosHistoryView.vue'),
    },
    {
      path: 'dashboard',
      name: 'ios26-dashboard',
      component: () => import('@/ios26/views/IosDashboardView.vue'),
    },
    {
      path: 'login',
      name: 'ios26-login',
      component: () => import('@/ios26/views/IosLoginView.vue'),
    },
    {
      path: 'pricing',
      name: 'ios26-pricing',
      component: () => import('@/ios26/views/IosPricingView.vue'),
    },
  ],
}
```

## 样式隔离策略

新版样式必须挂在 `.ios26-app` 根类下，避免影响旧版页面。

推荐入口：

```css
/* web/src/ios26/styles/index.css */
@import './tokens.css';
@import './typography.css';
@import './materials.css';
@import './animations.css';
```

推荐作用域：

```css
.ios26-app {
  min-height: 100dvh;
  font-family: "SF Pro", -apple-system, BlinkMacSystemFont, "Noto Sans SC", system-ui, sans-serif;
  background: var(--ios-color-bg-grouped-primary);
  color: var(--ios-color-label-primary);
}
```

所有新版变量使用 `--ios-*` 前缀，例如：

```css
.ios26-app {
  --ios-color-label-primary: #000;
  --ios-color-label-secondary: rgba(60, 60, 67, 0.6);
  --ios-color-bg-primary: #fff;
  --ios-color-bg-grouped-primary: #f2f2f7;
  --ios-color-separator: rgba(60, 60, 67, 0.29);
  --ios-radius-button: 12px;
  --ios-radius-sheet: 34px;
  --ios-glass-blur-md: 12px;
  --ios-glass-shadow-layer: 0 18px 40px rgba(0, 0, 0, 0.12);
}
```

暗色模式通过 media query 或 `data-theme` 扩展：

```css
@media (prefers-color-scheme: dark) {
  .ios26-app {
    --ios-color-label-primary: #fff;
    --ios-color-label-secondary: rgba(235, 235, 245, 0.6);
    --ios-color-bg-primary: #000;
    --ios-color-bg-grouped-primary: #1c1c1e;
    --ios-color-separator: rgba(84, 84, 88, 0.65);
  }
}
```

## iOS 26 设计规范落点

新版应遵循 `ios26-design` skill 中的核心规范：

- 字体：优先使用 SF Pro 字体栈，中文回退到 `Noto Sans SC` 或系统中文字体
- 布局：使用 8pt spacing grid
- 触控：按钮、tab、列表行最小触控区域不低于 44px
- 材质：使用 Liquid Glass 的半透明背景、模糊、边框高光和分层阴影
- 动效：使用 spring-like cubic-bezier，不使用线性 easing
- 暗色：默认支持 light / dark
- 颜色：优先语义 token，不在组件中散落硬编码颜色

推荐基础 token：

```css
.ios26-app {
  --ios-space-1: 4px;
  --ios-space-2: 8px;
  --ios-space-3: 12px;
  --ios-space-4: 16px;
  --ios-space-5: 20px;
  --ios-space-6: 24px;
  --ios-space-8: 32px;

  --ios-duration-tab: 0.3s;
  --ios-duration-sheet-present: 0.5s;
  --ios-ease-snappy: cubic-bezier(0.34, 1.56, 0.64, 1);
  --ios-ease-gentle: cubic-bezier(0.25, 0.46, 0.45, 0.94);
  --ios-ease-stiff: cubic-bezier(0.25, 0.1, 0.25, 1);
}
```

## 核心组件规划

### IosAppShell

新版应用壳，负责整体布局。

职责：

- 注入 `.ios26-app`
- 引入 `ios26/styles/index.css`
- 渲染顶部 `IosNavBar`
- 渲染底部 `IosTabBar`
- 处理 safe area padding
- 承载子路由页面

### IosNavBar

iOS 风格顶部导航。

规格：

- 常规高度 44px
- 页面大标题场景可扩展到 96px
- 横向 padding 16px
- 支持返回按钮、标题、右侧操作按钮
- 使用 glass 背景时需要保持内容可读性

### IosTabBar

底部 Liquid Glass tab bar。

规格：

- 常规高度 49px
- 含安全区时总高度约 83px
- 使用 fixed bottom
- 支持当前 tab 高亮
- 按钮触控区域不低于 44px

### IosGlassPanel

Liquid Glass 基础容器。

职责：

- 统一 blur、透明背景、边框、高光、阴影
- 支持 `size="sm" | "md" | "lg"`
- 支持 light / dark 自动适配

### IosButton

新版按钮组件。

规格：

- regular 高度 44px
- 默认圆角 12px
- 支持 filled、tinted、plain、glass 等 variant
- 使用 snappy spring 动效

### IosSheet

底部弹层组件。

规格：

- 顶部圆角 34px
- grabber 尺寸 36px x 5px
- 支持 25%、50%、100% detent 的后续扩展
- present 动效 0.5s gentle spring
- dismiss 动效 0.3s

### IosListRow

iOS 风格列表行。

规格：

- 常规高度 44px
- 小尺寸 36px
- 大尺寸 58px
- separator 左侧 inset 16px
- 支持 label、value、icon、disclosure

## 页面改造规划

### 第一批页面

优先开发能体现产品核心价值和新风格差异的页面：

```text
/ios26
/ios26/generator
/ios26/prompts
/ios26/prompts/:slug
```

目标：

- 建立整体视觉语言
- 验证 Liquid Glass 是否适合当前产品
- 验证移动端和桌面端布局
- 复用现有 store、api、types，降低业务风险

### 第二批页面

继续迁移账户、历史和商业化页面：

```text
/ios26/generator/history
/ios26/dashboard
/ios26/login
/ios26/pricing
```

目标：

- 补齐完整用户路径
- 验证表单、订阅卡片、个人中心等复杂组件
- 为后续灰度或正式切换做准备

## 数据和业务逻辑复用

新版前端只重做视觉和交互层，业务数据层尽量复用现有模块：

```text
web/src/api/
web/src/store/
web/src/types/
web/src/mock/
```

建议：

- API 请求不复制
- Pinia store 不复制
- 类型定义不复制
- mock 数据不复制
- 只新增 iOS26 风格的 view 和 component

这样可以避免新旧两套业务逻辑分叉。

## 开发阶段

### Phase 1: 基础设施

新增：

```text
web/src/ios26/styles/
web/src/ios26/components/
web/src/ios26/views/
web/src/ios26/router.ts
web/src/ios26/constants.ts
```

完成：

- iOS26 token
- typography
- Liquid Glass material
- animation variables
- `/ios26` 路由组
- `IosAppShell`

### Phase 2: 导航和首页

完成：

- `IosNavBar`
- `IosTabBar`
- `IosButton`
- `IosGlassPanel`
- `IosHomeView`

验收：

- `/ios26` 可以独立访问
- 旧版 `/` 无变化
- 亮色和暗色模式基础可用
- 移动端底部 tab 不遮挡内容

### Phase 3: 核心功能页

完成：

- `IosGeneratorView`
- `IosPromptsView`
- `IosPromptDetailView`
- 必要的 prompt card、filter、input、result panel 等新版组件

验收：

- 生成器核心流程可用
- 提示词列表可浏览
- 提示词详情可访问
- API / store 仍复用现有模块

### Phase 4: 用户路径补齐

完成：

- `IosHistoryView`
- `IosDashboardView`
- `IosLoginView`
- `IosPricingView`
- `IosSheet`
- `IosListRow`
- `IosSegmentedControl`

验收：

- 新版路径覆盖旧版主要用户旅程
- 登录、定价、历史、仪表盘页面视觉统一
- 表单和按钮满足 44px 触控要求

### Phase 5: 预览、灰度和切换

保留访问方式：

```text
/                 旧版正式入口
/ios26            新版预览入口
```

可选扩展：

```text
?theme=ios26      后续用于跳转或灰度
localStorage      后续用于记住用户选择
```

新版稳定后，再决定是否将 `/` 指向新版。

## 验收标准

功能验收：

- 旧版所有路径仍可访问
- 新版 `/ios26` 路由组可访问
- 新版核心页面可以复用现有 API 和 store
- 构建命令 `npm run build` 通过

样式验收：

- 旧版页面视觉无变化
- 新版样式不污染旧版页面
- 新版支持 light / dark
- 新版使用 Liquid Glass 材质
- 按钮、tab、列表行触控区域不低于 44px
- 页面在移动端和桌面端无明显重叠、溢出、遮挡

代码验收：

- 新版文件集中在 `web/src/ios26`
- 不复制 `api`、`store`、`types` 业务层
- 新版组件命名统一使用 `Ios` 前缀
- CSS 变量统一使用 `--ios-*` 前缀
- 组件内避免散落硬编码颜色

## 风险和控制

### 风险：新版样式影响旧版

控制方式：

- 所有新版样式挂在 `.ios26-app` 下
- CSS 变量使用 `--ios-*` 前缀
- 不修改旧版 `main.css` 中的现有变量

### 风险：新旧业务逻辑分叉

控制方式：

- 复用 `api`、`store`、`types`
- 只新增 view 和 component
- 页面逻辑尽量调用已有 store action

### 风险：Liquid Glass 可读性不足

控制方式：

- glass 容器必须有足够背景透明度
- 文本使用语义 label token
- 暗色模式单独调试
- 对复杂内容区域减少过强透明效果

### 风险：移动端底部 tab 遮挡内容

控制方式：

- `IosAppShell` 为主内容添加底部 safe area padding
- tab bar 高度按 83px 预留
- 使用 `env(safe-area-inset-bottom)` 兼容设备

## 推荐下一步

建议先实施 Phase 1 和 Phase 2，完成一个可访问的新风格入口：

```text
/ios26
```

随后再进入核心功能页：

```text
/ios26/generator
/ios26/prompts
/ios26/prompts/:slug
```

这样可以在不干扰旧版的前提下，快速验证 iOS 26 Liquid Glass 新风格是否适合 DevPrompt.AI。
