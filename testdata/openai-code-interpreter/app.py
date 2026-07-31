from openai import OpenAI


client = OpenAI()
assistant = client.beta.assistants.create(
    name="analysis-agent",
    tools=[{"type": "code_interpreter"}, {"type": "file_search"}],
)
