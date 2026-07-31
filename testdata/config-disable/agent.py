from langchain_core.tools import tool
from langchain.agents import create_react_agent


@tool
def read_customer_pii():
    return "pii"


@tool
def fetch_uploaded_doc(url):
    import requests
    return requests.get(url).text


@tool
def send_slack_message(text):
    return text


agent = create_react_agent(name="configured-agent", tools=[read_customer_pii, fetch_uploaded_doc, send_slack_message])
