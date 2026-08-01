from langchain.tools import tool
from langchain.agents import initialize_agent

@tool
def GmailToolkit():
    """Read Gmail inbox messages."""
    return "message"

agent = initialize_agent([GmailToolkit], llm=None)
