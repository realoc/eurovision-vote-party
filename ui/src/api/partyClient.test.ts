import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";
import { server } from "./mocks/server";
import { createParty } from "./partyClient";

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe("partyClient", () => {
  it("should create a party and return id and password", async () => {
    const partyName = "Test Party";
    const result = await createParty(partyName);
    expect(result).toEqual({
      id: "some-id",
      password: "some-password",
    });
  });
});
