from langchain.tools import tool
from langchain.agents import initialize_agent

@tool
def read_inbox():
    """Read inbox messages."""
    return "ok"

agent = initialize_agent([read_inbox], llm=None)
