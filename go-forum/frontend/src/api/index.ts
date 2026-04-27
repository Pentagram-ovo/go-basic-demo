import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api',  // 根据实际后端地址修改
    timeout: 10000,
});

// 请求拦截器：自动附加 token
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// 响应拦截器：统一处理错误
api.interceptors.response.use(
    (res) => res.data,           // 直接返回 data 部分（{ code, msg, data }）
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';   // 强制跳转登录页
        }
        return Promise.reject(error);
    }
);

export default api;