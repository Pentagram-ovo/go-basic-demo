import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getMyComments } from '../api/comment';
import { Comment } from '../types';

function MyComments() {
    const [comments, setComments] = useState<Comment[]>([]);

    useEffect(() => {
        const load = async () => {
            try {
                const res: any = await getMyComments();
                setComments(res.data.comments);
            } catch (err) {
                console.error(err);
            }
        };
        load();
    }, []);

    return (
        <div>
            <h2>我的评论</h2>
            {comments.length === 0 && <p>暂无评论</p>}
            {comments.map((c) => (
                <div key={c.id} style={{ border: '1px solid #eee', margin: 10, padding: 10 }}>
                    <p>{c.content}</p>
                    <small>
                        评论于帖子{' '}
                        <Link to={`/post/${c.post_id}`}>#{c.post_id}</Link>
                        {' '}· {c.created_at}
                    </small>
                </div>
            ))}
        </div>
    );
}

export default MyComments;