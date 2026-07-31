from openai import OpenAI


class ResponseClient:
    def run(self, create_kwargs):
        client = OpenAI()
        response = client.responses.create(**create_kwargs)
        return response
