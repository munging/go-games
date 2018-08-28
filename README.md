# Go Games

**`go-games`** is a sample GO application to scrape data from GitHub and Codewars based on usernames.

### Pulled Data

Data comes from GitHub 1-6 and Codewars 7-11.

1. User
2. Contributions
3. Repositories
4. Stars
5. Followers
6. Following
7. Rank
8. Honor
9. Leaderboard Position
10. Honor Percentile
11. Total Completed Kata

### Installing and Running
```
go get github.com/munging/go-games
cd $GOPATH/github.com/src/github.com/munging/go-games
go run main.go
```

Then open your browser point to http://localhost:9000
Demo application also available at [go-games.herokuapp.com](https://go-games.herokuapp.com)

