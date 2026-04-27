import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getPostList } from '../api/post';
import { Post } from '../types';

function Home() {
    const [posts, setPosts] = useState<Post[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const size = 10;

    const fetchPosts = async () => {
        try {
            const res: any = await getPostList(page, size);
            setPosts(res.data.list);
            setTotal(res.data.total);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetchPosts();
    }, [page]);

    return (
        <div>
            <h1>帖子列表</h1>
            {posts.map((post) => (
                <div key={post.id} style={{ border: '1px solid #ccc', margin: 10 }}>
                    <Link to={`/post/${post.id}`}>
                        <h3>{post.title}</h3>
                    </Link>
                    <p>{post.content.slice(0, 100)}...</p>
                </div>
            ))}
            <div>
                <button disabled={page === 1} onClick={() => setPage(page - 1)}>上一页</button>
                <span>第 {page} 页 / 共 {Math.ceil(total / size)} 页</span>
                <button disabled={page >= Math.ceil(total / size)} onClick={() => setPage(page + 1)}>下一页</button>
            </div>
        </div>
    );
}

export default Home;