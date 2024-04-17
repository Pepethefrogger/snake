snake is a fast web crawler written in go. It's easy to use and really fast.
Usage:
```
./snake -u http://crawl.website -fs js,css -fc 404,500,403
-fs|--filter-string : Filters url's containing this string.
-fc|--filter-code   : Filters url's that return these status codes.
-u |--url           : Url to use espaider against.
```
