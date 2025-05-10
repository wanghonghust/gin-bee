# gin-bee

#### 介绍
后台管理平台，包含用户管理、角色管理、权限管理,异步任务、菜单管理，以及工具箱等功能。
#### 演示网站
[bee](https://bee-front-beta.vercel.app/)
#### 软件架构
软件架构说明:
golang+gin+gorm+mysql/sqlite
#### 使用说明

1. 默认数据库为sqlite，若使用mysql作为数据库,在config/config.yaml中修改相关配置，并将config/config.go中52,53行取消注释，注释54行。
2. 初始化数据库（建表）:
    ```shell
   go run main.go init
   ```
3. 创建超级用户：
    ```shell
   go run main.go createsuperuser
   ```
   按提示输入账号密码即可
4. 菜单数据初始化
   ```shell
   sqlite3 dsn.db < ./sql/bee_menus.sql
   ```
5. 运行：
    ```shell
   go run main.go server
   ```
#### 相关链接
   前端项目地址：https://gitee.com/wang-honghust/vite-ts-bee