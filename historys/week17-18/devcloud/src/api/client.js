import axios from 'axios'

// http client
const client = axios.create({
    // API请求的base URL
    baseURL: "http://localhost:8050",
    // 超时时间
    timeout: 5000,
})

// request中间件
client.interceptors.request.use(
    // 成功的处理逻辑
    request => {
        return request
    },
    // 错误时的处理逻辑
    err => {
        console.log(err)
        return Promise.reject(err)
    }
)

// response中间件
client.interceptors.response.use(
    response => {
        const resp = response.data
        // 判断返回的error code是否为0, 如果为0请求成功
        if (resp.code === 0) {
            return resp
        }

        // 如果不为0, 请求失败
        console.log(resp)
    },
    err => {
        // 错误时的处理逻辑
        console.log(err)
        return Promise.reject(err)
    }
)

export default client