import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getPostDetail } from '../api/post';
import { getCommentsByPostId, createComment } from '../api/comment';
import { toggleLike, getLikeCount, getLikeStatus } from '../api/like';
import { Post, Comment } from '../types';
import { useAuth } from '../context/AuthContext';

function PostDetail() {
    const { id } = useParams<{ id: string }>();
    const postId = Number(id);
    const { isAuthenticated } = useAuth();
    const [post, setPost] = useState<Post | null>(null);
    const [comments, setComments] = useState<Comment[]>([]);
    const [commentPage, setCommentPage] = useState(1);
    const [newComment, setNewComment] = useState('');
    const [likeCount, setLikeCount] = useState(0);
    const [liked, setLiked] = useState(false);

    const loadPost = async () => {
        const res: any = await getPostDetail(postId);
        setPost(res.data);
    };

    const loadComments = async () => {
        const res: any = await getCommentsByPostId(postId, commentPage);
        setComments(res.data.comments);
    };

    const loadLike = async () => {
        const countRes: any = await getLikeCount(postId);
        setLikeCount(countRes.data.count);
        if (isAuthenticated) {
            const statusRes: any = await getLikeStatus(postId);
            setLiked(statusRes.data === '用户已点赞！');
        }
    };

    useEffect(() => {
        loadPost();
        loadComments();
        loadLike();
    }, [id, commentPage]);

    const handleToggleLike = async () => {
        if (!isAuthenticated) {
            alert('请先登录');
            return;
        }
        await toggleLike(postId);
        loadLike(); // 刷新点赞状态和数量
    };

    const handleAddComment = async () => {
        if (!isAuthenticated) {
            alert('请先登录');
            return;
        }
        await createComment(postId, newComment);
        setNewComment('');
        loadComments(); // 刷新评论列表
    };

    if (!post) return <div>加载中...</div>;

    return (
        <div>
            <h1>{post.title}</h1>
            <p>{post.content}</p>
            <div>
                <button onClick={handleToggleLike}>
                    {liked ? '❤️ 取消点赞' : '🤍 点赞'} ({likeCount})
                </button>
            </div>

            <h3>评论</h3>
            {comments.map((c) => (
                <div key={c.id} style={{ borderBottom: '1px solid #eee' }}>
                    <p>{c.content}</p>
                    <small>用户 {c.user_id} · {c.created_at}</small>
                </div>
            ))}
            <div>
                {/* 评论分页逻辑同帖子列表 */}
                <button onClick={() => setCommentPage(commentPage - 1)} disabled={commentPage === 1}>上一页</button>
                <button onClick={() => setCommentPage(commentPage + 1)}>下一页</button>
            </div>

            {isAuthenticated && (
                <div>
                    <textarea value={newComment} onChange={e => setNewComment(e.target.value)} />
                    <button onClick={handleAddComment}>发表评论</button>
                </div>
            )}
        </div>
    );
}

export default PostDetail;