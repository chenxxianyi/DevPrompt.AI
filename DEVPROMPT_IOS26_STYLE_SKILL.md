---
name: devprompt-ios26-style
description: Use when creating or modifying files under web/src/ios26/ for DevPrompt AI and the UI must match the existing iOS26 variant: Liquid Glass surfaces, macaron blue aurora palette, Apple-inspired spacing and typography, custom Ios* Vue components, and the current iOS26 route/layout conventions.
---

# DevPrompt iOS26 Style

Use this skill whenever work touches `web/src/ios26/` and the output must match the existing iOS26 frontend branch of DevPrompt AI.

This is not the same style as the main site under `web/src/views/`.
Do not mix the dark-purple main-site language into the iOS26 variant unless the existing iOS26 files already do so.

## Goal

Extend the current DevPrompt iOS26 implementation instead of creating a generic Apple clone.

This branch is:
- Apple-inspired, but customized
- Liquid Glass based
- Light, airy, translucent
- Built around a macaron blue palette
- Powered by custom `Ios*` Vue components
- More tactile and softer than the main desktop site

## Source Of Truth

Read these files first before making style decisions:
- `web/src/ios26/styles/tokens.css`
- `web/src/ios26/styles/materials.css`
- `web/src/ios26/styles/typography.css`
- `web/src/ios26/components/IosAppShell.vue`
- `web/src/ios26/components/IosButton.vue`
- `web/src/ios26/components/IosGlassPanel.vue`
- `web/src/ios26/views/IosHomeView.vue`

If the task is page-specific, also inspect a nearby page in `web/src/ios26/views/`.

## Core Identity

This iOS26 branch uses a specific visual direction:
- Macaron blue / cornflower blue / lavender tinting
- Frosted glass panels with strong highlights
- Soft aurora background fields
- SF-style hierarchy adapted to web
- High polish, low visual noise
- Mobile-first interaction assumptions, but still works on desktop

Do not reinterpret it as:
- the dark glass style from the main site
- a flat white admin panel
- a generic Tailwind SaaS landing page
- a pure unmodified Apple system clone

## Color Rules

Always prefer `--ios-*` variables from `tokens.css`.

Frequently used variables:
- `--ios-color-label-primary`
- `--ios-color-label-secondary`
- `--ios-color-label-tertiary`
- `--ios-color-fill-primary`
- `--ios-color-fill-secondary`
- `--ios-color-bg-primary`
- `--ios-color-bg-secondary`
- `--ios-color-bg-grouped-primary`
- `--ios-color-separator`
- `--ios-color-tint`
- `--ios-color-tint-hover`
- `--ios-color-tint-pressed`
- `--ios-color-tint-soft`
- `--ios-color-tint-soft-strong`
- `--ios-color-tint-gradient`
- `--ios-color-systemBlue`
- `--ios-color-systemGreen`
- `--ios-color-systemRed`
- `--ios-color-systemOrange`
- `--ios-color-systemPink`
- `--ios-color-systemTeal`
- `--ios-color-systemIndigo`

Use these instead of hardcoded colors whenever possible.

## Surface Rules

### 1. Prefer Ios components first

Before writing custom markup, check whether one of these should be used:
- `IosGlassPanel`
- `IosButton`
- `IosIcon`
- `IosSheet`
- `IosListRow`
- `IosSegmentedControl`
- `IosSelect`
- `IosToast`
- `IosTabBar`
- `IosNavBar`

If the UI matches one of these primitives, reuse it instead of rebuilding.

### 2. Glass over flat

The default surface style is not a flat card.

Prefer:
- `.ios-glass`
- `.ios-glass-sm`
- `.ios-glass-lg`
- `IosGlassPanel`

Only use solid grouped surfaces when the page already uses them for structure.

### 3. Highlights matter

Liquid Glass in this codebase includes:
- blur
- saturation
- border
- top highlight sheen
- soft layered shadow

Do not create a “glass” card with only opacity and blur. It will look off-brand.

## Typography Rules

Use the iOS26 typography system, not the main site typography.

Prefer the project’s iOS text classes where available, such as:
- `ios-text-large-title`
- `ios-text-title-1`
- `ios-text-title-2`
- `ios-text-title-3`
- `ios-text-headline`
- `ios-text-body`
- `ios-text-callout`
- `ios-text-subheadline`
- `ios-text-footnote`
- `ios-text-caption-1`
- `ios-text-caption-2`

Typography should feel:
- clean
- calm
- compact
- touch-native

Avoid heavy desktop-marketing headings unless the existing page already uses them.

## Layout Rules

### 1. Shell ownership

`IosAppShell.vue` owns the outer shell:
- aurora background
- page transition
- mobile bottom tab bar
- top padding
- theme switching scope

Do not recreate global shell behavior inside individual pages.

### 2. Page wrapper

Most pages should fit inside the existing iOS page structure, usually something like:

```vue
<div class="ios-page">
  <!-- sections -->
</div>
```

Keep spacing aligned with the existing token system:
- `--ios-page-padding`
- `--ios-web-page-padding`
- `--ios-space-*`

### 3. Mobile-first behavior

This branch is more mobile-native than the main site.

Design decisions should assume:
- touch targets at least `44px`
- stacked content on narrow screens
- compact but legible controls
- bottom tab bar present on mobile

Do not introduce dense desktop-only controls without a mobile fallback.

## Interaction Rules

### 1. Reuse the press language

If an element is tappable, prefer matching existing conventions like:
- `ios-press`
- `IosButton`
- subtle color change
- gentle spring-like transitions

Do not use aggressive hover transforms or exaggerated desktop card motion.

### 2. Animation tone

Use the existing iOS easing vocabulary:
- `--ios-ease-snappy`
- `--ios-ease-bouncy`
- `--ios-ease-gentle`
- `--ios-ease-stiff`

Use durations from tokens when possible.

Preferred motion:
- page transitions
- sheet transitions
- tab indicator motion
- soft button feedback

Avoid:
- generic `linear`
- oversized hover scaling
- dramatic bounce for ordinary UI

## Component Usage Patterns

### 1. Primary action

Use `IosButton` first.

Typical variants:
- `filled` for primary CTA
- `tinted` for secondary emphasis
- `glass` for subtle surfaced actions
- `plain` for lightweight inline actions
- `destructive` only for destructive intent

### 2. Cards and panels

Use `IosGlassPanel` for:
- feature cards
- quick-input panels
- CTA blocks
- lightweight grouped modules

### 3. Icons

Use `IosIcon` rather than mixing ad hoc icon libraries.

Stroke weights and sizes should stay visually consistent with existing pages.

## Copy Style

The copy tone should feel like a polished Chinese productivity app:
- concise
- refined
- clear
- slightly premium

Examples that fit:
- `立即生成`
- `查看全部`
- `申请试用`
- `选择适合你的方案`
- `补充一些信息，方便我们尽快联系并为你开通`

Avoid:
- noisy marketing slogans
- meme-like AI copy
- rigid backend/admin wording on end-user pages

## Implementation Rules

### 1. Prefer composition over raw CSS invention

Before adding custom CSS, ask:
- can this be built with existing `Ios*` components?
- can this be styled through existing `--ios-*` tokens?
- can this follow a pattern already used in `IosHomeView.vue` or similar pages?

### 2. Scoped styles are fine, but only for local structure

Use scoped styles for:
- page-specific spacing
- local arrangement
- custom sub-blocks

Do not redefine:
- global theme tokens
- core glass behavior
- shared animation primitives

### 3. Keep the branch self-consistent

When editing `web/src/ios26/`, do not pull in:
- main-site `glass`, `btn`, `form-input`, `tag`
- main-site text hierarchy
- main-site navbar/footer assumptions

This branch has its own system. Use it.

## Common Patterns To Reuse

### Glass section

```vue
<IosGlassPanel class="my-panel">
  <h3 class="ios-text-title-3">标题</h3>
  <p class="ios-text-subheadline">说明文案</p>
</IosGlassPanel>
```

### Section header

```vue
<header class="ios-section__header">
  <h3 class="ios-text-title-3">标题</h3>
  <router-link to="/ios26/target" class="ios-section__more">查看全部</router-link>
</header>
```

### Primary CTA

```vue
<IosButton variant="filled" size="lg">立即操作</IosButton>
```

### Input zone

Use glass container + native-feeling input/textarea styling, similar to `IosHomeView.vue`.

## Anti-Patterns

Do not introduce these unless explicitly requested:
- main-site dark purple panel styling inside `ios26`
- flat white cards with gray borders
- generic Tailwind button sets instead of `IosButton`
- arbitrary spacing unrelated to `--ios-space-*`
- hardcoded random gradients that ignore the existing macaron palette
- mixed icon languages and inconsistent stroke widths
- oversized desktop-only hero compositions

## Task Checklist

When editing `web/src/ios26/`, check:
1. Does this look like the same product family as `IosHomeView.vue`?
2. Are `Ios*` components reused before creating raw replacements?
3. Are `--ios-*` tokens used instead of unrelated custom colors?
4. Are touch targets and mobile layout still solid?
5. Does the motion feel gentle and tactile rather than flashy?
6. Is this clearly separate from the main-site style system?

## Output Expectation

When asked to build or restyle an iOS26 page, produce UI that feels:
- translucent
- soft
- tactile
- elegant
- Apple-inspired
- consistent with DevPrompt AI’s custom macaron Liquid Glass branch

Do not make it look like the main site. Do not make it look like a generic iOS mockup. Match the existing `web/src/ios26/` implementation.
