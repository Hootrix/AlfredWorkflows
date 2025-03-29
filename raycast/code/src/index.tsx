import { ActionPanel, List, Action, showToast, Toast, popToRoot, getPreferenceValues, openExtensionPreferences, Clipboard } from "@raycast/api";
import { useEffect, useState, useRef } from "react";
import { exec } from "child_process";
import { promisify } from "util";
import path from "path";

const execPromise = promisify(exec);

interface Item {
  title: string;
  subtitle: string;
  arg: string;
}

interface Preferences {
  binaryPath: string;
}

interface CommandProps {
  arguments?: {
    query?: string;
  };
}

// 防抖函数
function debounce<T extends (...args: any[]) => any>(func: T, wait: number): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null;
  
  return function(...args: Parameters<T>) {
    if (timeout) {
      clearTimeout(timeout);
    }
    
    timeout = setTimeout(() => {
      func(...args);
      timeout = null;
    }, wait);
  };
}

export default function Command(props: CommandProps) {
  const [items, setItems] = useState<Item[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchText, setSearchText] = useState(props.arguments?.query || "");

  // 使用 useRef 存储防抖后的函数
  const debouncedFetchRef = useRef<(text: string) => void>();
  
  // 初始化防抖函数
  useEffect(() => {
    async function fetchItems(text: string) {
      // 如果是空输入或只有空格，则不执行查询
      if (!text.trim()) {
        setItems([]);
        setIsLoading(false);
        return;
      }
      
      setIsLoading(true);
      try {
        const fs = require('fs');
        const preferences = getPreferenceValues<Preferences>();
        
        // 使用用户配置的路径
        const binaryPath = preferences.binaryPath;
        console.log(`使用用户配置路径: ${binaryPath}`);
        
        if (!fs.existsSync(binaryPath)) {
          console.error("错误: 找不到执行程序", binaryPath);
          // 显示错误提示
          showToast({
            style: Toast.Style.Failure,
            title: "错误",
            message: `找不到执行程序: ${binaryPath}，将打开设置页面`
          });
          
          setIsLoading(false);

          // 等待 3s
          await new Promise(resolve => setTimeout(resolve, 3 * 1000));
          //返回主页面
          popToRoot();
          // 打开扩展的偏好设置页面
          await openExtensionPreferences();
          return;
        }
        
        console.log(`使用二进制文件路径: ${binaryPath}`);
        
        // 使用传入的文本作为查询
        const query = text.trim();
        console.log(`使用查询: ${query}`);
        // 执行二进制文件并获取输出
        // 使用 shellEscape 函数处理参数中的特殊字符
        const shellEscape = (str: string) => {
          return `'${str.replace(/'/g, "'\\''")}'`;
        };
        const { stdout } = await execPromise(`${shellEscape(binaryPath)} ${shellEscape(query)}`);
        
        // 解析 JSON 结果
        const response = JSON.parse(stdout);
        setItems(response.items || []);
      } catch (error) {
        console.error(error);
        showToast({
          style: Toast.Style.Failure,
          title: "执行错误",
          message: String(error)
        });
      } finally {
        setIsLoading(false);
      }
    }

    // 创建防抖函数，延迟 300ms
    debouncedFetchRef.current = debounce(fetchItems, 300);
    
    // 如果有初始查询，立即执行
    if (props.arguments?.query) {
      fetchItems(props.arguments.query);
    } else {
      setIsLoading(false); // 如果没有初始查询，则不显示加载状态
    }
  }, [props.arguments?.query]);
  
  // 监听搜索文本变化
  useEffect(() => {
    if (debouncedFetchRef.current) {
      debouncedFetchRef.current(searchText);
    }
  }, [searchText]);

  return (
    <List isLoading={isLoading} onSearchTextChange={setSearchText} searchText={searchText} searchBarPlaceholder="输入时间戳或日期...">
      {items.map((item, index) => (
        <List.Item
          key={index}
          title={item.title}
          subtitle={item.subtitle}
          accessories={[
            { text: item.arg }
          ]}
          actions={
            <ActionPanel>
              <Action
                title="复制结果"
                onAction={() => {
                  Clipboard.copy(item.arg);
                  showToast({
                    style: Toast.Style.Success,
                    title: "已复制",
                    message: item.arg
                  });
                }}
                shortcut={{ modifiers: [], key: "return" }}
              />
              <Action.CopyToClipboard
                title="复制标题"
                content={item.title}
              />
              <Action.CopyToClipboard
                title="复制副标题"
                content={item.subtitle}
              />
            </ActionPanel>
          }
        />))}
    </List>
  );
}
