import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createPost } from '../api/post';

function CreatePost() {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await createPost(title, content);
            navigate('/');
        } catch (err: any) {
            if (err.response?.data?.msg) {
                setError(err.response.data.msg);
            } else {
                setError('发帖失败');
            }
        }
    };

    return (
        <div>
            <h2>发布新帖子</h2>
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
                {error && <p style={{ color: 'red' }}>{error}</p>}
                <button type="submit">发布</button>
            </form>
        </div>
    );
}

export default CreatePost;