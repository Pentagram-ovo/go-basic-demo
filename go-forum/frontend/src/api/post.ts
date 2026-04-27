import api from './index';

// 帖子列表（分页）
export const getPostList = (page = 1, size = 10) =>
    api.get('/post/list', { params: { page, size } });

// 帖子详情
export const getPostDetail = (id: number) => api.get(`/post/${id}`);

// 发布帖子
export const createPost = (title: string, content: string) =>
    api.post('/post/set', { title, content });

// 更新帖子
export const updatePost = (id: number, title: string, content: string) =>
    api.put(`/post/update/${id}`, { title, content });

// 删除帖子
export const deletePost = (id: number) =>
    api.delete(`/post/delete/${id}`);

// 热门帖子（top N）
export const getHotPosts = (top = 5) =>
    api.get('/post/hot', { params: { top } });