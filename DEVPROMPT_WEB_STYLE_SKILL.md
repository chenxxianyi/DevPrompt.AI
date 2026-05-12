---
name: devprompt-web-style
description: Use when creating or modifying files under web/ for DevPrompt AI and the UI must match the existing frontend style: dark glassmorphism, purple-cyan developer aesthetic, Vue 3 + Tailwind patterns, Chinese product copy, and the current component/layout conventions.
---

# DevPrompt Web Style

Use this skill whenever work touches `web/` and the output must visually and structurally match the existing DevPrompt AI frontend.

## Goal

Preserve the current design language instead of introducing a new one.

This frontend is:
- Dark-first
- Glassmorphism-heavy
- Purple + cyan as primary accents
- Developer-product oriented, not consumer-playful
- Dense but polished, with strong hierarchy and compact actions
- Built with Vue 3 SFCs, Tailwind utilities, and shared CSS tokens from `web/src/assets/main.css`

## Source Of Truth

Read these files first when making style decisions:
- `web/src/assets/main.css`
- `web/src/App.vue`
- `web/src/components/NavBar.vue`
- `web/src/components/PromptCard.vue`
- `web/src/views/HomeView.vue`

If a new page is being added, also inspect a nearby existing page in `web/src/views/` before designing.

## Visual Rules

### 1. Color system

Do not invent a new palette unless explicitly asked.

Prefer existing CSS variables:
- `--bg-deep`
- `--bg-base`
- `--bg-surface`
- `--bg-elevated`
- `--bg-card`
- `--border`
- `--border-active`
- `--text-primary`
- `--text-secondary`
- `--text-muted`
- `--accent`
- `--accent-hover`
- `--accent-glow`
- `--cyan`
- `--cyan-glow`
- `--green`
- `--amber`
- `--rose`

Use them via inline style or existing utility patterns, for example:

```vue
<div class="glass p-6" style="border-color:var(--border)">
  <p style="color:var(--text-secondary)">说明文案</p>
</div>
```

### 2. Surfaces

Primary containers should usually use `.glass`.

Use this pattern for cards, panels, dialogs, metric blocks, and content modules:
- rounded corners: `rounded-xl`, `rounded-2xl`, or inherited from `.glass`
- semi-transparent dark surface
- subtle purple border
- hover border enhancement, not exaggerated motion

Do not switch to flat white cards or bright enterprise dashboard styling.

### 3. Typography

Typography should feel technical and decisive.

Use:
- `'Noto Sans SC'` for UI and body
- `'JetBrains Mono'` for code, prompts, badges that feel technical, or compact brand accents

Common hierarchy already used in the project:
- page title: `text-[28px]` to `text-[36px]`, `font-extrabold` or `font-black`
- section title: `text-[22px]`, `font-bold`
- card title: `text-[15px]` to `text-[17px]`, `font-semibold` or `font-bold`
- body copy: `text-sm` or `text-base`
- helper copy: `text-xs` to `text-[13px]`

Avoid overly light text on dark backgrounds. Default copy should use `--text-primary` or `--text-secondary`.

### 4. Buttons and inputs

Reuse existing component classes whenever possible:
- `.btn`
- `.btn-primary`
- `.btn-ghost`
- `.btn-cyan`
- `.form-input`
- `.tag`
- `.tag-cyan`

Preferred interaction style:
- subtle lift on hover
- glow on primary actions
- borders on secondary actions
- no giant bounce, no noisy transforms

### 5. Motion

Motion should be restrained and meaningful.

Use:
- `transition-all duration-300`
- `transition-transform`
- mild hover translate like `hover:-translate-y-0.5`
- page transitions similar to current `App.vue`

Avoid:
- large spring animations
- rotating decorative UI
- scale-heavy hover states on layout containers

## Layout Rules

### 1. Page shell

Most pages should follow this structure:

```vue
<div class="relative z-[1] min-h-screen">
  <div class="max-w-screen-xl mx-auto px-6">
    <!-- page content -->
  </div>
</div>
```

Top spacing usually assumes the shared navbar already exists.

### 2. Responsive behavior

Default to desktop-first productivity layout, but ensure mobile still works.

Common grid patterns in this codebase:
- `grid-cols-[repeat(auto-fit,minmax(260px,1fr))]`
- `grid-cols-[repeat(auto-fill,minmax(300px,1fr))]`
- `grid grid-cols-1 md:grid-cols-2`
- `flex flex-wrap gap-*`

On mobile:
- stack panels vertically
- keep primary CTA visible
- avoid dense horizontal toolbars without wrapping

### 3. Information density

This product is for developers, so moderate density is good.

Prefer:
- compact controls
- clear grouping
- visible metadata
- concise labels

Avoid:
- overly spacious marketing-only layout
- giant empty hero sections on internal pages
- card designs with too little information

## Copy And Product Tone

The copy style should match a Chinese developer tool product:
- concise
- direct
- functional
- mildly premium

Good examples:
- `立即生成`
- `浏览模板库`
- `申请试用`
- `查看会员方案`
- `管理你的 Prompt 资产和生成记录`

Avoid:
- exaggerated marketing hype
- trendy consumer-social phrases
- fluffy AI jargon with no task meaning

## Vue/Tailwind Implementation Rules

### 1. Follow existing SFC style

Prefer:
- `<script setup lang="ts">`
- Tailwind utilities in template
- scoped styles only for page-specific sections
- CSS variables for custom colors

### 2. Do not over-componentize small one-off sections

If a pattern is clearly local to one page, keep it in that page.

Extract a component only when:
- reused in 2+ places
- structurally meaningful
- likely to stay stable

### 3. Reuse store and router patterns

When adding a page:
- match naming style in `web/src/router/index.ts`
- use existing stores before creating new ones
- prefer current API wrappers under `web/src/api/`

## Common UI Patterns To Reuse

### Panel

```vue
<div class="glass p-6">
  <h3 class="text-sm font-semibold mb-4" style="color:var(--text-secondary)">标题</h3>
  <div style="color:var(--text-primary)">内容</div>
</div>
```

### Card grid

```vue
<div class="grid grid-cols-[repeat(auto-fill,minmax(300px,1fr))] gap-4">
  <!-- cards -->
</div>
```

### Section header

```vue
<div class="flex items-center justify-between mb-6">
  <h2 class="text-[22px] font-bold">标题</h2>
  <router-link to="/target" class="text-sm no-underline" style="color:var(--accent)">查看全部</router-link>
</div>
```

### Form block

```vue
<div class="flex flex-col gap-1.5">
  <label class="text-[13px] font-semibold" style="color:var(--text-secondary)">字段名</label>
  <input class="form-input" type="text" placeholder="请输入..." />
</div>
```

## Anti-Patterns

Do not introduce these unless explicitly requested:
- light theme redesign
- neumorphism
- shadcn-style white/gray SaaS visuals
- excessive rounded blobs
- pastel consumer colors
- default browser blue
- mixed icon systems with inconsistent stroke weights
- random spacing scales unrelated to current pages
- new font stacks that conflict with `Noto Sans SC` and `JetBrains Mono`

## Task Checklist

When editing `web/`, check:
1. Does the page still look like the same product as `HomeView.vue` and `PromptCard.vue`?
2. Are existing CSS variables reused instead of hardcoding unrelated colors?
3. Are actions using `.btn-*` patterns where possible?
4. Is the layout wrapped in `max-w-screen-xl mx-auto px-6` unless there is a reason not to?
5. Does the copy sound like a Chinese developer productivity tool?
6. Is mobile still usable?

## Output Expectation

When asked to build or restyle a `web/` page, produce UI that feels:
- dark
- sharp
- premium
- technical
- calm
- consistent with the existing DevPrompt AI brand

Do not redesign the system from scratch. Extend the current language.
