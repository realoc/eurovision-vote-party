import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";

type PartyCreateRequest = {
  party_name: string;
};

export const handlers = [
  http.post("http://localhost:8080/party/create", async ({ request }) => {
    try {
      const body: PartyCreateRequest =
        (await request.json()) as PartyCreateRequest;
      if (body.party_name === "Test Party") {
        return HttpResponse.json({
          id: "some-id",
          password: "some-password",
        });
      }
      return new HttpResponse("Not Found", { status: 404 });
    } catch (e) {
      console.error("An error occurred mocking party creation endpoint ", e);
    }
  }),
];

export const server = setupServer(...handlers);
