# API

# REST

| Path                                  | Parameters    | Description    |
| --------------------------------------| ------------- | ------------------ |
| /energy/parameters                    |                     | get module params                          |
| /energy/{%s}/source_routes            | {sourceAddress}     | get routes from source                     |
| /energy/{%s}/destination_routes       | {destinationAddress}| get routes to destination                  |
| /energy/{%s}/source_routed_energy     | {sourceAddress}     | get routed total energy from source        |
| /energy/{%s}/destination_routed_energy| {destinationAddress}| get routed total energy to destination     |
| /energy/route/{%s}/{%s}               | {sourceAddress} and {destinationAddress} | get given route information |
| /energy/routes                        | page and limit      | get all routes                             |

# GRPC

| Path                                  | Parameters    | Description    |
| --------------------------------------| ------------- | ------------------ |
| /cyber/energy/v1beta1/energy/params                   |                     | get module params                          |
| /cyber/energy/v1beta1/energy/source_routes            | {sourceAddress}     | get routes from source                     |
| /cyber/energy/v1beta1/energy/destination_routes       | {destinationAddress}| get routes to destination                  |
| /cyber/energy/v1beta1/energy/destination_routed_energy| {sourceAddress}     | get routed total energy from source        |
| /cyber/energy/v1beta1/energy/source_routed_energy     | {destinationAddress}| get routed total energy to destination     |
| /cyber/energy/v1beta1/energy/route                    | {sourceAddress} and {destinationAddress} | get given route information |
| /cyber/energy/v1beta1/energy/routes                   | page and limit      | get all routes      