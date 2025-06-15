# Phoenix 项目状态记录

## 项目暂停时间
2025-06-15

## 项目概述
Phoenix 是一个用 Go 语言编写的文件恢复工具，专门用于恢复被 `rm` 命令删除的文件。通过低级别分析 ext4 文件系统结构来查找已删除的 inode 并恢复相关数据块。

## 对话记录

### 初始化阶段
1. **用户请求**: 使用 `init` 命令分析代码库并创建 CLAUDE.md 文件
2. **Claude 操作**: 
   - 分析了整个代码库结构
   - 读取了 README.md, go.mod, main.go 等核心文件
   - 读取了技术文档：ext4-official-docs.md, file-recovery-analysis.md, inode-knowledge.md
   - 创建了 CLAUDE.md 文件，包含项目架构、开发命令、技术基础等信息

### 项目暂停
3. **用户决定**: 暂时停止项目去做其他事情
4. **用户请求**: 记录所有对话和项目进度到文档中

## 当前项目状态

### 已完成部分
1. **基础结构搭建**
   - Go 模块初始化 (go.mod)
   - 基本的 main.go 文件框架
   - 命令行参数处理
   - Root 权限检查
   - 设备路径参数处理

2. **技术文档准备**
   - ext4-official-docs.md: 详细的 ext4 文件系统官方文档总结
   - inode-knowledge.md: inode 结构和删除机制知识
   - file-recovery-analysis.md: Go 语言实现文件恢复的技术要点
   - CLAUDE.md: 项目开发指导文档

3. **项目配置**
   - 依赖管理: golang.org/x/sys v0.15.0
   - Go 版本要求: 1.21+
   - 许可证: LICENSE 文件
   - 基本的 README.md

### 未完成部分 (核心功能)
1. **FilesystemAnalyzer 结构体**
   - 需要实现 `FilesystemAnalyzer` 类型
   - 需要实现 `NewFilesystemAnalyzer(devicePath string)` 构造函数
   - 需要实现 `Close()` 方法

2. **超级块读取功能**
   - 需要实现 `ReadSuperblock()` 方法
   - 需要实现 `PrintFilesystemInfo()` 方法
   - 需要定义 ext4 超级块数据结构

3. **核心恢复逻辑**
   - inode 表扫描
   - 已删除文件识别 (通过 dtime 字段)
   - 数据块恢复
   - 文件重建逻辑

### 技术要点记录
- **超级块位置**: 偏移 1024 字节，魔数 0xEF53
- **字节序**: 小端序 (除了日志)
- **关键字段**: inode 的 dtime (删除时间) 用于识别已删除文件
- **权限要求**: 必须 root 权限访问原始块设备
- **目标设备**: 如 /dev/sda1

### 代码结构
```
Phoenix/
├── main.go                    # 主程序入口 (框架已完成)
├── go.mod                     # Go 模块文件
├── CLAUDE.md                  # 开发指导文档
├── README.md                  # 项目说明
├── ext4-official-docs.md      # ext4 技术文档
├── inode-knowledge.md         # inode 结构知识
├── file-recovery-analysis.md  # 实现分析文档
└── LICENSE                    # 许可证
```

## 重启项目时的建议

1. **首先实现 FilesystemAnalyzer**
   - 创建结构体定义
   - 实现设备文件打开和关闭
   - 实现基本的磁盘读取功能

2. **然后实现超级块读取**
   - 定义 ext4_super_block 结构体
   - 实现从偏移 1024 字节读取超级块
   - 验证魔数 0xEF53

3. **最后实现核心恢复逻辑**
   - 扫描 inode 表
   - 识别已删除文件 (dtime != 0)
   - 恢复数据块

## 开发命令回顾
```bash
# 构建
go build -o phoenix-recovery main.go

# 运行 (需要 root)
sudo ./phoenix-recovery /dev/sda1

# 依赖管理
go mod tidy

# 代码格式化
go fmt ./...
```

## 项目目标
学习 Go 语言 + 实现文件恢复功能，通过底层文件系统操作来恢复被 rm 删除的文件。

---
*记录于 2025-06-15，项目暂停时的完整状态快照*