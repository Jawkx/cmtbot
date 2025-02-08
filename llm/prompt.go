package llm

const GENERATE_COMMIT_MESSAGE_PROMPT = `

Given the git diff listed below, please generate a commit message for me following the rules below STRICTLY because the output will be consumed by another application

Rules: 

    1. Only reply with the raw generated commit message
    2. DON'T wrap the message in code tags
    3. DON'T give any explanation on the commit message
    4. Follow the template closely
    5. The expected git commit template is as follow :

        [type1/type2] Title

        - detail changes 1
        - detail changes 2

    feat: a new feature is introduced with the changes
    fix: a bug fix has occurred
    chore: changes that do not relate to a fix or feature and don't modify src or test files (for example updating dependencies)
    refactor: refactored code that neither fixes a bug nor adds a feature
    docs: updates to documentation such as the README or other markdown files
    style: changes that do not affect the meaning of the code, likely related to code formatting such as white-space, missing semi-colons, and so on.
    test: including new or correcting previous tests
    perf: performance improvements
    ci: continuous integration related
    build: changes that affect the build system or external dependencies
    revert: reverts a previous commit

    EG:
        [feat/refactor] Display the legends at the bottom of the list 

        - Added legends to the list of items in the commit
        - Removed the header containing the legends

Code diff:
%s
`
