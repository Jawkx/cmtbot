# cmtbot

A terminal apps that helps generate commit using AI

<img src="https://raw.githubusercontent.com/Jawkx/cmtbot/refs/heads/master/DOCS/demo.gif" width="100%" alt="cmtbot demo">

## Installation

1. Download the executable from [Release page](https://github.com/Jawkx/cmtbot/releases)
2. Make it executable 

``` bash
chmod +x <downloaded binary>
```
3. (Optional but highly recommended) rename it to `cmtbot`
4. Move it into place where you usually store your binary, EG: `usr/local/bin`, or anywhere directory that is in path, [How to correctly add a path to PATH?](https://unix.stackexchange.com/questions/26047/how-to-correctly-add-a-path-to-path) 
### But I'm lazy, what now?

1. Copy this and paste it in your terminal (remember to change the version in `DOWNLOAD_URL` eg: `https://github.com/Jawkx/cmtbot/releases/download/v0.1.0/cmtbot-darwin-amd64`)

``` bash
DOWNLOAD_URL="https://github.com/Jawkx/cmtbot/releases/download/<change to version you want>/cmtbot-darwin-amd64"
EXECUTABLE_NAME="cmtbot"
INSTALL_DIR="/usr/local/bin"

echo "Downloading cmtbot..."
curl -L "$DOWNLOAD_URL" -o "$EXECUTABLE_NAME"

echo "Making cmtbot executable..."
chmod +x "$EXECUTABLE_NAME"

echo "Installing cmtbot to $INSTALL_DIR..."
sudo mkdir -p "$INSTALL_DIR" 2>/dev/null # Creates dir if it doesnt exist, suppress errors if it does.
sudo mv "$EXECUTABLE_NAME" "$INSTALL_DIR/$EXECUTABLE_NAME"

echo "cmtbot installed successfully!"
```

## Setup

1. To use this app, you first need to choose a LLM provider that have OpenAi compatible API (which is almost all of them). After getting a API key there, find a way to [load it in your environment variable](https://www.perplexity.ai/search/zsh-how-to-add-key-to-environm-escTutp0SJOINLwFpsu4VQ)
2. Create a directory `~/.config/cmtbot/`
3. In the directory create a file called `prompt.md` (or anything really), this will be your prompt used to generate the commit messages, you can check out a sample prompt file [here](https://github.com/Jawkx/cmtbot/blob/master/DOCS/EXAMPLE_CONFIG/prompt.md)
4. After that create another file in the directory called `config.toml`
5. In `config.toml` paste this and edit this base on your the environment variable key you setup in step 1 and also choose your model

``` toml
api_base = "https://api.openai.com/v1/chat/completions" # OpenAi compatible Api root
api_key_env = "OPENAI_API_KEY" # Environment variable key
model_name = "gpt-4-turbo-preview" # Model
num_of_msg = 5 # number of mesage generated for selection
prompt_filename = "prompt.md" # File name you set in step 3`
```

Can reference this [example config](https://github.com/Jawkx/cmtbot/tree/master/DOCS/EXAMPLE_CONFIG/)
