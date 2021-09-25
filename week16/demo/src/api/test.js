import request from '../utils/request'

export function GET_TEST_DATA(id, query) {
    return request({
      url: `/hosts/${id}`,
      method: 'get',
      params: query
    })
  }