# EXT4官方文档总结

## 官方文档来源
- **Linux内核官方文档**: https://www.kernel.org/doc/html/latest/filesystems/ext4/
- **主要章节**:
  - `overview.html` - 文件系统概述
  - `dynamic.html` - 动态结构（inode、目录等）
  - `globals.html` - 全局结构（超级块等）

## 1. EXT4文件系统整体结构

### 磁盘布局
```
[1024字节填充] [Block Group 0] [Block Group 1] ... [Block Group N]
```

### Block Group内部结构
```
[Superblock] [Group Descriptors] [Reserved GDT] [Data Block Bitmap] 
[Inode Bitmap] [Inode Table] [Data Blocks]
```

### 关键特性
- **字节序**: 小端序(Little-endian)，除了日志使用大端序
- **块大小**: 1-64 KiB，通常为4 KiB
- **Block Group大小**: 通常为128 MiB
- **魔数**: `0xEF53` 用于识别EXT4文件系统

## 2. 超级块(Superblock)结构

### 位置和大小
- **偏移**: 1024字节
- **大小**: 1024字节
- **魔数**: 0xEF53

### 关键字段 (基于官方文档)
```c
struct ext4_super_block {
    __le32  s_inodes_count;         /* Inodes count */
    __le32  s_blocks_count_lo;      /* Blocks count */
    __le32  s_r_blocks_count_lo;    /* Reserved blocks count */
    __le32  s_free_blocks_count_lo; /* Free blocks count */
    __le32  s_free_inodes_count;    /* Free inodes count */
    __le32  s_first_data_block;     /* First Data Block */
    __le32  s_log_block_size;       /* Block size */
    __le32  s_log_cluster_size;     /* Cluster size */
    __le32  s_blocks_per_group;     /* # Blocks per group */
    __le32  s_clusters_per_group;   /* # Clusters per group */
    __le32  s_inodes_per_group;     /* # Inodes per group */
    __le32  s_mtime;                /* Mount time */
    __le32  s_wtime;                /* Write time */
    __le16  s_mnt_count;            /* Mount count */
    __le16  s_max_mnt_count;        /* Maximal mount count */
    __le16  s_magic;                /* Magic signature */
    __le16  s_state;                /* File system state */
    // ... 更多字段
};
```

### 派生计算
```go
// 块大小计算
blockSize := 1024 << superblock.LogBlockSize

// 块组数量计算  
groupCount := (totalBlocks + blocksPerGroup - 1) / blocksPerGroup
```

## 3. inode结构详解

### inode基本信息
- **默认大小**: 256字节 (现代EXT4)
- **最小大小**: 128字节
- **位置计算**: `inode表偏移 + (inode号-1) * inode大小`

### 关键字段 (文件恢复相关)
```c
struct ext4_inode {
    __le16  i_mode;         /* File mode */
    __le16  i_uid;          /* Low 16 bits of Owner Uid */
    __le32  i_size_lo;      /* Size in bytes */
    __le32  i_atime;        /* Access time */
    __le32  i_ctime;        /* Inode Change time */
    __le32  i_mtime;        /* Modification time */
    __le32  i_dtime;        /* Deletion Time ← 关键! */
    __le16  i_gid;          /* Low 16 bits of Group Id */
    __le16  i_links_count;  /* Links count */
    __le32  i_blocks_lo;    /* Blocks count */
    __le32  i_flags;        /* File flags */
    // ...
    __le32  i_block[EXT4_N_BLOCKS]; /* Pointers to blocks */
    // ... 扩展字段
};
```

### 删除时间字段 (i_dtime) - 恢复的核心
**官方文档说明**:
- **偏移**: 0x14 (20字节)
- **含义**: "Deletion Time, in seconds since the epoch"
- **值为0**: 文件未删除
- **值非0**: 文件删除的Unix时间戳

**特殊情况** (orphan_file特性):
- 如果文件系统没有orphan_file特性，此字段可能被重载
- 用于构建孤儿inode链表

## 4. 文件恢复技术要点

### 恢复原理
1. **扫描所有inode**: 遍历每个块组的inode表
2. **检查删除时间**: `i_dtime != 0` 表示已删除文件
3. **验证完整性**: 检查数据块指针和链接计数
4. **重建文件**: 通过块指针读取数据块

### 关键验证点
```go
// 检查是否为已删除文件
if inode.Dtime != 0 && inode.LinksCount == 0 {
    // 可能的已删除文件
}

// 验证inode合理性
if inode.Mode != 0 && inode.Size > 0 {
    // inode看起来有效
}
```

### 数据块访问
- **直接块**: `i_block[0-11]` 直接指向数据块
- **间接块**: `i_block[12-14]` 指向间接块指针表
- **大文件**: 需要遍历间接块结构

## 5. 实现时的技术考虑

### 字节序处理
```go
import "encoding/binary"

// EXT4使用小端序
value := binary.LittleEndian.Uint32(data)
```

### 时间戳转换
```go
import "time"

// Unix时间戳转换
deleteTime := time.Unix(int64(inode.Dtime), 0)
```

### 内存对齐
```go
// 确保Go结构体与C结构体内存布局一致
type Ext4Inode struct {
    Mode     uint16
    Uid      uint16
    SizeLo   uint32
    // 可能需要padding字段
} // 总大小必须匹配EXT4规范
```

## 6. 调试工具参考

### 系统工具
```bash
# 查看文件系统信息
sudo dumpe2fs /dev/sda1

# 调试特定inode
sudo debugfs -R "stat <inode_number>" /dev/sda1

# 查看超级块
sudo debugfs -R "stats" /dev/sda1
```

### 官方源码参考
- **内核源码**: `fs/ext4/` 目录
- **关键文件**: 
  - `ext4.h` - 数据结构定义
  - `super.c` - 超级块操作
  - `inode.c` - inode操作

这些官方文档为我们的Go实现提供了准确的技术规范。