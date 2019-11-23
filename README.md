# Go Games

The Go Games app that is running on Heroku is now using a different codebase which is originally based on the code here.

**`go-games`** is a sample GO application which scrapes data from GitHub, Codewars, Codecademy and PyCheckio based on usernames and uses jQuery Tablesorter on the client side to display the data in an interactive grid.

### Calculated and Scraped Data

Data for columns 1 and 2 is calculated from the other data columns that contain the scraped data. 
Scraped data comes from GitHub in columns 4-8, Codewars in columns 9-13, Codecadmey in columns 14-16
and PyCheckio in column 17.

1. Go Games Points
2. Rank
3. GitHub User
4. Contributions
5. Repositories
6. Stars
7. Followers
8. Following
9. Codewars Rank
10. Honor
11. Leaderboard Position
12. Honor Percentile
13. Total Completed Kata
14. Codecademy Points
15. Badges
16. Day Streak
17. PyCheckio Level

### Installing and Running
```
go get github.com/munging/go-games
cd $GOPATH/src/github.com/munging/go-games
go run main.go
```

Then open your browser point to http://localhost:9000



