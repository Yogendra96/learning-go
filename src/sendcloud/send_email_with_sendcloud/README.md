# 从Excel中读取邮箱数据，进行自动发送

### 对Excel文件的格式要求： 现在简单支持 第一张sheet， 邮箱字段要在第一列

### 使用方式

- 登入sendcloud
- 选择 邮件=> 发送相关 => 邮件模板, 到邮件管理页面
- 编辑alan这个模板，这个模板-可修改邮件标题和邮件内容
- 你只需要自己编辑“邮件标题”，和“邮件内容” 这两项，其他的已经和代码兼容配置好，不要改动，改动会导致脚本程序不可用
- 邮件内容格式字体颜色等，可以自己编辑
- 点击保存，回到邮件管理列表，可以点击预览，预览邮件内容
- 进入 alan_sendcloud_email文件夹，文件夹中的文件都不能删除
- 打开文件夹中的 email_list.xlsx Excel文件, 将需要发的邮箱粘贴进去,竖着排列，要用第一列，不要有空行， 建议先粘一个你自己的邮箱做一次测试，没问题了，就需要发的邮箱
- 打开终端
  1. 输入 " cd alan_sendcloud_email " (引号不需要) ,然后按回车就进入了这个目录。你可以直接拷贝这个命令，粘贴到终端中， 引号不需要
  2. 输入 " ./alan_sendcloud_email " (引号不需要), 然后按回车，终端有成功信息打印出，说明就开始发送了(你可以直接拷贝这个命令，粘贴到终端中， 引号不需要)。 如果这时候要中断发送， 按 control + c 能中断程序

