# iOS26 Frontend Optimization Review

## Scope

本文档汇总 iOS26 前端的代码审查结果和修复记录，按时间倒序追加。

---

## Session 2025-05 — 滚动紧凑导航 & 代码审查修复

本次 session 新增了 `IosTabBar` 滚动紧凑行为，并在随后的代码审查中修复了 7 项问题。

### 新功能：IosTabBar 滚动紧凑模式

- 向下滚动超过 60px：导航胶囊收缩为 `width: 380px`，品牌名、Tab 标签、经典版文字通过 `max-width + opacity` 动画隐藏
- 向上滚动：恢复为 `width: 80vw`
- `IosNavBar`（页面级标题栏）改为 `position: static`，随页面内容滚动，不再吸附

### 已修复（优先级排序）

| 严重度 | 文件 | 问题 | 状态 |
| --- | --- | --- | --- |
| 🔴 High | `IosTabBar.vue:onUnmounted` | `ResizeObserver` 未在卸载时 `disconnect()`，内存泄漏 | ✅ 已修复 |
| 🔴 High | `IosTabBar.vue:~445` | CSS 原生嵌套（`.ios-topnav.is-compact { .ios-topnav__classic-btn {} }`）在 Safari <16.5 等旧版本静默失效 | ✅ 展开为平级选择器 |
| 🔴 Critical | `IosTabBar.vue:watch` | `watch(() => isCompact.value, …)` 声明在 `const isCompact` 之前，触发 temporal dead zone ReferenceError，导致整页白/黑屏 | ✅ 已修复（变量声明提前）|
| 🟠 Medium | `IosTabBar.vue:~53` | 指示器位置仅在 mount 和路由切换时重算，容器宽度过渡期间陈旧 | ✅ 改用 `ResizeObserver` 监听 `trackRef`，每次尺寸变化自动重算 |
| 🟠 Medium | `IosTabBar.vue:167` | `overflow: hidden` 参与 flex 布局计算，在宽度过渡中造成点击区域偏移 | ✅ 改为 `clip-path: inset(0 round 32px)`，仅裁剪视觉，不干扰布局 |
| 🟡 Low | `IosTabBar.vue:56` | `setTimeout(…) as unknown as number` 双重强制类型转换 | ✅ 随 ResizeObserver 重构一并消除 |
| 🟡 Low | `IosTabBar.vue:90` | 品牌区域使用 `div + @click`，无键盘支持、无语义 | ✅ 改为 `<router-link>`，补 `text-decoration: none; color: inherit` |

### 未修复 / 待跟进

- `IosTabBar.vue:tabRefs[i]` — 模板中 ref 自动解包，原写法正确；Code Review 中误判，已回滚
- `IosSelect` 键盘导航 & ARIA 状态仍未实现（见下方 Session 1 第 2 条）

---

## Session 初版 — 原始代码审查修复

## Priority Fixes

### 1. Restore system theme behavior

Severity: Medium

File: `web/src/ios26/components/IosAppShell.vue`

Problem:

The default theme was changed from `system` to `light`, and the app writes `light` to `localStorage` when no saved value exists. This permanently bypasses `prefers-color-scheme: dark` for first-time users and makes the existing `system` theme option ineffective.

Recommended fix:

```ts
const theme = ref<IosTheme>('system')

onMounted(() => {
  auth.initAuth()
  if (auth.isLoggedIn) auth.fetchProfile()

  try {
    const saved = localStorage.getItem(IOS_THEME_STORAGE_KEY) as IosTheme | null
    if (saved === 'light' || saved === 'dark' || saved === 'system') {
      theme.value = saved
    }
  } catch {
    /* noop */
  }
})
```

Acceptance criteria:

- First-time users follow the OS color scheme.
- Existing saved `light`, `dark`, and `system` preferences still work.
- The app does not write a default theme before the user explicitly chooses one.

### 2. Make `IosSelect` keyboard and screen-reader accessible

Severity: Medium

File: `web/src/ios26/components/IosSelect.vue`

Problem:

The custom select replaced a native `<select>`, but it does not provide equivalent keyboard navigation or ARIA state. Keyboard users cannot reliably open, navigate, select, or close the dropdown.

Recommended fix:

```ts
const listboxId = `ios-select-${Math.random().toString(36).slice(2)}`
const activeIndex = ref(-1)

function openDropdown() {
  calcPosition()
  open.value = true
  activeIndex.value = Math.max(
    0,
    props.options.findIndex(o => o.value === props.modelValue),
  )
}

function closeDropdown() {
  open.value = false
}

function moveActive(delta: number) {
  if (!open.value) {
    openDropdown()
    return
  }
  const count = props.options.length
  if (!count) return
  activeIndex.value = (activeIndex.value + delta + count) % count
}

function selectActive() {
  const option = props.options[activeIndex.value]
  if (option) select(option.value)
}
```

```vue
<button
  ref="triggerRef"
  type="button"
  class="ios-select__trigger ios-input"
  role="combobox"
  aria-haspopup="listbox"
  :aria-expanded="open"
  :aria-controls="listboxId"
  :aria-activedescendant="open && activeIndex >= 0 ? `${listboxId}-option-${activeIndex}` : undefined"
  @click="toggle"
  @keydown.down.prevent="moveActive(1)"
  @keydown.up.prevent="moveActive(-1)"
  @keydown.enter.prevent="open ? selectActive() : openDropdown()"
  @keydown.space.prevent="open ? selectActive() : openDropdown()"
  @keydown.escape.prevent="closeDropdown"
>
```

```vue
<ul :id="listboxId" class="ios-select__list" role="listbox">
  <li
    v-for="(opt, index) in options"
    :id="`${listboxId}-option-${index}`"
    :key="opt.value"
    class="ios-select__option"
    :class="{
      'is-selected': modelValue === opt.value,
      'is-active': activeIndex === index,
    }"
    role="option"
    :aria-selected="modelValue === opt.value"
    @mousemove="activeIndex = index"
    @mousedown.prevent="select(opt.value)"
  >
```

Acceptance criteria:

- `Enter` and `Space` open/select.
- `ArrowUp` and `ArrowDown` move through options.
- `Escape` closes the dropdown.
- Screen readers receive expanded state and active option information.

### 3. Keep `IosSelect` dropdown inside the viewport

Severity: Medium

File: `web/src/ios26/components/IosSelect.vue`

Problem:

The dropdown always opens below the trigger using absolute document coordinates. Near the bottom of the viewport, or with the mobile keyboard open, the dropdown can render off-screen.

Recommended fix:

```ts
const dropStyle = ref<Record<string, string>>({
  top: '0px',
  left: '0px',
  width: '0px',
  maxHeight: '260px',
})

function calcPosition() {
  const el = triggerRef.value
  if (!el) return

  const rect = el.getBoundingClientRect()
  const gap = 6
  const preferredMaxHeight = 260
  const spaceBelow = window.innerHeight - rect.bottom
  const spaceAbove = rect.top
  const openUp = spaceBelow < preferredMaxHeight && spaceAbove > spaceBelow
  const availableHeight = Math.max(
    120,
    Math.min(preferredMaxHeight, openUp ? spaceAbove - gap : spaceBelow - gap),
  )

  dropStyle.value = {
    top: `${(openUp ? rect.top - availableHeight - gap : rect.bottom + gap) + window.scrollY}px`,
    left: `${rect.left + window.scrollX}px`,
    width: `${rect.width}px`,
    maxHeight: `${availableHeight}px`,
  }
}
```

```css
.ios-select__list {
  max-height: inherit;
  overflow-y: auto;
}
```

Acceptance criteria:

- Dropdown opens upward when there is not enough space below.
- Dropdown remains usable on mobile and narrow viewports.
- Long option lists scroll inside the dropdown instead of pushing outside the viewport.

### 4. Recalculate top-nav indicator on layout changes ✅ 已实现

Severity: Low → **已修复**

File: `web/src/ios26/components/IosTabBar.vue`

实现方式：使用 `ResizeObserver` 监听 `trackRef` 元素，每次容器尺寸变化自动调用 `updateIndicator()`，并在 `onUnmounted` 中调用 `resizeObserver?.disconnect()` 防止内存泄漏。同时保留 `window.addEventListener('resize', updateIndicator)` 作为兜底。

Acceptance criteria — 验证状态：

- ✅ Indicator stays aligned after desktop resize.
- ✅ Indicator stays aligned when switching compact / expanded widths.
- ✅ No leaked ResizeObserver after component unmount.

## Implementation Order

| 步骤 | 项目 | 状态 |
| --- | --- | --- |
| 1 | 恢复 `system` 为默认主题，移除自动写入 `light` | ✅ 已完成 |
| 2 | `IosSelect` 增加 ARIA 状态与键盘导航 | ⏳ 待实现 |
| 3 | `IosSelect` 下拉菜单视口感知定位 | ⏳ 待实现 |
| 4 | `IosTabBar` 指示器 ResizeObserver 重算 | ✅ 已完成 |
| 5 | `IosTabBar` 滚动紧凑模式 & 相关 Bug 修复 | ✅ 已完成 |
| 6 | `npm run build`（`web/`）验证构建 | — |
| 7 | 手动验证 `/ios26/generator` 桌面端和移动端宽度 | — |

## Regression Checklist

- [ ] `web` build passes.
- [ ] `/ios26` home tab is active on `/ios26` and `/ios26/`.
- [ ] Navigation pill shows compact state on scroll down, restores on scroll up.
- [ ] Indicator position is accurate in both default and compact states.
- [ ] `/ios26/generator` required fields still enable the generate button correctly.
- [ ] Project type and code style selects work with mouse (keyboard nav pending).
- [ ] Dropdown closes on outside click, `Escape`, option selection, scroll, and route changes.
- [ ] Dark mode follows the OS when no explicit theme has been saved.
- [ ] Saved theme preferences still survive refresh.
- [ ] Brand logo in top nav navigates to `/ios26` and is keyboard-focusable.

