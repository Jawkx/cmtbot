# cmtbot

A terminal application that helps you generate and select commit messages based on your staged changes.

## Functionality

cmtbot performs the following actions:

1.  Reads the diff of your staged files.
2.  Generates multiple commit message options based on the diff using LLM
3.  Allows you to choose the desired commit message.
4.  Commits the changes using the selected message.

## Setup

In `~/.config/cmtbot/cmtbot.toml` create a toml file

``` toml

api_base = "https://openrouter.ai/api/v1/chat/completions" # Open Ai compatible api root 
api_key_env = "OPENROUTER_API_KEY" # Access token that will be called in the Bearer Token
model_name = "google/gemini-2.0-flash-001" # model name
num_of_msg = 4 # number of message generated for you to choose on
```
