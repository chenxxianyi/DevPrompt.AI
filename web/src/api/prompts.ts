import request from './request'
import type { PromptTemplate, PaginatedData } from '@/types'

export function getPrompts(params?: { keyword?: string; category?: string; page?: number; pageSize?: number; sort?: string }) {
  return request.get<{ data: PaginatedData<PromptTemplate> }>('/prompts', { params })
}

export function getPromptBySlug(slug: string) {
  return request.get<{ data: PromptTemplate }>(`/prompts/${slug}`)
}

export function likePrompt(id: number) {
  return request.post(`/prompts/${id}/like`)
}

export function favoritePrompt(id: number) {
  return request.post(`/prompts/${id}/favorite`)
}
