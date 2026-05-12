import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { PromptTemplate, PromptCategory } from '@/types'
import { getPrompts, getPromptBySlug, likePrompt, favoritePrompt, getFavoritePrompts } from '@/api/prompts'
import { getCategories } from '@/api/categories'

export const usePromptStore = defineStore('prompt', () => {
  const allTemplates = ref<PromptTemplate[]>([])
  const allCategories = ref<PromptCategory[]>([])
  const selectedCategory = ref('全部')
  const searchKeyword = ref('')
  const sortBy = ref('hot')
  const page = ref(1)
  const pageSize = ref(20)
  const total = ref(0)
  const loading = ref(false)
  const favoriteTemplates = ref<PromptTemplate[]>([])
  const favoriteTotal = ref(0)
  const favoriteLoading = ref(false)

  const filteredTemplates = computed(() => allTemplates.value)

  const totalPages = computed(() => Math.ceil(total.value / pageSize.value) || 1)

  async function fetchCategories() {
    const allCat: PromptCategory = { id: 0, name: '全部', slug: '', description: '', sort: 0, status: 'active' }
    try {
      const res = await getCategories()
      allCategories.value = [allCat, ...res.data.data.list]
    } catch {
      allCategories.value = [allCat]
    }
  }

  async function fetchTemplates() {
    loading.value = true
    try {
      const params: Record<string, any> = {
        page: page.value,
        pageSize: pageSize.value,
        sort: sortBy.value,
      }
      if (searchKeyword.value) params.keyword = searchKeyword.value
      if (selectedCategory.value !== '全部') {
        const cat = allCategories.value.find(c => c.name === selectedCategory.value)
        if (cat && cat.slug) params.category = cat.slug
      }
      const res = await getPrompts(params)
      allTemplates.value = res.data.data.list
      total.value = res.data.data.total
    } catch (e: any) {
      console.error('获取模板列表失败:', e)
    } finally {
      loading.value = false
    }
  }

  async function getTemplateBySlug(slug: string): Promise<PromptTemplate | null> {
    // 先从缓存找
    const cached = allTemplates.value.find(t => t.slug === slug)
    if (cached) return cached
    // 再从 API 获取
    try {
      const res = await getPromptBySlug(slug)
      return res.data.data
    } catch {
      return null
    }
  }

  async function toggleLike(id: number) {
    try {
      const res = await likePrompt(id)
      const t = allTemplates.value.find(t => t.id === id)
      if (t) {
        t.isLiked = res.data.data.liked
        t.likeCount += res.data.data.liked ? 1 : -1
      }
      return res.data.data.liked
    } catch {
      return null
    }
  }

  async function toggleFavorite(id: number) {
    try {
      const res = await favoritePrompt(id)
      const t = allTemplates.value.find(t => t.id === id)
      if (t) {
        t.isFavorited = res.data.data.favorited
        t.favoriteCount += res.data.data.favorited ? 1 : -1
      }
      const fav = favoriteTemplates.value.find(t => t.id === id)
      if (fav) {
        fav.isFavorited = res.data.data.favorited
        fav.favoriteCount += res.data.data.favorited ? 1 : -1
      }
      if (!res.data.data.favorited) {
        favoriteTemplates.value = favoriteTemplates.value.filter(t => t.id !== id)
        favoriteTotal.value = Math.max(0, favoriteTotal.value - 1)
      }
      return res.data.data.favorited
    } catch {
      return null
    }
  }

  async function fetchFavoriteTemplates(pageNumber = 1, size = 20) {
    favoriteLoading.value = true
    try {
      const res = await getFavoritePrompts({ page: pageNumber, pageSize: size })
      favoriteTemplates.value = res.data.data.list
      favoriteTotal.value = res.data.data.total
      return res.data.data
    } finally {
      favoriteLoading.value = false
    }
  }

  function setPage(p: number) {
    page.value = p
    fetchTemplates()
  }

  function setSort(s: string) {
    sortBy.value = s
    page.value = 1
    fetchTemplates()
  }

  function setCategory(cat: string) {
    selectedCategory.value = cat
    page.value = 1
    fetchTemplates()
  }

  function setKeyword(kw: string) {
    searchKeyword.value = kw
    page.value = 1
    fetchTemplates()
  }

  return {
    allTemplates, allCategories, selectedCategory, searchKeyword, sortBy,
    page, pageSize, total, totalPages, loading, favoriteTemplates, favoriteTotal, favoriteLoading,
    filteredTemplates, fetchCategories, fetchTemplates,
    getTemplateBySlug, toggleLike, toggleFavorite, fetchFavoriteTemplates,
    setPage, setSort, setCategory, setKeyword,
  }
})
