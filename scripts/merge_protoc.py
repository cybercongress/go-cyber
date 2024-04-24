#!/usr/bin/env python3

# Call this from the ./scripts/protoc_swagger_openapi_gen.sh script

# merged protoc definitions together into 1 JSON file without duplicate keys
# this is done AFTER swagger-merger has been run, merging the multiple name-#.json files into 1.

import json
import os
import random
import string

current_dir = os.path.dirname(os.path.realpath(__file__))
project_root = os.path.dirname(current_dir)

all_dir = os.path.join(project_root, "tmp-swagger-gen", "_all")

# get the go.mod file Version
version = ""
with open(os.path.join(project_root, "go.mod"), "r") as f:
    for line in f.readlines():
        if line.startswith("module"):
            version = line.split("/")[-1].strip()
            break

if not version:
    print("Could not find version in go.mod")
    exit(1)

# What we will save when all combined
output: dict
output = {
    "swagger": "2.0",
    "info": {"title": "Bostrom network", "version": version},
    "consumes": ["application/json"],
    "produces": ["application/json"],
    "paths": {},
    "definitions": {},
}

# Combine all individual files calls into 1 massive file.
for file in os.listdir(all_dir):
    if not file.endswith(".json"):
        continue

    # read file all_dir / file
    with open(os.path.join(all_dir, file), "r") as f:
        data = json.load(f)

    for key in data["paths"]:
        output["paths"][key] = data["paths"][key]

    for key in data["definitions"]:
        output["definitions"][key] = data["definitions"][key]


# loop through all paths, then alter any keys which are "operationId" to be a random string of 20 characters
# this is done to avoid duplicate keys in the final output (which opens 2 tabs in swagger-ui)
# current-random
for path in output["paths"]:
    for method in output["paths"][path]:
        if "operationId" in output["paths"][path][method]:
            output["paths"][path][method][
                "operationId"
            ] = f'{output["paths"][path][method]["operationId"]}_' + "".join(
                random.choices(string.ascii_uppercase + string.digits, k=5)
            )


# save output into 1 big json file
with open(
    os.path.join(project_root, "tmp-swagger-gen", "_all", "FINAL.json"), "w"
) as f:
    json.dump(output, f, indent=2)
