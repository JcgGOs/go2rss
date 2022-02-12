# go2rss

go2rss 是一款rss生成工具，功能类同与 [RSSHub](https://github.com/DIYgod/RSSHub),也是给奇怪的网站生成RSS.

不同于RSSHUB，你只需要在配置目录`/config` 中配置相应的抓取表达式，就可以生成一个新的rss数据源.

default.json 为公共配置， 可以被其他配置给覆盖


## Docker 

```shell
docker run -d \
 -v config2:/app/config/config2 \
 -p 8081:8081 \
 tantao700/go2rss:latest
```