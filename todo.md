# 迁移说明


执行以下正则替换
## 替换 1
```
github\.com/admpub/nging/v5/application/(library|cmd|dbschema|initialize|listener|middleware|model|registry|request|response|version)(/|")
```
替换为
```
github.com/coscms/webcore/$1$2
```

## 替换 2
```
handler\.(User|URLFor|FrontendURLFor|NoticeWriter|IsBackendAdmin|OnInstalled)\(
```
替换为
```
backend.$1(
```

handler\.(Err|PagingWithLister|SendOk|SendFail|SendErr|Paging[A-Za-z]+|NewLister)\(
common.$1(

handler\.(RegisterToGroup|Register|WithRequest|IRegister)\(
route.$1(