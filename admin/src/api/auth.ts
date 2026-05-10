import request from './request'
import type { ApiResponse, LoginInput, LoginResult, User } from '@/types'

export function login(data: LoginInput) {
  return request.post<ApiResponse<LoginResult>>('/auth/login', data)
}

export function getProfile() {
  return request.get<ApiResponse<User>>('/auth/profile')
}
