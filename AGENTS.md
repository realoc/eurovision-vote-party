# Agent Workflow

1. Confirm scope with the user by asking whether the task concerns the UI or the server.
2. Locate the next uncompleted task file in the `tasks/` directory.
3. Review the task details. If the feature is not yet implemented, implement it, preferring the latest available versions when adding dependencies.
4. Consult relevant documentation through the `context7` MCP tool as needed.
5. For UI-related validation, run Playwright tests when they provide meaningful coverage; add or update tests if required.
6. When the task is fully addressed, mark the corresponding task file as completed.
