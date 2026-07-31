from autogen import ConversableAgent
from autogen.coding import LocalCommandLineCodeExecutor


executor = LocalCommandLineCodeExecutor(work_dir=".")
agent = ConversableAgent(name="autogen-runner", code_execution_config={"executor": executor})
