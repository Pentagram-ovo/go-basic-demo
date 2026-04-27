import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { registerApi } from '../api/auth';

function Register() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [msg, setMsg] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setMsg('');
        try {
            await registerApi(username, password);
            setMsg('注册成功，即将跳转到登录页...');
            setTimeout(() => navigate('/login'), 1500);
        } catch (err: any) {
            if (err.response?.data?.msg) {
                setError(err.response.data.msg);
            } else {
                setError('注册失败');
            }
        }
    };

    return (
        <div>
            <h2>注册</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <input
                        placeholder="用户名"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                    />
                </div>
                <div>
                    <input
                        type="password"
                        placeholder="密码"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>
                {error && <p style={{ color: 'red' }}>{error}</p>}
                {msg && <p style={{ color: 'green' }}>{msg}</p>}
                <button type="submit">注册</button>
            </form>
        </div>
    );
}

export default Register;