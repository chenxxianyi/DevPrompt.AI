import request from './request'

export function createTrialRequest(data: { planCode: string; contact?: string; message?: string }) {
  return request.post('/trial-requests', data)
}
