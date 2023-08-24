# DumDum
*[English](/docs/README-en.md) ∙ [繁體中文](README.md) ∙ *
Dum Dum專案 用TiDB做一些無聊東西

做個網頁 可以輸入 星座 血型 找到適合自己的Nici娃娃

線上執行<http://34.80.44.135/>

首頁
![首頁](test1.png)

列出所有
![首頁](test2.png)

查詢
![首頁](test3.png)

結果
![首頁](test4.png)


[TiDB說明](https://docs.pingcap.com/zh/)

[GORM說明](https://gorm.io/zh_CN/)


[API文件](http://34.80.185.163/swagger/index.html/)

[API示範網站](https://shuming-yu.github.io/demo/dist/#/)

很不像 請[shuming-yu](https://github.com/shuming-yu)好好努力


Google雲建制

用tiup下載TiDB 用 tiup playground 部屬

導出資料 tiup dumpling -u root -P 4000 -h 127.0.0.1 --filetype sql -t 8 -o ~/ -F256MiB