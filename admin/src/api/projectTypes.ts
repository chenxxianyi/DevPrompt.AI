import request from './request'
import type { ApiResponse, PaginatedData, ProjectType } from '@/types'

export function getProjectTypes() {
  return request.get<ApiResponse<PaginatedData<ProjectType>>>('/admin/project-types')
}

export function createProjectType(data: Partial<ProjectType>) {
  return request.post<ApiResponse<ProjectType>>('/admin/project-types', data)
}

export function updateProjectType(id: number, data: Partial<ProjectType>) {
  return request.put<ApiResponse<ProjectType>>(`/admin/project-types/${id}`, data)
}

export function deleteProjectType(id: number) {
  return request.delete<ApiResponse<null>>(`/admin/project-types/${id}`)
}
