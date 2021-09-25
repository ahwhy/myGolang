import axios from 'axios'

// create an axios instance
const service = axios.create({
    baseURL: 'http://localhost:8050', // url = base url + request url
    // withCredentials: true, // send cookies when cross-domain requests
    timeout: 5000 // request timeout
  })


// request interceptor
service.interceptors.request.use(
    config => {
      return config
    },
    error => {
      // do something with request error
      console.log(error) // for debug
      return Promise.reject(error)
    }
)


// response interceptor
service.interceptors.response.use(
    /**
     * If you want to get http information such as headers or status
     * Please return  response => response
    */
  
    /**
     * Determine the request status by custom code
     * Here is just an example
     * You can also judge the status by HTTP Status Code
     */
    response => {
      const res = response.data
      // if the custom code is not 20000, it is judged as an error.
      if (res.code !== 0) {
        // 比如 token过期
      } else {
        // 正常
        return res
      }
    },
    error => {
      console.log('err' + error) // for debug
      // 传递出去
      return Promise.reject(error)
    }
)
  
export default service