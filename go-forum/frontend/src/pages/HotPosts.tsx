import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getHotPosts } from '../api/post';
import { Post } from '../types';

function HotPosts() {
    const [posts, setPosts] = useState<Post[]>([]);

    useEffect(() => {
        const load = async () => {
            try {
                const res: any = await getHotPosts(5);
                setPosts(res.data.posts);
            } catch (err) {
                console.error(err);
            }
        };
        load();
    }, []);

    return (
        <div>
            <h2>热门帖子 TOP 5</h2>
            {posts.length === 0 && <p>暂无热榜数据</p>}
            {posts.map((post) => (
                <div key={post.id} style={{ border: '1px solid #ccc', margin: 10, padding: 10 }}>
                    <Link to={`/post/${post.id}`}>
                        <h3>{post.title}</h3>
                    </Link>
                    <p>{post.content.slice(0, 100)}...</p>
                    <small>点赞数: {post.like_count ?? 0}</small>
                </div>
            ))}
        </div>
    );
}

export default HotPosts;