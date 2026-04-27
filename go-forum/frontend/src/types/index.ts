export interface User {
    id: number;
    username: string;
}

export interface Post {
    id: number;
    title: string;
    content: string;
    user_id: number;
    like_count: number;
    created_at: string;
    updated_at: string;
}

export interface Comment {
    id: number;
    user_id: number;
    post_id: number;
    content: string;
    created_at: string;
    updated_at: string;
}