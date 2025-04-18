# Reddit Data Compiler

Pulls data from Reddit's API and keeps track of the most upvoted posts and user with the most posts since the app began to run.

Additionally, it exposes http endpoints to retrieve those statistics.

## How to run locally
### Env vars to set
- `REDDIT_USERNAME`: reddit username
- `REDDIT_PASSWORD`: reddit password
- `REDDIT_APPID`: reddit app id given when adding your app to your reddit account
- `REDDIT_SECRET`: reddit secret key provided when adding app to your reddit account

### http server
http server will run on `localhost:8005`, but this can be modified in `main.go` as desired.
* Method: `GET`
* Path: `/user_with_most_posts`
  Example Response:
```json
{
  "AuthorID": "t2_82y4ryff",
  "AuthorUserName": "BirthdayCute5478",
  "NumPosts": 7
}
```


* Method: `GET`
* Path: `/most_upvoted_post`
  Example Response:
```json
{
  "Title": "Guy is not amused!",
  "NumUpvotes": 17012,
  "PostID": "t3_1k1fbpy"
}
```