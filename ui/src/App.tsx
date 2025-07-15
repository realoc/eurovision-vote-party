import { useState } from "react";
import { createParty } from "@/api/partyClient.ts";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

function App() {
  const [partyName, setPartyName] = useState<string>("");
  const [party, setParty] = useState<{ id: string; password: string } | null>(
    null,
  );

  const handleCreateParty = async () => {
    const result = await createParty(partyName);
    setParty(result);
  };

  return (
    <div>
      <div className="flex min-h-svh flex-col items-center justify-center">
        <Input
          id="party-name"
          size={2}
          placeholder="Enter party name ..."
          value={partyName}
          onChange={(e) => setPartyName(e.target.value)}
          onKeyUp={(e) => {
            if (e.key === "Enter") {
              (async () => {
                await handleCreateParty();
              })();
            }
          }}
        />
        <Button onClick={handleCreateParty}>Create Party</Button>
        {party && (
          <>
            <p>Party ID: {party.id}</p>
            <p>Party Password: {party.password}</p>
          </>
        )}
      </div>
    </div>
  );
}

export default App;
