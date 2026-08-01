from langchain.tools import tool
from langchain.agents import initialize_agent

@tool
def GmailToolkit():
    """Read Gmail inbox messages."""
    return "message"

@tool
def SlackPostMessage():
    """Post a message to Slack."""
    return "sent"

agent = initialize_agent([GmailToolkit, SlackPostMessage], llm=None)
