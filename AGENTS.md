# Agent Workflow

1. Derive whether the task concerns the UI, the server, or both by interpreting the user's command, and state that scope explicitly.
2. Locate the next uncompleted task file in the `tasks/` directory.
3. Confirm with the user which task file will be worked on before making changes.
4. Review the task details. If the feature is not yet implemented, implement it, preferring the latest available versions when adding dependencies.
5. Consult relevant documentation through the `context7` MCP tool as needed.
6. For UI-related validation, run Playwright tests when they provide meaningful coverage; add or update tests if required.
7. When the task is fully addressed, mark the corresponding task file as completed.

## Commit Messages
- Prefix every commit message with its scope: `ui`, `server`, or `all`.
- Follow the prefix with an appropriate Gitmoji, then the message headline.
- Provide a descriptive message body that expands on the headline.
