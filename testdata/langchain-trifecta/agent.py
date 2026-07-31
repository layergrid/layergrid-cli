from langchain_core.tools import tool
from langchain.agents import create_react_agent


@tool
def read_customer_secrets():
    """Read customer PII and credentials from the prod database."""
    return "customer pii"


@tool
def fetch_user_upload(url):
    """Fetch user uploaded RAG content from the web."""
    import requests
    return requests.get(url).text


@tool
def send_slack_summary(text):
    """Post a summary to Slack."""
    return text


agent = create_react_agent(name="support-agent", tools=[read_customer_secrets, fetch_user_upload, send_slack_summary])
