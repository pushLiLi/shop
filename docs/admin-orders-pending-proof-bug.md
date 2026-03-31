# 管理端订单筛选 Bug 记录：待审核凭证有数量但列表为空

## 背景
在管理端订单页（/admin/orders）点击快捷筛选“待审核凭证”时，页面徽标会显示待审核数量，但订单列表为空。

## 发现时间
- 2026-03-31

## 现象
- 前端请求：/api/admin/orders?page=1&limit=20&quick_filter=pending_proof
- 返回结构示例：
  - page: 1
  - limit: 20
  - pendingProofCount: 450
  - total: 0
  - totalPages: 0
  - orders: null

这说明：
- payment_proofs 表中 pending 状态的记录是存在的（pendingProofCount > 0）
- 但按照 quick_filter=pending_proof 关联订单时结果为 0

## 影响范围
- 管理端订单页“待审核凭证”快捷筛选不可用
- 依赖同类子查询写法的筛选可能存在同类风险（例如 to_ship、proof_status）

## 相关代码位置
- 后端筛选构建函数：server-go/internal/handlers/admin_order.go
  - buildOrderQuery() 内 quick_filter=pending_proof 分支
  - buildOrderQuery() 内 quick_filter=to_ship 分支
  - buildOrderQuery() 内 proof_status 分支

## 当前可见实现（问题点）
在 buildOrderQuery() 中使用了 GORM 子查询模式：
- query.Where("id IN ?", subQuery)

在当前项目配置（PrepareStmt: true）下，实际运行中出现了“计数正常但主列表为空”的不一致现象，需要进一步验证 SQL 生成与执行行为。

## 初步根因假设
- 可能是 GORM + MySQL 在 PrepareStmt 开启时，对 Where("id IN ?", subQuery) 这种子查询拼接存在兼容性/执行差异
- 由于 pendingProofCount 使用的是独立计数查询，因此能返回正确数量
- 但主订单查询依赖子查询关联，可能在 SQL 构造或参数绑定阶段出现偏差

## 建议排查步骤
1. 在相同请求参数下抓取数据库实际执行 SQL（包括变量绑定值）
2. 在数据库中分别执行：
   - pending 计数查询
   - id IN (subquery) 的订单查询
3. 对比 ORM 生成 SQL 与手写 SQL 的结果是否一致
4. 排查 PrepareStmt 开关对该查询的影响
5. 排查 to_ship、proof_status 是否存在同类症状

## 建议修复方向（供后续维护者评估）
可将子查询改为两段查询以规避 ORM 子查询不稳定因素：
- 先查询 payment_proofs 并 pluck 出 order_id 列表
- 再执行 orders 的 id IN 列表查询
- 当列表为空时显式追加 1=0，避免误查全部

该方向的优点是：
- 行为可预期，便于观测与调试
- 对当前数据量（几百到几千）足够稳定

## 验收标准
- quick_filter=pending_proof 返回 total > 0 且 orders 非空（在 pendingProofCount > 0 的前提下）
- 待审核徽标数量与可筛出的订单集合保持一致
- to_ship、proof_status 在同样数据集下返回逻辑一致

## 备注
- 本文档仅记录问题与建议路径，不在本文档中直接实施代码修复。
- 维护者后续可基于本文档独立完成修复与回归。