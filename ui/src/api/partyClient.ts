export async function createParty(
  partyName: string,
): Promise<{ id: string; password: string }> {
  const response = await fetch("http://localhost:8080/party/create", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ party_name: partyName }),
  });

  if (!response.ok) {
    throw new Error("Failed to create party");
  }

  return response.json();
}
