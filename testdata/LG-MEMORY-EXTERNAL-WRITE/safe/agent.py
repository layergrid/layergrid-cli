from langchain.agents import initialize_agent

agent = initialize_agent([], llm=None, memory=PineconeMemory(index="customer-memory", read_only=True))
