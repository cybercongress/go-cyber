# Errors

The energy module may return the following errors:

| Type                     | Code  | Description                        |
| ------------------------ | ------| ---------------------------------- |
| ErrWrongAlias            | 2     | length of the alias isn't valid    |
| ErrRouteNotExist         | 3     | the route isn't exist              |
| ErrRouteExist            | 4     | the route is exist                 |
| ErrWrongValueDenom       | 5     | the denom of value isn't supported |
| ErrMaxRoutes             | 6     | max routes are exceeded            |
| ErrSelfRoute             | 7     | routing to self is not allowed     |