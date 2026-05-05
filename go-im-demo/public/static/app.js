// 全局变量
let ws = null;
let currentUser = "";
let currentUserID = 0;
let token = localStorage.getItem("im_token") || "";
let heartbeatInterval = null;
let onlineRefreshInterval = null;

function formatTime(timestamp) {
    const date = new Date(timestamp * 1000);
    return date.toLocaleTimeString();
}

// 显示消息（根据发送者决定气泡样式）
function appendChatMessage(from, content, time) {
    const div = document.getElementById("messages");
    if (!div) return;
    const bubble = document.createElement("div");
    const isMine = from === currentUser;
    bubble.className = "message-bubble " + (isMine ? "mine" : "other");
    const sender = document.createElement("div");
    sender.className = "sender";
    sender.textContent = isMine ? "你" : from;
    const contentDiv = document.createElement("div");
    contentDiv.className = "content";
    contentDiv.textContent = content;
    const timeDiv = document.createElement("div");
    timeDiv.className = "time";
    timeDiv.textContent = formatTime(time);
    bubble.appendChild(sender);
    bubble.appendChild(contentDiv);
    bubble.appendChild(timeDiv);
    div.appendChild(bubble);
    div.scrollTop = div.scrollHeight;
}

function appendSystemMsg(text) {
    const div = document.getElementById("messages");
    if (!div) return;
    const p = document.createElement("p");
    p.style.color = "#999";
    p.style.fontSize = "12px";
    p.style.textAlign = "center";
    p.textContent = text;
    div.appendChild(p);
    div.scrollTop = div.scrollHeight;
}

// 登录逻辑
function switchToLogin() {
    document.getElementById("authTitle").textContent = "登录";
    document.getElementById("btnLogin").classList.remove("hidden");
    document.getElementById("btnToRegister").classList.remove("hidden");
    document.getElementById("btnDoRegister").classList.add("hidden");
    document.getElementById("btnBackToLogin").classList.add("hidden");
    document.getElementById("authPassword").value = "";
    document.getElementById("authMsg").textContent = "";
}

function switchToRegister() {
    document.getElementById("authTitle").textContent = "注册";
    document.getElementById("btnLogin").classList.add("hidden");
    document.getElementById("btnToRegister").classList.add("hidden");
    document.getElementById("btnDoRegister").classList.remove("hidden");
    document.getElementById("btnBackToLogin").classList.remove("hidden");
    document.getElementById("authPassword").value = "";
    document.getElementById("authMsg").textContent = "";
}

async function register() {
    const username = document.getElementById("authUsername").value.trim();
    const password = document.getElementById("authPassword").value.trim();
    if (!username || !password) return;
    const res = await fetch("/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });
    const data = await res.json();
    if (res.ok) {
        document.getElementById("authMsg").textContent = "注册成功，请登录";
        switchToLogin();
    } else {
        document.getElementById("authMsg").textContent = data.error;
    }
}

async function login() {
    const username = document.getElementById("authUsername").value.trim();
    const password = document.getElementById("authPassword").value.trim();
    if (!username || !password) return;
    const res = await fetch("/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });
    const data = await res.json();
    if (res.ok && data.token) {
        localStorage.setItem("im_token", data.token);
        token = data.token;
        currentUser = username;
        try {
            const payload = JSON.parse(atob(data.token.split('.')[1]));
            currentUserID = payload.user_id;
        } catch (e) {}
        window.location.href = "/";
    } else {
        document.getElementById("authMsg").textContent = data.error || "登录失败";
    }
}

// WebSocket 连接
function connectWebSocket() {
    if (!token) return;
    ws = new WebSocket("ws://localhost:8080/ws?token=" + token);
    ws.onopen = function() {
        appendSystemMsg("已连接到服务器");
        // 心跳
        heartbeatInterval = setInterval(() => {
            if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({ type: "ping" }));
            }
        }, 20000);
        // 定期刷新在线列表
        updateOnlineUsers();
        onlineRefreshInterval = setInterval(updateOnlineUsers, 10000);
    };
    ws.onmessage = function(event) {
        try {
            const msg = JSON.parse(event.data);
            if (msg.type === "pong") return;
            // 在线列表数组
            if (Array.isArray(msg)) {
                renderOnlineUsers(msg);
                return;
            }
            // 聊天消息
            if (msg.from && msg.content) {
                appendChatMessage(msg.from, msg.content, msg.time || Math.floor(Date.now()/1000));
                return;
            }
            // 系统消息
            appendSystemMsg(event.data);
        } catch (e) {
            appendSystemMsg(event.data);
        }
    };
    ws.onclose = function() {
        appendSystemMsg("连接已关闭");
        ws = null;
        clearInterval(heartbeatInterval);
        clearInterval(onlineRefreshInterval);
    };
    ws.onerror = function() {
        appendSystemMsg("连接错误");
    };
}

function updateOnlineUsers() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ to: "server", content: "/online" }));
    }
}

function renderOnlineUsers(users) {
    const ul = document.getElementById("onlineUsers");
    if (!ul) return;
    ul.innerHTML = users.map(u => `<li>🟢 ${u}</li>`).join("");
    const countEl = document.getElementById("onlineCount");
    if (countEl) countEl.textContent = `${users.length} 人在线`;
}

function disconnect() {
    if (ws) {
        ws.close();
        ws = null;
    }
    clearInterval(heartbeatInterval);
    clearInterval(onlineRefreshInterval);
    localStorage.removeItem("im_token");
    window.location.href = "/login.html";
}

function sendMsg() {
    const content = document.getElementById("msgInput").value.trim();
    const target = document.getElementById("targetUser").value.trim();
    if (!content || !ws || !target) return;
    if (target === currentUser) {
        appendSystemMsg("不能给自己发消息");
        document.getElementById("msgInput").value = "";
        return;
    }
    const now = Math.floor(Date.now() / 1000);
    ws.send(JSON.stringify({ to: target, content: content, time: now }));
    appendChatMessage(currentUser, content, now);
    document.getElementById("msgInput").value = "";
}

async function loadHistory() {
    const peer = document.getElementById("targetUser").value.trim();
    if (!peer) return alert("请输入对方用户名");
    const res = await fetch(`/messages?token=${token}&peer_name=${peer}&page=1&size=30`);
    const data = await res.json();
    if (!res.ok) return alert(data.error || "加载历史失败");
    document.getElementById("messages").innerHTML = "";
    const messages = data.messages.reverse();
    messages.forEach(msg => {
        const sender = msg.from_id === currentUserID ? currentUser : peer;
        appendChatMessage(sender, msg.content, msg.created_at_unix || new Date(msg.created_at).getTime()/1000);
    });
}

// 页面初始化
window.onload = function() {
    if (document.getElementById("authArea")) {
        switchToLogin();
    }
    if (document.getElementById("chatArea") || document.getElementById("messages")) {
        if (!token) {
            window.location.href = "/login.html";
            return;
        }
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            currentUser = payload.user_name;
            currentUserID = payload.user_id;
        } catch (e) {
            currentUser = "unknown";
        }
        document.getElementById("currentUserDisplay").textContent = currentUser;
        connectWebSocket();
    }
};