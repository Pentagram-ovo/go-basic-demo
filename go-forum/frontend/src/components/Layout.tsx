import { Link, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function Layout() {
    const { isAuthenticated, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/');
    };

    return (
        <div style={{ maxWidth: 800, margin: '0 auto', padding: 20 }}>
            <nav style={{ marginBottom: 20, display: 'flex', gap: 15, alignItems: 'center' }}>
                <Link to="/">首页</Link>
                <Link to="/hot">热榜</Link>
                {isAuthenticated ? (
                    <>
                        <Link to="/create-post">发帖</Link>
                        <Link to="/my-comments">我的评论</Link>
                        <button onClick={handleLogout}>退出登录</button>
                    </>
                ) : (
                    <>
                        <Link to="/login">登录</Link>
                        <Link to="/register">注册</Link>
                    </>
                )}
            </nav>
            <hr />
            <Outlet />
        </div>
    );
}

export default Layout;