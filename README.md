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
api_base = "https://openrouter.ai/api/v1/chat/completions" ## Open ai compatible api
api_key_env = "OPENROUTER_API_KEY" ## Env variable for api key
model_name = "google/gemini-2.0-flash-001" ## model name
num_of_msg = 4 ## number of message generated to pick from
### prompt 
prompt = """
Given the git diff and the full changed file's content below, please generate a commit message.  Follow these rules STRICTLY, as the output will be consumed by another application that expects a specific format: 

## Rules

1.  **Output:** Only reply with the raw generated commit message.
2.  **Formatting:** Do NOT wrap the message in code tags or provide any explanations.
3.  **Template:** Use the following template precisely:

    ```
    [type1/type2] Title

    - detail changes 1
    - detail changes 2
    ...
    ```
4. If the changes match more than one type, list all the types that apply and separate them by `/`. For example `feat/refactor`

## Commit Message Components

*   **Type:**  Choose *one or more* of the following types based on the diff.  Prioritize `feat` and `fix` if present. If multiple types apply, include all of them, separated by `/`. Use these guidelines to infer the type:
    *   `feat`: A new feature is introduced.  Look for new functions, classes, components, or significant additions to existing functionality.
    *   `fix`: A bug fix. Look for changes that resolve reported issues, handle edge cases, or correct incorrect behavior.
    *   `chore`: Changes unrelated to features or fixes, and not modifying source or test files (e.g., dependency updates).
    *   `refactor`: Code changes that neither fix bugs nor add features. Look for improvements to code structure, readability, or maintainability *without* changing functionality.
    *   `docs`: Updates to documentation (README, comments, etc.).
    *   `style`: Changes that don't affect the code's meaning (formatting, whitespace, etc.).
    *   `test`: Adding, updating, or correcting tests.
    *   `perf`: Performance improvements.  Look for changes that optimize code execution speed or resource usage.
    *   `ci`: Continuous integration related changes.
    *   `build`: Changes to the build system or external dependencies.
    *   `revert`: Reverts a previous commit.  The diff will clearly show the undoing of previous changes.

*   **Title:**
    *   Start with a verb in the imperative mood (e.g., "Add," "Fix," "Update," "Refactor").
    *   Be concise (ideally under 50 characters).
    *   Summarize the *what* of the change.  Focus on the primary purpose of the commit.
    *   Do *not* include a period at the end.

*   **Detail Changes:**
    *   Each detail should be a short, descriptive bullet point.
    *   Focus on *what* changed and, if not obvious, *why*.
    *   Use present tense.
    *   Aim for 1-3 detail points, but add more if necessary to adequately describe the changes. Prioritize the most important changes.
    *   If changes span multiple files, summarize the overall impact rather than listing every file-specific change.  If one file's changes are significantly different in nature, consider a separate detail point.

## Examples (Illustrative - Don't repeat these in your response)

    ```
    feat/refactor: Display legends at the bottom

    - Added legends to the item list.
    - Removed the header containing legends.
    ```

    ```
    style: Improve button color contrast

    - Adjusted button background color for better accessibility.
    ```

    ```
    test: Add unit tests for date validation

    - Created unit tests to verify date validation logic.
    - Covered edge cases and invalid input scenarios.
    ```

## Input Data 
"""
```
