# 时间戳转换工具 - Raycast 插件

用于时间戳与日期时间的互相转换


## 使用方法

1. 在 Raycast 中启动插件（快捷命令：`ts`）
2. 输入时间戳或日期时间格式

## 安装方法

### 方法一：直接安装

1. 确保已安装 [Raycast](https://raycast.com/)
2. 安装插件：
   ```bash
   npm install && npm run dev
   ```

### 方法二：开发模式

1. 克隆此仓库
2. 安装依赖：`npm install`
3. 启动开发模式：`npm run dev`

## 配置说明

插件需要配置 Go 可执行文件的路径。在首次运行时，会提示你设置二进制文件路径。

二进制文件构建：`cmd/timestamp_plus/main.go`


## 技术实现

- 前端：React + TypeScript + Raycast API
- 后端：Go 可执行文件处理时间转换逻辑

## 许可证

MIT
