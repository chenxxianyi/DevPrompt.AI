import request from './request'
import type { TrialRequestPayload } from '@/types'

export function createTrialRequest(data: TrialRequestPayload) {
  return request.post('/trial-requests', data)
}
