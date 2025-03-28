# 效率插件收集

## 一、Alfred Workflows

### 1. code.alfredworkflow

合并编码解码请求,添加unicode编码,\x十六进制字符解码，长度计算,数字求和（查看369幸运数字）
 
参考：<a href="https://github.com/willfarrell/alfred-encode-decode-workflow">encode-decode-workflow v1.8</a>


### 2. Timestamp+.alfredworkflow 
 
参考：http://jousch.com
修改添加当前时间的格式数据获取
 
时间戳转换和当前时间获取


### 3. OCR.alfredworkflow
  
参考：<a href="https://github.com/sillybun/alfred-workflow/blob/master/OCR.alfredworkflow">https://github.com/sillybun/alfred-workflow/blob/master/OCR.alfredworkflow</a>
  
需要安装`tesseract`开源库 就得配置语言
截图OCR图片信息提取 自己添加了触发关键字 `zh`,`cn`


## 二、Raycast 插件

### 1. Timestamp+（Timestamp+.alfredworkflow的迁移版本）

Raycast版本的时间戳转换插件，提供以下功能：
- 时间戳与日期时间的相互转换
- 多种格式的当前时间获取

安装与使用：
```
make raycast_timestamp
```

## 三、构建与使用

使用 Makefile 进行构建和安装：

```
# 构建所有 Alfred 工作流二进制文件
make all

# 安装并启动 Raycast 时间戳插件
make raycast_timestamp

# 查看所有可用命令
make help