import api from './index';

export const toggleLike = (post_id: number) =>
    api.post('/like/action', { post_id });

export const getLikeStatus = (post_id: number) =>
    api.get(`/like/status/${post_id}`);

export const getLikeCount = (post_id: number) =>
    api.get(`/like/count/${post_id}`);