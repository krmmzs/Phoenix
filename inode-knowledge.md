# inode结构解析前置知识

## 1. EXT4 inode基础概念

### 什么是inode？
- **Index Node**的缩写，是文件系统中存储文件元数据的数据结构
- 每个文件和目录都有一个唯一的inode编号
- inode包含文件的所有信息，**除了文件名**（文件名存储在目录项中）

### inode存储的信息
```
- 文件类型和权限 (mode)
- 所有者UID和GID
- 文件大小 (size)
- 时间戳：创建时间(ctime)、修改时间(mtime)、访问时间(atime)
- 删除时间 (dtime) ← 文件恢复的关键！
- 硬链接计数 (links_count)
- 数据块指针 (block pointers)
- 扩展属性信息
```

## 2. EXT4 inode在磁盘上的布局

### inode表位置计算
```
Block Group的组成:
[Super Block] [Group Descriptors] [Block Bitmap] [Inode Bitmap] [Inode Table] [Data Blocks]

inode表偏移 = 块组起始位置 + inode表块号 * 块大小
特定inode位置 = inode表偏移 + (inode号-1) * inode大小
```

### EXT4 inode结构 (128字节基本 + 128字节扩展)
参考：Linux内核 `include/linux/ext4_fs.h` 中的 `ext4_inode` 结构

```c
struct ext4_inode {
    __le16  i_mode;         /* File mode */
    __le16  i_uid;          /* Low 16 bits of Owner Uid */
    __le32  i_size_lo;      /* Size in bytes */
    __le32  i_atime;        /* Access time */
    __le32  i_ctime;        /* Inode Change time */
    __le32  i_mtime;        /* Modification time */
    __le32  i_dtime;        /* Deletion Time ← 关键字段! */
    __le16  i_gid;          /* Low 16 bits of Group Id */
    __le16  i_links_count;  /* Links count */
    __le32  i_blocks_lo;    /* Blocks count */
    __le32  i_flags;        /* File flags */
    // ... 更多字段
    __le32  i_block[15];    /* Pointers to blocks */
    // ... 扩展字段
};
```

## 3. Go语言实现关键技术

### 为什么使用unsafe.Pointer？
```go
// 磁盘上的数据是连续的字节序列
// 我们需要将这些字节直接解释为Go结构体
diskData := []byte{...} // 从磁盘读取的原始数据
inode := (*Ext4Inode)(unsafe.Pointer(&diskData[0]))
```

**原因：**
1. **零拷贝**：避免逐字段赋值的开销
2. **字节对齐**：保持与C结构体完全相同的内存布局
3. **小端序处理**：EXT4使用小端序存储多字节数据

### 字节序处理
EXT4使用小端序(Little Endian)，Go的`encoding/binary`包提供转换：
```go
// 读取32位小端序整数
value := binary.LittleEndian.Uint32(data[offset:offset+4])
```

### 时间戳处理
EXT4时间戳是Unix时间戳(自1970年1月1日的秒数)：
```go
import "time"

// 转换EXT4时间戳到Go时间
t := time.Unix(int64(inode.Dtime), 0)
```

## 4. 文件恢复相关的关键字段

### i_dtime (删除时间)
- **值为0**：文件未被删除
- **值非0**：文件被删除的Unix时间戳
- **恢复原理**：扫描所有inode，找到dtime≠0且数据块可能完整的inode

### i_links_count (硬链接计数)
- 正常文件：≥1
- 已删除文件：通常为0
- 用于验证文件是否真的被删除

### i_blocks (数据块计数)
- 文件占用的512字节块数量
- 用于计算文件实际占用的磁盘空间
- 帮助判断文件数据是否完整

### i_block[15] (数据块指针数组)
```
i_block[0-11]:  直接块指针 (Direct blocks)
i_block[12]:    一级间接块指针 (Single indirect)
i_block[13]:    二级间接块指针 (Double indirect) 
i_block[14]:    三级间接块指针 (Triple indirect)
```

## 5. 需要处理的技术挑战

### 内存对齐问题
```go
// Go结构体需要添加padding确保与C结构体布局一致
type Ext4Inode struct {
    Mode        uint16
    Uid         uint16
    SizeLo      uint32
    // 可能需要添加 padding 字段
}
```

### 大文件支持
EXT4支持大于4GB的文件，需要组合高低位：
```go
// 组合64位文件大小
size := uint64(inode.SizeLo) | (uint64(inode.SizeHi) << 32)
```

### 错误检测
- **魔数验证**：检查inode结构的合理性
- **校验和验证**：EXT4可能包含inode校验和
- **时间戳合理性**：删除时间不能早于创建时间

## 6. 相关文档和代码参考

### Linux内核源码
- `fs/ext4/ext4.h` - EXT4数据结构定义
- `fs/ext4/inode.c` - inode操作函数
- `include/linux/ext4_fs.h` - 磁盘格式定义

### 工具参考
- `debugfs` - EXT4调试工具，可以查看inode信息
- `dumpe2fs` - 显示EXT4文件系统信息
- `e2fsck` - 文件系统检查工具的源码

### 使用示例
```bash
# 查看特定inode信息
sudo debugfs -R "stat <inode_number>" /dev/sda1

# 查看文件系统超级块信息
sudo dumpe2fs /dev/sda1 | head -20
```

## 下一步实现目标

1. 定义与EXT4完全兼容的Go inode结构体
2. 实现从磁盘读取inode的函数
3. 实现inode数据验证和解析
4. 添加已删除文件检测逻辑