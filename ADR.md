## Architecture Decision Record

_This document should detail any notable decisions that were made during service implementation and the
justification behind them. Information around structure, general difficulties encountered, 3rd party
libraries used and possible future improvements are examples of what would be relevant._

### General Information

- The algorithms in `updateScoreAndPublish` and `updateWinnerAndPublish` functions were initially running on linear time (`O(n)` complexity). This would not be ideal in production since we'd probably have much more data (e.g. hundreds of tournaments and thousands of teams). I modified the code to sort the viewmodel's fixtures/teams on initialization (`O(n log n)` complexity) and then use binary search (`O(log n)` complexity) to find and update it faster.

- Added a `/livedata` endpoint to the API that returns the current state of the viewmodel on demand. This endpoint could be used by the frontend to update the UI on demand.

### Possible Improvements

- Publish the live data to Redis so it can be accessed by other services as well. The publisher would also trigger a Redis Pub/Sub message to notify subscribed consumers about updated data. The live server would listen to this channel and refresh the ViewModel based on these changes. This means the ViewModel would be updated in real time using Redis data instead of storing it in application memory.

- Organize the endpoints so they're all declared in just one place.

- Increase unit test coverage. Add integration tests.
