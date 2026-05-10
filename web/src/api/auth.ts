import request from './request'
import type { LoginParams, RegisterParams, User } from '@/types'

export function login(data: LoginParams) {
  return request.post<{ data: { token: string; user: User } }>('/auth/login', data)
}

export function register(data: RegisterParams) {
  return request.post<{ data: { token: string; user: User } }>('/auth/register', data)
}

export function getProfile() {
  return request.get<{ data: User }>('/auth/profile')
}
