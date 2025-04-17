import axios from 'axios';
import { ElMessage } from 'element-plus';
import { useTokenStore } from '/@/store/token';
import router from '/@/router';
// 创建axios实例
const service = axios.create({
    baseURL: "/api/v1",
    // baseURL: "mock", // mock
    timeout: 15000 // 请求超时时间
})

// request拦截器
service.interceptors.request.use(async config => {
    const tokenStore = useTokenStore()
    if (tokenStore.getToken && tokenStore.getToken != '') {
        config.headers['Authorization'] = tokenStore.getToken
    }
    return config
}, error => {
    // Do something with request error
    console.log(error) // for debug
    Promise.reject(error)
})

service.interceptors.response.use(
    response => {
        if(response.request.responseURL.indexOf('download') > -1){
            return Promise.resolve(response);
        }
        let res = response.data
        return Promise.resolve(res);
    },
    err => {
        console.log(err.response)
        if (err.response.status == 401 || err.response.status == 403) {
            const tokenStore = useTokenStore()
            tokenStore.setToken('')
            tokenStore.setRole('')
            tokenStore.setNickname('')
            router.push("/login")
        }
        ElMessage.error(err.response.data.msg)
        return Promise.reject(err.response.data)
    }
)


export default service