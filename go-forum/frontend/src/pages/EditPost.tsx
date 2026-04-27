import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getPostDetail, updatePost } from '../api/post';

function EditPost() {
    const { id } = useParams<{ id: string }>();
    const postId = Number(id);
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const load = async () => {
            try {
                const res: any = await getPostDetail(postId);
                setTitle(res.data.title);
                setContent(res.data.content);
            } catch (err) {
                setError('帖子不存在或加载失败');
            }
        };
        load();
    }, [postId]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await updatePost(postId, title, content);
            navigate(`/post/${postId}`);
        } catch (err: any) {
            if (err.response?.data?.msg) {
                setError(err.response.data.msg);
            } else {
                setError('更新失败');
            }
        }
    };

    return (
        <div>
            <h2>编辑帖子</h2>
            {error && <p style={{ color: 'red' }}>{error}</p>}
            <form onSubmit={handleSubmit}>
                <div>
                    <input
                        placeholder="标题"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        style={{ width: '100%' }}
                    />
                </div>
                <div>
          <textarea
              placeholder="内容"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={6}
              style={{ width: '100%' }}
          />
                </div>
                <button type="submit">保存修改</button>
            </form>
        </div>
    );
}

export default EditPost;