# iOS26 Frontend Optimization Review

## Scope

This document summarizes the concrete fixes from the latest code review for the iOS26 frontend changes. The focus is on behavior, accessibility, and maintainability risks that can affect real users.

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

### 4. Recalculate top-nav indicator on layout changes

Severity: Low

File: `web/src/ios26/components/IosTabBar.vue`

Problem:

The active indicator is recalculated only after mount and active route changes. Because the top navigation is now width-dependent, resizing the window or crossing responsive breakpoints can leave the indicator misaligned.

Recommended fix:

```ts
import { computed, ref, watch, nextTick, onMounted, onUnmounted } from 'vue'

onMounted(() => {
  updateIndicator()
  window.addEventListener('resize', updateIndicator)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateIndicator)
})
```

Optional improvement:

Use `ResizeObserver` on `trackRef` for more precise recalculation when fonts or content dimensions change without a window resize.

Acceptance criteria:

- Indicator stays aligned after desktop resize.
- Indicator stays aligned when switching between desktop and mobile widths.
- No leaked global resize listener after component unmount.

## Implementation Order

1. Restore `system` as the default theme and remove automatic `light` persistence.
2. Add ARIA state and keyboard handling to `IosSelect`.
3. Add viewport-aware dropdown positioning to `IosSelect`.
4. Add resize recalculation for `IosTabBar` indicator.
5. Run `npm run build` from `web/`.
6. Manually verify `/ios26/generator` on desktop and mobile widths.

## Regression Checklist

- `web` build passes.
- `/ios26` home tab is active on `/ios26` and `/ios26/`.
- `/ios26/generator` required fields still enable the generate button correctly.
- Project type and code style selects work with mouse and keyboard.
- Dropdown closes on outside click, `Escape`, option selection, scroll, and route changes.
- Dark mode follows the OS when no explicit theme has been saved.
- Saved theme preferences still survive refresh.

