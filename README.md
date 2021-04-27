gorm-generator
===

这是 [gorm.io/gorm](https://github.com/go-gorm/gorm) 的代码生成器，从项目 https://github.com/BigKuCha/model-generator 中变更而来

依赖
===
- [jennifer](https://github.com/dave/jennifer) --go的代码生成器

安装
===

```
# 使用远程源安装
go get https://github.com/ycoe/gorm-generator

# 使用本地源安装
git clone git@github.com:ycoe/gorm-generator.git
cd model-generator
go get .

# 成功的话，进入 $GOPATH/bin/ 可以看到有文件生成：gorm-generator
```

使用
===

```
# 首先确保你的GOPATH有配置
$GOPATH/bin/gorm-generator -u=root -p=(pwd of your mysql) -d=database -t=table -dir=./model -dd=./dao -appid=myapp -tp=finance_
```

参数
===
<table>
 <tr>
    <th>Flag</th>
    <th>Rule</th>
    <th>使用</th>
  </tr>
  <tr>
    <td>username, u</td>
    <td>非必填，默认：root</td>
    <td>数据库账号</td>
  </tr>
  
  <tr>
    <td>password, p</td>
    <td>必填，默认为null</td>
    <td>数据库密码</td>
  </tr>
  
  <tr>
    <td>database, d</td>
    <td>必填</td>
    <td>库名s</td>
  </tr>
  <tr>
    <td>table, t</td>
    <td>非必填，默认值：ALL，表示所有表</td>
    <td>需要创建的表名，多个使用半角逗号分隔。ALL时，会生成dao.go</td>
  </tr>
  
  <tr>
    <td>dir</td>
    <td>非必填，默认值：model</td>
    <td>model文件的存储路径</td>
  </tr>
  
  <tr>
    <td>daodir, dd</td>
    <td>必填</td>
    <td>dao文件的存储目录</td>
  </tr>
  
  <tr>
    <td>appid</td>
    <td>必填</td>
    <td>appId，用于生成引用路径</td>
  </tr>
  
  <tr>
    <td>tablePrefix, tp</td>
    <td>可选</td>
    <td>表前缀</td>
  </tr>
  
  <tr>
    <td>dp</td>
    <td>可选</td>
    <td>dao.go的包路径,默认值：gitee.com/inngke/proto/common</td>
  </tr>
</table>
