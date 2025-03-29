import { ActionPanel, List, Action, showToast, Toast, popToRoot, getPreferenceValues, openExtensionPreferences, Clipboard } from "@raycast/api";
import { useEffect, useState } from "react";
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

export default function Command(props: CommandProps) {
  const [items, setItems] = useState<Item[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchText, setSearchText] = useState(props.arguments?.query || "");

  useEffect(() => {
    async function fetchItems() {
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
        
        // 如果没有输入，使用当前时间戳
        const query = searchText || "";
        
        // 执行二进制文件并获取输出
        let shell  = binaryPath
        if(query){
          shell = `"${binaryPath}" "${query}"`
        }
        const { stdout } = await execPromise(shell);
        
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

    fetchItems();
  }, [searchText]);

  return (
    <List isLoading={isLoading} onSearchTextChange={setSearchText} searchBarPlaceholder="输入时间戳或日期...">
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
