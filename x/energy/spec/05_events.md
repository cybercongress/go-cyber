# Events

The energy module emits the following events:

## Msg's

### MsgCreateRoute
| Type             | Attribute Key | Attribute Value    |
| ---------------- | ------------- | ------------------ |
| message          | module        | energy             |
| message          | action        | create_route       |
| create_route     | source        | {sourceAddress}    |
| create_route     | destination   | {destinationAddress} |
| create_route     | alias         | {alias}            |

### MsgEditRoute
| Type             | Attribute Key | Attribute Value    |
| ---------------- | ------------- | ------------------ |
| message          | module        | energy             |
| message          | action        | edit_route         |
| edit_route       | source        | {sourceAddress}    |
| edit_route       | destination   | {destinationAddress} |
| edit_route       | value         | {value}            |

### MsgDeleteRoute
| Type             | Attribute Key | Attribute Value    |
| ---------------- | ------------- | ------------------ |
| message          | module        | energy             |
| message          | action        | delete_route       |
| delete_route     | source        | {sourceAddress}    |
| delete_route     | destination   | {destinationAddress} |

### MsgEditRouteAlias
| Type             | Attribute Key | Attribute Value    |
| ---------------- | ------------- | ------------------ |
| message          | module        | energy             |
| message          | action        | edit_route_alias   |
| edit_route_alias | source        | {sourceAddress}    |
| edit_route_alias | destination   | {destinationAddress} |
| edit_route_alias | alias         | {alias}            |