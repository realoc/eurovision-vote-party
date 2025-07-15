import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import * as partyClient from "@/api/partyClient.ts";
import App from "./App";

vi.mock("@/api/partyClient.ts", () => ({
  createParty: vi.fn(),
}));

describe("App", () => {
  it("should create a party and display the id and password", async () => {
    const mockParty = { id: "test-id", password: "test-password" };
    vi.mocked(partyClient.createParty).mockResolvedValue(mockParty);

    render(<App />);

    const input = screen.getByPlaceholderText("Enter party name ...");
    const button = screen.getByText("Create Party");

    fireEvent.change(input, { target: { value: "Test Party" } });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText("Party ID: test-id")).toBeTruthy();
      expect(screen.getByText("Party Password: test-password")).toBeTruthy();

      expect(screen.getByText("Party ID: test-id")).toBeInstanceOf(HTMLElement);
      expect(screen.getByText("Party Password: test-password")).toBeInstanceOf(
        HTMLElement,
      );
    });
  });
});
