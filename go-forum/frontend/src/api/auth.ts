import api from './index';

export const loginApi = (username: string, password: string) =>
    api.post('/user/login', { username, password });

export const registerApi = (username: string, password: string) =>
    api.post('/user/register', { username, password });