<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import IosIcon from './IosIcon.vue'

interface Option {
  value: string
  label: string
}

const props = withDefaults(defineProps<{
  modelValue: string
  options: Option[]
  placeholder?: string
}>(), {
  placeholder: '请选择',
})

const emit = defineEmits<{
  (e: 'update:modelValue', val: string): void
}>()

const listboxId = `ios-select-${Math.random().toString(36).slice(2)}`
const open = ref(false)
const activeIndex = ref(-1)
const triggerRef = ref<HTMLButtonElement | null>(null)
const dropdownRef = ref<HTMLDivElement | null>(null)

const dropStyle = ref<Record<string, string>>({ top: '0px', left: '0px', width: '0px', maxHeight: '260px' })

const selectedLabel = computed(() => {
  const match = props.options.find(o => o.value === props.modelValue)
  return match ? match.label : ''
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
    top:       `${(openUp ? rect.top - availableHeight - gap : rect.bottom + gap) + window.scrollY}px`,
    left:      `${rect.left + window.scrollX}px`,
    width:     `${rect.width}px`,
    maxHeight: `${availableHeight}px`,
  }
}

function openDropdown() {
  calcPosition()
  open.value = true
  activeIndex.value = Math.max(0, props.options.findIndex(o => o.value === props.modelValue))
}

function closeDropdown() {
  open.value = false
}

function toggle() {
  open.value ? closeDropdown() : openDropdown()
}

function moveActive(delta: number) {
  if (!open.value) { openDropdown(); return }
  const count = props.options.length
  if (!count) return
  activeIndex.value = (activeIndex.value + delta + count) % count
}

function selectActive() {
  const option = props.options[activeIndex.value]
  if (option) select(option.value)
}

function select(val: string) {
  emit('update:modelValue', val)
  closeDropdown()
}

function onOutsideClick(e: MouseEvent) {
  const t = triggerRef.value
  const d = dropdownRef.value
  if (t && !t.contains(e.target as Node) && d && !d.contains(e.target as Node)) {
    closeDropdown()
  }
}

function onScroll() {
  if (open.value) calcPosition()
}

onMounted(() => {
  document.addEventListener('mousedown', onOutsideClick)
  window.addEventListener('scroll', onScroll, true)
  window.addEventListener('resize', onScroll)
})
onUnmounted(() => {
  document.removeEventListener('mousedown', onOutsideClick)
  window.removeEventListener('scroll', onScroll, true)
  window.removeEventListener('resize', onScroll)
})
</script>

<template>
  <div class="ios-select" :class="{ 'is-open': open }">
    <!-- Trigger -->
    <button
      ref="triggerRef"
      type="button"
      class="ios-select__trigger ios-input"
      :class="{ 'is-placeholder': !modelValue }"
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
      <span class="ios-select__value">
        {{ selectedLabel || placeholder }}
      </span>
      <span class="ios-select__arrow" :class="{ 'is-open': open }">
        <IosIcon
          path="M6 9l6 6 6-6"
          :size="16"
          :stroke="2"
        />
      </span>
    </button>

    <!-- Dropdown — teleport to body to escape overflow:hidden -->
    <Teleport to="body">
      <Transition name="ios-select-drop">
        <div
          v-if="open"
          ref="dropdownRef"
          class="ios-select__dropdown"
          :style="dropStyle"
        >
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
              <span class="ios-select__option-label">{{ opt.label }}</span>
              <span v-if="modelValue === opt.value" class="ios-select__check">
                <IosIcon path="M20 6 9 17l-5-5" :size="15" :stroke="2.5" />
              </span>
            </li>
          </ul>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.ios-select {
  position: relative;
  width: 100%;
}

/* ── Trigger ── */
.ios-select__trigger {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  cursor: pointer;
  text-align: left;
  padding: 11px 14px;
  height: auto;
  transition:
    border-color 0.2s var(--ios-ease-stiff),
    background  0.2s var(--ios-ease-stiff),
    box-shadow  0.2s var(--ios-ease-stiff);
}

.ios-select__trigger.is-placeholder .ios-select__value {
  color: var(--ios-color-label-tertiary);
}

.ios-select.is-open .ios-select__trigger {
  border-color: var(--ios-color-tint);
  box-shadow: 0 0 0 4px var(--ios-color-tint-soft);
  background: var(--ios-color-bg-primary);
}

.ios-select__value {
  font-size: 15px;
  color: var(--ios-color-label-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.ios-select__arrow {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  color: var(--ios-color-label-tertiary);
  transition: transform 0.28s cubic-bezier(0.34, 1.56, 0.64, 1),
              color    0.18s var(--ios-ease-stiff);
}

.ios-select__arrow.is-open {
  transform: rotate(180deg);
  color: var(--ios-color-tint);
}

</style>

<!-- Non-scoped: styles for teleported dropdown (outside component DOM tree) -->
<style>
.ios-select__dropdown {
  position: absolute;
  z-index: 9999;
  background: rgba(248, 252, 255, 0.97);
  backdrop-filter: blur(24px) saturate(2.2);
  -webkit-backdrop-filter: blur(24px) saturate(2.2);
  border: 0.5px solid rgba(125, 211, 252, 0.40);
  border-radius: 16px;
  box-shadow:
    0 8px 32px rgba(63, 164, 232, 0.18),
    0 2px 8px  rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.95);
  overflow: hidden;
  padding: 5px;
}

.ios-select__list {
  list-style: none;
  margin: 0;
  padding: 0;
  max-height: inherit;
  overflow-y: auto;
}

.ios-select__list::-webkit-scrollbar { width: 4px; }
.ios-select__list::-webkit-scrollbar-thumb {
  background: var(--ios-color-fill-tertiary);
  border-radius: 2px;
}

.ios-select__option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  border-radius: 11px;
  cursor: pointer;
  font-size: 15px;
  color: #1c1c1e;
  letter-spacing: -0.15px;
  transition: background 0.15s ease;
  user-select: none;
}

.ios-select__option:hover,
.ios-select__option.is-active {
  background: rgba(63, 164, 232, 0.12);
  color: #0a0a0a;
}

.ios-select__option.is-selected {
  color: #3fa4e8;
  font-weight: 600;
  background: rgba(63, 164, 232, 0.12);
}

.ios-select__option-label { flex: 1; }

.ios-select__check {
  color: #3fa4e8;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.ios-select-drop-enter-active {
  transition: opacity 0.2s var(--ios-ease-stiff),
              transform 0.28s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.ios-select-drop-leave-active {
  transition: opacity 0.16s var(--ios-ease-stiff),
              transform 0.18s var(--ios-ease-stiff);
}
.ios-select-drop-enter-from,
.ios-select-drop-leave-to {
  opacity: 0;
  transform: translateY(-6px) scale(0.97);
}
</style>
