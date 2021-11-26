import request from '@/api/client'

export function LIST_HOST(params) {
    return request({
        url: '/hosts',
        method: 'get',
        params: params
    })
}