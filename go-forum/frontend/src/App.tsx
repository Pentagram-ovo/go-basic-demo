import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import PostDetail from './pages/PostDetail';
import CreatePost from './pages/CreatePost';
import EditPost from './pages/EditPost';
import HotPosts from './pages/HotPosts';
import MyComments from './pages/MyComments';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route element={<Layout />}>
                    <Route path="/" element={<Home />} />
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/post/:id" element={<PostDetail />} />
                    <Route path="/hot" element={<HotPosts />} />

                    {/* 需要登录的路由 */}
                    <Route element={<ProtectedRoute />}>
                        <Route path="/create-post" element={<CreatePost />} />
                        <Route path="/edit-post/:id" element={<EditPost />} />
                        <Route path="/my-comments" element={<MyComments />} />
                    </Route>
                </Route>
            </Routes>
        </BrowserRouter>
    );
}

export default App;