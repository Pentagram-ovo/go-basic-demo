一、代码核心知识点
表格
知识点分类	具体内容
JWT 核心技术	1. 自定义 Claims 结构体（嵌入 jwt.RegisteredClaims 标准声明扩展业务字段）
2. HS256 对称加密算法实现 Token 签名与验签
3. JWT 标准字段配置（过期时间 ExpiresAt、签发者 Issuer、签发时间 IssuedAt、唯一 ID ID）
4. Token 解析、算法合法性校验、有效性验证全流程
   Go 基础语法	1. 结构体定义与嵌入（组合式编程，复用标准 Claims 功能）
2. 类型断言（token.Claims.(*CustomClaims) 提取自定义业务数据）
3. 错误处理最佳实践（fmt.Errorf 包装底层错误、errors.New 定义业务错误）
4. 函数返回值类型规范（明确 string/error，避免 interface{} 隐式类型风险）
5. 全局变量初始化（包加载阶段执行 getJWTSecret 初始化密钥）
6. map 模拟数据存储（实现简易用户数据库）
   时间处理	1. time.Now() 获取系统当前时间
2. time.Duration 实现时间加减（设置 Token 2 小时过期、10 分钟刷新宽限期）
3. jwt.NewNumericDate 转换时间为 JWT 标准格式
4. time.Time.After 方法实现时间比较（判断 Token 是否过期）
   环境变量操作	1. os.Getenv 读取系统环境变量（JWT_SECRET_KEY/GO_ENV）
2. 环境差异化处理（开发环境默认密钥、生产环境强制配置密钥）
   安全校验机制	1. 签名算法精准校验（先验证 HMAC 类型，再验证 HS256 子算法，防算法混淆攻击）
2. 用户名 / 密码合法性校验（模拟用户身份认证）
3. Token 双重有效性验证（类型断言 + token.Valid 字段）
   二、代码实现的功能
   表格
   功能模块	具体实现
   用户登录认证	1. 接收用户名和密码参数，基于模拟用户数据库校验身份合法性
2. 为合法用户生成唯一 UserID（基于用户名长度），分配角色（admin/user）
3. 构建包含用户信息、标准字段的 CustomClaims 结构体
4. 使用 HS256 算法签名生成 JWT Token，返回给用户
   Token 验证解析	1. 接收前端传入的 Token 字符串，调用 jwt.ParseWithClaims 解析
2. 校验 Token 签名算法（仅允许 HMAC/HS256）
3. 验证 Token 整体有效性（签名、过期时间等）
4. 解析 Token 中的用户信息（UserID、Username、Role）并返回
   Token 刷新机制	1. 解析旧 Token（提取用户核心信息，暂不校验过期状态）
2. 判断旧 Token 是否在 10 分钟宽限期内（过期≤10 分钟允许刷新）
3. 复用用户信息生成新 Token，延长过期时间至当前时间 + 2 小时
4. 签名并返回新 Token，避免用户频繁登录
   环境适配处理	1. 开发环境（GO_ENV≠production）：JWT 密钥为空时使用默认值 dev_default_secret_987654321
2. 生产环境（GO_ENV=production）：JWT 密钥为空时直接 panic，强制配置环境变量保证安全性
   完整测试用例	1. 正常登录测试：验证合法账号密码生成 Token
2. Token 验证测试：解析 Token 并输出用户信息
3. Token 刷新测试：基于旧 Token 生成新 Token
4. 异常场景测试：错误密码登录验证失败逻辑