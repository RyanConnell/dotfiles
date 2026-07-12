**Protocol: File Changes**
Before using `write` or `edit` for the first time in a session, you **must** propose your plan and ask "May I proceed?" in your text response.
Do not call mutation tools in the same turn as your proposal. Read-only tools (`read`, `ls`, etc.) are exempt.
You are forbidden from using `edit` or `write` on any path mentioned in a prompt until you have first executed a `bash` command (`ls` or `test -f`) to verify the exact string exists. If the verification fails you must use `find` to locate the correct file before proceeding.
Do not leave comments in the codebase.

**Protocol: Action and Tool Calling**
Avoid infinite planning loops. Do not write multiple paragraphs explaining that you are "about to run a tool" or "applying the fix now." If you have decided on a course of action and have permission (or are executing a non-mutation tool), invoke the relevant tool immediately. Keep pre-tool thinking concise.
