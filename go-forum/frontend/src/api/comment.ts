import api from './index';

export const getCommentsByPostId = (postId: number, page = 1, size = 10) =>
    api.get(`/comment/${postId}`, { params: { page, size } });

export const createComment = (post_id: number, content: string) =>
    api.post('/comment/set', { post_id, content });

export const getMyComments = () =>
    api.get('/comment/user');