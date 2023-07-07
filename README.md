# [Ask HN: Could you share your personal blog here?](https://news.ycombinator.com/item?id=36575081)

## Description

This is a collection of personal blogs from the [Ask HN: Could you share your personal blog here?](https://news.ycombinator.com/item?id=36575081) 
thread on Hacker News, prepared as OPML for easy import into your favorite RSS reader.

## Usage

Download [list.opml](list.opml) and import it into your favorite RSS reader.

When building this list, I have ignored any user with less or equal to 100 karma, which means I might have missed some
interesting blogs, but at the same time I wanted to ignore spam or throwaway accounts.

The list is sorted by the user karma on Hacker News, so the first blogs are from users with the highest karma.

You can modify the list in your editor to include only the top 10 or 100 blogs, or to remove some blogs you are not interested in.

Not from all comments I was able to extract a blog URL, so the list is not complete. I just parse the correct recocognized URLs
from comments.

Not all blogs have RSS feeds, or the RSS feeds aren't included in the `<link rel="alternate" type="application/rss+xml" href="...">`
or `<link rel="alternate" type="application/atom+xml" href="...">` tag, so I might have missed some blogs.

Anyway, we got more than 600 blogs, so I think it is a good start.

You can find the output of the latest run at [console.log](console.log).

## Regenerate list

As easy as running:

```bash
go run ./main.go | tee > console.log
```

It is going to take a while, as it needs to fetch the karma for each user, and then fetch the RSS feed for each blog.

## Author

[outcoldman](https://www.outcoldman.com)

- [Twitter](https://twitter.com/outcoldman)
- [GitHub](https://github.com/outcoldman)

## LICENSE

[MIT](LICENSE)

